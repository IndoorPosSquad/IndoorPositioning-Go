package device

// #cgo LDFLAGS: -lusb-1.0
// #include <libusb-1.0/libusb.h>
import "C"

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"
	"strconv"

	"github.com/JohnFarmer/gousb/usb"
)

var (
	device       = flag.String("device", "0483:5741", "Device to which to connect")
	config       = flag.Int("config", 1, "Endpoint to which to connect")
	iface        = flag.Int("interface", 0, "Endpoint to which to connect")
	setup        = flag.Int("setup", 0, "Endpoint to which to connect")
	endpoint_bulk_read     = flag.Int("endpoint_bulk_read", 1, "Endpoint to which to connect")
	endpoint_bulk_write     = flag.Int("endpoint_bulk_write", 1, "Endpoint to which to connect")
	debug        = flag.Int("debug", 3, "Debug level for libusb")

	ctx          *usb.Context
	devs         []*usb.Device
	ep_bulk_read usb.Endpoint
	ep_bulk_write usb.Endpoint
)

func InitUSB() {
	var err error
	flag.Parse()

	// Only one context should be needed for an application.  It should always be closed.
	ctx = usb.NewContext()

	ctx.Debug(*debug)

	log.Printf("Scanning for device %q...", *device)

	dev, _ := ctx.GetDeviceWithVidPid(*device)

	// Open up two ep for read and write
	ep_bulk_read, err = dev.OpenEndpoint(
		uint8(*config),
		uint8(*iface),
		uint8(*setup),
		uint8(*endpoint_bulk_read)|uint8(usb.ENDPOINT_DIR_IN))

	ep_bulk_write, err = dev.OpenEndpoint(
		uint8(*config),
		uint8(*iface),
		uint8(*setup),
		uint8(*endpoint_bulk_write)|uint8(usb.ENDPOINT_DIR_OUT))
	_ = ep_bulk_write

	if err != nil {
		log.Fatalf("open: %s", err)
	}

	log.Printf("Init Done\n")
}

func CloseUSB() {
	defer ctx.Close()

	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()
}

func GetDistanceUSB() (float64, float64) {
	// read distance from device
	buf := make([]byte, 64)
	ep_bulk_read.Read(buf)

	// get the message out of buffer
	n := bytes.IndexByte(buf, byte(0))
	buf = buf[:n]
	
	distances := string(buf)

	d1_str := strings.Split(distances, " ")[0]
	d2_str := strings.Split(distances, " ")[1]

	d1_flt, _ := strconv.ParseFloat(d1_str, 64)
	d2_flt, _ := strconv.ParseFloat(d2_str, 64)

	fmt.Println(string(buf))
	//fmt.Printf("%c\n", buf)
	return d1_flt, d2_flt
}

func SendCommnadUSB(command string) {
	// TODO add a command hash/table
}
