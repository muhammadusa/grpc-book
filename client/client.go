package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"app/models"
	"app/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var c proto.BookServiceClient

func AllBook(ctx *gin.Context) {

	req := proto.GetAllBookRequest{
		TableName: "book",
	}

	res, err := c.AllBook(context.Background(), &req)

	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, res)
}

func getBook(ctx *gin.Context) {

	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Print(id)
	}

	req := proto.GetBookRequest{
		Id: int32(id),
	}

	res, err := c.GetBook(context.Background(), &req)

	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, res)
}

func postBook(ctx *gin.Context) {

	pBook := models.PostBook{}

	ctx.BindJSON(&pBook)

	req := proto.PostBookRequest{
		Book: &proto.Book{
			Name:       pBook.Name,
			Texts:      pBook.Texts,
			AuthorId:   pBook.AuthorId,
			CategoryId: pBook.CategoryId,
		},
	}

	res, err := c.PostBook(context.Background(), &req)

	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(201, res)
}

func deleteBook(ctx *gin.Context) {

	id := models.Id{}

	ctx.BindJSON(&id)

	req := proto.DeleteBookRequest{
		Id: id.Id,
	}

	res, err := c.DeleteBook(context.Background(), &req)

	if err != nil {
		log.Print(err)
	}

	ctx.JSON(202, res)
}

func updateBook(ctx *gin.Context) {

	book := models.Book{}

	ctx.BindJSON(&book)

	req := proto.UpdateBookRequest{
		Id: book.Id,
		Book: &proto.Book{
			Name:       book.Name,
			Texts:      book.Texts,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
		},
	}

	res, err := c.UpdateBook(context.Background(), &req)

	if err != nil {
		log.Print(err)
	}

	ctx.JSON(202, res)
}

func main() {

	fmt.Println("Connection Client RPC...")

	//gin.SetMode(gin.releaseMode)
	//httpRouter.Use(gin.Recovery())
	httpRouter := gin.Default()

	conn, err := grpc.Dial(":50050", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c = proto.NewBookServiceClient(conn)

	httpRouter.GET("/getBook", AllBook)
	httpRouter.GET("/getBook/:id", getBook)
	httpRouter.POST("/postBook", postBook)
	httpRouter.DELETE("/deleteBook", deleteBook)
	httpRouter.PUT("/putBook", updateBook)

	httpRouter.Run()
}
