package main

import (
	"flag"
	"log"
	"os"

	"answer.io/cmd/handler"
	"answer.io/pkg/bolt"
	"answer.io/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var path string

func main() {

	flag.Parse()

	flag.StringVar(&path, "path", "/tmpanswer.db", "path of the database to store the data")

	e := echo.New()
	utils.Generator = func() string {
		return uuid.NewString()
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	db, err := utils.Open(path)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	manager, err := bolt.NewService(db)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	handler.NewQuestionHandler(e, manager)

	e.Logger.Fatal(e.Start(":1323"))
}
