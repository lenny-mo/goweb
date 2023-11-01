package controllers

import (
	"go_web_app/logic"
	"go_web_app/models"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func PostVoteHandler(c *gin.Context) {
	// 参数校验
	votedata := new(models.VoteData)
	if err := c.ShouldBindJSON(votedata); err != nil {
		if errcheck, ok := err.(validator.ValidationErrors); ok {
			zap.L().Error("c.ShouldBindJSON(votedata) failed", zap.Error(errcheck))
			ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		}
		zap.L().Error("c.ShouldBindJSON(votedata) failed", zap.Error(err))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}

	// 业务逻辑: 投票
	userid := c.GetInt64("userid")
	err := logic.VoteForPost(votedata, userid)
	if err != nil {
		ReturnResponse(c, http.StatusInternalServerError, InvalidParamCode)
		return
	}

	ReturnResponse(c, http.StatusOK, SuccessCode)
}

// SortedPostHandler 获取社区下的帖子列表, 并且根据时间或者投票分数进行排序
func SortedPostHandler(c *gin.Context) {
	// get请求，从url 中获取参数
	// 跳过offset前面的行数，从offset+1行开始取limit 行数据
	querydata := new(models.PostListParam)
	if err := c.ShouldBind(querydata); err != nil {
		zap.L().Error("c.ShouldBindJSON(querydata) failed", zap.Error(err))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}

	// 2. 根据上述的offset和limit 查询redis 对应的数据，下放到logic层
	data, err := logic.GetSortedPost(querydata)
	if err != nil {
		zap.L().Error("logic.GetSortedPost(querydata) failed", zap.Error(err))
		ReturnResponse(c, http.StatusInternalServerError, InvalidParamCode)
	}

	// 3.
	ReturnResponse(c, http.StatusOK, SuccessCode, data)

}
