package handlers

import (
	"fmt"
	"log"
	"os"

	"diikstra.fr/letterboxd-statistics/app-client/models"
	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func Init() {
	fmt.Println("[INIT] Startup sequence starting...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[INIT] WARNING : Failed to load .env file")
		log.Println(err)
	}

	fmt.Printf("[INIT] Env : %s\n", os.Getenv("ENV"))

	models.Rdb.Init()

	fmt.Println("[INIT] Startup sequence done.")
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
