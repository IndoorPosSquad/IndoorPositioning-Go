package socket-server

var (
	JSON          = websocket.JSON           // codec for JSON
	Message       = websocket.Message        // codec for string, []byte
	ActiveClients = make(map[ClientConn]int) // map containing clients
)

type Page struct {
	Msg  string
	Xpos float64
	Ypos float64
}

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}
