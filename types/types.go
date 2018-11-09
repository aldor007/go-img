package types

import "image"

type Transformation struct {
	Type string `json:"type"`
	Parameters interface{} `json:"parameters"`

}

type TransformData struct {
	File string `json:"file"`
	Transformations []Transformation `json:"transformations"`
}

type JobMsg struct {
	Buf []byte
	Rransformation Transformation
}

type ResMsg struct {
	Err error
	Img image.Image
	Name string

}
