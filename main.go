package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dixonwille/wmenu"
	_ "github.com/mattn/go-sqlite3"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	checkErr(err)

	choice := "y"
	for choice == "y" {
		fmt.Printf("\n")
		menu := wmenu.NewMenu("Choose what to do: ")
		menu.Action(func(opts []wmenu.Opt) error { handleFunc(&choice, db, opts); return nil })
		menu.Option("Add new task", 0, false, nil)
		menu.Option("view tasks", 1, true, nil)
		menu.Option("update task's status", 2, false, nil)
		menu.Option("delete a task", 3, false, nil)
		menu.Option("exit", 4, false, nil)
		menu.LoopOnInvalid()

		menuerrr := menu.Run()
		checkErr(menuerrr)
	}

	db.Close()
}
