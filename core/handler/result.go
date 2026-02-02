package handler

type Result struct {
	Label    string `json:"label"`
	Subtitle string `json:"subtitle"`
	Rank     int    `json:"rank"`
	// Icon  string `json:"icon"`
	Action Action `json:"action"`
}
