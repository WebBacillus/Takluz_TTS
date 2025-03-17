package sound

import (
	"fmt"
	"os/exec"
)

func FFplayAudio(filePath string) error {
	args := []string{
		"-nodisp",
		"-autoexit",
		filePath,
	}
	cmd := exec.Command("ffplay", args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to play audio file: %v\nOutput: %s", err, string(output))
	}

	return nil
}
