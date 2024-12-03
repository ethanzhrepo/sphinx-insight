package processor

import (
	"context"
	"log"

	"github.com/openai/openai-go"
)

type ChatgptProcessor struct {
	client       *openai.Client
	conversation []string
}

func NewChatgptProcessor() *ChatgptProcessor {
	// openai_token := os.Getenv("OPENAI_API_KEY")
	return &ChatgptProcessor{
		client: openai.NewClient(
		// defaults to os.Getenv("OPENAI_API_KEY")
		// option.WithAPIKey(openai_token),
		),
		conversation: []string{
			"Based on the following text, find possible cryptocurrency symbols from the text and list them separated by commas. Analyze their impact on the market and divide them into three levels: negative, neutral, and positive. The answer should not exceed 200 characters.",
		},
	}
}

func (p *ChatgptProcessor) Process(data *ProcessorData) (string, error) {

	params := openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelGPT3_5Turbo),
		Messages: openai.F(
			[]openai.ChatCompletionMessageParamUnion{
				openai.AssistantMessage(p.conversation[0]),
				openai.UserMessage(data.Content),
			},
		),
	}

	completion, err := p.client.Chat.Completions.New(context.TODO(), params)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return data.Content + "\n<pre>" + completion.Choices[0].Message.Content + "</pre>", nil
}

func (p *ChatgptProcessor) Name() string {
	return "ChatgptProcessor"
}
