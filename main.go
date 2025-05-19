package main

import (
	command "Takluz_TTS/audio/prefix"
	"Takluz_TTS/audio/sound"
	cfg "Takluz_TTS/cfg"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getPath(nextPath string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	newPath := filepath.Join(currentDir, nextPath)
	return filepath.ToSlash(newPath)
}

func checkNewRelease() {
	resp, err := http.Get("https://api.github.com/repos/webbacillus/Takluz_TTS/releases/latest")
	if err != nil {
		log.Println("Error fetching latest release:", err)
		return
	}
	defer resp.Body.Close()

	var release struct {
		HTMLURL string `json:"html_url"`
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Println("Error decoding release response:", err)
		return
	}

	fmt.Println("---------------------------------------------------")
	fmt.Println("  Takluz TTS - Text-to-Speech Application")
	fmt.Println("  Made by:  WebBacillus (https://github.com/webbacillus)")
	fmt.Println("  Contact:  Web.pasit.kh@gmail.com")
	fmt.Println("  Phone:    +66918683540")
	fmt.Println("---------------------------------------------------")

	currentVersion := "v1.0.4"
	if release.TagName == currentVersion {
		fmt.Println(color.GreenString("You are using the latest version:"), currentVersion)
	} else {
		fmt.Println(color.RedString("New version available:"), release.HTMLURL)
	}

}

func waitForExit() {
	command.CreateSilentAudio()
	fmt.Println("Press 'Enter' to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Job struct for queueing POST requests
// Each job holds a function to execute and a channel for the result

type PostJob struct {
	Ctx     *fiber.Ctx
	Handler func(*fiber.Ctx) error
	Result  chan error
}

var postJobQueue = make(chan PostJob, 100)

func startPostJobWorker() {
	go func() {
		for job := range postJobQueue {
			err := job.Handler(job.Ctx)
			job.Result <- err
		}
	}()
}

func main() {
	checkNewRelease()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("fatal error config file: %s", err))
		waitForExit()
		return
	}

	OBS_Config, err := cfg.InitOBSConfig()
	if err != nil {
		log.Println(err)
		waitForExit()
		return
	}

	General_Config, err := cfg.InitGeneralConfig()
	if err != nil {
		log.Println(err)
		waitForExit()
		return
	}

	var goobsClient *goobs.Client
	if General_Config.Player == "OBS" {
		goobsClient, err = goobs.New(OBS_Config.URL, goobs.WithPassword(OBS_Config.Key))
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				log.Println("Error: Unable to connect to OBS. Please ensure OBS is running and the WebSocket server is enabled.")
			} else {
				log.Println("Error: Invalid OBS WebSocket URL or password.")
			}
			waitForExit()
			return
		}
		defer goobsClient.Disconnect()
		err = sound.ObsInitSound(goobsClient, OBS_Config.Media, getPath("speech.wav"))
		if err != nil {
			log.Println(err.Error())
		}

		fmt.Println("Successfully connected to OBS")
	}

	var collection *mongo.Collection

	if General_Config.DataCollect {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		uri := "mongodb+srv://user:ymRLJfpzc5Hy9whv@takluz-tts.y1hqc.mongodb.net/?retryWrites=true&w=majority&appName=takluz-tts"
		opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		mongoClient, err := mongo.Connect(ctx, opts)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err = mongoClient.Disconnect(context.TODO()); err != nil {
				log.Println(err.Error())
			}
		}()

		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Println(err.Error())
		} else {
			fmt.Println("Successfully connected to MongoDB")
		}

		database := mongoClient.Database("takluz")
		collection = database.Collection("takluz-tts")
	}

	command.CreateSilentAudio()

	type Message struct {
		UserName string `json:"userName"`
		Message  string `json:"message"`
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use("/", func(c *fiber.Ctx) error {

		var message Message
		if err := c.BodyParser(&message); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if len(message.Message) >= General_Config.LimitToken {
			message.Message = message.Message[:General_Config.LimitToken]
		}

		c.Locals("message", message)
		return c.Next()
	})

	startPostJobWorker()

	// Helper to wrap POST handlers into the queue
	wrapPost := func(handler func(*fiber.Ctx) error) fiber.Handler {
		return func(c *fiber.Ctx) error {
			result := make(chan error)
			postJobQueue <- PostJob{Ctx: c, Handler: handler, Result: result}
			err := <-result
			return err
		}
	}

	app.Post("/open_ai/:id", wrapPost(func(c *fiber.Ctx) error {
		id := c.Params("id")
		Open_AI_Config, err := cfg.InitOpenAIConfig()
		if id == "default" {
			Open_AI_Config.Model = "tts-1"
		} else {
			Open_AI_Config.Model = "gpt-4o-mini-tts"
		}
		if err != nil {
			log.Println(err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		message := c.Locals("message").(Message)
		msgCollectoin := struct {
			UserName  string    `json:"userName"`
			Message   string    `json:"message"`
			ID        string    `json:"id"`
			Timestamp time.Time `json:"timestamp"`
		}{
			UserName:  message.UserName,
			Message:   message.Message,
			ID:        id,
			Timestamp: time.Now(),
		}

		// Use the correct context for MongoDB insert
		if General_Config.DataCollect && collection != nil {
			_, err = collection.InsertOne(context.TODO(), msgCollectoin)
			if err != nil {
				log.Println("Error inserting message into collection:", err)
			}
		}
		err = sound.GetSoundOpenAI(message.Message, Open_AI_Config, id, "speech.wav")
		if err != nil {
			log.Println(err.Error())
		}

		if General_Config.Player == "OBS" {
			err = sound.ObsPlaySound(goobsClient, OBS_Config.Media, General_Config.TimeLimit, getPath("speech.wav"))
		} else if General_Config.Player == "FFPLAY" {
			err = sound.FFplayAudio(getPath("speech.wav"))
		}
		if err != nil {
			log.Println(err.Error())
		}

		fmt.Println(color.GreenString(message.UserName), "used", color.RedString(fmt.Sprintf("%d", len(message.Message))), "characters", message.Message, "with ID:", id)
		return c.Status(200).SendString(fmt.Sprintf("Message: %s, ID: %s", message.Message, id))
	}))

	app.Post("/bot_noi", wrapPost(func(c *fiber.Ctx) error {
		BOT_NOI_Config, err := cfg.InitBotNoiConfig()
		if err != nil {
			log.Println(err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		message := c.Locals("message").(Message)
		err = sound.GetSoundBotNoi(message.Message, BOT_NOI_Config, "speech.wav")
		if err != nil {
			log.Println(err.Error())
		}

		if General_Config.Player == "OBS" {
			err = sound.ObsPlaySound(goobsClient, OBS_Config.Media, General_Config.TimeLimit, getPath("speech.wav"))
		} else if General_Config.Player == "FFPLAY" {
			err = sound.FFplayAudio(getPath("speech.wav"))
		}
		if err != nil {
			log.Println(err.Error())
		}

		fmt.Println(color.GreenString(message.UserName), "used", color.RedString(fmt.Sprintf("%d", len(message.Message))), "characters", message.Message)
		return c.Status(200).SendString(message.Message)
	}))

	app.Post("/resemble", wrapPost(func(c *fiber.Ctx) error {
		Resemble_config, err := cfg.InitResembleConfig()
		if err != nil {
			log.Println(err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		message := c.Locals("message").(Message)
		text := fmt.Sprintf(`<speak>
                <prosody rate="%s" pitch="x-high">
                    <lang xml:lang="th-th">
                            %s
                    </lang>
                </prosody>
        </speak>`, Resemble_config.Speed, message.Message)
		err = sound.GetSoundResemble(text, Resemble_config, "speech.wav")
		if err != nil {
			log.Println(err.Error())
		}

		if General_Config.Player == "OBS" {
			err = sound.ObsPlaySound(goobsClient, OBS_Config.Media, General_Config.TimeLimit, getPath("speech.wav"))
		} else if General_Config.Player == "FFPLAY" {
			err = sound.FFplayAudio(getPath("speech.wav"))
		}
		if err != nil {
			log.Println(err.Error())
		}

		fmt.Println(color.GreenString(message.UserName), "used", color.RedString(fmt.Sprintf("%d", len(message.Message))), "characters", message.Message)
		return c.Status(200).SendString(message.Message)
	}))

	app.Post("/azure", wrapPost(func(c *fiber.Ctx) error {
		Microsoft_Config, err := cfg.InitMicrosoftConfig()
		if err != nil {
			log.Println(err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		message := c.Locals("message").(Message)
		err = sound.GetSoundAzure(message.Message, Microsoft_Config, "speech.wav")
		if err != nil {
			log.Println(err)
		}

		if General_Config.Player == "OBS" {
			err = sound.ObsPlaySound(goobsClient, OBS_Config.Media, General_Config.TimeLimit, getPath("speech.wav"))
		} else if General_Config.Player == "FFPLAY" {
			err = sound.FFplayAudio(getPath("speech.wav"))
		}
		if err != nil {
			log.Println(err)
		}

		fmt.Println(color.GreenString(message.UserName), "used", color.RedString(fmt.Sprintf("%d", len(message.Message))), "characters", message.Message)
		return c.Status(200).SendString(message.Message)
	}))

	app.Post("/google", wrapPost(func(c *fiber.Ctx) error {
		Google_Config, err := cfg.InitGoogleConfig()
		if err != nil {
			log.Println(err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		message := c.Locals("message").(Message)
		err = sound.GetSoundGoogle(message.Message, Google_Config, "speech.wav")
		if err != nil {
			log.Println(err)
		}

		if General_Config.Player == "OBS" {
			err = sound.ObsPlaySound(goobsClient, OBS_Config.Media, General_Config.TimeLimit, getPath("speech.wav"))
		} else if General_Config.Player == "FFPLAY" {
			err = sound.FFplayAudio(getPath("speech.wav"))
		}
		if err != nil {
			log.Println(err)
		}

		fmt.Println(color.GreenString(message.UserName), "used", color.RedString(fmt.Sprintf("%d", len(message.Message))), "characters", message.Message)
		return c.Status(200).SendString(message.Message)
	}))

	err = app.Listen("localhost:4444")
	if err != nil {
		log.Println(err)
		waitForExit()
		return
	}
}
