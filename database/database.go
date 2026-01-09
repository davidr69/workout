package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // go get github.com/lib/pq
	"workout.lavacro.net/models"
)

type Dao struct {
	conn *sql.DB
}

func (db *Dao) Init() {
	pw := os.Getenv("PGPASSWORD")
	uri := fmt.Sprintf("postgres://david:%s@dev-db:5432/workout?sslmode=disable", pw)

	var err error
	db.conn, err = sql.Open("postgres", uri)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
}

func (db *Dao) AllProgress() ([]models.Progress, error) {
	var allProg []models.Progress

	query := `
		SELECT progid, allprogress.exerciseid, exercise, muscle, mydate, weight, rep1, rep2
		FROM app.allprogress
		WHERE progid IS NOT NULL
		ORDER BY muscle, exercise, mydate
	`

	rows, sqlErr := db.conn.Query(query)
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
		var p models.Progress
		p = getProgress(rows)

		allProg = append(allProg, p)
	}

	return allProg, nil
}

func (db *Dao) Exercises() ([]models.Exercises, error) {
	var exercises []models.Exercises

	query := `
		SELECT m.description AS muscle, e.id, e.description AS exercise_name
		FROM app.exercise e
		JOIN app.muscle m ON e.muscle = m.id
		ORDER BY m.description, exercise_name
	`
	rows, err := db.conn.Query(query)
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

func (db *Dao) YearMonths() ([]string, error) {
	var months []string

	query := `
		WITH x AS (
			SELECT DISTINCT mydate
			FROM app.progress
		)
		SELECT CAST(date_part('year', mydate) AS varchar) || LPAD(CAST(date_part('month', mydate) AS varchar), 2, '0') AS yrmon
		FROM x
		ORDER BY mydate
	`
	rows, err := db.conn.Query(query)
	if err != nil {
		log.Fatal("Problem executing query ...", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal("Problem closing 'rows' resource")
		}
	}()

	for rows.Next() {
		var yrmon string
		if err := rows.Scan(&yrmon); err != nil {
			log.Fatal("Problem scanning row ...", err)
		}
		months = append(months, yrmon)
	}

	return months, nil
}

func (db *Dao) Progress(year int, month int) ([]models.Progress, error) {
	var resp []models.Progress

	query := `
		SELECT e.id, m.description AS muscle, m.id AS muscle_id, e.description AS exercise, p.weight, p.rep1, p.rep2,
				p.id AS progress_id
		FROM app.exercise e
		JOIN app.muscle m ON e.muscle = m.id
		LEFT JOIN app.progress p ON e.id = p.exercise
		    AND DATE_PART('year', mydate) = $1
		    AND DATE_PART('month', mydate) = $2
		ORDER BY muscle, exercise
	`

	rows, sqlErr := db.conn.Query(query, year, month)
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
		var p models.Progress
		var weight sql.NullFloat64
		var rep1 sql.NullInt16
		var rep2 sql.NullInt16

		if err := rows.Scan(
			&p.ExerciseId,
			&p.Muscle,
			&p.MuscleId,
			&p.Exercise,
			&weight,
			&rep1,
			&rep2,
			&p.ProgressId,
		); err != nil {
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

		resp = append(resp, p)
	}

	return resp, nil
}

func (db *Dao) Activity(id int) (models.Progress, error) {
	query := `
		SELECT progid, allprogress.exerciseid, exercise, muscle, mydate, weight, rep1, rep2
		FROM app.allprogress
		WHERE progid = $1
	`

	rows, sqlErr := db.conn.Query(query, id)
	if sqlErr != nil {
		log.Fatal("Problem executing query ...", sqlErr)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal("Problem closing 'rows' resource")
		}
	}()

	if !rows.Next() {
		log.Fatal("No rows returned")
	}

	return getProgress(rows), nil
}

func getProgress(rows *sql.Rows) models.Progress {
	var prog models.Progress

	var weight sql.NullFloat64
	var rep1 sql.NullInt16
	var rep2 sql.NullInt16

	if err := rows.Scan(
		&prog.ProgressId,
		&prog.ExerciseId,
		&prog.Exercise,
		&prog.Muscle,
		&prog.Mydate,
		&weight,
		&rep1,
		&rep2); err != nil {
		log.Fatal("Problem scanning row ...", err)
	}

	if weight.Valid {
		w := float32(weight.Float64)
		prog.Weight = &w
	}
	if rep1.Valid {
		r := int(rep1.Int16)
		prog.Rep1 = &r
	}
	if rep2.Valid {
		r := int(rep2.Int16)
		prog.Rep2 = &r
	}

	return prog
}

func (db *Dao) NewActivity(act models.NewActivity) (int64, error) {
	query := `
		INSERT INTO app.progress (exercise, mydate, weight, rep1, rep2)
		VALUES ($1, $2, $3, $4, $5)
	`

	res, err := db.conn.Exec(query, act.ExerciseID, act.MyDate, act.Weight, act.Rep1, act.Rep2)
	if err != nil {
		log.Println("Error inserting activity: ", err)
		return 0, err
	}

	return res.LastInsertId()
}
