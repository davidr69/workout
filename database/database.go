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

	query := `
		SELECT progid, allprogress.exerciseid, exercise, muscle, mydate, weight, rep1, rep2
		FROM app.allprogress
		WHERE progid IS NOT NULL
		ORDER BY muscle, exercise, mydate
	`

	rows, sqlErr := db.Query(query)
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
			&p.ExerciseId,
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

func Exercises() ([]models.Exercises, error) {
	var exercises []models.Exercises

	query := `
		SELECT m.description AS muscle, e.id, e.description AS exercise_name
		FROM app.exercise e
		JOIN app.muscle m ON e.muscle = m.id
		ORDER BY m.description, exercise_name
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Problem executing query ...", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal("Problem closing 'rows' resource")
		}
	}()

	muscle := ""
	var exercise models.Exercises

	for rows.Next() {
		var e models.Exercise
		if err := rows.Scan(
			&e.Muscle,
			&e.Id,
			&e.ExerciseName); err != nil {
			log.Fatal("Problem scanning row ...", err)
		}

		if muscle == "" || muscle != *e.Muscle {
			if muscle != "" {
				exercises = append(exercises, exercise)
			}
			muscle = *e.Muscle
			exercise = models.Exercises{}
			exercise.Muscle = muscle
			exercise.Exercises = make([]models.Exercise, 0)
		}
		exercise.Exercises = append(exercise.Exercises, e)
	}
	exercises = append(exercises, exercise)

	return exercises, nil
}
