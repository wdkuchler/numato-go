package main

// para teste -> http://localhost:8080/?device=com6&command=pulse
// para teste -> http://localhost:8080/?device=/dev/ttyACM0&command=pulse

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

const wakeUp = "\r"
const openIt = "relay on 0\r"
const closeIt = "relay off 0\r"

func doCommand(w http.ResponseWriter, device string, command string) {

	fmt.Fprintf(w, "Numato - openning device %s\n", device)

	var tty, err = os.OpenFile(device, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		fmt.Fprintf(w, "Numato - Error %s\n", err.Error())
		return
	}
	defer tty.Close()

	fmt.Fprintf(w, "Numato - writing %s on device %s\n", command, device)

	_, err = tty.WriteString(command)
	if err != nil {
		fmt.Fprintf(w, "Numato - Error %s\n", err.Error())
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	d := r.Form.Get("device")
	c := r.Form.Get("command")

	if len(d) == 0 || len(c) == 0 {
		fmt.Fprintf(w, "Numato - invalid parameters\n")
		return
	}
	doCommand(w, d, wakeUp)

	time.Sleep(50 * time.Millisecond)

	switch c {
	case "on":
		doCommand(w, d, openIt)
	case "off":
		doCommand(w, d, closeIt)
	case "pulse":
		doCommand(w, d, openIt)
		time.Sleep(1000 * time.Millisecond)
		doCommand(w, d, closeIt)
	default:
		fmt.Fprintf(w, "invalid command %s\n", c)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
