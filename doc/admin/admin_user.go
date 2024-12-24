package admin

type AdminUser struct {
}

type AdminUserLoginFormat struct {
	Code int32             `json:"code"`
	Msg  string            `json:"msg"`
	Data AdminUserLoginRes `json:"data"`
}

type AdminUserLoginReq struct {
	Username string `json:"username"  validate:"required"` //用户名
	Password string `json:"password"  validate:"required"` //密码
}

type AdminUserLoginRes struct {
	Token string `json:"token"` //token
	Name  string `json:"name"`  //name
}

// Login
// @Tags 后台用户
// @Summary 登陆
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body AdminUserLoginReq true "数据"
// @Success 200 {object} AdminUserLoginFormat
// @Router /admin/admin-user/login [POST]
func (AdminUser) Login() {

}

type AdminUserLogoutFormat struct {
	Code int32              `json:"code"`
	Msg  string             `json:"msg"`
	Data AdminUserLogoutRes `json:"data"`
}

type AdminUserLogoutReq struct {
}

type AdminUserLogoutRes struct {
}

// Logout
// @Tags 后台用户
// @Summary 登出
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body AdminUserLogoutReq true "数据"
// @Success 200 {object} AdminUserLogoutFormat
// @Router /admin/admin-user/logout [POST]
func (AdminUser) Logout() {

}

type AdminUserInfoFormat struct {
	Code int32            `json:"code"`
	Msg  string           `json:"msg"`
	Data AdminUserInfoRes `json:"data"`
}

type AdminUserInfoReq struct {
}

type AdminUserInfoRes struct {
	UserId string `json:"user_id"` //user_id
	Name   string `json:"name"`    //name
	Avatar string `json:"avatar"`  //name
}

// Info
// @Tags 后台用户
// @Summary 信息接口
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data query AdminUserInfoReq true "数据"
// @Success 200 {object} AdminUserInfoFormat
// @Router /admin/admin-user/info [GET]
func (AdminUser) Info() {

}
