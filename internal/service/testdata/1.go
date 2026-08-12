package service_test

import "fmt"

func a() {
	// some comment
	fmt.Println("This is test print for code with no comments // this is not a comment")
	fmt.Println("//This //is //not //a //comment //too")
	fmt.Println(`
	/* And
	this
	is
	not
	a
	multiline
	comment */
	`)
	// another
	// one
	// comment
}

/*
multi
line
comment
*/
