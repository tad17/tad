package fiber

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // для связи с mysql
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
)

type App struct {
	name  string
	fiber *fiber.App
	db    *sqlx.DB
}


type StringHandler func(db *sqlx.DB, query string) string
type HTMLHandler func(db *sqlx.DB, query string) map[string]interface{}

func OpenDB(database string) *sqlx.DB {
	// fmt.Printf("открываю БД nano-svelte\n")
	username := os.Getenv("USERNAMEDB")
	password := os.Getenv("PASSWORDDB")
	connstr := fmt.Sprintf("%s:%s@(xigmanas:3306)/%s", username, password, database)
	db, err := sqlx.Open("mysql", connstr)
	if err != nil {
		fmt.Errorf("невозможно открыть базу данных %s, %v\n", database, err)
		os.Exit(0)
	}
	return db
}

func NewApp(name string, db *sqlx.DB) *App {
	// Create a new engine
	//engine := html.New("./app/views", ".html")
	engine := html.NewFileSystem(http.Dir("./app/views"), ".html")

	engine.Reload(true) // Optional. Default: false
	// Debug will print each template that is parsed, good for debugging
	engine.Debug(true) // Optional. Default: false

	// инициализация fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// midlleware
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(favicon.New())
	app.Use(compress.New())

	// static
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

func (app *App) HTMLRouter(name string, funcRouter HTMLHandler) {
	app.fiber.Get(name, func(ctx *fiber.Ctx) error {
		q := ctx.Query("q", "")
		return ctx.Render(name, funcRouter(app.db, q))
	})
}
