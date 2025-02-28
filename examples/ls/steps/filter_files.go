package steps

import (
	"context"
	"path/filepath"

	"github.com/temirov/pipeline"
	"os"
)

// FilterFilesStep filters files based on the wildcard pattern in the context.
type FilterFilesStep struct{}

func NewFilterFilesStep() *FilterFilesStep {
	return &FilterFilesStep{}
}

func (step *FilterFilesStep) Name() string {
	return "FilterFiles"
}

func (step *FilterFilesStep) Process(
	executionContext context.Context,
	inputChannel <-chan pipeline.StepResult[os.DirEntry],
	outputChannel chan<- pipeline.StepResult[os.DirEntry],
	pipelineContext *LsContext,
) error {
	pattern := pipelineContext.Wildcard
	for fileResult := range inputChannel {
		fileEntry := fileResult.Item
		// If a non-trivial pattern is provided, filter based on it.
		if pattern != "" && pattern != "*" {
			matches, err := filepath.Match(pattern, fileEntry.Name())
			if err != nil || !matches {
				continue
			}
		}
		result := pipeline.ToStepResult(fileEntry)
		select {
		case <-executionContext.Done():
			return executionContext.Err()
		case outputChannel <- result:
			// continue
		}
	}
	return nil
}
