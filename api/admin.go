package api

const (
	SUCCESS               = 0
	INVALID_PARAMETERS    = 1
	UNSUPPORTED_OPERATION = 2

	// 内部错误
	INTERNAL_ERROR = 10

	// 外部错误，20开头
	EXTERNAL_ERROR = 20
)

type BaseResponse struct {
	Code int
	Msg  string
}
type BaseRequest struct {
	TraceID string
}

type TtsRequest struct {
	BaseRequest
	Text string
}

type TtsResponse struct {
	BaseResponse
	ID string
}

type MorningNoticeRequest struct {
	BaseRequest
}

type MorningNoticeResponse struct {
	BaseResponse
}
