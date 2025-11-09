package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// This file handles message events, and checks it for attachments, that includes "minecraft" and then uploads it to https://mclo.gs/, and returns it.

func HandleMCLogsMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Attachments) == 0 {
		return
	}

	for _, attachment := range m.Attachments {
		mcLogContent, isMCLog := isMCLog(attachment)
		if !isMCLog {
			continue
		}

		url, err := uploadToMCLogs(*mcLogContent)
		if err != nil {
			LOGGER.Println("Failed to attachment log to mclo.gs:", err)
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, "Uploaded "+attachment.Filename+" to mclo.gs: "+url)
	}
}

func isMCLog(attachment *discordgo.MessageAttachment) (*string, bool) {
	if attachment == nil {
		return nil, false
	}

	if !strings.HasPrefix(attachment.ContentType, "text/plain") {
		return nil, false
	}

	// Download the attachment from the URL
	httpResp, err := http.Get(attachment.URL)
	if err != nil {
		LOGGER.Println("Encountered an error while downloading attachment:", err)
		return nil, false
	}
	defer httpResp.Body.Close()
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		LOGGER.Println("Failed to read attachment body:", err)
		return nil, false
	}
	fileContents := string(bodyBytes)
	if httpResp.StatusCode != http.StatusOK {
		LOGGER.Println("Failed to download attachment, status code:", fileContents)
		return nil, false
	}

	if strings.Contains(fileContents, "minecraft") {
		return &fileContents, true
	}

	return nil, false
}

func uploadToMCLogs(logContent string) (uploadURL string, err error) {

	data := url.Values{}
	data.Set("content", logContent)

	resp, err := http.Post("https://api.mclo.gs/1/log", "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		LOGGER.Println("Failed to upload log to mclo.gs:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LOGGER.Println("Failed to read mclo.gs response body:", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		LOGGER.Println("Failed to upload log, status code:", body)
		return "", err
	}

	var respBody struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &respBody); err != nil {
		LOGGER.Println("Failed to decode mclo.gs response:", err)
		return "", err
	}

	return respBody.URL, nil
}
