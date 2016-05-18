package main

import (
	alexandriaOrm "github.com/maiconio/alexandria/orm"
)

//Book is on the table
type Book struct {
	ID    int
	Name  string
	Pages int
}

func main() {
	orm, err := alexandriaOrm.ConnectToPostgres()

	if err != nil {
		panic(err)
	}

	//find the first book
	orm.Find(Book{}).All()
}
