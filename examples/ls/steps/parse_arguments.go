package steps

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/temirov/pipeline"
)

// LsContext holds the shared context for the ls pipeline.
type LsContext struct {
	Wildcard string
}

// ParseArgumentsStep reads command-line arguments and sets the wildcard in the context.
type ParseArgumentsStep struct{}

func NewParseArgumentsStep() *ParseArgumentsStep {
	return &ParseArgumentsStep{}
}

func (step *ParseArgumentsStep) Name() string {
	return "ParseArguments"
}

func (step *ParseArgumentsStep) Process(
	executionContext context.Context,
	_ <-chan pipeline.StepResult[any],
	outputChannel chan<- pipeline.StepResult[string],
	pipelineContext *LsContext,
) error {
	if len(os.Args) > 1 {
		rawArgument := os.Args[1]
		// If the argument contains any wildcard characters, use it directly.
		if strings.ContainsAny(rawArgument, "*?[") {
			pipelineContext.Wildcard = rawArgument
		} else {
			// Otherwise, if the argument is an existing file, reconstruct the intended wildcard
			// by using its extension. For example, "content.sh" becomes "*.sh".
			if fileInfo, err := os.Stat(rawArgument); err == nil && !fileInfo.IsDir() {
				extension := filepath.Ext(rawArgument)
				if extension != "" {
					pipelineContext.Wildcard = "*" + extension
				} else {
					pipelineContext.Wildcard = rawArgument
				}
			} else {
				pipelineContext.Wildcard = rawArgument
			}
		}
	} else {
		pipelineContext.Wildcard = "*" // default: match all
	}
	result := pipeline.ToStepResult("arguments parsed")
	select {
	case <-executionContext.Done():
		return executionContext.Err()
	case outputChannel <- result:
		return nil
	}
}
