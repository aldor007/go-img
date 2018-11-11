package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/aldor007/transformer-go/fetch"
	"github.com/aldor007/transformer-go/operations"
	"github.com/aldor007/transformer-go/types"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const addr = "localhost:8081"

var configInstance types.Config

func main() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &configInstance)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/accept", handle)
	err = http.ListenAndServe(configInstance.Address, nil)
	fmt.Println(err.Error())
}

func handle(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			log.Print("Unable to parse request body")
			return
		}

		var data types.TransformData
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Json parse error", http.StatusBadRequest)
			log.Printf("Unable to parse json body %s", err)
			return
		}

		if data.File == "" {
			http.Error(w, "Invalid image address", http.StatusBadRequest)
			log.Print("Empty image address")
			return
		}

		buf, err := fetch.FetchFile(data.File)
		ct := http.DetectContentType(buf)

		if !strings.HasPrefix(ct, "image/") {
			http.Error(w, "Invalid file, image required", http.StatusBadRequest)
			log.Print("Wrong content type of file")
			return
		}
		ext := strings.Replace(ct, "image/", "", 1)

		total := 0
		jobs := make(chan types.JobMsg, 100)
		results := make(chan types.ResMsg, 100)

		for i := 1; i <= configInstance.WorkersCount; i++ {
			go operations.TransformWorker(jobs, results)
		}

		for _, trans := range data.Transformations {
			if trans.Parameters != nil {
				jobs <- types.JobMsg{buf, trans}
				total++
			}
		}
		close(jobs)
		w.WriteHeader(http.StatusOK)
		zipWriter := zip.NewWriter(w)
		if total > 0 {
			for j := 1; j <= total; j++ {
				res := <-results
				if res.Err == nil {
					f, err := zipWriter.Create(res.Name + "." + ext)
					if err != nil {
						continue
					}
					switch ext {
					case "png":
						png.Encode(f, res.Img)
					case "jpg", "jpeg":
						jpeg.Encode(f, res.Img, nil)
					default:
						log.Printf("Unknow ext %s", ext)

					}
				} else {
					log.Printf("Processing error %s", err)
				}
			}

		}
		zipWriter.Close()

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
