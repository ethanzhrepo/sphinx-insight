package pipeline

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ethanzhrepo/sphinx-insight/core/processor"
)

// each pipline has a set of processors
type Pipeline struct {
	Processors []processor.Processor
}

// NewPipeline creates a new pipeline
func NewPipeline(processors ...processor.Processor) *Pipeline {
	return &Pipeline{Processors: processors}
}

// Process processes the input string through the pipeline
func (p *Pipeline) Process(input string) (string, error) {
	var err error
	// input is json string, parse to ProcessorData
	var processorData processor.ProcessorData
	err = json.Unmarshal([]byte(input), &processorData)
	if err != nil {
		log.Println("Error parsing input to ProcessorData")
		return "", err
	}

	for _, processor := range p.Processors {
		if isDebug() {
			log.Println("Processing by processor: ", processor.Name())
		}
		input, err = processor.Process(&processorData)
		if err != nil {
			log.Println("Error processing by processor: ", processor.Name())
			return "", err
		}
	}
	return input, nil
}

// AddProcessor adds a processor to the pipeline
func (p *Pipeline) AddProcessor(processor processor.Processor) {
	p.Processors = append(p.Processors, processor)
}

// RemoveProcessor removes a processor from the pipeline
func (p *Pipeline) RemoveProcessor(processor processor.Processor) {
	for i, proc := range p.Processors {
		if proc == processor {
			p.Processors = append(p.Processors[:i], p.Processors[i+1:]...)
			return
		}
	}
}

// start
func (p *Pipeline) Start(ch <-chan string) {
	log.Println("Pipeline started")
	for input := range ch {
		if isDebug() {
			log.Println("Processing by pipeline: ", input)
		}
		go p.Process(input)
	}
}

func isDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
