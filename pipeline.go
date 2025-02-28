package pipeline

import (
	"context"
	"fmt"
)

const defaultBufferSize = 100

// StepResult represents a processed item with a generic type.
type StepResult[OutputType any] struct {
	Item       OutputType
	IsSkipped  bool
	Error      error
	ShouldExit bool
	ExitReason string
}

// Step is a unified interface for pipeline steps.
// When the input channel is nil, the step acts as a source.
type Step[InputType, OutputType, ContextType any] interface {
	// Process receives an input channel (or nil for sources) and an output channel provided by the pipeline framework.
	// The step writes its results to outputChannel and returns an error if processing fails.
	Process(
		executionContext context.Context,
		inputChannel <-chan StepResult[InputType],
		outputChannel chan<- StepResult[OutputType],
		pipelineContext ContextType,
	) error
	Name() string
}

// PipelineBuilder builds a pipeline while enforcing type consistency between steps.
// ContextType is the shared context type; CurrentOutputType is the output type of the last step.
type PipelineBuilder[ContextType, CurrentOutputType any] struct {
	stepFunctions []func(
		executionContext context.Context,
		inputData any,
		pipelineContext ContextType,
	) (any, error)
}

// NewPipelineBuilder creates a new builder starting with a source step.
// The source step ignores its input (it receives nil).
func NewPipelineBuilder[ContextType, OutputType any](
	sourceStep Step[any, OutputType, ContextType],
) *PipelineBuilder[ContextType, OutputType] {
	wrappedSourceFunction := func(
		executionContext context.Context,
		ignoredInput any,
		pipelineContext ContextType,
	) (any, error) {
		outputChannel := make(chan StepResult[OutputType], defaultBufferSize)
		go func() {
			_ = sourceStep.Process(executionContext, nil, outputChannel, pipelineContext)
			close(outputChannel)
		}()
		return outputChannel, nil
	}
	return &PipelineBuilder[ContextType, OutputType]{
		stepFunctions: []func(
			executionContext context.Context,
			inputData any,
			pipelineContext ContextType,
		) (any, error){wrappedSourceFunction},
	}
}

// AppendStep adds a new transformation step to the pipeline.
// NewOutputType is the output type of the new step.
func AppendStep[
	ContextType, CurrentOutputType, NewOutputType any,
](
	builder *PipelineBuilder[ContextType, CurrentOutputType],
	transformationStep Step[CurrentOutputType, NewOutputType, ContextType],
) *PipelineBuilder[ContextType, NewOutputType] {
	wrappedTransformationFunction := func(
		executionContext context.Context,
		inputData any,
		pipelineContext ContextType,
	) (any, error) {
		typedInputChannel, ok := inputData.(chan StepResult[CurrentOutputType])
		if !ok {
			return nil, fmt.Errorf("invalid input type for step %s", transformationStep.Name())
		}
		outputChannel := make(chan StepResult[NewOutputType], defaultBufferSize)
		go func() {
			_ = transformationStep.Process(executionContext, typedInputChannel, outputChannel, pipelineContext)
			close(outputChannel)
		}()
		return outputChannel, nil
	}
	newStepFunctions := make([]func(
		executionContext context.Context,
		inputData any,
		pipelineContext ContextType,
	) (any, error), len(builder.stepFunctions))
	copy(newStepFunctions, builder.stepFunctions)
	newStepFunctions = append(newStepFunctions, wrappedTransformationFunction)
	return &PipelineBuilder[ContextType, NewOutputType]{stepFunctions: newStepFunctions}
}

// Pipeline represents an executable pipeline.
// ContextType is the shared context and FinalOutputType is the final output type.
type Pipeline[ContextType, FinalOutputType any] struct {
	stepFunctions []func(
		executionContext context.Context,
		inputData any,
		pipelineContext ContextType,
	) (any, error)
}

// Build converts a PipelineBuilder into an executable Pipeline.
func (builder *PipelineBuilder[ContextType, FinalOutputType]) Build() *Pipeline[ContextType, FinalOutputType] {
	return &Pipeline[ContextType, FinalOutputType]{stepFunctions: builder.stepFunctions}
}

// Execute runs all the steps in the pipeline sequentially.
// It asserts the final output channel as a bidirectional channel and drains it.
func (pipelineInstance *Pipeline[ContextType, FinalOutputType]) Execute(
	executionContext context.Context,
	pipelineContext ContextType,
) error {
	var currentOutput any = nil // start with nil input for the source step
	var functionError error

	for _, stepFunction := range pipelineInstance.stepFunctions {
		currentOutput, functionError = stepFunction(executionContext, currentOutput, pipelineContext)
		if functionError != nil {
			return functionError
		}
	}

	// Assert the final output to be a bidirectional channel.
	finalOutputChannel, ok := currentOutput.(chan StepResult[FinalOutputType])
	if !ok {
		return fmt.Errorf("unexpected final output type")
	}

	// Drain the final output channel.
	for range finalOutputChannel {
		// Optionally process final items.
	}
	return nil
}

// ToStepResult wraps a given item into a StepResult.
// It can be extended to set default flags if needed.
func ToStepResult[ItemType any](item ItemType) StepResult[ItemType] {
	return StepResult[ItemType]{Item: item}
}
