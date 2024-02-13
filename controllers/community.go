package controllers

import (
	"go_web_app/contextcode"
	"go_web_app/logic"
	"go_web_app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	post := new(models.Post)
	// 要获取body 中的json 数据，必须通过结构体绑定的形式，shouldbindwith
	if err := c.ShouldBindJSON(post); err != nil {
		zap.L().Error("c.ShouldBindJSON(post) failed", zap.Error(err))
		contextcode.ReturnResponse(c, http.StatusBadRequest, contextcode.InvalidParamCode)
		return
	}

	id, ok := c.Get(contextcode.ContextUserIDKey) // 应该是JWT提供了信息
	if !ok {
		zap.L().Error("c.Get(ContextUserIDKey) failed")
		contextcode.ReturnResponse(c, http.StatusInternalServerError, contextcode.InvalidParamCode)
		return
	}
	post.AuthorID = id.(int64)

	// 2. 业务逻辑处理 创建帖子
	err := logic.CreatePost(post)
	// 3. 返回响应
	if err != nil {
		zap.L().Error("logic.CreatePost(post) failed", zap.Error(err))
		contextcode.ReturnResponse(c, http.StatusInternalServerError, contextcode.InvalidParamCode)
		return
	}
	contextcode.ReturnResponse(c, http.StatusOK, contextcode.SuccessCode, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数,post id 不需要community id 就可以访问post 详情
	// 从url 中获取post id
	// 要求获取url 中的数据，使用c.Param
	postIdStr := c.Param("post_id")
	if len(postIdStr) == 0 {
		zap.L().Error("invalid post id")
		contextcode.ReturnResponse(c, http.StatusBadRequest, contextcode.InvalidParamCode)
		return
	}
	// 转化成int64 类型
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(postIdStr, 10, 64) failed", zap.Error(err))
		contextcode.ReturnResponse(c, http.StatusBadRequest, contextcode.InvalidParamCode)
		return
	}

	// 2. 业务逻辑处理: 通过post id 查询帖子详情
	// 包含作者信息，社区信息， titile, content
	data, err := logic.GetPostDetailById(postId)
	if err != nil {
		zap.L().Error("logic.GetPostDetailById(postId) failed", zap.Error(err))
		contextcode.ReturnResponse(c, http.StatusInternalServerError, contextcode.InvalidParamCode)
		return
	}
	// 3. 返回响应
	contextcode.ReturnResponse(c, http.StatusOK, contextcode.SuccessCode, data)
}

func CommunitySortedPostHandler(c *gin.Context) {
	// 1. 获取社区id, 动态路由
	communityID := c.Param("id")
	if len(communityID) == 0 {
		zap.L().Error("invalid community id")
		contextcode.ReturnResponse(c, http.StatusBadRequest, contextcode.InvalidParamCode)
		return
	}

	// 2. 获取query string 参数
	querydata := new(models.PostListParam)
	if err := c.ShouldBind(querydata); err != nil {
		zap.L().Error("c.ShouldBindJSON(querydata) failed", zap.Error(err))
		contextcode.ReturnResponse(c, http.StatusBadRequest, contextcode.InvalidParamCode)
		return
	}

	// 3. 业务逻辑处理: 传递社区id, orderstr, offset, limit
	data, err := logic.CommunitySortedPost(communityID, querydata)
	if err != nil {
		zap.L().Error("logic.CommunitySortedPost(communityID, querydata) failed", zap.Error(err))
		contextcode.ReturnResponse(c, http.StatusInternalServerError, contextcode.InvalidParamCode)
		return
	}

	zap.L().Info("logic.CommunitySortedPost(communityID, querydata) success", zap.Any("data", data))
	contextcode.ReturnResponse(c, http.StatusOK, contextcode.SuccessCode, data)
	return
}
