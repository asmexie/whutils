package whutils

//HTTPReplyJSON ...
type HTTPReplyJSON struct {
	httpMsgInfo map[int]string
	Reply       map[string]interface{}
}

//SetStatus ...
func (reply *HTTPReplyJSON) SetStatus(code int) {
	reply.Reply["status"] = code
	reply.Reply["msg"] = reply.httpMsgInfo[code]
}

// GetStatus ...
func (reply *HTTPReplyJSON) GetStatus() int {
	if v, ok := reply.Reply["status"]; ok {
		return v.(int)
	}
	return 0
}

//SetStatusAndMsg ...
func (reply *HTTPReplyJSON) SetStatusAndMsg(code int, msg string) {
	reply.Reply["status"] = code
	reply.Reply["msg"] = msg
}

//SetStatusAndData ...
func (reply *HTTPReplyJSON) SetStatusAndData(code int, data interface{}) {
	reply.Reply["status"] = code
	reply.Reply["msg"] = reply.httpMsgInfo[code]
	reply.Reply["data"] = data
}

//SetStatusAndMsgAndData ...
func (reply *HTTPReplyJSON) SetStatusAndMsgAndData(code int, msg string, data interface{}) {
	reply.Reply["status"] = code
	reply.Reply["msg"] = msg
	reply.Reply["data"] = data
}

//SetStatusAndMsgAndData ...
func (reply *HTTPReplyJSON) SetCodeAndMsgAndData(code int, msg string, data interface{}) {
	reply.Reply["code"] = code
	reply.Reply["msg"] = msg
	reply.Reply["data"] = data
}

// NewHTTPReplyJSON ...
func NewHTTPReplyJSON(msgMap map[int]string, defaultErr int) *HTTPReplyJSON {
	return &HTTPReplyJSON{
		httpMsgInfo: msgMap,
		Reply: map[string]interface{}{
			"status": defaultErr,
			"msg":    msgMap[defaultErr],
		},
	}
}
