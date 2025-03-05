package prefix

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConcatAudio(inputFiles []string, outputFile string) error {
	for _, file := range inputFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", file)
		}
	}

	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	concatString := "concat:" + strings.Join(inputFiles, "|")

	args := []string{
		"-y",
		"-i", concatString,
		"-c:a", "libmp3lame",
		outputFile,
	}
	cmd := exec.Command("ffmpeg", args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to concatenate audio files: %v\nOutput: %s", err, string(output))
	}

	return nil
}
