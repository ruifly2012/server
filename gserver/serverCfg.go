package gserver

// ServerConfig  server cfg
type ServerConfig struct {
	ServerName string
	ServerID   int32
	Version    string

	Daemon     bool
	RestartNum int

	OpenHTTP bool
	HTTPPort int32

	StatsView     bool
	StatsViewPort int32

	NetType     string
	Port        int32
	Packet      int32
	Readtimeout int32 //读超时时间

	MsgTime int32
	MsgNum  int32

	ProtoPath string
	GoOut     string

	MongoConnStr string
	Mongodb      string

	RedisConnStr string
	RedisDB      int

	CfgPath string
	CfgType string

	LogWrite bool
	Loglevel string
	LogPath  string
	LogName  string

	ListenRangeBegin int
	ListenRangeEnd   int
	EPMDPort         int
	Cookie           string
}

// ServerCfg  Program overall configuration
var ServerCfg = ServerConfig{
	ServerName: "server",
	ServerID:   1,
	Version:    "1.0.0",

	Daemon:     false,
	RestartNum: 2,

	// http
	OpenHTTP: false,
	HTTPPort: 8080,

	StatsView:     false,
	StatsViewPort: 8081,

	// #network : tcp/udp
	NetType:     "tcp",
	Port:        3344,
	Packet:      2,
	Readtimeout: 0,

	MsgTime: 300,
	MsgNum:  500,

	// #protobuf path
	ProtoPath: "./proto",
	GoOut:     "./msgproto",

	MongoConnStr: "mongodb://localhost:27017",
	Mongodb:      "mygame",

	RedisConnStr: "localhost:6379",
	RedisDB:      0,

	CfgPath: "./config",
	CfgType: "json",

	Loglevel: "info",
	LogPath:  "./log",
	LogName:  "log",
	LogWrite: false,

	ListenRangeBegin: 15151,
	ListenRangeEnd:   25151,
	EPMDPort:         4369,
	Cookie:           "123",
}
