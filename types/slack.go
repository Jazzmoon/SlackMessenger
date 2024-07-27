package types

type SlackMessage struct {
	Channel string `json:"channel"`

	TS string `json:"ts"`

	Text string `json:"text"`
}

type SlackResponse struct {
	OK      bool   `json:"ok"`
	Channel string `json:"channel"`
	TS      string `json:"ts"`
}
