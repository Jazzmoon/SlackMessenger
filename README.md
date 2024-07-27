# SlackMessenger

A simple program that sends/edits messages saved in the config file

This program was created to help with sending messages to a slack channel within Task Scheduler or other automation tools that can run executables.

## Installation

1. Download the latest release from the [releases page](https://github.com/Jazzmoon/SlackMessenger/releases)
2. Clone `config.example.yaml` and rename it to `config.yaml` and fill in the necessary information
3. Run the program

## Usage

Provides a list of commands that can be used including the ones created in the config:

```bash
./SlackMessenger.exe -help
```

Run the program with the command you want to use:

```bash
./SlackMessenger.exe -<ActionName>
```

## Commands

The only default command provided is `-delete` which deletes the message saved in the config file

## Config Details

```yaml
token: XXXX # Slack API bot token
channel_id: XXXX # Slack channel ID
last_message_timestamp: "" # Timestamp shown in console when print_timestamp_for_config: true is set in an action
global_placeholders:
  - key: "value" # Global placeholders that can be used in all actions use them with {gp.<key>}
actions:
  - name: action1 # the name of the action used when calling the program
    description: "Description of the action" # Description of the action shown in the help command
    markdown_text: "Text to send" # Text to send to the channel
    print_timestamp_for_config: true # Prints the timestamp of this message sent to the console
    is_edit_only: true # If true this action can only edit the last message saved in the config
    is_send_only: true # If true this action can only send a message and not edit
urls:
  - SEND_MESSAGE: "https://slack.com/api/chat.postMessage" # Slack API URL to send a message
  - EDIT_MESSAGE: "https://slack.com/api/chat.update" # Slack API URL to edit a message
  - DELETE_MESSAGE: "https://slack.com/api/chat.delete" # Slack API URL to delete a message
```
