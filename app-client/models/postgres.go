package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type PDB struct {
	Client *sql.DB
}

var Pdb PDB = PDB{}

func (pdb *PDB) Init() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName))
	if err != nil {
		log.Default().Println("[INIT] fail to connect to postgres")
		log.Fatal(err)
	}

	pdb.Client = db
}

type MovieMeta struct {
	Id         string
	Slug       string
	Link       string
	Title      string
	Rating     float32
	Popularity int
	Poster     string
	Backdrop   string
}

func (pdb *PDB) GetMovieInfos(movieId string) (*MovieMeta, error) {
	stmt := `
		SELECT id, slug, link, title, rating, popularity, poster, backdrop FROM movies
		WHERE id = $1
	`
	row := pdb.Client.QueryRow(stmt, movieId)

	movieMeta := &MovieMeta{}
	err := row.Scan(&movieMeta.Id, &movieMeta.Slug, &movieMeta.Link, &movieMeta.Title, &movieMeta.Rating, &movieMeta.Popularity, &movieMeta.Poster, &movieMeta.Backdrop)
	if err != nil {
		return nil, err
	}

	return movieMeta, nil
}
