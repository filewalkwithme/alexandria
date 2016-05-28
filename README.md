# alexandria


## Usage

```
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
    ormBooks := orm.Handle(Book{})

    //Create Table
    ormBooks.CreateTable()

    //Insert/update
	ormBooks.Save(Book{Name: "The book is on the table", Pages: 198})

    //Select
    ormBooks.Find().Where("pages > 0")
	ormBooks.Find().ByID(9)
	ormBooks.Find().All()

    //Delete
    ormBooks.Delete().Where("id=9")
	ormBooks.Delete().ByID(9)
	ormBooks.Delete().All()

    //Drop Table
    ormBooks.DestroyTable()
}

```
