package device

import (
	"flag"
	"log"
	"github.com/kylelemons/gousb/usb"
)

var (
	device       = flag.String("device", "0483:5741", "Device to which to connect")
	config       = flag.Int("config", 1, "Endpoint to which to connect")
	iface        = flag.Int("interface", 0, "Endpoint to which to connect")
	setup        = flag.Int("setup", 0, "Endpoint to which to connect")
	endpoint_bulk_read     = flag.Int("endpoint_bulk_read", 1, "Endpoint to which to connect")
	endpoint_int_read     = flag.Int("endpoint_int_read", 2, "Endpoint to which to connect")
	endpoint_bulk_write     = flag.Int("endpoint_bulk_write", 3, "Endpoint to which to connect")
	debug        = flag.Int("debug", 3, "Debug level for libusb")

	ctx          *usb.Context
	devs         []*usb.Device
	ep_bulk_read usb.Endpoint
	ep_bulk_write usb.Endpoint
	ep_int_read usb.Endpoint
)

type DWM1k interface {
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
	ep_bulk_read usb.Endpoint
	ep_bulk_write usb.Endpoint
	ep_int_read usb.Endpoint
}

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
}

func SendCommnadUSB(command string) {
	// TODO add a command hash/table
}
