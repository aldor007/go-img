package operations

import (
	"fmt"
	"net/http"
	"strings"
	"bytes"
	"image"
	"github.com/disintegration/imaging"
)

type makeTrasform func (img image.Image, transformation Transformation)  resMsg
var Methods = map[string]makeTrasform{"crop": cropImage, "strip": removeExif, "rotate": rotateImage}

func isImage(buf []byte) bool {
	fmt.Println(http.DetectContentType(buf));
	return strings.HasPrefix(http.DetectContentType(buf), "imaage/")
}

func transformWorker(jobs <-chan jobMsg, results chan<-resMsg) { // define a worker that accepts its id, a read-only channel for jobs and a write-only channel for results
	for j := range jobs { // while the channel is open or has any messages
		imgReader := bytes.NewReader(j.buf)
		img, _ , err := image.Decode(imgReader)
		if err != nil {
			results <- resMsg{err, nil, ""}
		} else {

			if method, ok := methods[j.transformation.Type]; ok {
				results <- method(img, j.transformation)
			}

		}
	}
}

func creaeFileName(job jobMsg)  string {
	return job.transformation.Type
}



func cropImage(img image.Image, trans Transformation) resMsg {
	params, success := trans.Parameters.(map[string]interface{})
	if success == false {
		fmt.Println(params, success, trans)
		return resMsg{errors.New("invalid a params"), nil, ""}
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
				return resMsg{errors.New("unable to convert" + n), nil, ""}
			}
			paramsMap[n] = int(d)
			name = name + fmt.Sprintf("%s-%d", n, int(d));
		} else {
			return resMsg{errors.New("missing required params"), nil, ""}
		}
	}

	crop := image.Rect(paramsMap["x"], paramsMap["y"], paramsMap["x"] + paramsMap["width"], paramsMap["y"] + paramsMap["height"])
	crop = img.Bounds().Intersect(crop)
	result := image.NewRGBA(crop)
	draw.Draw(result, crop, img, crop.Min, draw.Src)
	return resMsg{nil, result, name}

}

func removeExif(img image.Image, trans Transformation) resMsg {
	return resMsg{nil, img, "exif"}
}

func rotateImage(img image.Image, trans Transformation) resMsg  {
	params, success := trans.Parameters.(map[string]interface{})
	if success == false {
		return resMsg{errors.New("invalid a params"), nil, ""}
	}


	d, ok := params["angle"].(float64)
	if !ok {
		fmt.Println(params, success, trans)
		return resMsg{errors.New("missing or invalid paramter"), nil, ""}
	}§

	result := imaging.Rotate(img, d * 90, color.Black)
	return resMsg{nil, result, fmt.Sprintf("rotate-%b", d)}
}


