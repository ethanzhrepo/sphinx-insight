package processor

type Processor interface {
	Process(data *ProcessorData) (string, error)
	Name() string
}

type ProcessorData struct {
	Content string `json:"content"`
	Link    string `json:"link"`
}
