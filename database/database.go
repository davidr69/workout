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
		var weight sql.NullFloat64
		var rep1 sql.NullInt16
		var rep2 sql.NullInt16

		if err := rows.Scan(
			&p.ProgressId,
			&p.Exercise,
			&p.Muscle,
			&p.Mydate,
			&weight,
			&rep1,
			&rep2); err != nil {
			log.Fatal("Problem scanning row ...", err)
		}

		if weight.Valid {
			p.Weight = float32(weight.Float64)
		}
		if rep1.Valid {
			p.Rep1 = int(rep1.Int16)
		}
		if rep2.Valid {
			p.Rep2 = int(rep2.Int16)
		}

		allProg = append(allProg, p)
	}

	return allProg, nil
}
