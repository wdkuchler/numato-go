package main

// para teste -> http://localhost:8080/?device=com6&command=pulse

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const wakeUp = "\r"
const openIt = "relay on 0\r"
const closeIt = "relay off 0\r"

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	d := r.Form.Get("device")
	c := r.Form.Get("command")

	if len(d) == 0 || len(c) == 0 {
		fmt.Fprintf(w, "Numato - invalid parameters\n")
		return
	}

	fmt.Fprintf(w, "Numato - openning device %s\n", d)

	var file, err = os.Create(d)
	if err != nil {
		fmt.Fprintf(w, "Numato - Error %s\n", err.Error())
		return
	}
	defer file.Close()

	_, err = file.WriteString(wakeUp)
	if err != nil {
		fmt.Fprintf(w, "Numato - Error %s\n", err.Error())
		return
	}
	time.Sleep(50 * time.Millisecond)

	fmt.Fprintf(w, "Numato - writing %s on device %s\n", c, d)

	switch c {
	case "on":
		_, err = file.WriteString(openIt)
		if err != nil {
			fmt.Fprintf(w, "Numato - Error %s\n", err.Error())
			return
		}
		fmt.Fprintf(w, "wrote %s command on device=%s\n", c, d)
	case "off":
		_, err = file.WriteString(closeIt)
		if err != nil {
			fmt.Fprintf(w, "Numato - Error %s\n", err.Error())
			return
		}
		fmt.Fprintf(w, "wrote %s command on device=%s\n", c, d)
	case "pulse":
		_, err = file.WriteString(openIt)
		if err != nil {
			fmt.Fprintf(w, "Numato - Error %s\n", err.Error())
			return
		}
		time.Sleep(1000 * time.Millisecond)
		_, err = file.WriteString(closeIt)
		if err != nil {
			fmt.Fprintf(w, "Numato - Error %s\n", err.Error())
			return
		}
		fmt.Fprintf(w, "wrote %s command on device=%s\n", c, d)
	default:
		fmt.Fprintf(w, "invalid command %s\n", c)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
