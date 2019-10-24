package communication

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"gofaces/controller"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	ID     string
	conn   *websocket.Conn
	cancel chan int
}

type Message struct {
	ID      string
	Content string
	SentAt  int64
	Type    int
}

type AliveList struct {
	ConnList  map[string]*Client
	register  chan *Client
	destroy   chan *Client
	broadcast chan Message
	receive   chan Message
	cancel    chan int
	Len       int
}

const (
	// SystemMessage 系统消息
	SystemMessage = iota
	// BroadcastMessage 广播消息(正常的消息)
	BroadcastMessage
	// HeartBeatMessage 心跳消息
	HeartBeatMessage
	// ConnectedMessage 上线通知
	ConnectedMessage
	// DisconnectedMessage 下线通知
	DisconnectedMessage
	// BreakMessage 服务断开链接通知(服务端关闭)
	BreakMessage
)

//var aliveList *AliveList
var upgrader = websocket.Upgrader{}
var addr = flag.String("addr", ":8080", "http service address")

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}

func NewAliveList() *AliveList {
	return &AliveList{
		ConnList:  make(map[string]*Client, 100),
		register:  make(chan *Client, 100),
		destroy:   make(chan *Client, 100),
		broadcast: make(chan Message, 100),
		receive:   make(chan Message, 100),
		cancel:    make(chan int),
		Len:       0,
	}
}

// Register 注册
func (al *AliveList) Register(client *Client) {
	al.register <- client
}

// Broadcast 个人广播消息
func (al *AliveList) Broadcast(message Message) {
	al.broadcast <- message
}

// Destroy 销毁
func (al *AliveList) Destroy(client *Client) {
	al.destroy <- client
}

//接收到消息
func (al *AliveList) ReciveMessage(message Message) {
	al.receive <- message
}

func (al *AliveList) sendMessage(id string, msg Message) error {
	if conn, ok := al.ConnList[id]; ok {
		return conn.SendMessage(msg.Type, msg.Content, al)
	}
	return fmt.Errorf("conn not found:%v", msg)
}

func (al *AliveList) run() {
	log.Println("开始监听注册事件")
	for {
		select {
		case client := <-al.register:
			log.Println("注册事件：", client.ID)
			al.ConnList[client.ID] = client
			al.Len++
		case client := <-al.destroy:
			log.Println("销毁事件:", client.ID)
			delete(al.ConnList, client.ID)
			al.Len--
		case message := <-al.broadcast:
			log.Printf("广播事件: %s %s %d \n", message.ID, message.Content, message.Type)
			for id := range al.ConnList {
				err := al.sendMessage(id, message)
				if err != nil {
					log.Println("broadcastError:", err)
				}
			}
		case message := <-al.receive:
			log.Print("消息接收事件:", message)
		case sign := <-al.cancel:
			log.Println("终止事件: ", sign)
			os.Exit(0)
		}
	}
}

func (cli *Client) Register(aliveList *AliveList) {
	aliveList.Register(cli)
}

func (cli *Client) Close(aliveList *AliveList) {
	aliveList.Destroy(cli)
}

func (cli *Client) SendMessage(messageType int, message string, aliveList *AliveList) error {
	msg := Message{
		ID:      cli.ID,
		Content: message,
		SentAt:  time.Now().Unix(),
		Type:    messageType,
	}
	err := cli.conn.WriteJSON(msg)
	if err != nil {
		log.Println("sendMessageError:%v", err)
		log.Println("message:%v", msg)
		log.Println("cli:%v", cli)
		cli.Close(aliveList)
	}
	return err
}

func (cli *Client) ReceiveMessage(message string, aliveList *AliveList) {
	log.Println("接收到一个消息：", message)
	msg := Message{
		ID:      cli.ID,
		Content: message,
		SentAt:  time.Now().Unix(),
		Type:    1,
	}
	aliveList.ReciveMessage(msg)
}

func NewWebSocket(id string, w http.ResponseWriter, r *http.Request, aliveList *AliveList) (client *Client, err error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client = &Client{
		ID:     id,
		conn:   conn,
		cancel: make(chan int, 1),
	}
	client.Register(aliveList)
	return
}

func ServerStart(aliveList *AliveList) {
	flag.Parse()
	go aliveList.run()
	http.HandleFunc("/api/v1/camera/getHistory", controller.GetHistory)
	http.HandleFunc("/api/v1/camera/ws", aliveList.socketServer)
	log.Println("端口监听：", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func (al *AliveList) socketServer(writer http.ResponseWriter, request *http.Request) {
	if websocket.IsWebSocketUpgrade(request) {
		log.Println("接收到websocket连接")
		id := strconv.Itoa(rand.Intn(65532))
		client, err := NewWebSocket(id, writer, request, al)
		if err == nil {
			defer client.Close(al)
			err = client.SendMessage(1, "WelCome!", al)
			/*加入登陆认证*/
			for {
				_, message, err := client.conn.ReadMessage()
				if websocket.IsCloseError(err, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
					log.Println("连接关闭：", client.ID)
					break
				}
				client.ReceiveMessage(string(message), al)
			}
		} else {
			log.Println("Web Socket 创建失败！")
		}
	} else {
		writer.Write([]byte("请使用Webocket连接"))
	}
}
