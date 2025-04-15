package pkg

import "errors"

// 网关错误
var (
	TCPErrInPackIsEmpty = errors.New("握手包没有内容")
	TCPErrAuthNotFirst  = errors.New("不是握手包")
	TCPErrAuthFail      = errors.New("鉴权失败")
	TCPErrAuthExpired   = errors.New("token过期")
)

// DB错误
var (
	DBParamInvalid      = errors.New("DB参数错误")
	DBRecordExist       = errors.New("DB数据已存在")
	DBRecordNotFound    = errors.New("DB未找到数据")
	DBRecordNotAffected = errors.New("DB没有数据行被改变")
)

// 客户端请求错误
var (
	APISessionErr             = errors.New("TCP上下文错误")
	APIPasswordIllegal        = errors.New("密码错误")
	APIPlatformIdTokenInvalid = errors.New("平台id_token错误")
	APIPlatformIdInvalid      = errors.New("平台id错误")
	APIPlatformVerifyInvalid  = errors.New("平台验证错误")
	APIServerErr              = errors.New("服务器内部错误")
	APIRequestIllegal         = errors.New("非法请求")
	APIStateIllegal           = errors.New("非法请求")
)
