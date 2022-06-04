package book

import "fmt"

type Book struct {
	ID              int     `json:"id"`
	Isbn            string  `json:"isbn"`
	Title           string  `json:"title"`
	Author          *string `json:"author"`
	PublicationYear *string `json:"publication_year"`
	Publisher       *string `json:"publisher"`
}

type NewBook struct {
	Isbn            string `json:"isbn"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear string `json:"publication_year"`
	Publisher       string `json:"publisher"`
}

type FieldError struct {
	field string
	err   string
}

func (fe FieldError) Error() string {
	return fmt.Sprintf("%s %s", fe.field, fe.err)
}
