package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PDB struct {
	Client *sql.DB
}

var Pdb PDB = PDB{}

func (pdb *PDB) Init() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName))
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

func (pdb *PDB) GetMoviesInfos(movieIds []string) ([]MovieMeta, error) {
	stmt := `
		SELECT id, slug, link, title, rating, popularity, poster, backdrop FROM movies
		WHERE id = ANY($1)
	`
	rows, err := pdb.Client.Query(stmt, pq.Array(movieIds))
	if err != nil {
		return nil, err
	}

	moviesMeta := make([]MovieMeta, len(movieIds))
	i := 0

	for rows.Next() {
		movieMeta := MovieMeta{}
		err := rows.Scan(&movieMeta.Id, &movieMeta.Slug, &movieMeta.Link, &movieMeta.Title, &movieMeta.Rating, &movieMeta.Popularity, &movieMeta.Poster, &movieMeta.Backdrop)
		if err != nil {
			return nil, err
		}

		moviesMeta[i] = movieMeta
		i++
	}

	return moviesMeta, nil
}
