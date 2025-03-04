package main

import (
	fetchsound "Takluz_TTS/fetchSound"
	"fmt"
	"os"

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
		"local_file": filePath, // Use "local_file" for Media Sources
	}))
	if err != nil {
		return err
	}

	_, err = client.MediaInputs.TriggerMediaInputAction(mediainputs.NewTriggerMediaInputActionParams().WithInputName(inputName).WithMediaAction("OBS_WEBSOCKET_MEDIA_INPUT_ACTION_RESTART"))

	return err
}

func main() {
	OPEN_API_KEY := os.Getenv("OPEN_API_KEY")
	OBS_KEY := os.Getenv("OBS_KEY")
	client, err := goobs.New("localhost:4455", goobs.WithPassword(OBS_KEY))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("GET SUCCESSFUL")
		var message struct {
			Message string `json:"message"`
		}
		message.Message = "hi"
		go fetchsound.GetSound("พี่กิ้ฟครับ ทราบหรือไม่ว่า การกินไก่ทำให้ตัวใหญ่ขึ้น", OPEN_API_KEY)
		return c.Status(200).JSON(message)
	})
	app.Post("/", func(c *fiber.Ctx) error {
		var message struct {
			Message string `json:"message"`
		}
		if err := c.BodyParser(&message); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		fetchsound.GetSound(message.Message, OPEN_API_KEY)
		err = playSound(client, "MediaSource", "F:/Project/Takluz/Takluz_TTS/output.mp3")
		if err != nil {
			panic(err)
		}
		fmt.Println("recieve post:", message)
		return c.Status(200).JSON(message)

	})

	app.Listen("localhost:4444")
}
