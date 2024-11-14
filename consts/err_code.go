package consts

type ErrorType int32

const (
	ErrorType_NoError       ErrorType = 0 // 没有错误
	ErrorType_ErrorNormal   ErrorType = 1 // 正常的错误，可以根据 err_msg 做具体提示
	ErrorType_ErrorProtocol ErrorType = 2 // 协议错误
	ErrorType_ErrorPanic    ErrorType = 3 // 服务器panic
	ErrorType_ErrorDB       ErrorType = 4 // 数据存储相关错误
)

const (
	InternalErrMsg string = "server_internal_error"
)
