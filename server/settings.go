package main

type Settings struct {
	Address    string
	DataFolder string
}

func InitSiteSettings(address string, dataFolder string) Settings {
	settings := Settings{
		Address:    "localhost:9001",
		DataFolder: "./data",
	}
	return settings
}
