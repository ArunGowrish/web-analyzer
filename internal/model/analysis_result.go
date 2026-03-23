package model

type AnalysisResult struct {
	HTMLVersion string
	Title       string
	Headings    map[string]int
	Link        Link
	LoginForm   bool
}
