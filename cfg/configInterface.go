package config

type Open_AI_Config struct {
	Key            string                    `json:"key"`
	Model          string                    `json:"model"`
	Speed          float64                   `json:"speed"`
	Voice          string                    `json:"voice"`
	InstructionSet map[string]InstructionSet `json:"instruction_set"`
}

type InstructionSet struct {
	Voice       string `json:"voice"`
	Instruction string `json:"instruction"`
}

type OBS_Config struct {
	URL     string `json:"url"`
	Key     string `json:"key"`
	Browser string `json:"browser"`
	Media   string `json:"media"`
}

type BOT_NOI_Config struct {
	Key       string  `json:"key"`
	Speaker   string  `json:"speaker"`
	Volume    float64 `json:"volume"`
	Speed     float64 `json:"speed"`
	TypeMedia string  `json:"type_media"`
	SaveFile  bool    `json:"save_file"`
	Language  string  `json:"language"`
}

type Resemble_Config struct {
	Key          string `json:"key"`
	VoiceUUID    string `json:"voice_uuid"`
	SampleRate   int    `json:"sample_rate"`
	OutputFormat string `json:"output_format"`
	Speed        string `json:"speed"`
}

type Microsoft_Config struct {
	Key    string `json:"key"`
	Region string `json:"region"`
	Voice  string `json:"voice"`
	Speed  string `json:"speed"`
}

type Google_Config struct {
	Key          string  `json:"key"`
	Name         string  `json:"name"`
	SpeakingRate float64 `json:"speaking_rate"`
	Pitch        float64 `json:"pitch"`
	VolumeGainDb float64 `json:"volume_gain_db"`
}

type General_Config struct {
	// AI         string `json:"ai"`
	LimitToken  int    `json:"limit_token"`
	TimeLimit   int    `json:"time_limit"`
	Player      string `json:"player"`
	DataCollect bool   `json:"data_collect"`
}
