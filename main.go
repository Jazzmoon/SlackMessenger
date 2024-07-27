package main

import (
	"Jazzmoon/SlackMessager/types"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var Config types.EnvConfig

var Flags map[string]*bool = make(map[string]*bool)

func main() {
	// Load the config.yaml file into the Config struct
	err := LoadConfig(&Config)
	if err != nil {
		err = fmt.Errorf("error loading config: %v", err)

		panic(err)
	}

	// Load the params passed in the command
	deleteMessage := flag.Bool("delete", false, "Delete the saved timestamp message in the config file")
	LoadParams()

	flag.Usage = func() {
		fmt.Println("Avialable messages:")
		for _, param := range Config.Actions {
			fmt.Printf("  -%s:\n      %s\n", param.Name, param.Description)
		}
		fmt.Println("Avialable actions:")
		fmt.Println("  -delete:\n      Delete the saved timestamp message in the config file")
	}

	flag.Parse()

	// Check if the delete flag is passed
	if *deleteMessage {
		URL := Config.URLS["DELETE_MESSAGE"]
		postBody := types.SlackMessage{
			Channel: Config.ChannelID,
			TS:      Config.LastMessageTimestamp,
		}

		res := MakeRequest(URL, &postBody, false)
		if !res {
			fmt.Println("Error deleting message")
			return
		}

		fmt.Println("Message deleted successfully")
		fmt.Println("Make sure to remove the timestamp from the config file")

		return
	}

	// Check if more than one action is passed
	actions := 0
	for _, value := range Flags {
		if *value {
			actions++
			if actions > 1 {
				fmt.Println("Only one action can be passed at a time")
				return
			}
		}
	}

	// Check if an action is passed
	if actions == 0 {
		fmt.Println("No action passed")
		return
	}

	// Fetch the action passed
	var action types.ParamAction
	for key, value := range Flags {
		if *value {
			for _, param := range Config.Actions {
				if param.Name == key {
					action = param
				}
			}
		}
	}

	// Perform the action
	PerformAction(&action)

}

// PerformAction performs the action passed
func PerformAction(action *types.ParamAction) {
	markdwn := action.MarkdownText

	// Replace the placeholders
	for key, value := range Config.GlobalPlaceholders {
		markdwn = strings.Replace(markdwn, fmt.Sprintf("{gp.%s}", key), value, -1)
	}

	timeNow := time.Now()
	// Replace {TZ} with unix timestamp
	markdwn = strings.Replace(markdwn, "{TZ}", fmt.Sprintf("%d", timeNow.Unix()), -1)
	// Replace {TZ_Locale} with unix timestamp in the format "February 18th, 2014 at 6:39 AM PST"
	markdwn = strings.Replace(markdwn, "{TZLocale}", timeNow.Format("January 2nd, 2006 at 3:04 PM MST"), -1)

	// Check if the action is edit only and the last message timestamp is not present
	if action.IsEditOnly && Config.LastMessageTimestamp == "" {
		fmt.Println("Last message timestamp not found, This action can't be performed without the last message timestamp")
		return
	}

	postBody := types.SlackMessage{
		Channel: Config.ChannelID,
		Text:    markdwn,
	}

	var URL string

	if Config.LastMessageTimestamp != "" && !action.IsSendOnly {
		// Edit the message
		postBody.TS = Config.LastMessageTimestamp
		URL = Config.URLS["UPDATE_MESSAGE"]
	} else {
		// Send the message
		URL = Config.URLS["SEND_MESSAGE"]
	}

	MakeRequest(URL, &postBody, action.PrintTimestampForConfig)
}

func MakeRequest(URL string, postBody *types.SlackMessage, PrintTimestampForConfig bool) bool {
	// Send the message
	jsonBody, err := json.Marshal(postBody)
	if err != nil {
		fmt.Printf("error marshalling json: %v", err)
		return false
	}

	// Create the request
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("error creating request: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Config.BotAuthToken))
	req.Header.Set("User-Agent", "Jazzmoon/SlackMessager.v1")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error sending request: %v", err)
		return false
	}

	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != 200 {
		fmt.Printf("error sending message: %v", resp.Status)
		return false
	}

	if PrintTimestampForConfig {
		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response body: %v", err)
			return false
		}

		// Pull the ts value from the response body
		var response types.SlackResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Printf("error unmarshalling response: %v", err)
			return false
		}

		// Print the ts value
		fmt.Printf("Message sent successfully with timestamp: %s\n", response.TS)
		fmt.Printf("Make sure to save this timestamp in the config file\n")

		// Save the timestamp to last_message_ts.txt
		file, err := os.Create("last_message_ts.txt")
		if err != nil {
			fmt.Printf("error creating file: %v", err)
			return false
		}
		defer file.Close()

		_, err = file.WriteString(response.TS)
		if err != nil {
			fmt.Printf("error writing to file: %v", err)
			return false
		}

		fmt.Println("Timestamp saved to last_message_ts.txt")

	}

	return true

}

// LoadConfig loads the config.yaml file into the Config struct
func LoadConfig(config *types.EnvConfig) error {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}

// LoadParams loads the params passed in the command
func LoadParams() {
	for _, param := range Config.Actions {
		flagPtr := flag.Bool(param.Name, false, param.Description)
		Flags[param.Name] = flagPtr
	}
}
