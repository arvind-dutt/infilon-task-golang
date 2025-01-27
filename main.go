package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	defer db.Close()

	r := gin.Default()

	r.GET("person/:person_id/info", func(ctx *gin.Context) {
		GetPersonInfo(ctx, db)
	})

	r.POST("person/create", func(ctx *gin.Context) {
		CreatePerson(ctx, db)
	})

	r.Run(":8080")
}
