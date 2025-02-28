# Pipeline

[![Go Reference](https://pkg.go.dev/badge/github.com/temirov/pipeline.svg)](https://pkg.go.dev/github.com/temirov/pipeline)
[![Go Report Card](https://goreportcard.com/badge/github.com/temirov/pipeline)](https://goreportcard.com/report/github.com/temirov/pipeline)

A flexible, composable, and type-safe data processing pipeline framework for Go.

## Features

- **Type Safety**: Uses Go generics to ensure type safety across the pipeline.
- **Composable**: Build pipelines by chaining independent processing steps.
- **Concurrent**: Built-in support for parallel processing.
- **Context-Aware**: Propagates cancellation signals and timeouts through the pipeline.
- **Flexibility**: Process any type of data through the pipeline.
- **Error Handling**: Comprehensive error handling with detailed logging.
- **Smart Wildcard Reconstruction**: In the provided examples (such as the Unix `ls` pipeline), if the shell expands a wildcard into a single file argument, the system reconstructs the intended pattern (e.g. turning `content.sh` into `*.sh`).

## Installation

```bash
go get github.com/temirov/pipeline
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/temirov/pipeline"
)

// Define your data type.
type Item struct {
	ID    int
	Value string
}

// Define your pipeline context.
type ProcessContext struct {
	Config map[string]string
}

// Example source step that generates data.
type SourceStep struct {
	pipeline.BaseStep[Item, ProcessContext]
}

func NewSourceStep() *SourceStep {
	return &SourceStep{
		BaseStep: pipeline.NewBaseStep[Item, ProcessContext]("Source"),
	}
}

func (s *SourceStep) Process(ctx context.Context, output chan<- pipeline.StepResult[Item], pipelineContext ProcessContext) {
	defer close(output)
	// Generate some sample data.
	for i := 1; i <= 5; i++ {
		output <- pipeline.StepResult[Item]{
			Item: Item{
				ID:    i,
				Value: fmt.Sprintf("Item %d", i),
			},
		}
	}
}

// Example processing step that transforms data.
type ProcessingStep struct {
	pipeline.BaseStep[Item, ProcessContext]
}

func NewProcessingStep() *ProcessingStep {
	return &ProcessingStep{
		BaseStep: pipeline.NewBaseStep[Item, ProcessContext]("Processor"),
	}
}

func (s *ProcessingStep) Process(ctx context.Context, input <-chan Item, output chan<- pipeline.StepResult[Item], pipelineContext ProcessContext) {
	for item := range input {
		result := s.ProcessItem(ctx, item, pipelineContext, func(ctx context.Context, item Item, pipelineContext ProcessContext) (bool, error) {
			// Transform the item.
			item.Value = fmt.Sprintf("%s - Processed", item.Value)
			return false, nil // Not skipped, no error.
		})
		output <- result
	}
	close(output)
}

func main() {
	// Set up logging.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	
	// Create pipeline.
	myPipeline := pipeline.NewPipeline[Item, ProcessContext]()
	
	// Add steps.
	myPipeline.SetFirstStep(NewSourceStep())
	myPipeline.AddStep(NewProcessingStep())
	
	// Create context.
	ctx := context.Background()
	pipelineContext := ProcessContext{
		Config: map[string]string{"key": "value"},
	}
	
	// Execute pipeline.
	if err := myPipeline.Execute(ctx, pipelineContext); err != nil {
		slog.Error("Pipeline execution failed", "error", err)
		os.Exit(1)
	}
	
	slog.Info("Pipeline completed successfully")
}
```

## Documentation

For detailed documentation and advanced usage examples, please see the [Go package documentation](https://pkg.go.dev/github.com/temirov/pipeline).

## Examples

Check out the [`examples`](examples) directory for more complete usage examples:

- ls: An implementation of a pipeline that mimics the Unix ls command. This example features smart wildcard reconstruction—if the shell expands a wildcard into a single file (e.g. content.sh), the system reconstructs the intended pattern (e.g. *.sh) so that matching is done correctly.
- More examples coming soon!

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 