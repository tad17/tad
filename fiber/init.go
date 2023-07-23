package fiber

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // для связи с mysql
	"github.com/jmoiron/sqlx"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type App struct {
	name  string
	fiber *fiber.App
	db    *sqlx.DB
}

type StringHandler func(db *sqlx.DB, query string) string

func OpenDB(database string) *sqlx.DB {
	// fmt.Printf("открываю БД nano-svelte\n")
	db, err := sqlx.Open("mysql", "itman:X753951x@(xigmanas:3306)/"+database)
	if err != nil {
		fmt.Errorf("невозможно открыть базу данных %s", database)
		os.Exit(0)
	}
	return db
}

func NewApp(name string, db *sqlx.DB) *App {
	// инициализация fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(favicon.New())
	app.Use(compress.New())

	app.Static("/", "app/dist/index.html")
	app.Static("/assets", "app/dist/assets", fiber.Static{
		Compress: true,
	})
	app.Static("/fonts", "app/dist/fonts", fiber.Static{
		Compress: true,
	})

	//fmt.Printf("fiber init\n")
	result := App{
		name:  name,
		fiber: app,
		db:    db,
	}
	return &result
}

func (app *App) Start(port int) error {
	err := app.fiber.Listen(fmt.Sprintf(":%d", port))
	return err
}

func (app *App) JSONRouter(name string, funcRouter StringHandler) {
	app.fiber.Get(name, func(ctx *fiber.Ctx) error {
		q := ctx.Query("q", "")
		return ctx.JSON(fiber.Map{
			"status": "ok",
			"result": funcRouter(app.db, q),
		})
	})
}
func (app *App) StringRouter(name string, funcRouter StringHandler) {
	app.fiber.Get(name, func(ctx *fiber.Ctx) error {
		q := ctx.Query("q", "")
		return ctx.SendString(funcRouter(app.db, q))
	})
}
