package types

import (
	"gopkg.in/yaml.v3"
)

type ParamAction struct {
	// Param passed in the command
	Name string `yaml:"name"`
	// Description of the param
	Description string `yaml:"description"`

	// Message to be sent to the channel/edited
	MarkdownText string `yaml:"markdown_text"`

	// Log the sent message timestamp for the config file
	PrintTimestampForConfig bool `yaml:"print_timestamp_for_config,omitempty"`

	// Is edit only
	IsEditOnly bool `yaml:"is_edit_only,omitempty"`

	// Is send only
	IsSendOnly bool `yaml:"is_send_only,omitempty"`
}

type EnvConfig struct {
	// Bot Auth Token for the bot
	BotAuthToken string `yaml:"bot_auth_token"`

	// Channel ID for the bot to post to
	ChannelID string `yaml:"channel_id"`

	// Timestamp for the last message sent (used to edit the message)
	LastMessageTimestamp string `yaml:"last_message_timestamp"`

	// Array of replacement strings
	GlobalPlaceholders GlobalPlaceholders `yaml:"global_placeholders"`

	// List of actions to be performed
	Actions []ParamAction `yaml:"actions"`

	URLS URLS `yaml:"urls"`
}

type GlobalPlaceholders map[string]string

func (r *GlobalPlaceholders) UnmarshalYAML(value *yaml.Node) error {
	// Check if the value is a sequence
	if value.Kind == yaml.SequenceNode {
		// Initialize the map
		*r = make(map[string]string)

		for _, val := range value.Content {
			// Add the key value pair to the map
			(*r)[val.Content[0].Value] = val.Content[1].Value
		}
	}

	return nil
}

type URLS map[string]string

func (r *URLS) UnmarshalYAML(value *yaml.Node) error {
	// Check if the value is a sequence
	if value.Kind == yaml.SequenceNode {
		// Initialize the map
		*r = make(map[string]string)

		for _, val := range value.Content {
			// Add the key value pair to the map
			(*r)[val.Content[0].Value] = val.Content[1].Value
		}
	}

	return nil
}
