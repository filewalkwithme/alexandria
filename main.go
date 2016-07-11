package main

import (
	"fmt"
	"time"

	alexandria "github.com/maiconio/alexandria/orm"
)

//Book is on the table
type Book struct {
	ID              int
	Name            string
	Pages           int
	HardCover       bool
	Price           float64
	PublicationDate time.Time
	Chapters        []*Chapter
	Author          *Author
}

//Chapter one
type Chapter struct {
	ID     int
	Name   string
	BookID int
}

//Author ...
type Author struct {
	ID        int
	FirstName string
	LastName  string
	BookID    int
}

func main() {

	//connect to Postgres
	orm, scream := alexandria.ConnectToPostgres("user=docker password=docker dbname=docker sslmode=disable")
	if scream != nil {
		panic(scream)
	}

	ormAuthor, scream := orm.NewHandler(Author{})
	ormAuthor.DropTable()
	ormAuthor.CreateTable()

	ormChapter, scream := orm.NewHandler(Chapter{})
	ormChapter.DropTable()
	ormChapter.CreateTable()

	//create the orm handler for Book objects
	ormBooks, scream := orm.NewHandler(Book{})
	if scream != nil {
		panic(scream)
	}

	//Create Table
	ormBooks.DropTable()
	ormBooks.CreateTable()

	//Insert
	book := Book{Name: "The book is on the table", Pages: 198, HardCover: true, Price: 99.99, PublicationDate: time.Date(2016, time.June, 1, 0, 0, 0, 0, time.UTC), Chapters: []*Chapter{&Chapter{0, "chapter-one", 0}, &Chapter{0, "chapter-two", 0}}, Author: &Author{0, "Maicon", "Costa", 0}}
	err := ormBooks.Save(&book)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	//Update
	book.Name = "The book is on fire"
	ormBooks.Save(&book)

	book.Author.FirstName = "maiconio"
	ormBooks.Save(&book)

	//Insert
	//book = Book{Name: "The book is on the table", Pages: 198, HardCover: true, Price: 99.99, PublicationDate: time.Date(2016, time.June, 1, 0, 0, 0, 0, time.UTC), Chapters: []*Chapter{&Chapter{0, "chapter-one", 0}, &Chapter{0, "chapter-two", 0}}}
	//err = ormBooks.Save(&book)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	return
	//}

	//Insert
	//book = Book{Name: "The book is on the table 2"}
	//err = ormBooks.Save(&book)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	return
	//}

	//Select
	//selBook, _ := ormBooks.Select().ByID(1)
	//selBooks, _ := ormBooks.Select().Where("pages > 0")
	//selBooks, _ = ormBooks.Select().All()

	//fmt.Printf("selBook: %v\n", selBook)
	//fmt.Printf("selBooks: %v\n", selBooks)

	//Delete
	//ormBooks.Delete().ByID(1)
	//ormBooks.Delete().Where("id=9")
	//ormBooks.Delete().All()

	//Drop Table
	//ormBooks.DropTable()

	//orm.FreeQuery("Select foo from bar")

}
