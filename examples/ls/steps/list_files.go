package steps

import (
	"context"
	"os"

	"github.com/temirov/pipeline"
)

// ListFilesStep reads the current directory and outputs each file entry.
type ListFilesStep struct{}

func NewListFilesStep() *ListFilesStep {
	return &ListFilesStep{}
}

func (step *ListFilesStep) Name() string {
	return "ListFiles"
}

func (step *ListFilesStep) Process(
	executionContext context.Context,
	inputChannel <-chan pipeline.StepResult[string],
	outputChannel chan<- pipeline.StepResult[os.DirEntry],
	pipelineContext *LsContext,
) error {
	// Consume the dummy input from the previous step.
	<-inputChannel

	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		result := pipeline.ToStepResult(entry)
		select {
		case <-executionContext.Done():
			return executionContext.Err()
		case outputChannel <- result:
			// continue
		}
	}
	return nil
}
