package models

type Envelope map[string]any

type AllProgress struct {
	ProgressId *int     `json:"progress_id"`
	ExerciseId *int     `json:"exercise_id"`
	Exercise   *string  `json:"exercise"`
	Muscle     *string  `json:"muscle"`
	Mydate     *string  `json:"mydate"`
	Weight     *float32 `json:"weight"`
	Rep1       *int     `json:"rep1"`
	Rep2       *int     `json:"rep2"`
}
