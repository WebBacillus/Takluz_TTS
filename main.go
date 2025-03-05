package main

import (
	"Takluz_TTS/audio/prefix"
	fetchsound "Takluz_TTS/audio/sound"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/mediainputs"
)

type Message struct {
	UserName string `json:"userName"`
	Message  string `json:"message"`
}

func playAnimation(client *goobs.Client, inputName, htmlDirectory string, userName string, message string) error {
	// if len(message) >= 50 {
	// 	message = message[:50] + "..."
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

	f, err := os.Create(htmlDirectory + "/print.html") //Create file named output.html
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer f.Close() // Ensure the file is closed, even if there's an error.

	err = tmpl.Execute(f, data) //write to a file
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

func main() {

	//เปลี่ยน Key
	OPEN_API_KEY := os.Getenv("OPENAI_API_KEY")
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

	htmlPath := filepath.Join(currentDir, "templates")
	htmlPath = filepath.ToSlash(htmlPath)
	err = playAnimation(client, "Browser", htmlPath, "WebBacillus", "Test Message")
	if err != nil {
		log.Println(err.Error())
	}

	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		var message Message
		if err := c.BodyParser(&message); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if len(message.Message) >= limitToken {
			message.Message = message.Message[:limitToken]
		}

		fetchsound.GetSound(message.Message, OPEN_API_KEY, "speech.mp3")
		prefix.ConcatAudio([]string{"sample-3s.mp3", "speech.mp3"}, "output.mp3")

		outPath := filepath.Join(currentDir, "output.mp3")
		outPath = filepath.ToSlash(outPath)
		err = playSound(client, "Media Source", outPath)
		if err != nil {
			panic(err)
		}

		htmlPath := filepath.Join(currentDir, "templates")
		htmlPath = filepath.ToSlash(htmlPath)
		err = playAnimation(client, "Browser", htmlPath, message.UserName, message.Message)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println(message.UserName, "used", len(message.Message), "characters", message.Message)
		return c.Status(200).SendString(message.Message)

	})

	app.Listen("localhost:4444")
}
