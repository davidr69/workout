package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // go get github.com/lib/pq
	"workout.lavacro.net/models"
)

var db *sql.DB

func Init() {
	pw := os.Getenv("PGPASSWORD")
	conn := fmt.Sprintf("postgres://david:%s@dev-db:5432/workout?sslmode=disable", pw)

	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
}

func AllProgress() ([]models.AllProgress, error) {
	var allProg []models.AllProgress

	rows, sqlErr := db.Query(`
SELECT progid, exercise, muscle, mydate, weight, rep1, rep2
FROM app.allprogress
WHERE progid IS NOT NULL
ORDER BY muscle, exercise, mydate
LIMIT 5
	`)
	if sqlErr != nil {
		log.Fatal("Problem executing query ...", sqlErr)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal("Problem closing 'rows' resource")
		}
	}()

	for rows.Next() {
		var p models.AllProgress

		if err := rows.Scan(
			&p.ProgressId,
			&p.Exercise,
			&p.Muscle,
			&p.Mydate,
			&p.Weight,
			&p.Rep1,
			&p.Rep2); err != nil {
			log.Fatal("Problem scanning row ...", err)
		}

		allProg = append(allProg, p)
	}

	return allProg, nil
}
