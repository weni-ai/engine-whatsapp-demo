package models

type Message struct {
	To   string `json:"to"`
	Type string `json:"type"`
	Text struct {
		Body string `json:"body"`
	} `json:"text"`
}
