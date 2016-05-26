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
	orm, _ := alexandria.ConnectToPostgres()
	ormBooks := orm.Handle(Book{})

	//ormBooks.CreateTable()
	err := ormBooks.DestroyTable()
	fmt.Printf("err: %v\n", err)

	//ormBooks.Save(Book{Name: "Fight Club", Pages: 198})
	//ormBooks.Find().Where("pages > 0")
	//ormBooks.Find().ByID(9)
	//ormBooks.Find().All()
	//ormBooks.Delete().Where("id=9")
	//ormBooks.Delete().ByID(10)
	//ormBooks.Delete().All()
	//ormBooks.DestroyTable()
}
