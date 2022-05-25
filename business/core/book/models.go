package book

type Book struct {
	ID              int     `json:"id"`
	Isbn            string  `json:"isbn"`
	Title           string  `json:"title"`
	Author          *string `json:"author"`
	PublicationYear *string `json:"publication_year"`
	Publisher       *string `json:"publisher"`
}
