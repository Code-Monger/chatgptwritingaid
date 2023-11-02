package main

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/neurosnap/sentences.v1"
)

func extractBetweenTags(text, startTag, endTag string) string {
	if strings.Contains(text, startTag) && strings.Contains(text, endTag) {
		start := strings.Index(text, startTag) + len(startTag)
		end := strings.Index(text, endTag)
		return text[start:end]
	}
	return ""
}

func getStringDelta(str1 string, str2 string) string {
	dmp := diffmatchpatch.New()

	words1 := strings.Split(str1, " ")
	words2 := strings.Split(str2, " ")

	diffs := dmp.DiffMain(strings.Join(words1, "\x00"), strings.Join(words2, "\x00"), false)

	// Process the diffs to replace null characters with spaces for human-readable output
	var result []string
	for _, diff := range diffs {
		text := strings.Replace(diff.Text, "\x00", " ", -1)
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			result = append(result, "[Insert: "+text+"]")
		case diffmatchpatch.DiffDelete:
			result = append(result, "[Delete: "+text+"]")
		case diffmatchpatch.DiffEqual:
			result = append(result, text)
		}
	}
	return strings.Join(result, "\n")
}

func splitIntoParagraphs(text string) (string, []string) {
	delim := "\r\n\r\n"
	if !strings.Contains(text, delim) {
		delim = "\n\n"
	}
	return delim, strings.Split(text, delim)
}

/*
	func splitIntoSentences(paragraph string) []string {
		dialogueRegex := regexp.MustCompile(`"([^"]+)"`)
		dialogues := dialogueRegex.FindAllString(paragraph, -1)

		for _, dialogue := range dialogues {
			// Replace dialogues temporarily with placeholders
			paragraph = strings.Replace(paragraph, dialogue, "<DIALOGUE>", 1)
		}

		sentences := regexp.MustCompile(`[.!?]\s`).Split(paragraph, -1)

		// Restore dialogues from placeholders
		for i, sentence := range sentences {
			if strings.Contains(sentence, "<DIALOGUE>") {
				sentence = strings.Replace(sentence, "<DIALOGUE>", dialogues[0], 1)
				dialogues = dialogues[1:]
			}

			sentencteIndex := strings.Index(paragraph, sentences[i])
			if sentencteIndex > -1 {
				sentencteIndex += len(sentences[i])
				if sentencteIndex < len(paragraph) {
					//fmt.Printf("addind to end %s\n", paragraph[sentencteIndex:sentencteIndex+1])
					sentence += paragraph[sentencteIndex : sentencteIndex+1]
				}
			}

			sentences[i] = strings.TrimSpace(sentence)
		}

		return sentences
	}
*/
func splitIntoSentences(paragraph string) []string {
	tokenizer := sentences.NewSentenceTokenizer(sentences.NewStorage())

	// Break the paragraph into sentences.
	sents := tokenizer.Tokenize(paragraph)
	var results []string

	// Display the individual sentences.

	for _, s := range sents {
		results = append(results, s.Text)
	}
	return results
}
