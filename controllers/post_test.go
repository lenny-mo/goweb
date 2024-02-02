package controllers

import (
	"bytes"
	"encoding/json"
	"go_web_app/models"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/smartystreets/goconvey/convey"
)

func TestPostVoteHandler(t *testing.T) {
	convey.Convey("Given a request to vote for a post", t, func() {
		// Set up a fake Gin context with a valid vote data and user ID
		c, _ := gin.CreateTestContext(nil)
		c.Set("userid", int64(123)) // Replace with a valid user ID
		voteData := &models.VoteData{
			// Set valid vote data here
			PostID: 12,
			Vote:   1,
		}
		rawdata, _ := json.Marshal(voteData)
		// Set up the request with raw JSON data and content type
		c.Request = &http.Request{
			Method: "POST",
			URL: &url.URL{
				Path: "/your-endpoint", // Replace with your actual endpoint
			},
			Header: make(http.Header),
			Host:   "localhost",
			Body:   io.NopCloser(bytes.NewBuffer(rawdata)),
		}
		c.Request.Header.Set("Content-Type", "application/json")

		PostVoteHandler(c)

		convey.Convey("It should return a successful response", func() {
			convey.So(c.Writer.Status(), convey.ShouldEqual, http.StatusOK)
			// Add more assertions based on your expected response
		})
	})
}

func TestSortedPostHandler(t *testing.T) {
	convey.Convey("Given a request to retrieve sorted posts", t, func() {
		// Set up a fake Gin context with valid query data
		c, _ := gin.CreateTestContext(nil)
		queryData := &models.PostListParam{
			// Set valid query data here
			Offset: 1,
			Limit:  1,
			Order:  "time",
		}
		rawdata, _ := json.Marshal(queryData)

		c.Request = &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path:     "/your-endpoint", // Replace with your actual endpoint
				RawQuery: string(rawdata),  // Encode the query parameters
			},
			Header: make(http.Header),
			Host:   "localhost",
		}

		SortedPostHandler(c)

		convey.Convey("It should return a successful response with sorted post data", func() {
			convey.So(c.Writer.Status(), convey.ShouldEqual, http.StatusOK)
			// Add more assertions based on your expected response
		})
	})
}
