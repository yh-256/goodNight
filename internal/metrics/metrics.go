package metrics
type OverallStats struct {
	TotalLinesAdded, TotalLinesDeleted, FunctionsOverThreshold int
	AverageComplexity float64
	FileStats map[string]*FileTypeStat
	ComplexityStats []ComplexityStat
}
type FileTypeStat struct { Extension string; Count int }
type ComplexityStat struct { Complexity int; Package, FunctionName, File string; Line int }
