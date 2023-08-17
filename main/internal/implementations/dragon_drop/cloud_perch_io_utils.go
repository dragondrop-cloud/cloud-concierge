package dragonDrop

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	terraformWorkspace "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_workspace"
)

// readFile reads a file and returns the data as a byte array.
func readFile(filename string) ([]byte, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, fmt.Errorf("[error reading output file]%w", err)
	}

	return fileBytes, nil
}

// readOutputFileAsMap reads a file and returns the data as a map.
func readOutputFileAsMap(filename string) (map[string]interface{}, error) {
	fileBytes, err := readFile(fmt.Sprintf("outputs/%s", filename))
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("[error reading output file]%w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("[error unmarshalling output file]%w", err)
	}

	return data, nil
}

// readOutputFileAsSlice reads a file and returns the data as a slice.
func readOutputFileAsSlice(filename string) ([]interface{}, error) {
	fileBytes, err := readFile(fmt.Sprintf("outputs/%s", filename))
	if err != nil {
		return []interface{}{}, fmt.Errorf("[error reading output file]%w", err)
	}

	var data []interface{}
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return []interface{}{}, fmt.Errorf("[error unmarshalling output file]%w", err)
	}

	return data, nil
}

// getAllTFFiles searches a directory for all terraform files within the user workspace directories.
func getAllTFFiles(ctx context.Context, directories terraformWorkspace.WorkspaceDirectoriesDecoder) []string {
	tfFiles := make([]string, 0)

	for _, directory := range directories {
		files, err := os.ReadDir(fmt.Sprintf("repo/%s", directory))
		if err != nil {
			return make([]string, 0)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if strings.HasSuffix(file.Name(), ".tf") {
				tfFiles = append(tfFiles, fmt.Sprintf("repo/%s/%s", directory, file.Name()))
			}
		}
	}

	return tfFiles
}
