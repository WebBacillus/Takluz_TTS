package fetchsound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetSound(message string, OPEN_API_KEY string, outputPath string) {

	url := "https://api.openai.com/v1/audio/speech"
	body := map[string]string{
		"model": "tts-1-hd",
		"input": message,
		"speed": "0.8",   //เปลียนความเร็ว
		"voice": "alloy", //เปลี่ยนเสียงคนพูด
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
	req.Header.Set("Authorization", "Bearer "+OPEN_API_KEY)

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

func GetSoundBotNoi(message string, OPEN_API_KEY string, outputPath string) {
	url := "https://api-voice.botnoi.ai/openapi/v1/generate_audio"
	body := map[string]any{
		"text":       message,
		"speaker":    "8",
		"volume":     0.8,
		"speed":      0.8,
		"type_media": "mp3",
		"save_file":  true,
		"language":   "th",
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
	req.Header.Set("Botnoi-Token", "WWVQcHVjb1pHRFB2OEFCa2ZOMlVHU24xQ2hWMjU2MTg5NA==")

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
