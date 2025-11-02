package discord

import "sync/atomic"

type PollMedia struct {
	Text string `json:"text"`
	// TODO: emoji https://discord.com/developers/docs/resources/poll#poll-media-object
}

type PollAnswer struct {
	Id        int       `json:"answer_id"`
	PollMedia PollMedia `json:"poll_media"`
}

type Poll struct {

	// The question of the poll. Only text is supported.
	Question PollMedia `json:"question"`

	lastAnswerId int64

	// Each of the answers available in the poll, up to 10.
	Answers []PollAnswer `json:"answers"`

	// Number of hours the poll should be open for, up to 32 days. Defaults to 24.
	Duration int `json:"duration"`

	// Whether a user can select multiple answers. Defaults to false.
	AllowMultiselect bool `json:"allow_multiselect"`
}

func (p *Poll) nextAnswerId() int {
	return int(atomic.AddInt64(&p.lastAnswerId, 1))
}

// AddAnswer simply add an answer to the poll, respecting the id.
func (p *Poll) AddAnswer(text string) {
	answer := PollAnswer{
		Id:        p.nextAnswerId(),
		PollMedia: PollMedia{Text: text},
	}
	p.Answers = append(p.Answers, answer)
}

type WebhookParams struct {
	Content   string `json:"content,omitempty"`
	Username  string `json:"username,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	Poll      Poll   `json:"poll"`
}
