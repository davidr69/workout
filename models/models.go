package models

import "database/sql"

type AllProgress struct {
	ProgressId int             `json:"progress_id"`
	ExerciseId int             `json:"exercise_id"`
	Exercise   string          `json:"exercise"`
	Muscle     string          `json:"muscle"`
	Mydate     string          `json:"mydate"`
	Weight     sql.NullFloat64 `json:"weight"`
	Rep1       sql.NullInt16   `json:"rep1"`
	Rep2       sql.NullInt16   `json:"rep2"`
}
