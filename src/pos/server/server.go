package server

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"os"

	"encoding/json"
	"html/template"
	"net/http"
	"pkg/websocket"
	
	"pos/positioning"
	"pos/device"
)

var (
	JSON          = websocket.JSON           // codec for JSON
	Message       = websocket.Message        // codec for string, []byte
	ActiveClients = make(map[ClientConn]int) // map containing clients

	ps  = [][]float64{{0.0, 0.0}, {500.0, 0.0}}
	rec = [][]float64{{0.0, 0.0}, {0.0, 0.0}}
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

type TxMsg struct {
	Xpos float64 `json:"x"`
	Ypos float64 `json:"y"`
}

func SockServer(ws *websocket.Conn) {
	var err error
	// use []byte if websocket binary type is blob or arraybuffer
	// var clientMessage []byte
	// cleanup on server side

	defer func() {
		if err = ws.Close(); err != nil {
			log.Println("Websocket could not be closed",
				err.Error())
		}
	}()

	client := ws.Request().RemoteAddr
	sockCli := ClientConn{ws, client}
	ActiveClients[sockCli] = 0

	for {
		time.Sleep(50 * time.Millisecond)

		p1_flt, p2_flt := device.GetDistanceUSB()
		positioning.Solve2d(rec, ps, p1_flt, p2_flt)

		txMsgStruct := TxMsg{rec[0][0], rec[0][1]}
		byteTxMsgJSON, err := json.Marshal(txMsgStruct)
		// []byte -> string
		if err != nil {
			log.Println("ERROR parsing: ", err)
			log.Println("         data: ", p1_flt, p2_flt)
			continue;
		}

		txMsg := string(byteTxMsgJSON)
		log.Println(txMsg)

		for cs, _ := range ActiveClients {
			if err = Message.Send(
				cs.websocket,
				txMsg); err != nil {
				// we could not send the message to a peer
				log.Println(
					"Could not send message to ",
					cs.clientIP,
					err.Error())
			}
		}
	}
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(os.Getwd())
	fmt.Fprintf(w, "<h1>%s</h1>", "test")
	param := r.URL.Path
	log.Println(param)
	t, _ := template.ParseFiles("./src/pos/server/index.html")
	p := &Page{Msg: strconv.Itoa(1234), Xpos: 0, Ypos: 0}
	t.Execute(w, p)
}

func Init(port int) {
	http.Handle("/js/",
		http.StripPrefix("/js/",
			http.FileServer(http.Dir("./src/pos/server/js"))))
	http.Handle("/css/",
		http.StripPrefix("/css/",
			http.FileServer(http.Dir("./src/pos/server/css"))))
	http.Handle("/sock", websocket.Handler(SockServer))
	http.HandleFunc("/", RequestHandler)
	log.Println("localhost:" + strconv.Itoa(port))
	http.ListenAndServe("localhost:" + strconv.Itoa(port), nil)
	log.Println("end of socket server init")
	//for {}
}
