package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"diikstra.fr/letterboxd-statistics/app-cron/src/letterboxd"
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

func (pdb *PDB) SetMovieInfos(movie letterboxd.MovieMeta) error {
	stmt := `
		INSERT INTO movies (id, slug, link, title, rating, popularity, poster, backdrop)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := pdb.Client.Exec(stmt,
		movie.Id,
		movie.Slug,
		movie.Link,
		movie.Title,
		movie.Rating,
		0,
		movie.Poster,
		movie.Backdrop,
	)
	return err
}

func (pdb *PDB) GetMovieIds() (letterboxd.MovieIds, error) {
	stmt := `SELECT id FROM movies`
	rows, err := pdb.Client.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movieIds []string
	for rows.Next() {
		var movieId string
		err := rows.Scan(&movieId)
		if err != nil {
			return nil, err
		}
		movieIds = append(movieIds, movieId)
	}

	return movieIds, nil
}
