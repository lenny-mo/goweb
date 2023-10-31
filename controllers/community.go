package controllers

import (
	"go_web_app/logic"
	"go_web_app/models"
	"go_web_app/pkg/snowflake"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityListHandler 获取社区的所有分类
func CommunityListHandler(c *gin.Context) {
	// 查询所有社区的信息

	// 2. 业务逻辑：查询所有社区的信息
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err)) // 不要把服务端报错返回给前端
		ReturnResponse(c, http.StatusInternalServerError, InvalidParamCode)
	}

	// 3. 返回响应
	ReturnResponse(c, http.StatusOK, SuccessCode, data)
}

// CommunityDetailHandler 获取某个社区的详细信息
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区id
	communityIdstr := c.Param("id")
	// 做参数校验
	if len(communityIdstr) == 0 {
		zap.L().Error("invalid community id")
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}
	// 转化成int64
	communityId, err := strconv.ParseInt(communityIdstr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(communityIdstr, 10, 64) failed", zap.Error(err))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
	}

	// 2. 业务逻辑 通过id查询社区详情
	data, err := logic.GetCommunityDetailById(communityId)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetailById(communityId) failed", zap.Error(err))
		ReturnResponse(c, http.StatusInternalServerError, InvalidParamCode)
	}
	ReturnResponse(c, http.StatusOK, SuccessCode, data)
}

func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	post := new(models.Post)
	// 要获取body 中的json 数据，必须通过结构体绑定的形式，shouldbindwith
	if err := c.ShouldBindJSON(post); err != nil {
		zap.L().Error("c.ShouldBindJSON(post) failed", zap.Error(err))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}

	id, ok := c.Get(ContextUserIDKey)
	if !ok {
		zap.L().Error("c.Get(ContextUserIDKey) failed")
		ReturnResponse(c, http.StatusInternalServerError, InvalidParamCode)
		return
	}
	post.PostID = snowflake.GetId()
	post.AuthorID = id.(int64)

	// 2. 业务逻辑处理 创建帖子
	err := logic.CreatePost(post)
	// 3. 返回响应
	if err != nil {
		zap.L().Error("logic.CreatePost(post) failed", zap.Error(err))
		ReturnResponse(c, http.StatusInternalServerError, InvalidParamCode)
		return
	}
	ReturnResponse(c, http.StatusOK, SuccessCode, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数,post id 不需要community id 就可以访问post 详情
	// 从url 中获取post id
	// 要求获取url 中的数据，使用c.Param
	postIdStr := c.Param("post_id")
	if len(postIdStr) == 0 {
		zap.L().Error("invalid post id")
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}
	// 转化成int64 类型
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(postIdStr, 10, 64) failed", zap.Error(err))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}

	// 2. 业务逻辑处理: 通过post id 查询帖子详情
	// 包含作者信息，社区信息， titile, content
	data, err := logic.GetPostDetailById(postId)
	if err != nil {
		zap.L().Error("logic.GetPostDetailById(postId) failed", zap.Error(err))
		ReturnResponse(c, http.StatusInternalServerError, InvalidParamCode)
		return
	}
	// 3. 返回响应
	ReturnResponse(c, http.StatusOK, SuccessCode, data)
}

// GetPostListHandler 获取社区下的帖子列表
func GetPostListHandler(c *gin.Context) {

	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	if len(offsetStr) == 0 || len(limitStr) == 0 {
		zap.L().Error("invalid offset or limit")
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}
	offset, err1 := strconv.ParseInt(offsetStr, 10, 64)
	limit, err2 := strconv.ParseInt(limitStr, 10, 64)
	if err1 != nil || err2 != nil {
		zap.L().Error("strconv.ParseInt(offsetStr, 10, 64) failed", zap.Error(err1))
		zap.L().Error("strconv.ParseInt(limitStr, 10, 64) failed", zap.Error(err2))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}

	// 获取社区的id, 从url 中获取，查询所有的帖子，要求community id= 指定的id
	communityIdStr := c.Param("id")
	if len(communityIdStr) == 0 {
		zap.L().Error("invalid community id", zap.String("communityIdStr", communityIdStr))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}
	// 转化成int64
	communityId, err := strconv.ParseInt(communityIdStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(communityIdStr, 10, 64) failed", zap.Error(err))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}

	// 2. 业务逻辑处理: 通过community id 查询帖子列表
	data, err := logic.GetPostListByCommunityId(communityId, offset, limit)
	if err != nil {
		zap.L().Error("logic.GetPostListByCommunityId(communityId) failed", zap.Int64("communityId", communityId), zap.Error(err))
		ReturnResponse(c, http.StatusInternalServerError, InvalidParamCode)
		return
	}

	ReturnResponse(c, http.StatusOK, SuccessCode, data)
}
