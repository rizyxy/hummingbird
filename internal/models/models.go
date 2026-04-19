package models

type Match struct {
	Name         string `json:"name"`
	Category     string `json:"category"` // "table" or "function"
	FileName     string `json:"file_name"`
	FunctionName string `json:"function_name"`
	LineNumber   int    `json:"line_number"`
	Snippet      string `json:"snippet"`
}

type Summary struct {
	Name          string `json:"name"`
	TotalRefs     int    `json:"total_refs"`
	UniqueFuncs   int    `json:"unique_functions_count"`
	FrictionScore int    `json:"friction_score"`
}

type ImpactReport struct {
	TargetTable    string   `json:"target_table"`
	DirectImpact   []string `json:"direct_impact_functions"`   // Functions touching the table
	IndirectImpact []string `json:"indirect_impact_functions"` // Functions calling those functions
	TotalRiskScore int      `json:"total_risk_score"`
}
