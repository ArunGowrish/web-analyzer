package model

type AnalysisResult struct {
	HTMLVersion string
	Title       string
	Headings    map[string]int
	Links       []Link
	Link        Link
}
