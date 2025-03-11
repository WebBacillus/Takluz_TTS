package main

import (
	command "Takluz_TTS/audio/prefix"
	"Takluz_TTS/audio/sound"
	myInterface "Takluz_TTS/myInterface"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"github.com/andreykaipov/goobs"
)

func getPath(nextPath string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	newPath := filepath.Join(currentDir, nextPath)
	return filepath.ToSlash(newPath)
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	//init OPEN_API
	Open_AI_Config := myInterface.Open_AI_Config{
		Key:   viper.GetString("OPEN_AI.KEY"),
		Model: viper.GetString("OPEN_AI.MODEL"),
		Speed: viper.GetString("OPEN_AI.SPEED"),
		Voice: viper.GetString("OPEN_AI.VOICE"),
	}

	//init OBS
	OBS_Config := myInterface.OBS_Config{
		URL:     viper.GetString("OBS.URL"),
		Key:     viper.GetString("OBS.KEY"),
		Browser: viper.GetString("OBS.BROWSER"),
		Media:   viper.GetString("OBS.MEDIA"),
	}

	//init BOT_NOI
	BOT_NOI_Config := myInterface.BOT_NOI_Config{
		Key:       viper.GetString("BOT_NOI.KEY"),
		Speaker:   viper.GetString("BOT_NOI.SPEAKER"),
		Volume:    viper.GetFloat64("BOT_NOI.VOLUME"),
		Speed:     viper.GetFloat64("BOT_NOI.SPEED"),
		TypeMedia: viper.GetString("BOT_NOI.TYPE_MEDIA"),
		SaveFile:  viper.GetBool("BOT_NOI.SAVE_FILE"),
		Language:  viper.GetString("BOT_NOI.LANGUAGE"),
	}

	//init Resemble
	Resemble_config := myInterface.Resemble_Config{
		Key:          viper.GetString("RESEMBLE.KEY"),
		VoiceUUID:    viper.GetString("RESEMBLE.VOICE_UUID"),
		SampleRate:   viper.GetInt("RESEMBLE.SAMPLE_RATE"),
		OutputFormat: viper.GetString("RESEMBLE.OUTPUT_FORMAT"),
		Speed:        viper.GetString("RESEMBLE.SPEED"),
	}

	AI := viper.GetString("AI")
	limitToken := viper.GetInt("LIMIT")

	client, err := goobs.New(OBS_Config.URL, goobs.WithPassword(OBS_Config.Key))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect()
	command.CreateSilentAudio()
	// sound.InitSound(client, OBS_Config.Media, getPath("speech.mp3"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("END")
		command.CreateSilentAudio()
		os.Exit(0)
	}()

	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		var message myInterface.Message
		if err := c.BodyParser(&message); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if len(message.Message) >= limitToken {
			message.Message = message.Message[:limitToken]

		}

		if AI == "BOT_NOI" {
			sound.GetSoundBotNoi(message.Message, BOT_NOI_Config, "speech.mp3")
		} else if AI == "OPEN_AI" {
			sound.GetSound(message.Message, Open_AI_Config, "speech.mp3")
		} else if AI == "RESEMBLE" {
			text := fmt.Sprintf(`<speak>
					<prosody rate="%s" pitch="x-high">
						<lang xml:lang="th-th">
								%s
						</lang>
					</prosody>
			</speak>`, Resemble_config.Speed, message.Message)
			sound.GetSoundResemble(text, Resemble_config, "speech.mp3")
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid AI configuration"})
		}
		// prefix.ConcatAudio([]string{"sample-3s.mp3", "speech.mp3"}, "output.mp3")

		err = sound.PlaySound(client, OBS_Config.Media, getPath("speech.mp3"))
		if err != nil {
			log.Println(err.Error())
		}

		// err = sound.PlayAnimation(client, OBS_Config.Browser, getPath("templates"), message.UserName, message.Message)
		// if err != nil {
		// 	log.Println(err.Error())
		// }

		fmt.Println(color.GreenString(message.UserName), "used", color.RedString(fmt.Sprintf("%d", len(message.Message))), "characters", message.Message)
		return c.Status(200).SendString(message.Message)

	})

	app.Listen("localhost:4444")
}
