package report

import (
	"fmt"
	"html/template" // Using html/template for Markdown to be safe, though text/template is often fine for MD
	"os"
	"path/filepath"
	"strings"

	"github.com/user/zenwatch/internal/git"
	"github.com/user/zenwatch/internal/metrics"
)

const markdownTemplate = `
# ZenWatch Analysis Report

**Repository:** {{.RepoURL}}
**Analyzed At:** {{.ReportDate}}

{{if .BadgeURL}}
![ZenWatch Stats]({{.BadgeURL}})
{{end}}

## Latest Commit Analyzed
- **Hash:** {{.Commit.Hash}}
- **Author:** {{.Commit.Author}} <{{.Commit.Email}}>
- **Date:** {{.Commit.Date}}
- **Message:** {{.Commit.Message}}

## Code Statistics
- **Total Lines Added:** {{.Stats.TotalLinesAdded}}
- **Total Lines Deleted:** {{.Stats.TotalLinesDeleted}}
  *Note: Line counts are overall for the commit. Per-file line counts were not available with current git analysis settings.*

### File Type Distribution
| Extension | Count |
|-----------|-------|
{{range $ext, $stat := .Stats.FileStats -}}
| {{$ext}} | {{$stat.Count}} |
{{end}}

## Cyclomatic Complexity Analysis (Threshold > {{.ComplexityThreshold}})
- **Average Complexity (of functions over threshold):** {{printf "%.2f" .Stats.AverageComplexity}}
- **Functions Over Threshold:** {{.Stats.FunctionsOverThreshold}}

{{if gt .Stats.FunctionsOverThreshold 0 -}}
### Functions Over Complexity Threshold
| Complexity | Function                               | File:Line        | Package        |
|------------|----------------------------------------|------------------|----------------|
{{range .Stats.ComplexityStats -}}
| {{.Complexity}} | {{.FunctionName}}                     | {{.File}}:{{.Line}} | {{.Package}}    |
{{end}}
{{else -}}
No functions found with cyclomatic complexity greater than {{.ComplexityThreshold}}.
{{end}}
`

// ReportData holds all necessary data for rendering the Markdown report.
type ReportData struct {
	RepoURL             string
	ReportDate          string
	BadgeURL            string // Optional: URL for the status badge
	Commit              *git.CommitInfo
	Stats               *metrics.OverallStats
	ComplexityThreshold int
}

// GenerateMarkdownReport creates a Markdown report from the analysis data.
func GenerateMarkdownReport(data ReportData, outputPath string) error {
	tmpl, err := template.New("markdownReport").Parse(markdownTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse markdown template: %w", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create report file %s: %w", outputPath, err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	fmt.Printf("Markdown report generated at %s\n", outputPath)
	return nil
}

// GenerateBadgeURL creates a URL for a shields.io badge.
// Example: Total Changes: 150, Avg Complexity: 8.5
func GenerateBadgeURL(totalChangedLines int, avgComplexity float64) string {
	label := "ZenWatch"
	// Ensure avgComplexity is formatted nicely for the URL, e.g., "8.5" not "8.500000"
	message := fmt.Sprintf("changes %d | avg complx %.1f", totalChangedLines, avgComplexity)
	color := "blue"

	// URL encode message
	safeMessage := strings.ReplaceAll(message, " ", "%20")
	safeMessage = strings.ReplaceAll(safeMessage, "|", "%7C")

	return fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s", label, safeMessage, color)
}
