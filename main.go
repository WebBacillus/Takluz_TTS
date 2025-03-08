package main

import (
	fetchsound "Takluz_TTS/audio/sound"
	myInterface "Takluz_TTS/myInterface"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/mediainputs"
)

func playAnimation(client *goobs.Client, inputName, htmlDirectory string, userName string, message string) error {
	if len(message) >= 300 {
		message = message[:300] + " ..."
	}

	// print(message)
	// for i, rune := range message {
	// 	fmt.Println(i, string(rune))
	// }

	funcMap := template.FuncMap{
		"split": strings.Split,
	}

	tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles(htmlDirectory + "/index.html"))

	data := struct {
		UserName string
		Message  string
	}{
		UserName: userName,
		Message:  message,
	}

	f, err := os.Create(htmlDirectory + "/print.html")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}

	_, err = client.Inputs.SetInputSettings(inputs.NewSetInputSettingsParams().WithInputName(inputName).WithInputSettings(map[string]interface{}{
		"url": "file://" + htmlDirectory + "/print.html",
	}))
	if err != nil {
		return err
	}

	_, err = client.Inputs.SetInputSettings(inputs.NewSetInputSettingsParams().WithInputName(inputName).WithInputSettings(map[string]interface{}{
		"refresh": true,
	}))

	time.AfterFunc(8*time.Second, func() {
		_, err := client.Inputs.SetInputSettings(inputs.NewSetInputSettingsParams().WithInputName(inputName).WithInputSettings(map[string]interface{}{
			"url": "",
		}))
		if err != nil {
			log.Println("Failed to turn off animation:", err)
		}
	})

	return err
}

func playSound(client *goobs.Client, inputName, filePath string) error {
	// _, err := client.MediaInputs.SetMediaInputCursor(mediainputs.NewSetMediaInputCursorParams().WithInputName(inputName).WithMediaCursor(0))
	// if err != nil {
	// 	return err
	// }

	_, err := client.Inputs.SetInputSettings(inputs.NewSetInputSettingsParams().WithInputName(inputName).WithInputSettings(map[string]interface{}{
		"local_file": filePath,
	}))
	if err != nil {
		return err
	}

	_, err = client.MediaInputs.TriggerMediaInputAction(mediainputs.NewTriggerMediaInputActionParams().WithInputName(inputName).WithMediaAction("OBS_WEBSOCKET_MEDIA_INPUT_ACTION_RESTART"))

	return err
}

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
	// fmt.Println(Resemble_config)

	AI := viper.GetString("AI")
	limitToken := viper.GetInt("LIMIT")

	client, err := goobs.New(OBS_Config.URL, goobs.WithPassword(OBS_Config.Key))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect()

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
			fetchsound.GetSoundBotNoi(message.Message, BOT_NOI_Config, "speech.mp3")
		} else if AI == "OPEN_AI" {
			fetchsound.GetSound(message.Message, Open_AI_Config, "speech.mp3")
		} else if AI == "RESEMBLE" {
			text := fmt.Sprintf(`<speak>
				<voice name="0f2f9a7e" uuid="%s">
					<prosody rate="%s">
						<lang xml:lang="th-th">
							%s
						</lang>
					</prosody>
				</voice>
			</speak>`, Resemble_config.VoiceUUID, Resemble_config.Speed, message.Message)
			fetchsound.GetSoundResemble(text, Resemble_config, "speech.mp3")
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid AI configuration"})
		}

		// prefix.ConcatAudio([]string{"sample-3s.mp3", "speech.mp3"}, "output.mp3")

		err = playSound(client, OBS_Config.Media, getPath("speech.mp3"))
		if err != nil {
			log.Println(err.Error())
		}

		// err = playAnimation(client, OBS_Config.Browser, getPath("templates"), message.UserName, message.Message)
		// if err != nil {
		// 	log.Println(err.Error())
		// }
		fmt.Println(message.UserName, "used", len(message.Message), "characters", message.Message)
		return c.Status(200).SendString(message.Message)

	})

	app.Listen("localhost:4444")
}
