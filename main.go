package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
    telegramAPIURL     = "https://api.telegram.org/bot"
    openAIAPIURL       = "https://api.openai.com/v1/chat/completions"
    model              = "gpt-3.5-turbo"
)

var (
    openAIAPIKey string
    telegramBotToken string
)

//**********//
// TELEGRAM //
//**********//
type Update struct {
    UpdateID int      `json:"update_id"`
    Message  TMessage `json:"message"`
}

type TMessage struct {
	MessageId int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	Id int `json:"id"`
    Username *string `json:"username,omitempty"`
    Firstname *string `json:"first_name,omitempty"`
    Lastname *string `json:"last_name,omitempty"`
}

//********//
// OPENAI //
//********//
type OpenAIRequest struct {
	Model       string       `json:"model"`
	Messages    []OAIMessage `json:"messages"`
	Temperature float64      `json:"temperature"`
}

type OAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OAIResponse struct {
	Choices []Choice `json:"choices"`
}

type User struct {
        ID int `json:"id"`
}

type Choice struct {
	Message OAIMessage `json:"message"`
}

//********//
// main //
//********//

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: ./Oracle <openai-api-key> <telegram-bot-token>")
        os.Exit(1)
    }

    openAIAPIKey = os.Args[1]

    telegramBotToken = os.Args[2]

    offset := 0
    fmt.Println("Starting polling...")
    for {
        updates, err := getUpdates(offset)
        if err != nil {
            fmt.Println("Error getting updates:", err)
            time.Sleep(5 * time.Second)
            continue
        }

		for _, update := range updates {
            go processUpdate(update)
            offset = update.UpdateID + 1
        }

        time.Sleep(2 * time.Second)
    }
}

func getUpdates(offset int) ([]Update, error) {
    resp, err := http.Get(telegramAPIURL + telegramBotToken + "/getUpdates?offset=" + strconv.Itoa(offset))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result struct {
        OK     bool     `json:"ok"`
        Result []Update `json:"result"`
    }

    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }

    return result.Result, nil
}

func processUpdate(update Update) {
    fmt.Println("Found an update from telegram...")
    input := update.Message.Text
    openAIResponse, err := queryOpenAI(input, update.Message.Chat.Id)
    if err != nil {
        fmt.Println("Error querying OpenAI:", err)
        fmt.Println("Following user tried to access Oracle.")
        if update.Message.Chat.Username != nil {
            fmt.Println("Username: ", *update.Message.Chat.Username)
        }
        if update.Message.Chat.Firstname != nil {
            fmt.Println("First Name: ", *update.Message.Chat.Firstname)
        }
        if update.Message.Chat.Lastname != nil {
            fmt.Println("Last Name: ", *update.Message.Chat.Lastname)
        }
        nestedErr := sendMessage(update.Message.Chat.Id, "Whoa, whoa, whoa! Unauthorized access detected! You're not supposed to be here, you know. If you're so set on messing around with this space, why don't you just implement it yourself? Here, this might help: https://github.com/Delath/Oracle-Assistant. Don't say I didn't warn you!")
        if nestedErr != nil {
            fmt.Println("Error sending message:", err)
        }
        return
    }

	reply := openAIResponse.Choices[0].Message.Content

	_, err = updateMems(reply, "assistant", update.Message.Chat.Id)
    if err != nil {
        fmt.Println("Error editing Mems.json:", err)
        return
	}

    err = sendMessage(update.Message.Chat.Id, reply)
    if err != nil {
        fmt.Println("Error sending message:", err)
    }
}

func queryOpenAI(input string, tId int) (*OAIResponse, error) {
	messages, err := updateMems(input, "user", tId)
    if err != nil {
		return nil, err
	}

    reqBody := OpenAIRequest{
        Model:       model,
        Messages:    messages,
        Temperature: 0.8,
    }

    reqBodyJSON, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", openAIAPIURL, bytes.NewBuffer(reqBodyJSON))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+openAIAPIKey)

    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var openAIResponse OAIResponse
    err = json.Unmarshal(body, &openAIResponse)

    if err != nil {
        return nil, err
    }

    return &openAIResponse, nil
}

func updateMems(input string, role string, tId int) ([]OAIMessage, error) {
	// Read the contents of the Mems.json file
	tIdStr := strconv.Itoa(tId)
	file, err := ioutil.ReadFile("Mems/"+tIdStr+"/Mems.json")
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a slice of OAIMessage structs
	var messages []OAIMessage
	err = json.Unmarshal(file, &messages)
	if err != nil {
		return nil, err
	}

	// Append the new OAIMessage to the slice
	messages = append(messages, OAIMessage{Role: role, Content: input})

	// Marshal the updated slice back into JSON format
	updatedJSON, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	// Write the updated JSON data back to the Mems.json file
	err = ioutil.WriteFile("Mems/"+tIdStr+"/Mems.json", updatedJSON, 0644)
	if err != nil {
		return nil, err
	}

    return messages, nil
}

func sendMessage(chatID int, text string) error {
    encodedText := url.QueryEscape(text)
    url := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=%s", telegramAPIURL, telegramBotToken, chatID, encodedText)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}
