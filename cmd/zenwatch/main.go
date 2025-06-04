package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	analyzeCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	outFilePath := analyzeCmd.String("out", "reports/latest.md", "Path to save the output Markdown report")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'analyze' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "analyze":
		analyzeCmd.Parse(os.Args[2:])
		if analyzeCmd.NArg() < 1 {
			fmt.Println("Usage: zenwatch analyze <repo-url> --out <output-file>")
			analyzeCmd.Usage()
			os.Exit(1)
		}
		repoURL := analyzeCmd.Arg(0)

		fmt.Printf("Repository URL: %s\n", repoURL)
		fmt.Printf("Output File: %s\n", *outFilePath)
	default:
		fmt.Println("Expected 'analyze' subcommand")
		os.Exit(1)
	}

	// Further implementation will follow in subsequent steps
}
