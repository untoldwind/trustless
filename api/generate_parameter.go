package api

type CharsParameter struct {
	NumChars         int  `json:"num_chars"`
	IncludeUpper     bool `json:"include_upper"`
	IncludeNumbers   bool `json:"include_numbers"`
	IncludeSymbols   bool `json:"include_symbols"`
	RequireUpper     bool `json:"require_upper"`
	RequireNumber    bool `json:"require_number"`
	RequireSymbol    bool `json:"require_symbol"`
	ExcludeSimilar   bool `json:"exclude_similar"`
	ExcludeAmbiguous bool `json:"exclude_ambiguous"`
}

type WordsParameter struct {
	NumWords int    `json:"num_words"`
	Delim    string `json:"delim"`
}

type GenerateParameter struct {
	Chars *CharsParameter `json:"chars,omitempty"`
	Words *WordsParameter `json:"words,omitempty"`
}
