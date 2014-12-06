// Copyright 2013 Google Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// rawread attempts to read from the specified USB device.
package main

// #cgo LDFLAGS: -lusb-1.0
// #include <libusb-1.0/libusb.h>
import "C"

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	. "math"
	"strconv"
	"strings"
	"time"

	"html/template"
	"net/http"
	"websocket"

	"gousb/usb"

//"usbid"
)

var (
	device   = flag.String("device", "0483:5741", "Device to which to connect")
	config   = flag.Int("config", 1, "Endpoint to which to connect")
	iface    = flag.Int("interface", 0, "Endpoint to which to connect")
	setup    = flag.Int("setup", 0, "Endpoint to which to connect")
	endpoint = flag.Int("endpoint", 1, "Endpoint to which to connect")
	debug    = flag.Int("debug", 3, "Debug level for libusb")
	ep       usb.Endpoint
	ctx      *usb.Context
	devs     []usb.Device

	JSON          = websocket.JSON           // codec for JSON
	Message       = websocket.Message        // codec for string, []byte
	ActiveClients = make(map[ClientConn]int) // map containing clients

	x_tmp = float64(0)
	y_tmp = float64(0)

	prange1 = float64(0)
	prange2 = float64(0)

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

func usb_getmsg() []byte {
	flag.Parse()

	// Only one context should be needed for an application.  It should always be closed.
	ctx = usb.NewContext()

	ctx.Debug(*debug)
	defer ctx.Close()

	log.Printf("Scanning for device %q...", *device)

	// ListDevices is used to find the devices to open.
	devs, err := ctx.ListDevices(func(desc *usb.Descriptor) bool {
		if fmt.Sprintf("%s:%s", desc.Vendor, desc.Product) != *device {
			return false
		}
		return true
	})

	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	if err != nil {
		log.Fatalf("list: %s", err)
	}

	if len(devs) == 0 {
		log.Fatalf("no devices found")
	}

	dev := devs[0]

	//log.Printf("Connecting to endpoint...")
	//log.Printf("- %#v", dev.Descriptor)
	ep, err := dev.OpenEndpoint(
		uint8(*config),
		uint8(*iface),
		uint8(*setup),
		uint8(*endpoint)|uint8(usb.ENDPOINT_DIR_IN))

	if err != nil {
		log.Fatalf("open: %s", err)
	}
	_ = ep

	buf := make([]byte, 64)
	for i := 0; i < 64; i++ {
		buf[i] = 0
	}

	ep.Read(buf)
	n := bytes.IndexByte(buf, byte(0))
	buf = buf[:n]

	fmt.Println(string(buf))
	//fmt.Printf("%c\n", buf)
	return buf
}

func sgn(x float64) float64 {
	if x >= 0.0 {
		return 1.0
	} else {
		return -1.0
	}
}

func fabs(x float64) float64 {
	if x >= 0.0 {
		return x
	} else {
		return -x
	}
}

func solve_2d(
	reciever [][]float64,
	pseudolites [][]float64,
	pranges1 float64,
	pranges2 float64) {

	var origin [2]float64
	var len float64
	var tan_theta float64
	var sin_theta float64
	var cos_theta float64
	var d1 float64
	var h1 float64
	var invrotation [2][2]float64
	var pranges [2]float64

	fmt.Printf("\nPseudolites\n%f %f %f %f\n",
		pseudolites[0][0],
		pseudolites[0][1],
		pseudolites[1][0],
		pseudolites[1][1])
	fmt.Printf("\npr1 %f pr2 %f\n", pranges1, pranges2)

	pranges[0] = pranges1
	pranges[1] = pranges2

	origin[0] = pseudolites[0][0]
	origin[1] = pseudolites[0][1]

	pseudolites[0][0] = 0
	pseudolites[0][1] = 0
	pseudolites[1][0] = pseudolites[1][0] - origin[0]
	pseudolites[1][1] = pseudolites[1][1] - origin[1]

	len = Sqrt(Pow(pseudolites[1][0], 2) + Pow(pseudolites[1][1], 2))

	tan_theta = pseudolites[1][1] / pseudolites[1][0]
	cos_theta = sgn(pseudolites[1][0]) / Sqrt(Pow(tan_theta, 2)+1)
	sin_theta = sgn(pseudolites[1][1]) * fabs(tan_theta) / Sqrt(Pow(tan_theta, 2)+1)

	invrotation[0][0] = cos_theta
	invrotation[0][1] = -sin_theta
	invrotation[1][0] = sin_theta
	invrotation[1][1] = cos_theta

	d1 = ((Pow(pranges[0], 2)-Pow(pranges[1], 2))/len + len) / 2

	h1 = Sqrt(Pow(pranges[0], 2) - Pow(d1, 2))

	reciever[0][0] = d1
	reciever[0][1] = h1
	reciever[1][0] = d1
	reciever[1][1] = -h1

	reciever[0][0] = invrotation[0][0]*d1 + invrotation[0][1]*h1
	reciever[0][1] = invrotation[1][0]*d1 + invrotation[1][1]*h1
	reciever[0][0] += origin[0]
	reciever[0][1] += origin[1]

	reciever[1][0] = invrotation[0][0]*d1 + invrotation[0][1]*-h1
	reciever[1][1] = invrotation[1][0]*d1 + invrotation[1][1]*-h1
	reciever[1][0] += origin[0]
	reciever[1][1] += origin[1]
}

func SockServer(ws *websocket.Conn) {
	var err error
	var clientMessage string
	// use []byte if websocket binary type is blob or arraybuffer
	// var clientMessage []byte
	// cleanup on server side
	defer func() {
		if err = ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()
	client := ws.Request().RemoteAddr
	log.Println("Client connected:", client)
	sockCli := ClientConn{ws, client}
	ActiveClients[sockCli] = 0
	log.Println(
		"Number of clients connected ...",
		len(ActiveClients))

	for {
		time.Sleep(1000 * time.Millisecond)

		msg := string(usb_getmsg())
		p1_str := strings.Split(msg, " ")[0]
		p2_str := strings.Split(msg, " ")[1]
		fmt.Println("str:")
		fmt.Println(p1_str, p2_str)

		p1_flt, _ := strconv.ParseFloat(p1_str, 64)
		p2_flt, _ := strconv.ParseFloat(p2_str, 64)

		p1_flt += 60
		p2_flt += 00

		fmt.Println("flt:")
		fmt.Println(p1_flt, p2_flt)

		fmt.Println("rec")
		solve_2d(rec, ps, p1_flt, p2_flt)

		fmt.Println(rec)

		clientMessage =
			strconv.FormatFloat(rec[0][0], 'g', 6, 64) +
				"," +
				strconv.FormatFloat(rec[0][1], 'g', 6, 64)

		for cs, _ := range ActiveClients {
			if err = Message.Send(
				cs.websocket,
				clientMessage); err != nil {
				// we could not send the message to a peer
				log.Println(
					"Could not send message to ",
					cs.clientIP,
					err.Error())
			}
		}
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calllll")
	param := r.URL.Path
	fmt.Println(param)
	t, _ := template.ParseFiles("index.html")
	p := &Page{Msg: string(usb_getmsg()), Xpos: 0, Ypos: 0}
	t.Execute(w, p)
}

func main() {
	for i := 0; i < 5; i++ {
		fmt.Println(usb_getmsg())
	}
	http.Handle("/js/",
		http.StripPrefix("/js/",
			http.FileServer(http.Dir("./js"))))
	http.Handle("/css/",
		http.StripPrefix("/css/",
			http.FileServer(http.Dir("./css"))))
	http.Handle("/sock", websocket.Handler(SockServer))
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe("localhost:2000", nil)
}
