package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"net/http"
	"time"
)

type Data struct {
	User User `json:"user"`
}
type User struct {
	ID            string    `json:"id" gorm:"column:id;primaryKey;type:bigint"`
	CreatedTime   string    `gorm:"column:created_time;type:bigint"`
	UpdatedTime   string    `gorm:"column:updated_time;type:bigint"`
	Username      string    `json:"username" gorm:"column:username;type:varchar(255);comment:用户名"`
	Password      string    `json:"password" gorm:"column:password;type:varchar(255);comment:密码"`
	Email         string    `json:"email" gorm:"column:email;type:varchar(255);comment:邮箱"`
	Avatar        string    `json:"avatar" gorm:"column:avatar;type:varchar(255);comment:头像URL"`
	Role          string    `json:"role" gorm:"column:role;type:int;comment:权限等级"`
	Phone         string    `gorm:"column:phone;type:varchar(20);comment:手机号" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	ClientIp      string    `gorm:"column:client_ip;type:varchar(50);comment:客户端IP"`
	ClientPort    string    `gorm:"column:client_port;type:varchar(20);comment:客户端端口"`
	LoginTime     time.Time `gorm:"column:login_time;comment:登录时间"`
	HeartbeatTime time.Time `gorm:"column:heartbeat_time;comment:心跳时间"`
	LoginOutTime  time.Time `gorm:"column:login_out_time;comment:登出时间" json:"login_out_time"`
	IsLogout      bool      `gorm:"column:is_logout;comment:是否登出"`
	DeviceInfo    string    `gorm:"column:device_info;type:varchar(255);comment:设备信息"`
	Bio           string    `gorm:"column:bio;type:varchar(255);comment:个人简介"`
}
type GetUserInfoInputParams struct {
	ID string `json:"id" jsonschema:"description=the id of the user"`
}

func GetUserInfo(ctx context.Context, params *GetUserInfoInputParams) (map[string]any, error) {

	// 1. 目标接口
	url := "http://localhost:8080/api/user/get-user-info?user_id=" + params.ID
	//获取 JWT
	//jwt, _ := aiUtils.GetToken()
	meta := ctx.Value("Authorization")
	jwt := meta.(string)
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
	u := d.User

	respMap := map[string]any{
		"id":             u.ID,
		"username":       u.Username,
		"password":       u.Password,
		"email":          u.Email,
		"avatar":         u.Avatar,
		"role":           u.Role,
		"phone":          u.Phone,
		"client_ip":      u.ClientIp,
		"client_port":    u.ClientPort,
		"login_time":     u.LoginTime.Format("2006-01-02 15:04:05"),
		"heartbeat_time": u.HeartbeatTime.Format("2006-01-02 15:04:05"),
		"login_out_time": u.LoginOutTime.Format("2006-01-02 15:04:05"),
		"is_logout":      u.IsLogout,
		"bio":            u.Bio,
		"device_info":    u.DeviceInfo,
	}
	return respMap, nil
}
func CreateGetUserInfoTool() tool.InvokableTool {
	GetUserInfoTool := utils.NewTool(&schema.ToolInfo{
		Name: "GetUserInfo",
		Desc: "Get User Info By ID",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"id": &schema.ParameterInfo{
					Type:     schema.String,
					Desc:     "the id of the user",
					Required: true,
				},
			},
		),
	}, GetUserInfo)
	return GetUserInfoTool
}
