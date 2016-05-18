package main

import (
	"fmt"
	alexandriaOrm "github.com/maiconio/alexandria/orm"
)

//Book is on the table
type Book struct {
	ID    int
	Name  string
	Pages int
}

func main() {
	//connect on postgres
	orm, err := alexandriaOrm.ConnectToPostgres()

	if err != nil {
		panic(err)
	}

	//finds all books
	books := orm.Find(Book{}).All()

	fmt.Printf("%v\n", books)
}
