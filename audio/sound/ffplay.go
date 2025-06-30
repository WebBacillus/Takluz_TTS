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

func FFmpegDecode(intputFilePath string, outputFilePath string) error {
	args := []string{
		"-f",
		"s16le",
		"-ar",
		"24000",
		"-ac",
		"1",
		"-i",
		intputFilePath,
		outputFilePath,
		"-y",
	}
	cmd := exec.Command("ffmpeg", args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to decode audio file: %v\nOutput: %s", err, string(output))
	}

	return nil
}
