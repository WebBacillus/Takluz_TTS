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
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

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
	fmt.Println("---------------------------------------------------")

	currentVersion := "v1.0.1"
	if release.TagName == currentVersion {
		fmt.Println(color.GreenString("You are using the latest version:"), currentVersion)
	} else {
		fmt.Println(color.RedString("New version available:"), release.HTMLURL)
	}

}

func waitForExit() {
	fmt.Println("Press 'Enter' to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
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

	Open_AI_Config, err := cfg.InitOpenAIConfig()
	if err != nil {
		log.Println(err)
		waitForExit()
		return
	}

	// OBS_Config, err := cfg.InitOBSConfig()
	// if err != nil {
	// 	log.Println(err)
	// 	waitForExit()
	// 	return
	// }

	BOT_NOI_Config, err := cfg.InitBotNoiConfig()
	if err != nil {
		log.Println(err)
		waitForExit()
		return
	}

	Resemble_config, err := cfg.InitResembleConfig()
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

	/*
		goobsClient, err := goobs.New(OBS_Config.URL, goobs.WithPassword(OBS_Config.Key))
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

		fmt.Println("Successfully connected to OBS")
	*/

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
	collection := database.Collection("takluz-tts")

	command.CreateSilentAudio()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("END")
		command.CreateSilentAudio()
		os.Exit(0)
	}()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Post("/", func(c *fiber.Ctx) error {

		type Message struct {
			UserName string `json:"userName"`
			Message  string `json:"message"`
		}

		var message Message
		if err := c.BodyParser(&message); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		collection.InsertOne(ctx, message)

		if len(message.Message) >= General_Config.LimitToken {
			message.Message = message.Message[:General_Config.LimitToken]
		}

		if General_Config.AI == "BOT_NOI" {
			sound.GetSoundBotNoi(message.Message, BOT_NOI_Config, "speech.mp3")
		} else if General_Config.AI == "OPEN_AI" {
			sound.GetSound(message.Message, Open_AI_Config, "speech.mp3")
		} else if General_Config.AI == "RESEMBLE" {
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

		// err = sound.ObsPlaySound(goobsClient, OBS_Config.Media, General_Config.TimeLimit, getPath("speech.mp3"))
		err = sound.FFplayAudio(getPath("speech.mp3"))
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
	if err != nil {
		log.Println(err)
		waitForExit()
		return
	}
}
