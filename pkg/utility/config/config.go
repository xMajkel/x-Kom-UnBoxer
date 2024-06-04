package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tidwall/pretty"
	"github.com/xMajkel/x-kom-unboxer/pkg/utility"
)

const DEFAULT_API_KEY = "ushoh9OoY7eerae8aiGh"

var GlobalConfig Config

type Config struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	WebhookURL        string `json:"webhook_url"`
	PreferredRollTime string `json:"preferred_roll_time"`
	ApiKey            string `json:"x-api-key"`
}

func ConfigInit() error {
	if _, err := os.Stat("config.json"); err != nil {
		f, err := os.OpenFile("config.json", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.New("could not create config.json")
		}

		var email, password, prefRollTime, webhookUrl string

		fmt.Print("Enter an email: ")
		fmt.Scanln(&email)

		fmt.Print("Enter a password: ")
		fmt.Scanln(&password)

		fmt.Print("Enter a preferred roll time (hh:mm): ")
		fmt.Scanln(&prefRollTime)

		fmt.Print("Enter a Discord Webhook Url (optional): ")
		fmt.Scanln(&webhookUrl)

		h, m := utility.ParsePreferredRollTime(prefRollTime)

		b, err := json.Marshal(Config{
			Email:             email,
			Password:          password,
			WebhookURL:        webhookUrl,
			PreferredRollTime: fmt.Sprintf("%02d:%02d", h, m),
			ApiKey:            DEFAULT_API_KEY,
		})
		if err != nil {
			return errors.New("could not write config.json")
		}
		_, err = f.Write(pretty.Pretty(b))
		if err != nil {
			return errors.New("could not write config.json")
		}
		f.Close()
	}
	err := loadConfing()
	if err != nil {
		return err
	}

	return nil
}

func loadConfing() error {
	f, err := os.OpenFile("config.json", os.O_RDONLY, 0644)
	if err != nil {
		return errors.New("could not load config.json")
	}
	json.NewDecoder(f).Decode(&GlobalConfig)
	f.Close()

	return nil
}

func WriteConfig() (err error) {
	f, err := os.OpenFile("config.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("could not write config.json")
	}
	b, err := json.Marshal(GlobalConfig)
	if err != nil {
		return errors.New("could not write config.json")
	}
	_, err = f.Write(pretty.Pretty(b))
	if err != nil {
		return errors.New("could not write config.json")
	}

	f.Close()

	return err
}
