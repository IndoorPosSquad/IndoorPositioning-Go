# Indoor Positioning [USB Host Program]

This is the USB host part of an UWB based indoor positioning program.
~~~
         SPI           USB
DWM1000 <===> STM32F1 <===> Host Software
~~~

~~~
Host Software:

Device Interface <===> Processing <===> Front End
|----------- Golang ------------|      |- Web -|
~~~

Build:

1. install libusb
2. install golang and gousb/websocket
3. add the root of this repo to GOPATH, then install the main package

