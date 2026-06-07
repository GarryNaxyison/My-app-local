package main

import "english-coach-bot/internal/appconfig"

func loadDotEnv(path string) error {
	return appconfig.LoadDotEnv(path)
}
