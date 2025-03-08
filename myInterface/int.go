package myInterface

type Message struct {
	UserName string `json:"userName"`
	Message  string `json:"message"`
}

type Open_AI_Config struct {
	Key   string `json:"key"`
	Model string `json:"model"`
	Speed string `json:"speed"`
	Voice string `json:"voice"`
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
