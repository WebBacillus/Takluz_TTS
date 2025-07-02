package config

type Open_AI_Config struct {
	Key            string                    `yaml:"KEY"`
	Model          string                    `yaml:"MODEL"`
	Speed          float64                   `yaml:"SPEED"`
	Voice          string                    `yaml:"VOICE"`
	InstructionSet map[string]InstructionSet `yaml:"INSTRUCTION_SET"`
}

type InstructionSet struct {
	Voice       string `yaml:"VOICE"`
	Instruction string `yaml:"INSTRUCTION"`
}

type CharacterProfile struct {
	Voice       string `yaml:"VOICE"`
	Instruction string `yaml:"INSTRUCTION"`
}

type OBS_Config struct {
	URL     string `yaml:"url"`
	Key     string `yaml:"key"`
	Browser string `yaml:"browser"`
	Media   string `yaml:"media"`
}

type BOT_NOI_Config struct {
	Key       string  `yaml:"key"`
	Speaker   string  `yaml:"speaker"`
	Volume    float64 `yaml:"volume"`
	Speed     float64 `yaml:"speed"`
	TypeMedia string  `yaml:"type_media"`
	SaveFile  bool    `yaml:"save_file"`
	Language  string  `yaml:"language"`
}

type Resemble_Config struct {
	Key          string `yaml:"key"`
	VoiceUUID    string `yaml:"voice_uuid"`
	SampleRate   int    `yaml:"sample_rate"`
	OutputFormat string `yaml:"output_format"`
	Speed        string `yaml:"speed"`
}

type Microsoft_Config struct {
	Key    string `yaml:"key"`
	Region string `yaml:"region"`
	Voice  string `yaml:"voice"`
	Speed  string `yaml:"speed"`
}

type Google_Config struct {
	Key          string  `yaml:"key"`
	Model        string  `yaml:"model"`
	SpeakingRate float64 `yaml:"speaking_rate"`
	Pitch        float64 `yaml:"pitch"`
	VolumeGainDb float64 `yaml:"volume_gain_db"`
	// GoogleInstructionSet GoogleInstructionSet
}

type General_Config struct {
	// AI         string `yaml:"ai"`
	LimitToken int    `yaml:"limit_token"`
	TimeLimit  int    `yaml:"time_limit"`
	Player     string `yaml:"player"`
	MongoDBKey string `yaml:"mongo_db_key"`
}
