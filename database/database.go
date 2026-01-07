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
			w := float32(weight.Float64)
			p.Weight = &w
		}
		if rep1.Valid {
			r := int(rep1.Int16)
			p.Rep1 = &r
		}
		if rep2.Valid {
			r := int(rep2.Int16)
			p.Rep2 = &r
		}

		allProg = append(allProg, p)
	}

	return allProg, nil
}
