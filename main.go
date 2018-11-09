package main // define package

import (
	"fmt"
"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"image"
	"image/draw"
	"strings"
	"archive/zip"
	"image/jpeg"
	"github.com/aldor007/transformer-go/fetch"
	"github.com/kataras/iris/core/errors"
	"image/png"
	"github.com/disintegration/imaging"
	"image/color"
)

const addr = "localhost:8081"

func main() { // define main function
	http.HandleFunc("/accept", handle) // add a handler to the default ServeMux
	err := http.ListenAndServe(addr, nil) // start listening on the addres and instruct to use the default ServeMux
	fmt.Println(err.Error()) // ListenAndServe blocks execution unless an error occurs, so we log that here
}


func handle (w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var data TransformData
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Json parse error", http.StatusBadRequest)
			return
		}

		fmt.Println(string(body), data)

		if data.File == "" {
			http.Error(w, "Invalid image address", http.StatusBadRequest)
			return
		}

		buf, err := fetch.FetchFile(data.File)
		ct := http.DetectContentType(buf)
		fmt.Println(ct)
		if !strings.HasPrefix(ct, "image/") {
			http.Error(w, "Invalid file, image required", http.StatusBadRequest)
			return
		}
		ext := strings.Replace(ct, "image/", "", 1)

		total := 0
		jobs := make(chan jobMsg, 100)
		results := make(chan resMsg, 100)

		for i := 1; i <= 3; i++ { // spawn 3 workers
			go transformWorker(jobs, results) // every worker runs in a separate goroutine
		}

		fmt.Println("AAAAAAA", data.Transformations)
		for _, trans := range data.Transformations {
			if trans.Parameters != nil {
				jobs <- jobMsg{buf, trans}
				total++
			}
		}
		close(jobs)
		w.WriteHeader(http.StatusOK)  // write the header to the outgoing socket with 200 status code
		zipWriter := zip.NewWriter(w)
		if total > 0 {
			for j := 1; j <= total; j++ {
				res := <-results
				if res.err == nil {
					fmt.Println("result", res)
					f, err := zipWriter.Create(res.name + "." + ext)
					if err != nil {
						continue
					}
					switch ext {
					case "png":
						png.Encode(f, res.img)
					case "jpg", "jpeg":
						jpeg.Encode(f, res.img, nil)
					default:
						fmt.Println(ext)

					}
				} else {
					fmt.Println(res.err)
				}
			}

		}
		zipWriter.Close()


	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type makeTrasform func (img image.Image, transformation Transformation)  resMsg

var methods = map[string]makeTrasform{"crop": cropImage, "strip": removeExif, "rotate": rotateImage}

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


