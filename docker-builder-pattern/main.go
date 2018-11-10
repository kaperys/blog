package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/gin-gonic/gin"
)

func main() {
	svc := rekognition.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))

	r := gin.Default()

	r.StaticFS("/", http.Dir("html"))
	r.POST("/analyse", func(c *gin.Context) {
		var req struct{ Image string }

		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		img, err := base64.StdEncoding.DecodeString(req.Image)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		labels, err := svc.DetectLabels(&rekognition.DetectLabelsInput{Image: &rekognition.Image{Bytes: img}})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if len(labels.Labels) == 0 {
			c.JSON(http.StatusOK, nil)
			return
		}

		type label struct {
			Name       string  `json:"label"`
			Confidence float64 `json:"confidence"`
		}

		var ll []label
		for _, l := range labels.Labels {
			ll = append(ll, label{Name: *l.Name, Confidence: *l.Confidence})
		}

		c.JSON(http.StatusOK, ll)
		return
	})

	log.Fatal(r.Run())
}
