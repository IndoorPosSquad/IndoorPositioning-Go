package device

// #cgo LDFLAGS: -lusb-1.0
// #include <libusb-1.0/libusb.h>
import "C"

import (
	"bytes"
	"flag"
	"fmt"
	"log"

	"gousb/usb"
)

var (
	device       = flag.String("device", "0483:5741", "Device to which to connect")
	config       = flag.Int("config", 1, "Endpoint to which to connect")
	iface        = flag.Int("interface", 0, "Endpoint to which to connect")
	setup        = flag.Int("setup", 0, "Endpoint to which to connect")
	endpoint     = flag.Int("endpoint", 1, "Endpoint to which to connect")
	debug        = flag.Int("debug", 3, "Debug level for libusb")

	ctx          *usb.Context
	devs         []*usb.Device
	ep_bulk_read usb.Endpoint
)

func InitUSB() {
	var err error
	flag.Parse()

	// Only one context should be needed for an application.  It should always be closed.
	ctx = usb.NewContext()

	ctx.Debug(*debug)

	log.Printf("Scanning for device %q...", *device)

	// ListDevices is used to find the devices to open.
	devs, err = ctx.ListDevices(func(desc *usb.Descriptor) bool {
		if fmt.Sprintf("%s:%s", desc.Vendor, desc.Product) != *device {
			return false
		}
		return true
	})

	if err != nil {
		log.Fatalf("list: %s", err)
	}

	if len(devs) == 0 {
		log.Fatalf("no devices found")
	}

	dev := devs[0]

	//log.Printf("Connecting to endpoint...")
	//log.Printf("- %#v", dev.Descriptor)
	ep_bulk_read, err = dev.OpenEndpoint(
		uint8(*config),
		uint8(*iface),
		uint8(*setup),
		uint8(*endpoint)|uint8(usb.ENDPOINT_DIR_IN))

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

	d1_str := strings.Split(msg, " ")[0]
	d2_str := strings.Split(msg, " ")[1]

	d1_flt, _ := strconv.ParseFloat(d1_str, 64)
	d2_flt, _ := strconv.ParseFloat(d2_str, 64)

	fmt.Println(string(buf))
	//fmt.Printf("%c\n", buf)
	return d1_flt, d2_flt
}
