package types

import "image"

type Transformation struct {
	Type       string      `json:"type"`
	Parameters interface{} `json:"parameters"`
}

type TransformData struct {
	File            string           `json:"file"`
	Transformations []Transformation `json:"transformations"`
}

type JobMsg struct {
	Buf            []byte
	Transformation Transformation
}

type ResMsg struct {
	Err  error
	Img  image.Image
	Name string
}

type Config struct {
	WorkersCount int    `json:"workers"`
	Address      string `json:"address"`
}
