package gserver

import "fmt"

type client struct {
	//OnMessage func(module int, method int, buf []byte)
	//handlers map[int32]interface{}
}

func (c *client) OnConnect() {
	fmt.Println("client object connect msg")
}

func (c *client) OnMessage(module int, method int, buf []byte) {

}

func (c *client) OnClose() {

}

//UserLogin 用户登陆
func (c *client) UserLogin(buf []byte) {

}
