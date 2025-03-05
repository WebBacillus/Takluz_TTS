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
