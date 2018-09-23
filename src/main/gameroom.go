package main

import (
	"net/http"
    "html/template"
	"fmt"
	"time"
	"github.com/gorilla/websocket"
	"regexp"
	//"strings"
	"log"
	"strconv"
)

type tmessage struct {
	content    []byte
	fromuser   []byte
	touser     []byte
	mtype      int
	createtime string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	//注册连接
	connections map[*connection]bool

	// Inbound messages from the connections.
	//连接中的绑定消息
	broadcast chan *tmessage

	// Register requests from the connections.
	//添加新连接
	register chan *connection

	// Unregister requests from connections.
	//删除连接
	unregister chan *connection
}

var h = hub{
	//广播slice
	broadcast: make(chan *tmessage),
	//注册者slice
	register: make(chan *connection),
	//未注册者sclie
	unregister: make(chan *connection),
	//连接map
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		//注册者有数据，则插入连接map
		case c := <-h.register:
			h.connections[c] = true
		//非注册者有数据，则删除连接map
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		//广播有数据
		case m := <-h.broadcast:
			//递归所有广播连接
			for c := range h.connections {
				var send_flag = false

				//根据广播消息标识记录
				/*
					text2 := string(m.content)
					reg2 := regexp.MustCompile(`^@.*? `)
					s2 := reg2.FindAllString(text2, -1)
				*/
				var send_msg []byte
				if m.mtype == 1 { //系统消息
					//send_msg = []byte(" system: " + string(m.content))
					send_msg = []byte("[_updateremain_]" + string(m.content))
				} else if m.mtype == 2 { //用户消息
					//send_msg = []byte(string(m.fromuser) + " say: " + string(m.content))
					t := time.Now()
					send_msg = []byte(t.Format("15:04:05") + ": " + string(m.content))
				} else {
					send_msg = []byte(string(m.content))
				}
				if string(m.touser) != "BANK" {
					if string(c.username) == string(m.touser) || string(c.username) == string(m.fromuser) {
						send_flag = true
					}
					if send_flag {
						select {
						//发送数据给连接
						case c.send <- send_msg:
						//关闭连接
						default:
							close(c.send)
							delete(h.connections, c)
						}
					}
				} else {
					select {
					//发送数据给连接
					case c.send <- send_msg:
					//关闭连接
					default:
						close(c.send)
						delete(h.connections, c)
					}
				}

			}
		}
	}
}

const (
	//对方写入会话等待时间
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	//对方读取下次消息等待时间
	// Time allowed to read the next pong message from the peer.
	pongWait = 1200 * time.Second

	//对方ping周期
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	//对方最大写入字节数
	// Maximum message size allowed from peer.
	maxMessageSize = 512

	//验证字符串
	authToken = "123456"
)

//服务器配置信息
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection 是websocket的conntion和hub的中间人
// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	//websocket的连接
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	//出站消息缓存通道
	send chan []byte

	//验证状态
	auth bool

	//验证状态
	username []byte
}

func transferMoney(roomId string, fromPlayer string, toPlayer string, cashType string, count int) bool{
	cashInAccount,_ := roomToPlayer[roomId][fromPlayer].Accounts[cashType];
	if(cashInAccount < count){
		return false;
	}
	
	if(count < 0){
		if(toPlayer != "BANK"){
			return false;
		}
	}
	if(toPlayer == "BANK"){
		roomToPlayer[roomId][fromPlayer].Accounts[cashType] = cashInAccount - count;
		return true;
	}
	roomToPlayer[roomId][fromPlayer].Accounts[cashType] = cashInAccount - count;
	otherUserAccount, _ := roomToPlayer[roomId][toPlayer].Accounts[cashType];
	roomToPlayer[roomId][toPlayer].Accounts[cashType] = otherUserAccount + count;
	return true;
}

func getTransferReceipt(fromPlayer string, toPlayer string, cashType string, count int) string{
	res := fromPlayer + "->";
	res += toPlayer + " [" + cashType + "] " + fmt.Sprintf("%d",count);
	return res;
}

func getRemainMoney(roomId string, theUser string, cashType string) string{
	remainmoney, _ := roomToPlayer[roomId][theUser].Accounts[cashType];
	return fmt.Sprintf("%d", remainmoney);
}

func sendCurrentMoney(roomId string, theUser string, cashType string){
	if(theUser == "BANK"){
		return
	}
	message := []byte(getRemainMoney(roomId, theUser,cashType));
	t := time.Now().Unix()
	h.broadcast <- &tmessage{content: message, fromuser: []byte(theUser), touser: []byte(theUser), mtype: 1, createtime: time.Unix(t, 0).String()}
}

//读取connection中的数据导入到hub中，实则发广播消息
//服务器读取的所有客户端的发来的消息
// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump(currentUsername string) {
	fmt.Println("Enter reed") //get request method
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		text := string(message)
		fmt.Println(currentUsername, text) //get request method
		reg := regexp.MustCompile(`[^&]+`)
		s := reg.FindAllString(text, -1)
		transferCount, _ := strconv.Atoi(s[0]);
		currentRoom, _ := playerToRoom[currentUsername];
		if(transferMoney(currentRoom, currentUsername, s[1], "money", transferCount)){
			fmt.Println(roomToPlayer[currentRoom][currentUsername].GetAccountString());
			mtype := 2 //用户消息
			message = []byte(getTransferReceipt(currentUsername, s[1], "money", transferCount));
			//c.username = []byte(currentUsername)
			//c.auth = true
			touser := []byte(s[1])
			t := time.Now().Unix()
			h.broadcast <- &tmessage{content: message, fromuser: c.username, touser: touser, mtype: mtype, createtime: time.Unix(t, 0).String()}
			
			mtype = 1 //message update remain money
			message = []byte(getRemainMoney(currentRoom, currentUsername,"money"));
			//c.username = []byte(currentUsername)
			//c.auth = true
			//touser := []byte(s[1])
			t = time.Now().Unix()
			h.broadcast <- &tmessage{content: message, fromuser: []byte(currentUsername), touser: []byte(currentUsername), mtype: mtype, createtime: time.Unix(t, 0).String()}
			
			if(s[1] != "BANK"){
				mtype = 1 //message update remain money
				message = []byte(getRemainMoney(currentRoom, s[1],"money"));
				//c.username = []byte(currentUsername)
				//c.auth = true
				//touser := []byte(s[1])
				t = time.Now().Unix()
				h.broadcast <- &tmessage{content: message, fromuser: []byte(s[1]), touser: []byte(s[1]), mtype: mtype, createtime: time.Unix(t, 0).String()}
			}
		}
		
		/*
		//默认all
		if len(s) == 2 {
			fromuser := strings.Replace(s[0], "=", "", 1)
			token := strings.Replace(s[1], "=", "", 1)
			if token == authToken {
				c.username = []byte(fromuser)
				c.auth = true
				message = []byte(fromuser + " join")
				mtype = 1 //系统消息
			}
		}

		touser := []byte("all")
		reg2 := regexp.MustCompile(`^@.*? `)
		s2 := reg2.FindAllString(text, -1)
		if len(s2) == 1 {
			s2[0] = strings.Replace(s2[0], "@", "", 1)
			s2[0] = strings.Replace(s2[0], " ", "", 1)
			touser = []byte(s2[0])
		}

		if c.auth == true {
			t := time.Now().Unix()
			h.broadcast <- &tmessage{content: message, fromuser: c.username, touser: touser, mtype: mtype, createtime: time.Unix(t, 0).String()}
		}
		*/
	}
}

//给消息，指定消息类型和荷载
// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

//从hub到connection写数据
//服务器端发送消息给客户端
// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	//定时执行
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

//处理客户端对websocket请求
// serveWs handles websocket requests from the peer.
func resWS(w http.ResponseWriter, r *http.Request) {
	if(!isLogedIn(r)){
		return;
	}
	//设定环境变量
    fmt.Println("Connect built") //get request method
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		fmt.Println("Error exist") //get request method
		return
	}
	nameInSession := getSessionUsername(r);
	//初始化connection
	c := &connection{send: make(chan []byte, 256), ws: ws, auth: false, username: []byte(nameInSession)}
	//加入注册通道，意思是只要连接的人都加入register通道
	h.register <- c
	fmt.Println("No Error") //get request method
	go c.writePump() //服务器端发送消息给客户端
	currentRoom, _ := playerToRoom[nameInSession];
	sendCurrentMoney(currentRoom, nameInSession, "money")
	c.readPump(nameInSession)     //服务器读取的所有客户端的发来的消息
}

type GameRoomData struct{
	PlayerName string;
	TheHost string;
	PlayerList map[string]*PlayerAccount;
}

func resGameRoom(w http.ResponseWriter, r *http.Request) {
	if(!isLogedIn(r)){
		http.Redirect(w, r, "/login", http.StatusFound);
	}
	//fmt.Fprintf(w, r.Form["lan_mode"][0])

    fmt.Println("gameroom method:", r.Method) //get request method
    //session, _ := store.Get(r, "cookie-name")

    if r.Method == "GET" {
    	nameInSession := getSessionUsername(r);
        t := template.Must(template.ParseFiles("gameroom.gtpl"))
        grd := GameRoomData{nameInSession, r.Host, roomToPlayer[playerToRoom[getSessionUsername(r)]]};
        t.Execute(w, grd)
    } else {
    	
    }
}