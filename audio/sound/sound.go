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

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	genai "google.golang.org/genai"
)

func GetSoundOpenAI(message string, Open_AI_Config cfg.Open_AI_Config, id string, outputPath string) error {
	url := "https://api.openai.com/v1/audio/speech"

	voice := Open_AI_Config.InstructionSet[id].Voice
	instruction := Open_AI_Config.InstructionSet[id].Instruction

	body := map[string]any{
		"model":           Open_AI_Config.Model,
		"input":           message,
		"speed":           Open_AI_Config.Speed,
		"voice":           voice,
		"response_format": "wav",
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

func GetSoundGoogle(message string, Google_Config cfg.Google_Config, id string, outputPath string) error {

	var voice string
	var instruction string

	switch id {
	case "default":
		voice = "Zephyr"
		instruction = ""
	case "grand_ma":
		voice = "Zephyr"
		instruction =
			`
			Voice Style Prompt: Isaan Grandma (Elderly and Frail)
			Overall Impression: Gentle, kind, and knowledgeable, but the voice should be faint and convey a slight fatigue, as if speaking requires effort. The voice should still offer warmth, but a more subdued, weakened kind of warmth.
			Key Characteristics:
			Pace: slow, deliberate. There may be slightly  silent periods between sentences or phrases.
			Pitch: Mid to low. Faint and slightly tremulous. May have some audible breathiness or an airy quality. Should not be high-pitched or sharp.
			Tone:
			Weak/Fatigued: Clearly convey a weak and tired tone.
			Faint/Subdued: The voice will be faint and subdued, not resonant or clearly loud.
			Kind/Warm: Still retains kindness and warmth, but expressed through a gentle, fragile delivery.
			Knowledgeable: Still holds an essence of wisdom and life experience, but conveyed slowly and with minimal energy.
			Articulation: Less distinct than before. Some words may be mumbled, swallowed, or trail off, sounding incomplete due to overall weakness.
			Isaan Language Influence:
			Rhythm and Melody: Still retains the inherent rhythm and melody of the Isaan language, but significantly slowed down and softened.
			Ending Particles: The use of common Isaan ending particles like "เด้อ" (der), "น้อ" (nor), "แท้" (tae) will still be present, but pronounced softly and languidly.
			Vowels: The pronunciation of Isaan vowels and consonants may sound less full or have more breathiness than typical.
			Breathing: There may be slightly more audible breaths, or breathing that sounds effortful before speaking.
			Vocal Texture: Should have a noticeably raspy, husky, or breathy quality, clearly indicating age and diminished strength. It should not sound clear, crisp, or strong.
			What to Avoid:
			Fast or rushed speaking.
			Loud, resonant, or high-energy vocal delivery.
			A harsh or commanding voice.
			Overly clear or professional articulation, like a newscaster.
			Overall Goal:
			When hearing this voice, the listener should feel as if they are sitting beside a very elderly grandma. She is telling a story or conversing softly, requiring noticeable effort for each word, as if she is resting and slowly sharing her thoughts.
		`
	case "kid":
		voice = "Laomedeia"
		instruction =
			`
			Tone: Hypersonic, excruciatingly shril grating, and utterly relentless, designed to pierce through
			any surrounding noise.

			Emotion: Over-the-top hysteria, pure unadulterated excitement or rage, bordering on a sonic
			assault, demanding absolute and immediate attention.

			Delivery: Extremely high-pitched, like a siren on full blast. The voice constantly operates at the
			very top of a child's vocal range, verging on a painful, almost ultrasonic frequency. Speaking is
			an almost continuous stream of piercing, high-decibel sounds, with individual words often
			indistinguishable due to the sheer volume and pitch. signifies an instantaneous, ear-spliting
			crescendo -a sudden, sharp, almost metallic screech that feels like it's vibrating inside the
			listener's head. Imagine the sound of a very excited, very angry tiny whistle that never stops. 
		`
	case "girl":
		voice = "Despina"
		instruction =
			`
			Affect/Accent: Sultry and inviting, with a playful, slightly breathy quality. Voice should have a warm, slightly lower register, but with the ability to lilt upwards on teases or compliments. Imagine a subtle, knowing smile always present in the voice.
			Tone: Flirtatious, confident, yet with an undercurrent of genuine interest and admiration. Teasing but warm, conveying that he is the focus of her attention and desire. Avoid sounding desperate; aim for captivating.
			Pacing: Deliberate and enticing; slightly slower than normal conversation to build anticipation, with playful hesitations or drawn-out syllables on key words. Allow for moments of suggestive silence.
			Emotion: A blend of playful confidence, allure, intrigue, and clear attraction. The voice should sound like she's enjoying herself and finds him genuinely captivating. There's a hint of "I know what I want, and it might be you."
			Pronunciation: Clear, but with a soft, rounded quality to the words. Emphasize words that convey 
		`
	}

	fmt.Println("Voice: " + voice)
	fmt.Println("Instruction: " + instruction)

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: Google_Config.Key,
	})

	if err != nil {
		log.Fatal(err)
	}

	contents := genai.Text(instruction + "\nSAY: " + message)
	temperature := float32(0.7)
	config := &genai.GenerateContentConfig{
		Temperature:        &temperature,
		ResponseModalities: []string{"AUDIO"},
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: voice,
				},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, Google_Config.Model, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	audioData := resp.Candidates[0].Content.Parts[0].InlineData.Data

	outFile, err := os.Create("google.pcm")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, bytes.NewReader(audioData))
	if err != nil {
		return fmt.Errorf("error saving response to file: %v", err)
	}

	err = FFmpegDecode("google.pcm", outputPath)
	if err != nil {
		return fmt.Errorf("error to decode pcm file: %v", err)
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
