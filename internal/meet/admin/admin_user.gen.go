// DO NOT EDIT. DO NOT EDIT. DO NOT EDIT.

package admin

type AdminUserLoginReq struct {
	Username string `json:"username" form:"username" validate:"required"` //用户名
	Password string `json:"password" form:"password" validate:"required"` //密码
}

type AdminUserLoginRes struct {
	Token string `json:"token"` //token
	Name  string `json:"name"`  //name
}

type AdminUserLogoutReq struct {
}

type AdminUserLogoutRes struct {
}

type AdminUserInfoReq struct {
}

type AdminUserInfoRes struct {
	UserId string `json:"user_id"` //user_id
	Name   string `json:"name"`    //name
	Avatar string `json:"avatar"`  //name
}
