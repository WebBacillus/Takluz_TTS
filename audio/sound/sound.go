package sound

import (
	cfg "Takluz_TTS/cfg"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"google.golang.org/api/option"
)

func GetSoundOpenAI(message string, Open_AI_Config cfg.Open_AI_Config, id string, outputPath string) error {
	url := "https://api.openai.com/v1/audio/speech"

	voice := Open_AI_Config.InstructionSet[id].Voice
	instruction := Open_AI_Config.InstructionSet[id].Instruction

	body := map[string]string{
		"model": Open_AI_Config.Model,
		"input": message,
		"speed": Open_AI_Config.Speed,
		"voice": voice,
	}
	if Open_AI_Config.Model == "gpt-4o-mini-tts" {
		body["voice"] = voice
		body["instructions"] = instruction
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Open_AI_Config.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error: received non-200 response code: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving response to file: %v", err)
	}

	return nil
}

func GetSoundBotNoi(message string, BOT_NOI_Config cfg.BOT_NOI_Config, outputPath string) error {
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
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Botnoi-Token", BOT_NOI_Config.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: received non-200 response code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("error decoding JSON response: %v", err)
	}

	audioURL, ok := result["audio_url"].(string)
	if !ok {
		return fmt.Errorf("error: audio_url not found in response")
	}

	audioResp, err := http.Get(audioURL)
	if err != nil {
		return fmt.Errorf("error downloading audio file: %v", err)
	}
	defer audioResp.Body.Close()

	if audioResp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: received non-200 response code while downloading audio: %d", audioResp.StatusCode)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, audioResp.Body)
	if err != nil {
		return fmt.Errorf("error saving audio to file: %v", err)
	}

	return nil
}

func GetSoundResemble(message string, Resemble_Config cfg.Resemble_Config, outputPath string) error {
	url := "https://f.cluster.resemble.ai/synthesize"
	body := map[string]any{
		"voice_uuid":    Resemble_Config.VoiceUUID,
		"data":          message,
		"sample_rate":   Resemble_Config.SampleRate,
		"output_format": Resemble_Config.OutputFormat,
	}
	// print(message)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Resemble_Config.Key)
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: received non-200 response code: %d", resp.StatusCode)
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("error creating gzip reader: %v", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var responseData map[string]interface{}
	err = json.NewDecoder(reader).Decode(&responseData)
	if err != nil {
		return fmt.Errorf("error decoding JSON response: %v", err)
	}

	if success, ok := responseData["success"].(bool); ok && success {
		audioContent, ok := responseData["audio_content"].(string)
		if !ok {
			return fmt.Errorf("error: 'audio_content' not found in the response")
		}

		audioBytes, err := base64.StdEncoding.DecodeString(audioContent)
		if err != nil {
			return fmt.Errorf("error: invalid base64 data in audio_content: %v", err)
		}

		outFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer outFile.Close()

		_, err = outFile.Write(audioBytes)
		if err != nil {
			return fmt.Errorf("error saving audio to file: %v", err)
		}

		return nil
	} else {
		return fmt.Errorf("error: Resemble API returned success=false. Issues: %v", responseData["issues"])
	}
}

func GetSoundAzure(message string, Microsoft_Config cfg.Microsoft_Config, outputPath string) error {
	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", Microsoft_Config.Region)

	ssml := fmt.Sprintf(`<speak version='1.0' xml:lang='th-TH'>
		<voice xml:lang='th-TH' xml:gender='Female' name='%s'>
			<prosody rate='%s'>
				%s
			</prosody>
		</voice>
	</speak>`, Microsoft_Config.Voice, Microsoft_Config.Speed, message)

	jsonBody := []byte(ssml)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", Microsoft_Config.Key)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("X-Microsoft-OutputFormat", "audio-16khz-128kbitrate-mono-mp3")
	req.Header.Set("User-Agent", "curl")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: received non-200 response code: %s", resp.Status)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving response to file: %v", err)
	}

	return nil

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return fmt.Errorf("error reading response body: %v", err)
	// }

	// err = ioutil.WriteFile(outputPath, body, 0644)
	// if err != nil {
	// 	return fmt.Errorf("error writing to file: %v", err)
	// }

	// fmt.Println("Audio content written to file:", outputPath)
	// return nil
}

func GetSoundGoogle(message string, Google_Config cfg.Google_Config, outputPath string) error {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx, option.WithAPIKey(Google_Config.Key))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: message},
		},
		// Build the voice request, select the language code ("th-TH") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "th-TH",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
			Name:         Google_Config.Name,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:   texttospeechpb.AudioEncoding_LINEAR16,
			SpeakingRate:    Google_Config.SpeakingRate,
			Pitch:           Google_Config.Pitch,
			VolumeGainDb:    Google_Config.VolumeGainDb,
			SampleRateHertz: 48000,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return fmt.Errorf("error synthesizing speech: %v", err)
	}

	err = os.WriteFile(outputPath, resp.AudioContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing audio content to file: %v", err)
	}

	return nil
}

func PlayAnimation(client *goobs.Client, inputName, htmlDirectory string, userName string, message string) error {
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
		return fmt.Errorf("error creating file: %v", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	_, err = client.Inputs.SetInputSettings(inputs.NewSetInputSettingsParams().WithInputName(inputName).WithInputSettings(map[string]interface{}{
		"url": "file://" + htmlDirectory + "/print.html",
	}))
	if err != nil {
		return fmt.Errorf("error setting input settings: %v", err)
	}

	_, err = client.Inputs.SetInputSettings(inputs.NewSetInputSettingsParams().WithInputName(inputName).WithInputSettings(map[string]interface{}{
		"refresh": true,
	}))
	if err != nil {
		return fmt.Errorf("error refreshing input settings: %v", err)
	}

	time.AfterFunc(8*time.Second, func() {
		_, err := client.Inputs.SetInputSettings(inputs.NewSetInputSettingsParams().WithInputName(inputName).WithInputSettings(map[string]interface{}{
			"url": "",
		}))
		if err != nil {
			log.Println("Failed to turn off animation:", err)
		}
	})

	return nil
}
