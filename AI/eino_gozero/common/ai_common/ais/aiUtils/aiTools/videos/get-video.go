package videos

import (
	"context"
	"eino_gozero/common/ai_common/ais/aiUtils"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"net/http"
	"time"
)

type Data struct {
	Video Video `json:"Video"`
}
type Video struct {
	ID          string `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	CreatedTime string `gorm:"column:created_time;type:bigint"`
	UpdatedTime string `gorm:"column:updated_time;type:bigint"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Cover       string `json:"cover"`
	Likes       string `json:"likes" gorm:"not null;column:likes;type:int;comment:点赞数"`
	Comments    string `json:"comments" gorm:"not null;column:comments;type:int;comment:评论数"`
	Shares      string `json:"shares" gorm:"not null;column:shares;type:int;comment:分享数"`
	VideoID     string `json:"Video_id" gorm:"column:Video_id;type:bigint;comment:用户ID"`
	PublishTime time.Time
	Type        string `json:"type" gorm:"column:type;type:varchar(63);comment:类型"`
	IsPrivate   bool   `json:"is_private" gorm:"not null;column:is_private;type:bool;comment:是否私密"`
}
type GetVideoInputParams struct {
}

func GetVideo(ctx context.Context, params *GetVideoInputParams) (map[string]any, error) {

	// 1. 目标接口
	url := "http://localhost:8080/api/videos"
	//获取 JWT
	jwt, _ := aiUtils.GetToken()
	// 4. 创建请求
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(">>> 发送失败:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var tmp struct {
		Code int             `json:"code"`
		Msg  string          `json:"message"`
		Data json.RawMessage `json:"data"` // 先拿原始字节
	}
	if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
		return nil, err
	}
	// 再二次反序列化到真正的结构
	var d Data
	if err := json.Unmarshal(tmp.Data, &d); err != nil {
		return nil, err
	}
	u := d.Video
	fmt.Println("$%^%$#$%$#$%^%$##$%^%$#$%^%$#@")
	fmt.Println(u)
	fmt.Println("$%^%$#$%$#$%^%$##$%^%$#$%^%$#@")

	respMap := map[string]any{
		"id":           u.ID,
		"created_time": u.CreatedTime,
		"updated_time": u.UpdatedTime,
		"title":        u.Title,
		"description":  u.Description,
		"url":          u.URL,
		"cover":        u.Cover,
		"likes":        u.Likes,
		"comments":     u.Comments,
		"shares":       u.Shares,
		"video_id":     u.VideoID,
		"publish_time": u.PublishTime,
		"type":         u.Type,
		"is_private":   u.IsPrivate,
	}
	return respMap, nil
}
func CreateGetVideoTool() tool.InvokableTool {
	GetVideoTool := utils.NewTool(&schema.ToolInfo{
		Name: "GetVideo",
		Desc: "Get Videos",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{},
		),
	}, GetVideo)
	return GetVideoTool
}
