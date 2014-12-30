package main

type TransResult struct {
	Error int
	From  string
	To    string
	Data  TransData
}

type TransData struct {
	Symbols   []TransSymbol
	Word_name string
}

type TransSymbol struct {
	Parts []TransPart
	Ph_am string
	Ph_en string
}

type TransPart struct {
	Means []string
	Part  string
}
