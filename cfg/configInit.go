package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitOpenAIConfig() (Open_AI_Config, error) {
	// var config Open_AI_Config
	// err := viper.Unmarshal(&config)
	// if err != nil {
	// 	return config, fmt.Errorf("UNMARSHAL ERROR")
	// }

	config := Open_AI_Config{
		Key:            viper.GetString("OPEN_AI.KEY"),
		Model:          viper.GetString("OPEN_AI.MODEL"),
		Speed:          viper.GetFloat64("OPEN_AI.SPEED"),
		Voice:          viper.GetString("OPEN_AI.VOICE"),
		InstructionSet: convertToInstructionSet(viper.GetStringMap("OPEN_AI.INSTRUCTION_SET")),
	}
	viper.UnmarshalKey("OPEN_AI.INSTRUCTION_SET", &config.InstructionSet)

	if config.Key == "ADD_YOUR_OWN_KEY_HERE" || config.Key == "" {
		return config, fmt.Errorf("OPEN_AI.KEY is required")
	}
	if config.Model == "ADD_YOUR_OWN_MODEL_HERE" || config.Model == "" {
		config.Model = "tts-1-hd"
	}
	if config.Speed == 0 {
		config.Speed = 1.0
	}
	if config.Voice == "ADD_YOUR_OWN_VOICE_HERE" || config.Voice == "" {
		config.Voice = "alloy"
	}

	return config, nil
}

func convertToInstructionSet(input map[string]any) map[string]InstructionSet {
	result := make(map[string]InstructionSet)
	for key, value := range input {
		if instruction, ok := value.(InstructionSet); ok {
			result[key] = instruction
		}
	}
	return result
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
		config.Browser = ""
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
		config.TypeMedia = "wav"
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
		config.OutputFormat = "wav"
	}
	if config.Speed == "ADD_YOUR_OWN_SPEED_HERE" || config.Speed == "" {
		config.Speed = "80%"
	}

	return config, nil
}

func InitGeneralConfig() (General_Config, error) {
	config := General_Config{
		LimitToken:  viper.GetInt("GENERAL.LIMIT"),
		TimeLimit:   viper.GetInt("GENERAL.TIME_LIMIT"),
		Player:      viper.GetString("GENERAL.PLAYER"),
		DataCollect: viper.GetBool("GENERAL.DATA_COLLECT"),
	}

	if config.LimitToken == 0 {
		config.LimitToken = 1000
	}
	if config.TimeLimit == 0 {
		config.TimeLimit = 10
	}
	if config.Player == "" {
		config.Player = "FFPLAY"
	}

	return config, nil
}

func InitMicrosoftConfig() (Microsoft_Config, error) {
	config := Microsoft_Config{
		Key:    viper.GetString("MICROSOFT.KEY"),
		Region: viper.GetString("MICROSOFT.REGION"),
		Voice:  viper.GetString("MICROSOFT.VOICE"),
		Speed:  viper.GetString("MICROSOFT.SPEED"),
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
	if config.Speed == "" {
		config.Speed = "1"
	}

	return config, nil
}

func InitGoogleConfig() (Google_Config, error) {
	config := Google_Config{
		Key:          viper.GetString("GOOGLE.KEY"),
		Name:         viper.GetString("GOOGLE.NAME"),
		SpeakingRate: viper.GetFloat64("GOOGLE.SPEAKING_RATE"),
		Pitch:        viper.GetFloat64("GOOGLE.PITCH"),
		VolumeGainDb: viper.GetFloat64("GOOGLE.VOLUME_GAIN_DB"),
	}

	if config.Name == "ADD_YOUR_OWN_KEY_HERE" || config.Name == "" {
		return config, fmt.Errorf("GOOGLE.KEY is required")
	}

	if config.Name == "ADD_YOUR_OWN_NAME_HERE" || config.Name == "" {
		return config, fmt.Errorf("GOOGLE.NAME is required")
	}
	if config.SpeakingRate == 0 {
		config.SpeakingRate = 1.0
	}
	if config.Pitch == 0 {
		config.Pitch = 0.0
	}
	if config.VolumeGainDb == 0 {
		config.VolumeGainDb = 0.0
	}

	return config, nil
}
