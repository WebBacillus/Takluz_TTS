package sound

import (
	"log"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/mediainputs"
)

func ObsInitSound(client *goobs.Client, inputName, filePath string) error {
	_, err := client.Inputs.SetInputSettings(inputs.NewSetInputSettingsParams().WithInputName(inputName).WithInputSettings(map[string]interface{}{
		"local_file": filePath,
	}))
	if err != nil {
		log.Println("Error setting input settings:", err)
		return err
	}
	return nil
}

func ObsPlaySound(client *goobs.Client, inputName string, timeLimit int, filePath string) error {
	_, err := client.MediaInputs.TriggerMediaInputAction(mediainputs.NewTriggerMediaInputActionParams().WithInputName(inputName).WithMediaAction("OBS_WEBSOCKET_MEDIA_INPUT_ACTION_RESTART"))
	if err != nil {
		log.Println("Error triggering media input action:", err)
		return err
	}
	return nil
}
