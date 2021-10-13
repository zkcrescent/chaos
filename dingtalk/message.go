package dingtalk

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Link is a message with a link
type Link struct {
	Title      string `json:"title,omitempty"`
	Text       string `json:"text,omitempty"`
	MessageURL string `json:"messageUrl,omitempty"`
	PicURL     string `json:"picUrl,omitempty"`
}

// Message is a message send to bot Webhook
type Message interface {
	Message() string
	ToDingTalkMsg() DingTalkMsg
}

type DingTalkMsg struct {
	Msgtype    string      `json:"msgtype"`
	Text       *text       `json:"text,omitempty"`
	Link       *Link       `json:"link,omitempty"`
	Markdown   *markdown   `json:"markdown,omitempty"`
	At         *at         `json:"at,omitempty"`
	ActionCard *actionCard `json:"actionCard,omitempty"`
	FeedCard   *feedCard   `json:"feedCard,omitempty"`
}

type message struct {
	Msgtype    string      `json:"msgtype"`
	Text       *text       `json:"text,omitempty"`
	Link       *Link       `json:"link,omitempty"`
	Markdown   *markdown   `json:"markdown,omitempty"`
	At         *at         `json:"at,omitempty"`
	ActionCard *actionCard `json:"actionCard,omitempty"`
	FeedCard   *feedCard   `json:"feedCard,omitempty"`
}

func (t message) Message() string {
	body, _ := json.Marshal(t)
	return string(body)
}
func (t message) ToDingTalkMsg() DingTalkMsg {
	return DingTalkMsg{
		Msgtype:    t.Msgtype,
		Text:       t.Text,
		Link:       t.Link,
		Markdown:   t.Markdown,
		At:         t.At,
		ActionCard: t.ActionCard,
		FeedCard:   t.FeedCard,
	}
}

type text struct {
	Content string `json:"content"`
}
type markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type at struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}
type actionCard struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	HideAvatar     string `json:"hideAvatar"`
	BtnOrientation string `json:"btnOrientation"`

	SingleTitle string `json:"singleTitle"`
	SingleURL   string `json:"singleURL"`

	Btns []*actionBtn `json:"btns"`
}
type feedCard struct {
	Links []*Link `json:"links"`
}
type actionBtn struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

// NewLink creates a linked message
func NewLink(title string, text string, messageURL string, picURL string) *Link {
	return &Link{
		Title:      title,
		Text:       text,
		MessageURL: messageURL,
		PicURL:     picURL,
	}
}

// NewTextMsg creates a text message
func NewTextMsg(content string, atMobiles []string, atAll bool) Message {
	return &message{
		Msgtype: "text",
		Text: &text{
			Content: content,
		},
		At: &at{
			AtMobiles: atMobiles,
			IsAtAll:   atAll,
		},
	}
}

// NewLinkMsg creates a linked text message
func NewLinkMsg(title string, text string, messageURL string, picURL string) Message {
	return &message{
		Msgtype: "link",
		Link:    NewLink(title, text, messageURL, picURL),
	}
}

// NewMarkdownMsg creates a markdown message
func NewMarkdownMsg(title string, atMobiles []string, atAll bool, markdownParts ...string) Message {
	return &message{
		Msgtype: "markdown",
		Markdown: &markdown{
			Title: title,
			Text:  strings.Join(markdownParts, "\n"),
		},
		At: &at{
			AtMobiles: atMobiles,
			IsAtAll:   atAll,
		},
	}
}

// NewSigleActionCardMsg creates a single action card message
func NewSigleActionCardMsg(title string, text string, hideAvatar bool, actionTitle string, actionURL string) Message {
	hideAvatarStr := "0"
	if hideAvatar {
		hideAvatarStr = "1"
	}

	return &message{
		Msgtype: "actionCard",
		ActionCard: &actionCard{
			Title:          title,
			Text:           text,
			HideAvatar:     hideAvatarStr,
			BtnOrientation: "0",
			SingleTitle:    actionTitle,
			SingleURL:      actionURL,
		},
	}
}

// NewMultiActionCardMsg creates a multi action card message
func NewMultiActionCardMsg(title string, text string, hideAvatar bool, horiz bool, actionAndURLs ...string) Message {
	hideAvatarStr := "0"
	if hideAvatar {
		hideAvatarStr = "1"
	}
	orientation := "0"
	if horiz {
		orientation = "1"
	}
	var btns []*actionBtn
	for i := 1; i < len(actionAndURLs); i += 2 {
		btns = append(
			btns,
			&actionBtn{
				Title:     actionAndURLs[i-1],
				ActionURL: actionAndURLs[i],
			},
		)
	}

	return &message{
		Msgtype: "actionCard",
		ActionCard: &actionCard{
			Title:          title,
			Text:           text,
			HideAvatar:     hideAvatarStr,
			BtnOrientation: orientation,
			Btns:           btns,
		},
	}
}

// NewFeedCardMsg creates a feed card message
func NewFeedCardMsg(links ...*Link) Message {
	return &message{
		Msgtype: "feedCard",
		FeedCard: &feedCard{
			Links: links,
		},
	}
}

// MarkdownInline returns markdown inline block
func MarkdownInline(block string) string {
	return fmt.Sprintf("```\n%s\n```", block)
}

// MarkdownBold returns markdown bold text
func MarkdownBold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

// MarkdownItalic returns markdown italic text
func MarkdownItalic(text string) string {
	return fmt.Sprintf("*%s*", text)
}

// MarkdownLink returns markdown linked text
func MarkdownLink(text string, href string) string {
	return fmt.Sprintf("[%s](%s)", text, href)
}

// MarkdownImage returns markdown images
func MarkdownImage(image string) string {
	return fmt.Sprintf("![image](%s)", image)
}

// MarkdownHeader returns markdown header
func MarkdownHeader(level int, header string) string {
	return fmt.Sprintf("%s %s", strings.Repeat("#", level), header)
}

// MarkdownRefer returns markdown referer text
func MarkdownRefer(text string) string {
	return fmt.Sprintf("> %s", text)
}

// MarkdownList returns markdown list block
func MarkdownList(items ...string) string {
	var list []string
	for _, i := range items {
		list = append(list, fmt.Sprintf("- %s", i))
	}
	return strings.Join(list, "\n")
}

// MarkdownOrderList returns markdown order list
func MarkdownOrderList(items ...string) string {
	var list []string
	for o, i := range items {
		list = append(list, fmt.Sprintf("%d. %s", o+1, i))
	}
	return strings.Join(list, "\n")
}

// MarkdownQuote returns markdown quote text
func MarkdownQuote(block string) string {
	items := strings.Split(block, "\n")
	var list []string
	for _, i := range items {
		list = append(list, fmt.Sprintf("%s", i))
	}
	return strings.Join(list, "\n")
}
