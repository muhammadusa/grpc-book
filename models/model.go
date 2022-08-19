package models

type Book struct {
	Id         int32
	Name       string
	Texts      string
	AuthorId   string
	CategoryId string
}

type PostBook struct {
	Name       string `json:"book_name"`
	Texts      string `json:"texts"`
	AuthorId   string `json:"author_id"`
	CategoryId string `json:"category_id"`
}

type Id struct {
	Id int32 `json:"id"`
}
