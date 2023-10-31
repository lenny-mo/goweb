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
