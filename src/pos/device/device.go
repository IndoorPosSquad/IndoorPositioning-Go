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
	ep_bulk_read usb.Endpoint
	ctx          *usb.Context
	devs         []*usb.Device
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

func GetmsgUSB() []byte {
	//InitUSB()
	//defer CloseUSB()

	buf := make([]byte, 64)
	for i := 0; i < 64; i++ {
		buf[i] = 0
	}

	log.Println("Ready to read from Device")
	ep_bulk_read.Read(buf)
	n := bytes.IndexByte(buf, byte(0))
	buf = buf[:n]

	fmt.Println(string(buf))
	//fmt.Printf("%c\n", buf)
	return buf
}
