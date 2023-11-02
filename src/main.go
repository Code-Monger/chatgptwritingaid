package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/gtkit/go-openai"
)

func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() && strings.HasSuffix(path, ".txt") {
		fmt.Println("Processing file:", path)
		content, _ := os.ReadFile(path)
		sourceStirng := string(content)
		/*
			if len(sourceStirng) > 0 {
				if !strings.Contains(sourceStirng, "<$p>") {
					result := testTextParsing(sourceStirng)
					if result != sourceStirng {
						delta := getStringDelta(sourceStirng, result)
						fmt.Printf("delta\n")
						fmt.Printf("%s\n", delta)
					}
				}
			}
		*/
		if len(sourceStirng) > 0 {
			if !strings.Contains(sourceStirng, "<$p>") {
				result, _ := runTextCleanupBySentence(openai.GPT3Dot5Turbo, sourceStirng)
				fmt.Printf("delta:\n")
				delta := getStringDelta(sourceStirng, result)
				fmt.Printf("Corrected Text for %s :\n", path)
				fmt.Printf("%s\n", delta)
				fmt.Printf("Corrected Text for %s :\n", path)
				fmt.Printf("%s\n", result)
				//os.WriteFile("delta.log", []byte(delta), fs.FileMode(0644))
				//os.WriteFile("corrected.log", []byte(result), fs.FileMode(0644))
				os.WriteFile(path, []byte(result), fs.FileMode(0644))

			}
		}
	}
	return nil
}

func main() {
	apiKeyF := flag.String("apikey", "", "API key for processing (required)")
	root := flag.String("path", ".", "Starting path for directory traversal")

	// Parse command-line flags
	flag.Parse()

	// Check if API key is provided
	if *apiKeyF == "" {
		fmt.Println("Error: API key not provided")
		flag.PrintDefaults()
		os.Exit(1)
	}
	apiKey = *apiKeyF

	// Walk the directory specified in the path flag (or the default if not provided)
	err := filepath.Walk(*root, processFile)
	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}
	//			Model: openai.GPT432K,

	//Model: openai.GPT4,

	/*
		iteration, ctp3Result, err := runTextCleanup(openai.GPT3Dot5Turbo ,string(content), inputFileName)
		if (err != nil) {
			fmt.Printf("Error running chatgtp3 model %v", err)
		}
		fmt.Printf("Original Text:\n")
		fmt.Printf("%s\n",string(content))

		fmt.Printf("b3.5 iterations: %d\n",iteration);
		fmt.Printf("v3.5 Corrected Text:\n")
		fmt.Printf("%s\n",ctp3Result)
	*/
	/*
	   //	if (len(ctp3Result) > 0) {
	   	iterationTotal := 0
	   	iteration := 2
	   	result := string(content)
	   	for (iteration > 1 && iterationTotal < 3) {
	   		iterationTotal++
	   		iteration, result, err = runTextCleanup(openai.GPT4 ,result, inputFileName)
	   		if (err != nil) {
	   			fmt.Printf("Error running chatgtp4 model %v", err)
	   		}
	   		fmt.Printf("v4 iterations: %d\n",iteration);
	   		fmt.Printf("v4 Corrected Text:\n")
	   		fmt.Printf("%s\n", result)
	   	}
	*/
}
