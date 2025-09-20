package types

type UploadFileResponse struct {
	Url string `json:"url"`
}
type MQUploadResultRequest struct {
	ID string `json:"id"` // 请求唯一标识
}
type MQUploadResultResponse struct {
}
type UploadSSERequest struct {
	ID string `json:"id"` // 请求唯一标识
}

// UploadTask 塞进 MQ 的消息体
type UploadTask struct {
	ID       string `json:"id"`       // 请求唯一标识
	TmpPath  string `json:"tmp_path"` // 本地临时文件路径
	FileName string `json:"file_name"`
	Ext      string `json:"ext"`
}

// UploadResult 上传最终结果（写到 Redis）
type UploadResult struct {
	ID    string `json:"id"`
	OK    bool   `json:"ok"`
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}
