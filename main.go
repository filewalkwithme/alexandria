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
	//the new fabulous alexandria command set
	//
	//orm := alexandria.ConnectToPostgres("127.0.0.1:3452", "postgres")
	//ormBooks := orm.Handle(book)
	//ormBooks.CreateTable()
	//ormBooks.DestroyTable()
	//ormBooks.Save(book)
	//ormBooks.Find().Where("id=abcd")
	//ormBooks.Find().ID(123)
	//ormBooks.Find().All()
	//ormBooks.Delete().Where("id=abcd")
	//ormBooks.Delete().ID(123)
	//ormBooks.Delete().All()

	//connect on postgres
	orm, err := alexandria.ConnectToPostgres()
	if err != nil {
		panic(err)
	}

	ormBooks := orm.Handle(Book{})

	ormBooks.Save(Book{Name: "Fight Club", Pages: 198})

	books := ormBooks.Find().Where("pages > 0")

	fmt.Printf("book: %v\n", books)

}
