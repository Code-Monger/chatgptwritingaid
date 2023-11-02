package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func runTextCleanup(model string, content string, inputFileName string) (int, string, error) {
	const initialQuestion = `can you fix only punctuation and report the number of errors corrected and list them (and the reason for them) 
be very careful about it and double check each condition, 
repeat the process until there are no more errors with the corrected text 
if there are no errors then there is no need to provide the corrected text. 
If there are any corrections provide the full corrected text at the end
never provide the original text
when giving the corrected ADD a tag to the begining "---BEGIN---" and at the end add "---END---"
when giving the corrections made ADD a tag to the begining "---BEGIN ERRORS---" and at the end add "---END ERRORS---"
if no corrections are needed only produce the tag "--COMPLETE--" 
for the following:`

	iteration := 0
	cleanText := ""

	question := initialQuestion + "\n\n" + string(content)
	//response, promptTokens, completionTokens, totalTokens ,err := GetChatResult(question)

	//fmt.Println(response)

	//os.Exit(0);
	total_promptTokens := 0
	total_completionTokens := 0
	total_totalTokens := 0
	correctedText := string(content)
	for iteration < 20 {
		iteration++
		response, promptTokens, completionTokens, totalTokens, err := GetChatResult(model, question)
		for err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 65)
		}
		total_promptTokens += promptTokens
		total_completionTokens += completionTokens
		total_totalTokens += totalTokens
		errorLog := extractBetweenTags(response, "---BEGIN ERRORS---", "---END ERRORS---")
		if len(errorLog) > 0 {
			// Open the file in append mode, or create it if it doesn't exist
			file, err := os.OpenFile(inputFileName+".errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return iteration, correctedText, err
			}

			// Append text to file
			_, err = file.WriteString(errorLog)
			if err != nil {
				return iteration, correctedText, err
			}
			file.Close()
			fmt.Println("errors:" + errorLog)
		} else if !strings.Contains(response, "--COMPLETE--") {
			fmt.Printf("missing completion tag in\n")
			fmt.Printf(response)
			fmt.Printf("\nforcing cmpletion\n")
			response = "--COMPLETE--"
		}

		if strings.Contains(response, "--COMPLETE--") {
			//fmt.Println(response)
			out := []byte(strings.TrimSpace(strings.ReplaceAll(correctedText, "\n\n", "\n")))
			ioutil.WriteFile(inputFileName+".corrected.txt", out, 0644)
			cleanText = correctedText
			break
		} else {
			correctedText = extractBetweenTags(response, "---BEGIN---", "---END---")
			if len(correctedText) == 0 {
				fmt.Printf("No text found?\n")
				fmt.Println(response)
				continue
			} else {
				cleanText = correctedText
			}
		}
		question = initialQuestion + "\n\n" + strings.TrimSpace(correctedText)
	}
	return iteration, cleanText, nil
}

func testTextParsing(content string) string {
	outtext := ""
	content = strings.TrimSpace(content)
	delm, paragraphs := splitIntoParagraphs(content)
	for _, paragraph := range paragraphs {
		// Convert each paragraph to sentences
		sentences := splitIntoSentences(paragraph)

		for _, sentence := range sentences {
			if sentence != "" {
				outtext += sentence
				//outtext += sentence + " "
			}
		}
		outtext += delm
	}
	outtext = strings.TrimSpace(outtext)
	if outtext != content {
		//delta := getStringDelta(content, outtext)
		//fmt.Printf("parse error:delta %s\n", delta)

		fmt.Printf("parse error:text\n%s\nnewtext\n%s \n", content, outtext)

		//os.Exit((1))
	}
	return outtext
}

func runTextCleanupBySentence(model string, content string) (string, error) {
	//const initialQuestion = `correct any grammatical errors in the following:
	const initialQuestion = `correct spelling pr punctuation errors in the following:
	`
	total := 0
	delm, paragraphs := splitIntoParagraphs(content)
	for _, paragraph := range paragraphs {
		// Convert each paragraph to sentences
		sentences := splitIntoSentences(paragraph)
		for _, sentence := range sentences {
			if sentence != "" {
				total++
			}
		}
	}
	outtext := ""
	progress := 0
	for _, paragraph := range paragraphs {
		// Convert each paragraph to sentences
		sentences := splitIntoSentences(paragraph)

		for _, sentence := range sentences {
			updated := sentence
			if sentence != "" {
				//sentence += "."
				i := 0
				progress++
				correct := false
				responseMap := make(map[string]int)
				for i < 8 && !correct {
					i++
					question := initialQuestion + sentence
					response, _, _, _, err := GetChatResult(model, question)
					errcnt := 0
					for err != nil && errcnt < 3 {
						errcnt++
						fmt.Println(err)
						time.Sleep(time.Second * 65)
						response, _, _, _, err = GetChatResult(model, question)
					}
					// skip first line
					pos := strings.Index(response, "\n")
					if pos > 0 {
						response = response[pos:]
					}
					response = strings.TrimSpace(response)
					//. (No errors found)
					if !correct {
						//quoteFix := "He couldn’t imagine how she’d pulled all this off."
						quoteFix := strings.ReplaceAll(sentence, "’", "'")
						quoteFix = strings.ReplaceAll(sentence, "“", "\"")
						quoteFix = strings.ReplaceAll(sentence, "”", "\"")
						//
						//here are no spelling or punctuation errors in the given sentence.
						correct = response == quoteFix
					}
					if !correct {
						correct = strings.Contains(response, "no grammatical errors")
					}
					if !correct {
						correct = strings.Contains(response, "grammatically correct")
					}
					if !correct {
						correct = strings.Contains(response, "No errors found")
					}
					if !correct {
						correct = strings.Contains(response, "There are no spelling or punctuation errors in the sentence provided")
					}
					if !correct {
						correct = strings.Contains(response, "here are no spelling or punctuation errors in the given sentence.")
					}
					if !correct {
						correct = strings.Contains(response, "No corrections needed")
					}
					if !correct {
						correct = strings.Contains(response, "The sentence appears to be correctly spelled and punctuated")
					}
					if !correct {
						correct = strings.Contains(response, "already correct")
					}
					if !correct {
						correct = strings.Contains(response, "already correct")
					}
					if !correct {
						//No spelling or punctuation error
						correct = strings.Contains(response, "No spelling or punctuation error")
					}

					if correct {
						response = sentence
						updated = sentence
						correct = true
						break
					} else {
						if existing, ok := responseMap[response]; ok {
							responseMap[response] = existing + 1
						} else {
							responseMap[response] = 1
						}

						//fmt.Printf("Old: %s\n", sentence)
						//fmt.Printf("New: %s\n", response)
					}
				}

				percentage := float64(progress) / float64(total) * 100
				fmt.Printf("\rProcessing... %.2f%% complete    ", percentage)

				if !correct {
					maxFound := 0
					//if len(responseMap) > 0 {
					//	fmt.Printf("%v\n", responseMap)
					//}

					for k, v := range responseMap {
						if v > maxFound {
							maxFound = v
							updated = k
						}
					}
				}
				outtext += updated
			}
		}
		outtext += delm
	}
	fmt.Printf("\r paragraph:%d sentences: %d delta: %d ", len(paragraphs), total, len(outtext)-len(content))
	fmt.Printf("\n")

	return strings.TrimSpace(outtext), nil
}
