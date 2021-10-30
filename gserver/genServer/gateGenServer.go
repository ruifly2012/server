package genServer

import (
	"server/network"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

//接收处理socket 发送过来的信息
// 处理玩家独立无交互的游戏逻辑
// 在socket中断后 此进程会保留一段时间以便于重新建立连接

type GateGenServer struct {
	gen.Server
	process  *gen.ServerProcess
	sendChan chan []byte
}

func (gateGS *GateGenServer) Unregister() {
	//gateGS.process.Node.Unregister(gateGS.process.Name())
}

func (gateGS *GateGenServer) Register(name string) {
	//gateGS.process.Node.Register(name, gateGS.process.Self())
}

func (gateGS *GateGenServer) Init(process *gen.ServerProcess, args ...etf.Term) error {
	log.Infof("Init (%v): args %v ", process.Name(), args)
	gateGS.process = process
	gateGS.sendChan = args[0].(chan []byte)
	return nil
}

func (gateGS *GateGenServer) HandleCast(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Infof("HandleCast (%v): %v", gateGS.process.Name(), message)

	switch info := message.(type) {
	case etf.Atom:
		switch info {
		case "SocketStop":
			return gen.ServerStatusStopWithReason("stop normal")
		case "timeloop":
			log.Debug("time loop")
		}
	case etf.Tuple:
		module := info[0].(int)
		method := info[1].(int)
		buf := info[2].([]byte)

		log.Debug("socket info ", module, method, buf)
		//gateGS.sendChan <- []byte("send msg test")
	case []byte:
		log.Debug("[]byte:", info)
	}
	return gen.ServerStatusOK
}

func (gateGS *GateGenServer) HandleCall(process *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	log.Infof("HandleCall (%v): %v, From: %v", gateGS.process.Name(), message, from)

	reply := etf.Term(etf.Tuple{etf.Atom("error"), etf.Atom("unknown_request")})

	return reply, gen.ServerStatusOK
}

func (gateGS *GateGenServer) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Infof("HandleInfo (%v): %v", gateGS.process.Name(), message)
	return gen.ServerStatusOK
}

// Terminate called when process died
func (gateGS *GateGenServer) Terminate(process *gen.ServerProcess, reason string) {
	gateGS.Unregister()
	log.Infof("Terminate (%v): %v", gateGS.process.Name(), reason)
}

// //Send 发送消息
func (gateGS *GateGenServer) Send(module int32, method int32, pb proto.Message) {
	//log.Debugf("client send msg [%v] [%v] [%v]", module, method, pb)
	data, err := proto.Marshal(pb)
	if err != nil {
		log.Errorf("proto encode error[%v] [%v][%v] [%v]", err.Error(), module, method, pb)
		return
	}
	// msginfo := &common.NetworkMsg{}
	// msginfo.Module = module
	// msginfo.Method = method
	// msginfo.MsgBytes = data
	// msgdata, err := proto.Marshal(msginfo)
	// if err != nil {
	// 	log.Errorf("msg encode error[%s]\n", err.Error())
	// }
	// gateGS.sendChan <- msgdata

	mldulebuf := network.IntToBytes(module, 2)
	methodbuf := network.IntToBytes(method, 2)
	gateGS.sendChan <- network.BytesCombine(mldulebuf, methodbuf, data)

}

//protobuf 解码
// func decode(pb proto.Message, buf []byte) bool {
// 	if e := proto.Unmarshal(buf, pb); e != nil {
// 		log.Error(e)
// 		return false
// 	}
// 	return true
// }
