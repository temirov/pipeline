package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/temirov/pipeline"
	"github.com/temirov/pipeline/examples/ls/steps"
)

func main() {
	lsContext := &steps.LsContext{}

	executionContext, cancelFunction := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunction()

	parseArgsStep := steps.NewParseArgumentsStep()
	listFilesStep := steps.NewListFilesStep()
	filterFilesStep := steps.NewFilterFilesStep()
	printFilesStep := steps.NewPrintFilesStep()

	builderStep := pipeline.NewPipelineBuilder(parseArgsStep)
	builderList := pipeline.AppendStep(builderStep, listFilesStep)
	builderFilter := pipeline.AppendStep(builderList, filterFilesStep)
	builderPrint := pipeline.AppendStep(builderFilter, printFilesStep)
	pipelineInstance := builderPrint.Build()

	if err := pipelineInstance.Execute(executionContext, lsContext); err != nil {
		fmt.Fprintf(os.Stderr, "Pipeline execution error: %v\n", err)
		os.Exit(1)
	}
}
