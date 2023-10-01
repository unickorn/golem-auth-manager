package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var apiKey string
var apiSecret string

func Initialize(mailjetApiKey, mailjetApiSecret string) {
	apiKey = mailjetApiKey
	apiSecret = mailjetApiSecret
}

// Message is the message to send to the user.
type Message struct {
	From     map[string]string   `json:"From"`
	To       []map[string]string `json:"To"`
	Subject  string              `json:"Subject"`
	TextPart string              `json:"TextPart"`
	HTMLPart string              `json:"HTMLPart"`
	CustomID string              `json:"CustomID"`
}

// MailjetRequest is the request body for sending an email.
type MailjetRequest struct {
	Messages []Message `json:"Messages"`
}

// MailjetResponse is the response sent back by mailjet.
type MailjetResponse struct {
	StatusCode   int    `json:"StatusCode,omitempty"`
	ErrorMessage string `json:"ErrorMessage,omitempty"`
	Messages     []struct {
		Status string `json:"Status,omitempty"`
	} `json:"Messages,omitempty"`
}

// SendMailjetEmail sends an email to the user with their secret code.
func SendMailjetEmail(target string, code string) error {
	url := "https://api.mailjet.com/v3.1/send"
	request := MailjetRequest{
		Messages: []Message{
			{
				From: map[string]string{
					"Name":  "Golem Cloud Poll Demo",
					"Email": "golem-poll@ilhan.me",
				},
				To: []map[string]string{
					{
						"Email": target,
					},
				},
				Subject:  "Golem Poll Demo | Your Secret Code",
				TextPart: "Your secret authentication code to use golem-poll is:\r\n" + code + "\r\nIf you haven't signed up for this, you can safely ignore this message.\r\n",
				HTMLPart: "<!DOCTYPE html>\n<html>\n<head>\n<style>\n    body {\n        font-family: Arial, sans-serif;\n        background-color: #f0f0f0;\n    }\n\n    .container {\n        width: 80%;\n        margin: auto;\n        background-color: white;\n        padding: 20px;\n        border-radius: 10px;\n    }\n\n    .code {\n        font-size: 24px;\n        font-weight: bold;\n        color: #1DBF73;\n        padding: 20px;\n        border: 2px solid #1DBF73;\n        border-radius: 10px;\n        text-align: center;\n    }\n\n    .message {\n        font-size: 16px;\n        color: #333;\n        padding: 20px;\n    }\n\n    .footer {\n        font-size: 12px;\n        color: #777;\n        padding: 20px;\n        text-align: center;\n    }\n</style>\n</head>\n<body>\n\n<div class=\"container\">\n    <div class=\"code\">\n        Your secret authentication code to use golem-poll is: <br>\n        <span style=\"font-weight: normal;\">" + code + "</span>\n    </div>\n    <div class=\"message\">\n        If you haven't signed up for this, you can safely ignore this message.\n    </div>\n    <div class=\"footer\">\n        Made with ❤️ in Istanbul\n    </div>\n</div>\n\n</body>\n</html>\n",
				CustomID: "GolemPollDemo",
			},
		},
	}
	jsonStr, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apiKey, apiSecret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var response MailjetResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if len(response.Messages) == 0 || response.Messages[0].Status != "success" {
		return fmt.Errorf("mailjet error: %s", response.ErrorMessage)
	}
	return nil
}
