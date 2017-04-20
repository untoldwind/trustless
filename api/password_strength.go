package api

// PasswordEstimate options to estimate the strength of a password
type PasswordEstimate struct {
	Password string   `json:"password"`
	Inputs   []string `json:"inputs,omitempty"`
}

// PasswordStrength contains details about the strength of a password
type PasswordStrength struct {
	Entropy          float64 `json:"entropy"`
	CrackTime        float64 `json:"crackTime"`
	CrackTimeDisplay string  `json:"crackTimeDisplay"`
	Score            int     `json:"score"`
}
