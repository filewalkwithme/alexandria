package main

import (
	"fmt"
	alexandria "github.com/maiconio/alexandria/orm"
)

//Book is on the table
type Book struct {
	ID    int
	Name  string
	Pages int
}

func main() {
	//connect on postgres
	orm, err := alexandria.ConnectToPostgres()

	if err != nil {
		panic(err)
	}

	//Create a new book
	newBook, err := orm.Save(Book{Name: "Fight Club", Pages: 198})
	fmt.Printf("%v\n", newBook)
	fmt.Printf("%v\n", err)

	//update an existing one
	orm.Save(Book{ID: 1, Name: "Fight Club UPDATED!", Pages: 198})

	//finds all books
	books := orm.Find(Book{}).All()
	fmt.Printf("%v\n", books)

	//finds book with ID=1
	book := orm.Find(Book{ID: 1}).One()
	fmt.Printf("%v\n", book)

	//finds book with pages >= 3
	largeBooks := orm.Find(Book{}).Where("pages >= 3")
	fmt.Printf("%v\n", largeBooks)
}
