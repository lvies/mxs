package worldmanager

import (
	"mxs/api/iface"
	"mxs/api/mnet"
	"mxs/log"
	"mxs/mmo/proto/flat/sample/strupro"
)

// 世界聊天 路由业务
type WorldChat struct {
	 mnet.BaseRouter
}

func (*WorldManager) Handler(request iface.IRequest) {
	msg := strupro.GetSizePrefixedRootAsTalkMessage(request.GetData(), 32)

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		log.Error("GetProperty pid error %v", pid)
		request.GetConnection().Stop()
		return
	}

	player := WorldMgrObj.GetPlayerByEid(pid.(int32))

	player.BubbleTalk(msg.Table().Bytes)
}

