package operations

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aldor007/transformer-go/types"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/draw"
)

type makeTrasform func(img image.Image, transformation types.Transformation) types.ResMsg

var methods = map[string]makeTrasform{"crop": cropImage, "strip": removeExif, "rotate": rotateImage}

func TransformWorker(jobs <-chan types.JobMsg, results chan<- types.ResMsg) {
	for j := range jobs {
		imgReader := bytes.NewReader(j.Buf)
		img, _, err := image.Decode(imgReader)
		if err != nil {
			results <- types.ResMsg{err, nil, ""}
		} else {

			if method, ok := methods[j.Transformation.Type]; ok {
				results <- method(img, j.Transformation)
			}

		}
	}
}

func cropImage(img image.Image, trans types.Transformation) types.ResMsg {
	params, success := trans.Parameters.(map[string]interface{})
	if success == false {
		fmt.Println(params, success, trans)
		return types.ResMsg{errors.New("invalid params"), nil, ""}
	}

	paramsMap := make(map[string]int)
	requiredParams := []string{"x", "y", "width", "height"}
	name := "crop"
	for _, n := range requiredParams {
		if v, ok := params[n]; ok {
			success = true
			d, success := v.(float64)
			if success == false {
				fmt.Println(v, success, params, paramsMap)
				return types.ResMsg{errors.New("unable to convert" + n), nil, ""}
			}
			paramsMap[n] = int(d)
			name = name + fmt.Sprintf("%s-%d", n, int(d))
		} else {
			return types.ResMsg{errors.New("missing required params"), nil, ""}
		}
	}

	crop := image.Rect(paramsMap["x"], paramsMap["y"], paramsMap["x"]+paramsMap["width"], paramsMap["y"]+paramsMap["height"])
	crop = img.Bounds().Intersect(crop)
	result := image.NewRGBA(crop)
	draw.Draw(result, crop, img, crop.Min, draw.Src)
	return types.ResMsg{nil, result, name}

}

func removeExif(img image.Image, trans types.Transformation) types.ResMsg {
	return types.ResMsg{nil, img, "exif"}
}

func rotateImage(img image.Image, trans types.Transformation) types.ResMsg {
	params, success := trans.Parameters.(map[string]interface{})
	if success == false {
		return types.ResMsg{errors.New("invalid a params"), nil, ""}
	}

	d, ok := params["angle"].(float64)
	if !ok {
		fmt.Println(params, success, trans)
		return types.ResMsg{errors.New("missing or invalid paramter"), nil, ""}
	}

	result := imaging.Rotate(img, d*90, color.Black)
	return types.ResMsg{nil, result, fmt.Sprintf("rotate-%d", int(d))}
}
