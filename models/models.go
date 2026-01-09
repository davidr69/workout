package models

import "time"

type Envelope map[string]any

type Progress struct {
	ProgressId *int     `json:"progress_id"`
	ExerciseId *int     `json:"exercise_id"`
	Exercise   *string  `json:"exercise"`
	Muscle     *string  `json:"muscle"`
	MuscleId   *int     `json:"muscle_id"`
	Mydate     *string  `json:"mydate"`
	Weight     *float32 `json:"weight"`
	Rep1       *int     `json:"rep1"`
	Rep2       *int     `json:"rep2"`
}

type Exercises struct {
	Muscle    string     `json:"muscle"`
	Exercises []Exercise `json:"exercises"`
}

type Exercise struct {
	Id           *int    `json:"id"`
	Muscle       *string `json:"muscle"`
	ExerciseName *string `json:"exercise_name"`
}

type NewActivity struct {
	ExerciseID *int       `json:"exercise"`
	MyDate     *time.Time `json:"mydate"`
	Weight     *float32   `json:"weight"`
	Rep1       *int       `json:"rep1"`
	Rep2       *int       `json:"rep2"`
}
