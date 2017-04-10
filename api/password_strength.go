package api

type PasswordEstimate struct {
	Password string   `json:"password"`
	Inputs   []string `json:"inputs,omitempty"`
}

type PasswordStrength struct {
	Entropy          float64 `json:"entropy"`
	CrackTime        float64 `json:"crackTime"`
	CrackTimeDisplay string  `json:"crackTimeDisplay"`
	Score            int     `json:"score"`
}
