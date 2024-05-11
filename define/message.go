package define

// 发送消息结构体
type SendMsg struct {
	TargetId    int         `json:"target_id"`    //接收者
	Type        int         `json:"type"`         //1-私聊 2-群聊 3-广播
	MessageType int         `json:"message_type"` //消息类型  文字、图片、音频
	Content     interface{} `json:"content"`      //消息内容
	Cmd         string      `json:"cmd"`          //请求命令字
	Date        string      `json:date`
}

// 通用websocket请求数据格式
// type Request struct {
// 	Seq  string      `json:"seq"`            // 消息的唯一Id
// 	Cmd  string      `json:"cmd"`            // 请求命令字
// 	Data interface{} `json:"data,omitempty"` // 数据 json
// }
