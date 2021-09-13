package gserver

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"server/db"
	"server/gserver/cfg"
	"server/gserver/genServer"
	"server/gserver/nodeManange"
	"server/logger"
	"server/network"
	"server/web"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"

	//msg "server/proto"
	"github.com/facebookgo/pidfile"
	"github.com/halturin/ergo"
	"github.com/halturin/ergo/etf"
)

//GameServerInfo game info
var GameServerInfo *gameServer

type gameServer struct {
	nw      *network.NetWorkx
	command chan string
}

func (g *gameServer) Start() {
	nodeManange.Start(&ServerCfg, g.command)
	dbNode := nodeManange.GetNode(fmt.Sprintf("dbNode_%v@127.0.0.1", ServerCfg.ServerID))
	if dbNode == nil {
		panic("节点启动失败")
	}
	//启动网络
	g.nw.Start(dbNode)
}

func (g *gameServer) Close() {
	for _, node := range nodeManange.GetNodes() {
		for _, p := range node.GetProcessList() {
			p.Cast(p.Self(), etf.Atom("stop"))
		}
		// node.Stop()
		// node.Wait()
	}
	g.nw.Close()
}

//StartGServer 启动game server
//go run main.go start --config=E:/worke/server/cfg.yaml
func StartGServer() {
	if ServerRunState() {
		log.Infof("[%v][%v] runing", ServerCfg.ServerName, ServerCfg.ServerID)
	}

	log.Infof("============================= Begin Start [%v][%v] ===============================", ServerCfg.ServerName, ServerCfg.ServerID)
	if level, err := log.ParseLevel(ServerCfg.Loglevel); err == nil {
		logger.Init(level, ServerCfg.LogWrite, ServerCfg.LogName, ServerCfg.LogPath)
	} else {
		logger.Init(log.InfoLevel, ServerCfg.LogWrite, ServerCfg.LogName, ServerCfg.LogPath)
	}

	// if ServerCfg.Daemon {
	//https://github.com/takama/daemon
	// }

	cfg.InitViperConfig(ServerCfg.CfgPath, ServerCfg.CfgType)
	if ServerCfg.WatchConfig {
		cfg.WatchConfig(ServerCfg.CfgPath, func(in fsnotify.Event) {
			log.Debug("Config file changed: [%v]  ", in.Name)
			cfg.InitViperConfig(ServerCfg.CfgPath, ServerCfg.CfgType)
		})
	}

	//启动定时器
	//timedtasks.StartCronTasks()
	// //定时器
	// timedtasks.AddTasks("loop", "* * * * * ?", func() {
	// 	log.Info("server time:", time.Now())
	// })
	//defer timedtasks.RemoveTasks("loop")

	db.StartMongodb(ServerCfg.Mongodb, ServerCfg.MongoConnStr)
	if ok, err := db.MongodbPing(); ok {
		log.Info("mongodb conn success")
	} else {
		panic(err)
	}

	db.StartRedis(ServerCfg.RedisConnStr, ServerCfg.RedisDB)
	if ok, err := db.RedisConn(); ok {
		log.Info("redis conn success")
	} else {
		panic(err)
	}

	if ServerCfg.OpenHTTP {
		go web.Start(ServerCfg.HTTPPort)
	}

	if ServerCfg.StatsView {
		go web.StartStatsView(ServerCfg.StatsViewPort)
	}

	GameServerInfo = &gameServer{
		nw: network.NewNetWorkX(
			func() ergo.GenServerBehaviour {
				return &genServer.GateGenServer{}
			},
			// &sync.Pool{
			// New: func() interface{} {
			// 	return &genServer.GateGenServer{}
			// }},
			ServerCfg.Port,
			ServerCfg.Packet,
			ServerCfg.Readtimeout,
			ServerCfg.NetType,
			ServerCfg.MsgTime,
			ServerCfg.MsgNum,
			func() { SendGameServerMsg("StartSuccess") },
			func() { db.RedisExec("del", "ConnectNumber") },
			func() { log.Info("connect number: ", db.RedisINCRBY("ConnectNumber", 1)) },
			func() { log.Info("connect number: ", db.RedisINCRBY("ConnectNumber", -1)) },
		),
		command: make(chan string),
	}
	GameServerInfo.Start()

	defer ClonseServer()
	defer GameServerInfo.Close()

	//退出消息监控
	var exitChan = make(chan os.Signal)

	if runtime.GOOS == "linux" {
		//signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
		signal.Notify(exitChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTSTP)
	} else {
		signal.Notify(exitChan, os.Interrupt)
	}

	//create pid
	file, _ := ioutil.TempFile("", fmt.Sprintf("pid_%v_%v_", ServerCfg.ServerName, ServerCfg.ServerID))
	pidfile.SetPidfilePath(file.Name())

	for {
		select {
		case command := <-GameServerInfo.command:
			switch command {
			case "StartSuccess":
				pid := StartSuccess()
				log.Infof("====================== Start Game Server pid:[%v] Success =========================", pid)
			case "shutdown":
				log.Warn("Shut down the game server")
				return
			default:
				log.Warn("command:", command)
			}
		case s := <-exitChan:
			log.Info("收到信号: ", s)
			if runtime.GOOS == "linux" && s.String() == "quit" || s.String() == "terminated" {
				return
			} else if runtime.GOOS == "windows" && s.String() == "interrupt" {
				return
			}

			//case <-time.After(60 * time.Second):
			//log.Infof("time: [%v]  online:[%v]", time.Now().Format(tools.DateTimeFormat), db.RedisGetInt("ConnectNumber"))
		}
	}

}

//SendGameServerMsg game system msg
func SendGameServerMsg(msg string) {
	GameServerInfo.command <- msg
}

func ServerRunState() bool {
	return false
}

//成功启动
func StartSuccess() int {
	pidfile.Write()
	log.Infof("pidfile :%v", pidfile.GetPidfilePath())
	i, _ := pidfile.Read()
	return i
}

//关闭服务
func ClonseServer() {
	os.Remove(pidfile.GetPidfilePath())
}
