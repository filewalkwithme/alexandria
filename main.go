package main

import (
	"fmt"
	alexandria "github.com/maiconio/alexandria/orm"
	"time"
)

//Book is on the table
type Book struct {
	ID              int
	Name            string
	Pages           int
	HardCover       bool
	Price           float64
	PublicationDate time.Time
	Chapters        []Chapter
	Authors         []Author
}

//Chapter one
type Chapter struct {
	ID   int
	Name string
}

//Author ...
type Author struct {
	ID        int
	FirstName string
	LastName  string
}

func main() {

	//connect to Postgres
	orm, scream := alexandria.ConnectToPostgres("user=docker password=docker dbname=docker sslmode=disable")
	if scream != nil {
		panic(scream)
	}

	//create the orm handler for Book objects
	chap := Chapter{1, "chapter-name"}
	ormBooks, scream := orm.NewHandler(Book{Chapters: []Chapter{chap}})
	if scream != nil {
		panic(scream)
	}

	//Create Table
	ormBooks.DropTable()
	ormBooks.CreateTable()

	//Insert
	book := Book{Name: "The book is on the table", Pages: 198, HardCover: true, Price: 99.99, PublicationDate: time.Date(2016, time.June, 1, 0, 0, 0, 0, time.UTC)}
	ormBooks.Save(&book)

	//Update
	book.Name = "The book is on fire"
	ormBooks.Save(&book)

	//Select
	selBook, _ := ormBooks.Select().ByID(1)
	selBooks, _ := ormBooks.Select().Where("pages > 0")
	selBooks, _ = ormBooks.Select().All()

	fmt.Printf("%v\n", selBook)
	fmt.Printf("%v\n", selBooks)

	//Delete
	ormBooks.Delete().ByID(1)
	ormBooks.Delete().Where("id=9")
	ormBooks.Delete().All()

	//Drop Table
	ormBooks.DropTable()

	//orm.FreeQuery("Select foo from bar")

}
