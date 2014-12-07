package server

import (
	"fmt"
	"time"
	"strconv"
	"os"

	"html/template"
	"net/http"
	"websocket"
)

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

func SockServer(ws *websocket.Conn) {
	var err error
	var clientMessage string
	// use []byte if websocket binary type is blob or arraybuffer
	// var clientMessage []byte
	// cleanup on server side
	defer func() {
		if err = ws.Close(); err != nil {
			fmt.Println("Websocket could not be closed", err.Error())
		}
	}()
	client := ws.Request().RemoteAddr
	fmt.Println("Client connected:", client)
	sockCli := ClientConn{ws, client}
	ActiveClients[sockCli] = 0
	fmt.Println(
		"Number of clients connected ...",
		len(ActiveClients))

	for {
		time.Sleep(1000 * time.Millisecond)

		//msg := string(usb_getmsg())
		//p1_str := strings.Split(msg, " ")[0]
		//p2_str := strings.Split(msg, " ")[1]
		//fmt.Println("str:")
		//fmt.Println(p1_str, p2_str)
//
		//p1_flt, _ := strconv.ParseFloat(p1_str, 64)
		//p2_flt, _ := strconv.ParseFloat(p2_str, 64)
//
		//p1_flt += 60
		//p2_flt += 00
//
		//fmt.Println("flt:")
		//fmt.Println(p1_flt, p2_flt)
//
		//fmt.Println("rec")
		//solve_2d(rec, ps, p1_flt, p2_flt)
//
		//fmt.Println(rec)
//
		//clientMessage =
			//strconv.FormatFloat(rec[0][0], 'g', 6, 64) +
				//"," +
				//strconv.FormatFloat(rec[0][1], 'g', 6, 64)
		clientMessage = "test"

		for cs, _ := range ActiveClients {
			if err = Message.Send(
				cs.websocket,
				clientMessage); err != nil {
				// we could not send the message to a peer
				fmt.Println(
					"Could not send message to ",
					cs.clientIP,
					err.Error())
			}
		}
	}
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(os.Getwd())
	fmt.Fprintf(w, "<h1>%s</h1>", "test")
	param := r.URL.Path
	fmt.Println(param)
	t, _ := template.ParseFiles("./src/pos/server/index.html")
	p := &Page{Msg: strconv.Itoa(1234), Xpos: 0, Ypos: 0}
	t.Execute(w, p)
}

func Init(port int) {
	fmt.Println(port)
	http.Handle("/js/",
		http.StripPrefix("/js/",
			http.FileServer(http.Dir("./src/pos/server/js"))))
	http.Handle("/css/",
		http.StripPrefix("/css/",
			http.FileServer(http.Dir("./src/pos/server/css"))))
	http.Handle("/sock", websocket.Handler(SockServer))
	http.HandleFunc("/", RequestHandler)
	fmt.Println("localhost:" + strconv.Itoa(port))
	http.ListenAndServe("localhost:" + strconv.Itoa(port), nil)
	fmt.Println("end of socket server init")
	//for {}
}
