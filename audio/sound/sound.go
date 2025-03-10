package fetchsound

import (
	myInterface "Takluz_TTS/myInterface"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetSound(message string, Open_AI_Config myInterface.Open_AI_Config, outputPath string) {
	url := "https://api.openai.com/v1/audio/speech"
	body := map[string]string{
		"model": Open_AI_Config.Model,
		"input": message,
		"speed": Open_AI_Config.Speed,
		"voice": Open_AI_Config.Voice,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Open_AI_Config.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code:", resp.StatusCode)
		return
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Println("Error saving response to file:", err)
		return
	}

	fmt.Println("MP3 file saved as", outputPath)
}

func GetSoundBotNoi(message string, BOT_NOI_Config myInterface.BOT_NOI_Config, outputPath string) {
	url := "https://api-voice.botnoi.ai/openapi/v1/generate_audio"
	body := map[string]any{
		"text":       message,
		"speaker":    BOT_NOI_Config.Speaker,
		"volume":     BOT_NOI_Config.Volume,
		"speed":      BOT_NOI_Config.Speed,
		"type_media": BOT_NOI_Config.TypeMedia,
		"save_file":  BOT_NOI_Config.SaveFile,
		"language":   BOT_NOI_Config.Language,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Botnoi-Token", BOT_NOI_Config.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code:", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	audioURL, ok := result["audio_url"].(string)
	if !ok {
		fmt.Println("Error: audio_url not found in response")
		return
	}

	audioResp, err := http.Get(audioURL)
	if err != nil {
		fmt.Println("Error downloading audio file:", err)
		return
	}
	defer audioResp.Body.Close()

	if audioResp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code while downloading audio:", audioResp.StatusCode)
		return
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, audioResp.Body)
	if err != nil {
		fmt.Println("Error saving audio to file:", err)
		return
	}

	fmt.Println("MP3 file saved as", outputPath)
}

func GetSoundResemble(message string, Resemble_Config myInterface.Resemble_Config, outputPath string) bool {
	url := "https://f.cluster.resemble.ai/synthesize"
	body := map[string]any{
		"voice_uuid":    Resemble_Config.VoiceUUID,
		"data":          message,
		"sample_rate":   Resemble_Config.SampleRate,
		"output_format": Resemble_Config.OutputFormat,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Resemble_Config.Key)
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code:", resp.StatusCode)
		return false
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return false
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var responseData map[string]interface{}
	err = json.NewDecoder(reader).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return false
	}

	if success, ok := responseData["success"].(bool); ok && success {
		audioContent, ok := responseData["audio_content"].(string)
		if !ok {
			fmt.Println("Error: 'audio_content' not found in the response.")
			return false
		}

		audioBytes, err := base64.StdEncoding.DecodeString(audioContent)
		if err != nil {
			fmt.Println("Error: Invalid base64 data in audio_content.")
			return false
		}

		outFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return false
		}
		defer outFile.Close()

		_, err = outFile.Write(audioBytes)
		if err != nil {
			fmt.Println("Error saving audio to file:", err)
			return false
		}

		fmt.Println("Audio saved to", outputPath)
		return true
	} else {
		fmt.Printf("Error: Resemble API returned success=false. Issues: %v\n", responseData["issues"])
		return false
	}
}
