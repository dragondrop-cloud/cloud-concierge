package documentize

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/Jeffail/gabs/v2"
)

// AllWorkspaceStatesToDocuments converts all workspace states to documents of non-sensitive strings.
func (d *documentize) AllWorkspaceStatesToDocuments(workspaceToDirectory map[string]string) (map[Workspace][]byte, error) {
	outputWorkspaceToDocument := map[Workspace][]byte{}

	for workspace := range workspaceToDirectory {
		doc, err := d.WorkspaceStateToDocument(Workspace(workspace))
		if err != nil {
			return nil, fmt.Errorf("[WorkspaceStateToDocument] Error while documentizing %v: %v", workspace, err)
		}
		outputWorkspaceToDocument[Workspace(workspace)] = doc
	}

	return outputWorkspaceToDocument, nil
}

// ConvertWorkspaceDocumentsToJSON converts the output of AllWorkspaceStatesToDocuments to a json-format byte array.
func (d *documentize) ConvertWorkspaceDocumentsToJSON(workspaceDocMap map[Workspace][]byte) ([]byte, error) {
	jsonObj := gabs.New()
	for workspace, doc := range workspaceDocMap {
		_, err := jsonObj.Set(string(doc), string(workspace))
		if err != nil {
			return nil, fmt.Errorf(
				"[ConvertWorkpsaceDocumentsToJSON] error in jsonObj.Set() for Workspace %v and doc %v: %v",
				string(workspace), string(doc), err)
		}
	}
	return jsonObj.Bytes(), nil
}

// WorkspaceStateToDocument converts a workspace state to a document of non-sensitive strings.
func (d *documentize) WorkspaceStateToDocument(workspace Workspace) ([]byte, error) {
	tfState, err := os.ReadFile(fmt.Sprintf("state_files/%v.json", string(workspace)))
	if err != nil {
		return nil, fmt.Errorf("Error reading in terraform state file for workspace %v: %v", workspace, err)
	}

	tfStateParsed, err := gabs.ParseJSON(tfState)
	if err != nil {
		return nil, fmt.Errorf("Error parsing terraform state file for workspace %v via gabs: %v", workspace, err)
	}

	workspaceDoc, err := d.workspaceDocFromTFState(tfStateParsed)
	if err != nil {
		return nil, fmt.Errorf("[resourceDetailsFromTFState] Error extracting resource details: %v", err)
	}

	return []byte(workspaceDoc), nil
}

// workspaceDocFromTFState produces a document from tfStateParsed by scanning and extracting information
// from each resource.
func (d *documentize) workspaceDocFromTFState(tfStateParsed *gabs.Container) (string, error) {
	resourceDetails := ""

	i := 0
	for tfStateParsed.Exists("resources", strconv.Itoa(i)) {
		j := 0
		for tfStateParsed.Exists("resources", strconv.Itoa(i), "instances", strconv.Itoa(j)) {
			currentMode := tfStateParsed.Search("resources", strconv.Itoa(i), "mode").Data().(string)
			if currentMode != "managed" {
				j++
				continue
			}

			currentResourceDetails, err := d.extractResourceDocument(tfStateParsed, false, i, j)
			if err != nil {
				return "", fmt.Errorf("[d.extractResourceDocument] Error: %v", err)
			}

			resourceDetails += currentResourceDetails

			j++
		}

		i++
	}

	return resourceDetails, nil
}

// regexProviderName extracts the terraform provider name from a string.
func regexProviderName(rawProvider string) (string, error) {
	r, err := regexp.Compile(`\.?(provider.*]).*$`)
	if err != nil {
		return "", fmt.Errorf("Error in regexp.Compile: %v", err)
	}

	output := r.FindStringSubmatch(rawProvider)[1]

	return output, nil
}
