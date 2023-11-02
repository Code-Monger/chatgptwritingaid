package main

import (
	"context"

	openai "github.com/gtkit/go-openai"
)

var apiKey string = ""

func GetChatResult(model string, prompt string) (answer string, promotTokens int, completionTokens int, totalTokens int, err error) {
	/*
		if (model == openai.GPT3Dot5Turbo) {
			promptLen := len(prompt);
			if promptLen < (4000 * 4) {
				fmt.Println("using instruct.")
				model = "gpt-3.5-turbo-instruct"
			}
		}
	*/

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(

		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Temperature: 0,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	//
	if err != nil {
		//fmt.Printf("ChatCompletion error: %v\n", err)
		return "", 0, 0, 0, err
	}

	answer = resp.Choices[0].Message.Content

	//fmt.Println(answer)
	return answer, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens, nil

}
