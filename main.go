package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/souhoc/when-next/datepicker"
	"github.com/souhoc/when-next/discord"
)

func sendWebhook(cfg Config, dates []time.Time) error {
	poll := discord.Poll{
		Question:         discord.PollMedia{Text: cfg.Question},
		Duration:         cfg.Duration,
		AllowMultiselect: cfg.AllowMultiselect,
		Answers:          make([]discord.PollAnswer, 0, len(dates)),
	}
	for _, date := range dates {
		poll.AddAnswer(date.Format(cfg.TimeLayout))
	}

	params := discord.WebhookParams{
		Poll: poll,
	}
	b, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	resp, err := http.Post(cfg.WebhookUrl+"?wait=true", "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read resp body: %w", err)
		}

		return fmt.Errorf("failed to send webhook: %s", string(respBody))
	}

	return nil
}

func run(cfg *Config, args []string, stdout io.Writer) error {
	if err := cfg.Parse(args); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	datepickerModel := datepicker.New()
	datepickerModel.TimeLayout = cfg.TimeLayout
	p := tea.NewProgram(datepickerModel)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("Error running datepicker: %w", err)
	}
	dates, err := datepickerModel.GetSelected()
	if err != nil {
		return fmt.Errorf("failed to retrieve selected dates: %w", err)
	}

	if len(dates) != 0 {
		fmt.Fprintln(stdout, "Sending webhook...")
		return sendWebhook(*cfg, dates)
	}
	return nil
}

func main() {
	cfg := NewConfig("when-next")
	if err := run(cfg, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
