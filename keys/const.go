package keys

const (
	ErrorPhoneCode      = "手机验证码错误"
	ErrorPhoneEmpty     = "手机号不能为空"
	ErrorPhoneExists    = "该手机号已经被注册"
	ErrorUsernameEmpty  = "手机号或邮箱不能同时为空"
	ErrorUsernameExists = "该手机号或邮箱已经被注册"
	ErrorPasswordSafe   = "密码不能包含空格，长度为8-20个英文字符，必须同时包含数字和字母"
	ErrorCaptchaCode    = "验证码错误"
	ErrorLogin          = "用户名或密码错误"
	ErrorParam          = "参数错误"
	ErrorParamPage      = "分页参数错误"
	ErrorPermission     = "权限错误"
	ErrorSave           = "保存数据错误，请稍后重试"
	ErrorFile           = "读取上传文件错误"
	ErrorFileMaxSize    = "上传文件大小不能超过10M"
	ErrorFileInfo       = "无法读取文件"
	ErrorFileSave       = "保存文件错误"
	ErrorPassword       = "密码错误"
	ErrorUserNoExists   = "用户不存在"
	ErrorEmail          = "发送邮件错误，请稍后重试"
	ErrorActiveCode     = "激活码错误"
	ErrorNoActive       = "当前账号不可用，请联系我们"
	ErrorRead           = "数据不存在；或读取数据错误，请稍后重试"
	ErrorNeedSign       = "请先登录"
	ErrorProxy          = "网络错误，请稍后重试"
	ErrorPhoneSms       = "发送短信失败，请检查手机号是否正确或稍后重试"
	ErrorSmsTimer       = "发送短信过快，请稍后重试"
	ErrorServer         = "读取第三方服务错误，请稍后重试"
	ErrorID             = "无效的ID"
)

const (
	PageIndex = "pageIndex"
	PageCount = "pageCount"
)
