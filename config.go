package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Config struct {
	fs       *flag.FlagSet
	filePath string

	WebhookUrl       string
	TimeLayout       string
	Question         string
	Duration         int
	AllowMultiselect bool
}

func NewConfig(name string) *Config {
	c := Config{}

	c.fs = flag.NewFlagSet(name, flag.ExitOnError)
	c.fs.StringVar(&c.filePath, "config", "", "Path to config file.")

	c.fs.StringVar(&c.WebhookUrl, "webhook", "", "Discord webhooh url.")
	c.fs.StringVar(&c.TimeLayout, "layout", "Mon, 02 Jan 06", "Time layout.")
	c.fs.StringVar(&c.Question, "question", "WHEN?", "The question of the poll.")
	c.fs.IntVar(&c.Duration, "duration", 24, "Number of hours the poll should be open for, up to 32 days.")
	c.fs.BoolVar(&c.AllowMultiselect, "multi-select", true, "Whether a user can select multiple answers.")

	return &c
}

func (c *Config) parseFile() error {
	f, err := os.Open(c.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var args []string
	for scanner.Scan() {
		args = append(args, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	hasConfigPath := slices.ContainsFunc(args, func(s string) bool {
		return strings.HasPrefix(s, "-config") || strings.HasPrefix(s, "--config")
	})
	if hasConfigPath {
		return errors.New("config file cannot contains config path.")
	}

	c.filePath = ""
	return c.Parse(args)
}

func (c *Config) Parse(args []string) error {
	err := c.fs.Parse(args)
	if err != nil {
		return fmt.Errorf("failed to args: %w", err)
	}

	if c.filePath != "" {
		return c.parseFile()
	}

	// Check required.
	if c.WebhookUrl == "" {
		return errors.New("webhook required")
	}

	return nil
}

func (c *Config) Usage() {
	c.fs.Usage()
}
