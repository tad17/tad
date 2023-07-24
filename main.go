package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/tad17/tad/fiber"
	"os"
	"strings"
)

func testRouter(db *sqlx.DB, query string) string {
	//fmt.Printf("query: %s\n", query)
	var urls []string
	cmd := "select newurl from filemeta limit 5"
	if err := db.Select(&urls, cmd); err != nil {
		fmt.Errorf("ошибка выполнения: %s", cmd)
		os.Exit(0)
	}
	s := strings.Join(urls, " ")
	return s
}

func testHTML(db *sqlx.DB, query string) map[string]interface{} {
	result := make(map[string]interface{})
	result["Item"] = "Проверка"
	result["value"] = 456
	result["items"] = []string{"Первый", "Второй", "Третий"}
	return result
}

func main() {
	db := fiber.OpenDB("sea")
	app := fiber.NewApp("test-app", db)
	app.JSONRouter("/test", testRouter)
	app.HTMLRouter("test2", testHTML)
	app.Start(4444)
	fmt.Printf("init app fiber: %v\n", app)
}
