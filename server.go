package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"app/database"
	"app/models"
	"app/proto"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedBookServiceServer
}

func (*server) AllBook(c context.Context, req *proto.GetAllBookRequest) (*proto.GetAllBookResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	book := []models.Book{}

	rows := db.Table(req.TableName).Find(&book)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	fmt.Println(book)

	var result []*proto.Book

	for _, data := range book {

		res := &proto.Book{
			Name:       data.Name,
			Texts:      data.Texts,
			AuthorId:   data.AuthorId,
			CategoryId: data.CategoryId,
		}

		result = append(result, res)
	}

	return &proto.GetAllBookResponse{
		Book: result,
	}, nil
}

func (*server) GetBook(c context.Context, req *proto.GetBookRequest) (*proto.GetBookResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Print(err)
	}

	book := models.Book{}

	rows := db.Table("book").
		Where("id = ?", req.GetId()).
		Find(&book)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	res := &proto.GetBookResponse{
		Id: book.Id,
		Book: &proto.Book{
			Name:       book.Name,
			Texts:      book.Texts,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
		},
	}

	return res, nil
}

func (*server) PostBook(c context.Context, req *proto.PostBookRequest) (*proto.PostBookResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Print(err)
	}

	var id int32

	rows := db.Table("book").
		Create(&req.Book).
		Select("id").
		Where("book_name = ?", req.Book.Name).
		Find(&id)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	res := &proto.PostBookResponse{
		Id: id,
		Book: &proto.Book{
			Name:       req.Book.Name,
			Texts:      req.Book.Texts,
			AuthorId:   req.Book.AuthorId,
			CategoryId: req.Book.CategoryId,
		},
	}

	return res, nil
}

func (*server) DeleteBook(c context.Context, req *proto.DeleteBookRequest) (*proto.DeleteBookResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	book := models.PostBook{}
	fmt.Println(req.Id)

	rows := db.Table("book").
		Select("book_name, texts, author_id, category_id").
		Where("id = ?", req.Id).
		Find(&book)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	del := db.Table("books").
		Delete(&models.Book{}, req.Id)

	if database.IsNotFound(del) {
		log.Fatal(del.Error)
	}

	res := &proto.DeleteBookResponse{
		Book: &proto.Book{
			Name:       book.Name,
			Texts:      book.Texts,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
		}}

	return res, nil
}

func (*server) UpdateBook(c context.Context, req *proto.UpdateBookRequest) (*proto.UpdateBookResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	bookTemp := models.PostBook{}

	findBook := db.Table("book").
		Select("first_name, texts, author_id, category_id").
		Where("id = ?", req.Id).
		Find(&bookTemp)

	if database.IsNotFound(findBook) {
		log.Fatal(findBook.Error)
	}

	if req.Book.Name == "" {
		req.Book.Name = bookTemp.Name
	}

	if req.Book.Texts == "" {
		req.Book.Texts = bookTemp.Texts
	}

	if req.Book.AuthorId == "" {
		req.Book.AuthorId = bookTemp.AuthorId
	}

	if req.Book.CategoryId == "" {
		req.Book.CategoryId = bookTemp.CategoryId
	}

	rows := db.Table("book").
		Where("id = ?", req.Id).
		Updates(
			map[string]interface{}{
				"book_name":   req.Book.Name,
				"texts":       req.Book.Texts,
				"author_id":   req.Book.AuthorId,
				"category_id": req.Book.CategoryId,
			},
		)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	return &proto.UpdateBookResponse{}, nil
}

func main() {

	fmt.Println("Connection Server RPC...")

	lis, err := net.Listen("tcp", ":50050")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	proto.RegisterBookServiceServer(s, &server{})

	err = s.Serve(lis)

	if err != nil {
		log.Fatal(err)
	}
}
