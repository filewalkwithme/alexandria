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
	//connect to Postgres
	orm, scream := alexandria.ConnectToPostgres()
	if scream != nil {
		panic(scream)
	}

	//create the orm handler for Book objects
	ormBooks, scream := orm.NewHandler(Book{})
	if scream != nil {
		panic(scream)
	}

	//Create Table
	ormBooks.DropTable()
	ormBooks.CreateTable()
	book := Book{Name: "The book is on the table", Pages: 198}
	err := ormBooks.Insert(&book)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("book: %v\n", book)

	/*
		//Insert/update
		ormBooks.Insert(Book{Name: "The book is on the table", Pages: 198})

		//Update
		ormBooks.Update().Where("pages > 0")
		ormBooks.Update().ByID(9)
		ormBooks.Update().All()


		//Select
		ormBooks.Select().Where("pages > 0")
		ormBooks.Select().ByID(9)
		ormBooks.Select().All()

		//Delete
		ormBooks.Delete().Where("id=9")
		ormBooks.Delete().ByID(9)
		ormBooks.Delete().All()

		//Drop Table
		ormBooks.DropTable()
	*/
}
