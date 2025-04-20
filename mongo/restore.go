package mongo

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Naman-B-Parlecha/BullStash/util"
)

func Restore(mongoURI, inputPath string, dropBeforeRestore, isCompressed bool) error {
	if mongoURI == "" {
		errorMsg := "MongoDB URI cannot be empty"
		util.CallWebHook(errorMsg, true)
		return fmt.Errorf("%s", errorMsg)
	}

	if inputPath == "" {
		errorMsg := "input path cannot be empty"
		util.CallWebHook(errorMsg, true)
		return fmt.Errorf("%s", errorMsg)
	}

	fileInfo, err := os.Stat(inputPath)
	if os.IsNotExist(err) {
		errorMsg := fmt.Sprintf("input path does not exist: %s", inputPath)
		util.CallWebHook(errorMsg, true)
		return fmt.Errorf("%s", errorMsg)
	}

	var cmd *exec.Cmd
	args := []string{
		"--uri=" + mongoURI,
	}

	if dropBeforeRestore {
		args = append(args, "--drop")
	}

	if isCompressed {
		args = append(args, "--gzip")
	}

	if fileInfo.IsDir() {
		args = append(args, "--dir="+inputPath)
	} else {
		args = append(args, "--archive="+inputPath)
	}

	cmd = exec.Command("mongorestore", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMsg := fmt.Sprintf("mongorestore failed: %v\nOutput: %s", err, string(output))
		util.CallWebHook(errorMsg, true)
		return fmt.Errorf("%s", errorMsg)

	}

	successMsg := fmt.Sprintf("MongoDB restoration completed successfully from: %s", inputPath)
	util.CallWebHook(successMsg, false)
	return nil
}
