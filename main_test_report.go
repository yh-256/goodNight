package main

import (
	"fmt"
	"os"
	"time"

	"github.com/user/zenwatch/internal/git"
	"github.com/user/zenwatch/internal/metrics"
	"github.com/user/zenwatch/internal/report"
)

func main() {
	// Dummy data for testing
	commitInfo := &git.CommitInfo{
		Hash:    "a1b2c3d4e5f6",
		Author:  "Jules Verne",
		Email:   "jules@example.com",
		Date:    time.Now().Format(time.RFC1123),
		Message: "feat: implement amazing new features",
	}

	overallStats := &metrics.OverallStats{
		TotalLinesAdded:   150,
		TotalLinesDeleted: 30,
		FileStats: map[string]*metrics.FileTypeStat{
			".go": {Extension: ".go", Count: 5},
			".md": {Extension: ".md", Count: 2},
		},
		ComplexityStats: []metrics.ComplexityStat{
			{Complexity: 20, Package: "main", FunctionName: "complexFunc", File: "main.go", Line: 42},
			{Complexity: 16, Package: "helper", FunctionName: "anotherComplex", File: "utils/helper.go", Line: 101},
		},
		AverageComplexity: 18.0,
		FunctionsOverThreshold: 2,
	}

	complexityThreshold := 15
	repoURL := "https://github.com/user/testrepo"
	reportDate := time.Now().Format("2006-01-02 15:04:05 MST")

	// Test badge URL generation
	totalChanges := overallStats.TotalLinesAdded + overallStats.TotalLinesDeleted
	badgeURL := report.GenerateBadgeURL(totalChanges, overallStats.AverageComplexity)
	fmt.Println("Generated Badge URL:", badgeURL)

	reportData := report.ReportData{
		RepoURL:             repoURL,
		ReportDate:          reportDate,
		BadgeURL:            badgeURL, // Include the badge
		Commit:              commitInfo,
		Stats:               overallStats,
		ComplexityThreshold: complexityThreshold,
	}

	outputFilePath := "test_report.md"
	err := report.GenerateMarkdownReport(reportData, outputFilePath)
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Report generation test complete. Check %s\n", outputFilePath)

	// You can optionally cat the file here to see its content in the subtask output
	// but it might be long.
	// content, _ := os.ReadFile(outputFilePath)
	// fmt.Println("\n--- Report Content ---")
	// fmt.Println(string(content))
	// fmt.Println("--- End Report Content ---")

	// Test without badge
	reportDataNoBadge := report.ReportData{
		RepoURL:             repoURL,
		ReportDate:          reportDate,
		BadgeURL:            "", // No badge
		Commit:              commitInfo,
		Stats:               overallStats,
		ComplexityThreshold: complexityThreshold,
	}
	outputFilePathNoBadge := "test_report_no_badge.md"
	err = report.GenerateMarkdownReport(reportDataNoBadge, outputFilePathNoBadge)
	if err != nil {
		fmt.Printf("Error generating report without badge: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Report generation (no badge) test complete. Check %s\n", outputFilePathNoBadge)


}
