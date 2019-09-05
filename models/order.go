package models

type Order struct {
	Model
	MemberId int
	Amount   float64
}
