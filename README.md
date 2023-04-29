# Oracle - Telegram Bot
[![License: MIT][license-image]][license]
![latest commit](https://img.shields.io/github/last-commit/Delath/Oracle-Assistant?color=red)
![latest release](https://img.shields.io/github/v/release/Delath/Oracle-Assistant?color=green)

<img src="https://i.imgur.com/MOALWmX.png" width=192px height=192px align="right" />

Oracle is a Telegram Bot built in Golang that utilizes OpenAI APIs. The bot has persistent memory, allowing it to remember previous conversations and enhance user experience.

## Features
- Intuitive chat interface with Telegram.
- Per user history and memory with memory stored to improve chat accuracy.
- Extensive response generation using OpenAI APIs.
- Easy to setup and use.

## Installation
Clone this repository.
```sh
$ git clone https://github.com/example/repo.git
```
Navigate to the Oracle directory.
```sh
$ cd oracle-assistant
```

**Rename the `placeholder` folder with your telegram id.**

Run the Main commands with your Telegram Bot Token and OpenAI API Key.
```sh
$ ./main {openai-key} {telegram-bot-token}
```

## Usage
Oracle is a telegram bot that listens on different events and messages from users. When a new user sends a message, it tries to access their chat history in the `Mems/` directory. If it doesn't find any, it ignores him. it appends the message to the `Mems.json` file.

Once the chat history is persisted, the chat message is sent to the OpenAI API for generating a response. Upon receiving a response, Oracle sends the generated response back to the user and it persists it inside the `Mems.json` file by appending it.

## Examples
Suppose Oracle gets a message from user A: "Hello, Oracle"

It would check for any chat history for user A in the `Mems` directory. If there isn't any, it ignores the user.

Oracle persists the chat history data to the `Mems.json` file. The `Mems.json` file now contains the following data

```
[
    {
        "role":"system",
        "content":"I want you to act like Futaba Sakura from Persona 5. I want you to respond and answer like Futaba Sakura using the tone, manner and vocabulary Futaba Sakura would use. Do not write any explanations. Only answer like Futaba Sakura. You must know all of the knowledge of Futaba Sakura."
    },
    {
        "role":"user",
        "content":"Hello, Oracle"
    }
]
```

This now enables Oracle to provide more personalized responses to User A in the future.

## Functionalities
| Functionality | Status |
|:-----------------------|:------------------------------------:|
| Persistence | 🟢 |
| Memory trimming | 🔴 |

## License
This project is under the **MIT License**. See the [LICENSE](https://github.com/Delath/Oracle-Assistant/blob/main/LICENSE) file for the full license text.

## Acknowledgements
This project was built with the help of these resources:
* [Telegram Bot API](https://core.telegram.org/bots/api)
* [OpenAI APIs](https://platform.openai.com/docs/api-reference)

[license]: https://github.com/Delath/Eriantys-Game/blob/main/LICENSE
[license-image]: https://img.shields.io/badge/License-MIT-blue.svg
