package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBooksByID)
	router.POST("/books", postBooks)

	router.Run("localhost:8080")
}

func getBooks(c *gin.Context) {
	books := apiendpoints.getBooks(c)

	c.IndentedJSON(http.StatusOK, books)
}

func getBooksByID(c *gin.Context) {
	db, err := sql.Open("mysql", "admin:password@tcp(parker-database.cfhfkqv5cjrl.us-east-1.rds.amazonaws.com:3306)/book_schema")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()

	query := "SELECT * FROM books WHERE id = ?"

	queryResult, err := db.Query(query, c.Param("id"))
	if err != nil {
		fmt.Println(err)
	}

	var books []Book

	for queryResult.Next() {
		var bk Book
		if err := queryResult.Scan(&bk.ID, &bk.Title); err != nil {
			fmt.Println(books, err)
		}
		books = append(books, bk)

	}

	id := c.Param("id")

	for _, a := range books {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

func postBooks(c *gin.Context) {
	db, err := sql.Open("mysql", "admin:password@tcp(parker-database.cfhfkqv5cjrl.us-east-1.rds.amazonaws.com:3306)/book_schema")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()

	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	insert := "INSERT INTO books (title) VALUES (?)"
	insertResult, err := db.ExecContext(context.Background(), insert, &newBook.Title)
	if err != nil {
		fmt.Println(err)
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	c.IndentedJSON(http.StatusCreated, id)
}
