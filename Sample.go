package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// the strings in the back tick are used to define how this would get converted into
// json notation
// basically we are saying that book struct would have id in json instead of ID
// we are capitalizing the field names for some reason so that, that field becomes
// ready to be exported outside.
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int32  `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsaby", Author: "Frank Fitzgerald", Quantity: 20},
	{ID: "3", Title: "War & Peace", Author: "Leo Tolstoy", Quantity: 3},
}

// gin.Context is essentially all information about request.
// return books in the form of a json object.
func getBooks(c *gin.Context) {
	// return json version of all books.
	c.IndentedJSON(http.StatusOK, books)
}
func addBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, books)
}
func bookById(c *gin.Context) {
	id := c.Param("id")
	boo, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, boo)
}
func getBookByID(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book Not Found :-(")
}
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}
	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not Found"})
		return
	}

	if book.Quantity == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not available"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
	return
}
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}
	boo, err := getBookByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not found"})
		return
	}
	boo.Quantity += 1
	c.IndentedJSON(http.StatusOK, gin.H{"message": "book returned successfully"})
	return

}
func main() {

	// create a router and handle endpoints of our route.
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", addBook)
	router.GET("/books/:id", bookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
