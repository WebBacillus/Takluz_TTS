package main

import (
	prefix "Takluz_TTS/audio/prefix"
	audio "Takluz_TTS/audio/sound"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/mediainputs"
)

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

func main() {

	// prefix.ConcatAudio([]string{"sample-3s.mp3", "speech.mp3"}, "xdd.mp3")
	fmt.Println("done")
	//เปลี่ยน Key
	OPEN_API_KEY := os.Getenv("OPEN_API_KEY")
	OBS_KEY := os.Getenv("OBS_KEY")
	limitToken := 200
	//

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	client, err := goobs.New("localhost:4455", goobs.WithPassword(OBS_KEY))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect()

	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		var message struct {
			Message string `json:"message"`
		}
		if err := c.BodyParser(&message); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if len(message.Message) >= limitToken {
			message.Message = message.Message[:limitToken]
		}

		audio.GetSound(message.Message, OPEN_API_KEY, "speech.mp3")
		prefix.ConcatAudio([]string{"sample-3s.mp3", "speech.mp3"}, "output.mp3")

		outPath := filepath.Join(currentDir, "output.mp3")
		outPath = filepath.ToSlash(outPath)
		err = playSound(client, "Media Source", outPath)
		if err != nil {
			panic(err)
		}
		fmt.Println(len(message.Message), "characters:", message.Message)
		return c.Status(200).SendString(message.Message)

	})

	app.Listen("localhost:4444")
}
