package workflow

import (
	"encoding/json"
	"io"
	"os"
)

type Transformer struct {
	Name   string `json:"transformer"`
	Output string `json:"output,omitempty"`
	Input  *struct {
		Item   string `json:"item,omitempty"`
		ItemNo *int   `json:"itemNo,omitempty"`
	} `json:"input,omitempty"`
}

type Workflow struct {
	Name            string        `json:"name"`
	Transformations []Transformer `json:"transformations"`
}

func NewWorkflow(filename string) (*Workflow, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var workflow Workflow
	json.Unmarshal(bytes, &workflow)

	return &workflow, nil
}
