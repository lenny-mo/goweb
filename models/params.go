package models

// 请求的参数结构体，在logic层, controllers层都会用到，
// 所以放在models层
type SignupParam struct {
	Username string `json:"username" binding:"required" ` // binding:"required"表示必须要传，gin的validator使用
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Gender   uint8  `json:"gender" binding:"required"`
}

type LoginParam struct {
	Username string `json:"username" binding:"required" ` // binding:"required"表示必须要传
	Password string `json:"password" binding:"required"`
}

type VoteData struct {
	PostID      int64  `json:"post_id,string" binding:"required"`
	Vote        int8   `json:"vote,string" binding:"oneof=-1 0 1"` // 1赞成 0取消赞成 -1反对, 不要设置required
	CommunityID string `json:"community_id" binding:"required"`
}

// PostListParam 获取帖子列表的请求参数
//
// 支持json和form两种方式，form包含了query string和post form
type PostListParam struct {
	Offset int64  `json:"offset,string" form:"offset,string"`
	Limit  int64  `json:"limit,string" form:"limit,string"`
	Order  string `json:"order,string" form:"order,string" binding:"oneof=time vote"` //
}
