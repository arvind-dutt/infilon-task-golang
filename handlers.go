package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func GetPersonInfo will fetch the person information from the database and send it in response
func GetPersonInfo(ctx *gin.Context, db *sql.DB) {
	var personInfo PersonInfo

	// fetching the path parameter person_id from the request url
	personId := ctx.Param("person_id")

	query := `SELECT p.name, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code
				FROM person p
				JOIN phone ph ON p.id = ph.person_id
				JOIN address_join aj ON p.id = aj.person_id
				JOIN address a ON aj.address_id = a.id
				WHERE p.id = ?`

	err := db.QueryRow(query, personId).Scan(&personInfo.Name, &personInfo.PhoneNumber, &personInfo.City,
		&personInfo.State, &personInfo.Street1, &personInfo.Street2, &personInfo.ZipCode)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Person info not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, personInfo)
}

// func CreatePerson will create a person in Database
func CreatePerson(ctx *gin.Context, db *sql.DB) {
	var newPerson PersonInfo
	if err := ctx.ShouldBindJSON(&newPerson); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Begin of a transaction
	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// inserting into the person table
	resp, err := tx.Exec("INSERT INTO person (name, age) VALUES (?,0)", newPerson.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	personID, err := resp.LastInsertId()
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// inserting into the phone table
	_, err = tx.Exec("INSERT INTO phone (person_id, number) VALUES (?,?)", personID, newPerson.PhoneNumber)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// inserting into the address table
	_, err = tx.Exec("INSERT INTO address (city, state, street1, street2, zip_code) VALUES (?, ?, ?, ?, ?)",
		newPerson.City, newPerson.State, newPerson.Street1, newPerson.Street2, newPerson.ZipCode)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	addressID, err := resp.LastInsertId()
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// inserting into the address_join table
	_, err = tx.Exec("INSERT INTO address_join (person_id, address_id) VALUES (?,?)", personID, addressID)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Person Info Created successfully"})
}

// custome function for making error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
