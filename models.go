package main

type PersonInfo struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	PhoneNumber string `json:"phone_number" binding:"required,numeric,len=10"`
	City        string `json:"city" binding:"required,alpha"`
	State       string `json:"state" binding:"required"`
	Street1     string `json:"stree1" binding:"required"`
	Street2     string `json:"street2" binding:"required"`
	ZipCode     string `json:"zip_code" binding:"required,len=5"`
}
