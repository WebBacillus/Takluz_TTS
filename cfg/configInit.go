package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitOpenAIConfig() (Open_AI_Config, error) {
	config := Open_AI_Config{
		Key:   viper.GetString("OPEN_AI.KEY"),
		Model: viper.GetString("OPEN_AI.MODEL"),
		Speed: viper.GetString("OPEN_AI.SPEED"),
		Voice: viper.GetString("OPEN_AI.VOICE"),
	}

	if config.Key == "ADD_YOUR_OWN_KEY_HERE" || config.Key == "" {
		return config, fmt.Errorf("OPEN_AI.KEY is required")
	}
	if config.Model == "ADD_YOUR_OWN_MODEL_HERE" || config.Model == "" {
		config.Model = "tts-1-hd"
	}
	if config.Speed == "ADD_YOUR_OWN_SPEED_HERE" || config.Speed == "" {
		config.Speed = "0.8"
	}
	if config.Voice == "ADD_YOUR_OWN_VOICE_HERE" || config.Voice == "" {
		config.Voice = "alloy"
	}

	return config, nil
}

func InitOBSConfig() (OBS_Config, error) {
	config := OBS_Config{
		URL:     viper.GetString("OBS.URL"),
		Key:     viper.GetString("OBS.KEY"),
		Browser: viper.GetString("OBS.BROWSER"),
		Media:   viper.GetString("OBS.MEDIA"),
	}

	if config.URL == "ADD_YOUR_OWN_URL_HERE" || config.URL == "" {
		config.URL = "localhost:4455"
	}
	if config.Key == "ADD_YOUR_OWN_KEY_HERE" || config.Key == "" {
		return config, fmt.Errorf("OBS.KEY is required")
	}
	if config.Browser == "ADD_YOUR_OWN_BROWSER_HERE" || config.Browser == "" {
		return config, fmt.Errorf("OBS.BROWSER is required")
	}
	if config.Media == "ADD_YOUR_OWN_MEDIA_HERE" || config.Media == "" {
		return config, fmt.Errorf("OBS.MEDIA is required")
	}

	return config, nil
}

func InitBotNoiConfig() (BOT_NOI_Config, error) {
	config := BOT_NOI_Config{
		Key:       viper.GetString("BOT_NOI.KEY"),
		Speaker:   viper.GetString("BOT_NOI.SPEAKER"),
		Volume:    viper.GetFloat64("BOT_NOI.VOLUME"),
		Speed:     viper.GetFloat64("BOT_NOI.SPEED"),
		TypeMedia: viper.GetString("BOT_NOI.TYPE_MEDIA"),
		SaveFile:  viper.GetBool("BOT_NOI.SAVE_FILE"),
		Language:  viper.GetString("BOT_NOI.LANGUAGE"),
	}

	if config.Key == "ADD_YOUR_OWN_KEY_HERE" || config.Key == "" {
		return config, fmt.Errorf("BOT_NOI.KEY is required")
	}
	if config.Speaker == "ADD_YOUR_OWN_SPEAKER_HERE" || config.Speaker == "" {
		config.Speaker = "8"
	}
	if config.Volume == 0 {
		config.Volume = 0.8
	}
	if config.Speed == 0 {
		config.Speed = 0.8
	}
	if config.TypeMedia == "ADD_YOUR_OWN_TYPE_MEDIA_HERE" || config.TypeMedia == "" {
		config.TypeMedia = "mp3"
	}
	if config.Language == "ADD_YOUR_OWN_LANGUAGE_HERE" || config.Language == "" {
		config.Language = "th"
	}

	return config, nil
}

func InitResembleConfig() (Resemble_Config, error) {
	config := Resemble_Config{
		Key:          viper.GetString("RESEMBLE.KEY"),
		VoiceUUID:    viper.GetString("RESEMBLE.VOICE_UUID"),
		SampleRate:   viper.GetInt("RESEMBLE.SAMPLE_RATE"),
		OutputFormat: viper.GetString("RESEMBLE.OUTPUT_FORMAT"),
		Speed:        viper.GetString("RESEMBLE.SPEED"),
	}

	if config.Key == "ADD_YOUR_OWN_KEY_HERE" || config.Key == "" {
		return config, fmt.Errorf("RESEMBLE.KEY is required")
	}
	if config.VoiceUUID == "ADD_YOUR_OWN_VOICE_UUID_HERE" || config.VoiceUUID == "" {
		config.VoiceUUID = "55592656"
	}
	if config.SampleRate == 0 {
		config.SampleRate = 48000
	}
	if config.OutputFormat == "ADD_YOUR_OWN_OUTPUT_FORMAT_HERE" || config.OutputFormat == "" {
		config.OutputFormat = "mp3"
	}
	if config.Speed == "ADD_YOUR_OWN_SPEED_HERE" || config.Speed == "" {
		config.Speed = "80%"
	}

	return config, nil
}

func InitGeneralConfig() (General_Config, error) {
	config := General_Config{
		AI:         viper.GetString("AI"),
		LimitToken: viper.GetInt("LIMIT"),
		TimeLimit:  viper.GetInt("TIME_LIMIT"),
	}

	if config.AI == "ADD_YOUR_OWN_AI_HERE" || config.AI == "" {
		config.AI = "default_ai"
	}
	if config.LimitToken == 0 {
		config.LimitToken = 1000
	}
	if config.TimeLimit == 0 {
		config.TimeLimit = 60
	}

	return config, nil
}

func InitMicrosoftConfig() (Microsoft_Config, error) {
	config := Microsoft_Config{
		Key:    viper.GetString("MICROSOFT.KEY"),
		Region: viper.GetString("MICROSOFT.REGION"),
		Voice:  viper.GetString("MICROSOFT.VOICE"),
	}

	if config.Key == "ADD_YOUR_OWN_KEY_HERE" || config.Key == "" {
		return config, fmt.Errorf("MICROSOFT.KEY is required")
	}
	if config.Region == "ADD_YOUR_OWN_REGION_HERE" || config.Region == "" {
		config.Region = "eastus"
	}
	if config.Voice == "ADD_YOUR_OWN_VOICE_HERE" || config.Voice == "" {
		config.Voice = "th-TH-SuchadaNeural"
	}

	return config, nil
}
