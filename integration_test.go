package pipeline

import (
	"context"
	"testing"
	"time"
)

// TestContext is a simple context type for our test pipelines.
type TestContext struct {
	Value string
}

// HelloWorld outputs a single message: "Hello, World!".
// As a source, it ignores its input.
type HelloWorld struct {
	sourceName string
}

// NewHelloWorldSource creates a new HelloWorldSource.
func NewHelloWorldSource() *HelloWorld {
	return &HelloWorld{sourceName: "HelloWorldSource"}
}

// Name returns the name of the source step.
func (helloWorldSource *HelloWorld) Name() string {
	return helloWorldSource.sourceName
}

func (helloWorldSource *HelloWorld) Hi() string {
	return "Hi, World!"
}

// Process implements the Step interface for a source step.
// It sends a single "Hello, World!" message to the provided output channel.
func (helloWorldSource *HelloWorld) Process(
	executionContext context.Context,
	_ <-chan StepResult[any],
	outputChannel chan<- StepResult[string],
	pipelineContext TestContext,
) error {
	// result := StepResult[string]{Item: helloWorldSource.Hi()}
	result := ToStepResult(helloWorldSource.Hi())
	select {
	case <-executionContext.Done():
		return executionContext.Err()
	case outputChannel <- result:
		return nil
	}
}

// TestSingleStepHelloWorld verifies that a pipeline with a single source step executes without error.
func TestSingleStepHelloWorld(testInstance *testing.T) {
	// Create a context with timeout for the pipeline execution.
	executionContext, cancelFunction := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunction()

	// Create the HelloWorld source step.
	helloWorldSource := NewHelloWorldSource()

	// Build a pipeline starting with the HelloWorld source.
	builderInstance := NewPipelineBuilder(helloWorldSource)
	pipelineInstance := builderInstance.Build()

	// Define a test pipeline context.
	pipelineExecutionContext := TestContext{Value: "test"}

	// Execute the pipeline.
	executionError := pipelineInstance.Execute(executionContext, pipelineExecutionContext)
	if executionError != nil {
		testInstance.Fatalf("Failed to execute pipeline: %v", executionError)
	}

	testInstance.Log("Pipeline executed successfully with Hello World source")
}
