package controllers

import (
	"go_web_app/contextcode"
	"go_web_app/logic"
	"go_web_app/models"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// PostVoteHandler 投票功能
func PostVoteHandler(c *gin.Context) {
	// 参数校验
	votedata := new(models.VoteData)
	if err := c.ShouldBindJSON(votedata); err != nil {
		if errcheck, ok := err.(validator.ValidationErrors); ok {
			zap.L().Error("c.ShouldBindJSON(votedata) failed", zap.Error(errcheck))
			contextcode.ReturnResponse(c, http.StatusBadRequest, contextcode.InvalidParamCode)
		}
		zap.L().Error("c.ShouldBindJSON(votedata) failed", zap.Error(err))
		contextcode.ReturnResponse(c, http.StatusBadRequest, contextcode.InvalidParamCode)
		return
	}

	// 业务逻辑: 投票
	userid := c.GetInt64("userid") // JWT 获取到用户信息
	err := logic.VoteForPost(votedata, userid)
	if err != nil {
		contextcode.ReturnResponse(c, http.StatusInternalServerError, contextcode.InvalidParamCode)
		return
	}

	contextcode.ReturnResponse(c, http.StatusOK, contextcode.SuccessCode)
}
