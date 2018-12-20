package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-xray-sdk-go/xray"
)

const imageCount = 2

func main() {
	http.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("VisionService.Analyse"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		svc := rekognition.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))
		xray.AWS(svc.Client)

		w.Header().Set("Content-type", "text/html")
		for _, v := range getImages(r.Context(), imageCount) {
			labels, err := svc.DetectLabelsWithContext(r.Context(), &rekognition.DetectLabelsInput{Image: &rekognition.Image{Bytes: v}})
			if err != nil {
				xray.GetSegment(r.Context()).Close(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "<img src=\"data:image/jpeg;base64,%s\"> <br>", base64.StdEncoding.EncodeToString(v))
			for _, l := range labels.Labels {
				fmt.Fprintln(w, *l.Name, *l.Confidence, "<br>")
			}
		}
	})))

	http.ListenAndServe(":8000", nil)
}

func getImages(ctx context.Context, c int) [][]byte {
	images := [][]byte{}

	var wg sync.WaitGroup
	for i := 0; i < c; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			client := xray.Client(&http.Client{})
			req, _ := http.NewRequest(http.MethodGet, "https://picsum.photos/500/500/?random", nil)
			res, _ := client.Do(req.WithContext(ctx))
			bod, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()

			images = append(images, bod)
		}(&wg)
	}

	wg.Wait()

	return images
}
