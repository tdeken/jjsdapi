package code

const OK = 200
const OKMsg = "ok"

const (
	SystemErrorCode = -1  //系统级别错误
	CommonErrorCode = 0   //通用错误类型的都用这个
	VerifyErrorCode = 400 //表单验证不通过
	AuthErrorCode   = 401 // 身份校验不通过
)

var (
	SystemError = NewError(SystemErrorCode, "存在一些问题，请稍后再试") //系统异常
	AuthError   = NewError(AuthErrorCode, "无效token")        //权限错误（例如：数据的企业ID与登录者的企业ID不一致，登录者身份信息异常等）
)

const (
	AdminUserLoginErrCode = 10001 // 登陆失败错误码
)

var (
	AdminUserLoginErr = NewError(AdminUserLoginErrCode, "用户名或密码错误")
)
