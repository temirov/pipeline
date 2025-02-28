package steps

import (
	"context"
	"fmt"

	"os"

	"github.com/temirov/pipeline"
)

// PrintFilesStep prints each file's name to stdout.
type PrintFilesStep struct{}

func NewPrintFilesStep() *PrintFilesStep {
	return &PrintFilesStep{}
}

func (step *PrintFilesStep) Name() string {
	return "PrintFiles"
}

func (step *PrintFilesStep) Process(
	executionContext context.Context,
	inputChannel <-chan pipeline.StepResult[os.DirEntry],
	outputChannel chan<- pipeline.StepResult[string],
	pipelineContext *LsContext,
) error {
	for fileResult := range inputChannel {
		fileEntry := fileResult.Item
		// Print file name (you can extend to output more info if desired).
		fmt.Println(fileEntry.Name())
		result := pipeline.ToStepResult(fileEntry.Name())
		select {
		case <-executionContext.Done():
			return executionContext.Err()
		case outputChannel <- result:
			// continue
		}
	}
	return nil
}
