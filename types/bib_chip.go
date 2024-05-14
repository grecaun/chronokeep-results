package types

type BibChip struct {
	Identifier int64  `json:"-"`
	Bib        string `json:"bib"`
	Chip       string `json:"chip"`
}
