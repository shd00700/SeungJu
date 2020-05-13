package main

import (
        "fmt"
        "github.com/stianeikeland/go-rpio"
        "github.com/tarm/serial"
        "time"
        "os"
        "log"
        "flag"
        "net"
        "strconv"
)
const (
        gwState = rpio.Pin(17)
        lan9514 = rpio.Pin(20)
        lan9512 = rpio.Pin(21)
        rs485A = rpio.Pin(22)
        rs485B = rpio.Pin(23)
        rs485C = rpio.Pin(24)
        wireless = rpio.Pin(25)
)

func pinset(){
        if err := rpio.Open(); err!= nil{
                fmt.Println(err)
                os.Exit(1)
                }
        gwState.Output()
        lan9514.Output()
        lan9512.Output()
        rs485A.Output()
        rs485B.Output()
        rs485C.Output()
        wireless.Output()
}
func pintoggle(){
        gwState.Toggle()
        lan9514.Toggle()
        lan9512.Toggle()
        rs485A.Toggle()
        rs485B.Toggle()
        rs485C.Toggle()
}
func LedToggle(){

        for{    //GPIO LED Toggle
                pintoggle()
                time.Sleep(time.Second)
        }
}

func Rs485(){
        rs485a := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600, StopBits: 1, Parity: 'N'}
        rs485b := &serial.Config{Name: "/dev/ttyUSB1", Baud: 9600, StopBits: 1, Parity: 'N'}
        rs485c := &serial.Config{Name: "/dev/ttyUSB3", Baud: 9600, StopBits: 1, Parity: 'N'}

        a, err := serial.OpenPort(rs485a)
        if err != nil{
                log.Fatal(err)
        }
        b, err := serial.OpenPort(rs485b)
        if err != nil{
                log.Fatal(err)
        }
        c, err := serial.OpenPort(rs485c)
        if err != nil{
                log.Fatal(err)
        }
        fmt.Println("Serial port Open")
        for{
                n, err := a.Write([]byte("test"))
                if err != nil {
                        log.Fatal(n)
                }
                m, err := b.Write([]byte("test"))
                if err != nil {
                        log.Fatal(m)
                }
                l, err := c.Write([]byte("test"))
                if err != nil {
                        log.Fatal(l)
                }
        }
}

func handleRequest(conn net.Conn) {
        log.Println("Accepted new connection.")
        for{
                buf := make([]byte, 1024)
                size, err := conn.Read(buf)
                if err != nil {
                        return
                }
                data := buf[:size]
                log.Println("Read new data from connection", data)
                conn.Write(data)

        }
}
func main() {
        //gpio pin setting
        pinset()
        //serial Rs-485 setting
        /*a, err := serial.OpenPort(rs485a)
        if err != nil{
                log.Fatal(err)
        }
        b, err := serial.OpenPort(rs485b)
        if err != nil{
                log.Fatal(err)
        }
        c, err := serial.OpenPort(rs485c)
        if err != nil{
                log.Fatal(err)
        }*/
        //Ethernet TCP setting
        port := flag.Int("port", 3322, "Port to accept connections on.")
        flag.Parse()

        l, err := net.Listen("tcp",":"+strconv.Itoa(*port))
        if err != nil {
                log.Panicln(err)
        }
        log.Println("Listening to connections at on port", strconv.Itoa(*port))
        fmt.Println(l)
        //Board Test
                //Ethernet TCP
        for{
                fmt.Println("ee")
                conn, err := l.Accept()
                if err != nil {
                        log.Panicln(err)
                }
                //Led start
                go LedToggle()
                //Rs485 start
                go Rs485()
                go handleRequest(conn)
        }
}
