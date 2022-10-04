package notify

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/httpclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"strings"
)

type WebHook struct {
	//飞书还是其他
	Kind string
	Url  string
}

var _defaultWebHook *WebHook

func (w *WebHook) SendMsg(msg *Message) {
	//todo
	switch _defaultWebHook.Kind {
	case "feishu":
		var sendData = feiShuTemplateCard
		sendData = strings.Replace(sendData, "timeSlot", msg.OccurTime.String(), 1)
		userSlot := ""
		for _, to := range msg.To {
			userSlot += fmt.Sprintf("<at>%s</at>", to)
		}
		sendData = strings.Replace(sendData, "userSlot", userSlot, 1)
		sendData = strings.Replace(sendData, "msgSlot", msg.Body, 1)
		err := httpclient.PostJson(w.Url, sendData, 0)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("feishu  send msg[%+v] err: %s", msg, err.Error()))
		}
	default:
		err := httpclient.PostJson(w.Url, msg, 0)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("web hook api send msg[%+v] err: %s", msg, err.Error()))
		}
	}
	return
}
