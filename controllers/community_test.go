package controllers

import (
	"fmt"
	"go_web_app/dao/mysql"
	"go_web_app/dao/redis"
	"go_web_app/logger"
	"go_web_app/middleware"
	"go_web_app/pkg/snowflake"
	"go_web_app/settings"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/smartystreets/goconvey/convey"
)

func Init() {
	settings.Init("../conf/config.yaml")
	logger.Init(settings.Config.LogConfig, settings.Config.Mode)
	mysql.Init(settings.Config.MySQLConfig)
	redis.Init(settings.Config.RedisConfig)
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	if err := snowflake.Init(formattedTime, settings.Config.MachineId); err != nil {
		fmt.Println("Init snowflake failed, err: ", err)
		panic(err)
	}
}

func TestCreatePostHandler(t *testing.T) {

	Init()
	// 测试路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.JWT())
	router.POST("/createpost", CreatePostHandler)

	// 顶层convey
	convey.Convey("CreatePostHandler start test", t, func() {

		convey.Convey("with valid JSON input and valid Bearer Token", func() {
			jsonInput := `{"title": "snow", "content": "2024最好的游戏是幻兽帕鲁", "community_id": 2}`
			req, _ := http.NewRequest("POST", "/createpost", strings.NewReader(jsonInput))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZzcHBsZSIsInVzZXJfaWQiOi0xMTk0MTU2MjY1MzA4MTYsImlzcyI6ImdvX3dlYl9hcHAiLCJleHAiOjE3MDcyMjI2MTh9.NVDLEALm4lh8DIGeKM97-zk6kFwJTxQMTaT_RaUOsFk")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			convey.Convey("should return status 200 OK", func() {
				convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestGetPostDetailHandler(t *testing.T) {
	Init()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/community/post/:post_id", GetPostDetailHandler)

	convey.Convey("GetPostDetailHandler start test", t, func() {

		convey.Convey("with valid post ID", func() {
			postID := -120728552734720
			req, _ := http.NewRequest("GET", fmt.Sprintf("/community/post/%d", postID), nil)
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZzcHBsZSIsInVzZXJfaWQiOi0xMTk0MTU2MjY1MzA4MTYsImlzcyI6ImdvX3dlYl9hcHAiLCJleHAiOjE3MDcyMjI2MTh9.NVDLEALm4lh8DIGeKM97-zk6kFwJTxQMTaT_RaUOsFk")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			convey.Convey("should return status 200 OK and valid response", func() {
				convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestCommunitySortedPostHandler(t *testing.T) {
	Init()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/community/:id/sortedpost", CommunitySortedPostHandler)

	convey.Convey("CommunitySortedPostHandler start test", t, func() {

		convey.Convey("with valid community ID and query parameters", func() {
			//communityID := "2"
			//queryParams := map[string]string{
			//	"orderstr": "time",
			//	"offset":   "0",
			//	"limit":    "10",
			//}
			// req, _ := http.NewRequest("GET", fmt.Sprintf("/community/%s/sortedpost?order=%s&offset=%s&limit=%s", communityID, queryParams["orderstr"], queryParams["offset"], queryParams["limit"]), nil)
			req, _ := http.NewRequest("GET", "http://localhost:8083/community/2/sortedpost?offset=0&limit=10&order=time", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			convey.Convey("should return status 200 OK and valid response", func() {
				convey.So(resp.Code, convey.ShouldEqual, http.StatusOK)
				convey.Println(resp.Body)
			})
		})
	})
}
