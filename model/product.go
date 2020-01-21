package model

import alira "github.com/ivohutasoit/alira/model"

type Product struct {
	alira.BaseModel
	Barcode     string
	Name        string
	LongName    string
	Manufacture string
}
