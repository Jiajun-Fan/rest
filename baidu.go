package main

type TransResult struct {
	From         string
	To           string
	Trans_result []TransElem
}

type TransElem struct {
	Src string
	Dst string
}
