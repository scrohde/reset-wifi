/*
reset-wifi continuously checks if the connection to the internet is down and
resets the network if necessary.

This works around an issue I'm having with my laptop.

	Usage of reset-wifi:
		-dst string
			destination to ping (default "google.com")
		-monitor duration
			(duration) delayed start to monitoring (default 1m0s)
		-ping duration
			(duration) interval between pings (default 2s)
		-pwd
			ask for root password

Duration flags are a sequence of positive decimal numbers, each with optional
fraction and a unit suffix,such as "300ms", "1.5h" or "2h45m". Valid time
units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
*/
package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var pingDuration, monitorDuration time.Duration

func main() {
	dst := flag.String("dst", "google.com", "destination to ping")
	askForPwd := flag.Bool("pwd", false, "ask for root password")
	flag.DurationVar(&pingDuration, "ping", 2*time.Second, "(duration) interval between pings")
	flag.DurationVar(&monitorDuration, "monitor", time.Minute, "(duration) delayed start to monitoring")
	flag.Parse()

	if pingDuration < 0 {
		log.Println("--ping must be positive")
	}
	if monitorDuration < 0 {
		log.Println("--monitor must be positive")
	}

	password := ""
	if *askForPwd {
		password = getPassword("password: ")
	}
	monitorNetwork(password, *dst)
}

func monitorNetwork(password string, dst string) {
	for {
		time.Sleep(monitorDuration)
		waitForPingToFail(dst)
		err := restartNetwork(password)
		if err != nil {
			log.Printf("failed to restart network service: %v\n", err)
			os.Exit(1)
		}
		log.Println("network service restarted successfully.")
	}
}

func waitForPingToFail(dst string) {
	for {
		if !ping(dst, true) {
			return
		}
		time.Sleep(pingDuration)
	}
}

// restartNetwork calls `systemctl restart NetworkManager.service`
func restartNetwork(password string) error {
	var restartCmd *exec.Cmd
	if password == "" {
		restartCmd = exec.Command("systemctl", "restart", "NetworkManager.service")
		restartCmd.Stdin = os.Stdin
	} else {
		restartCmd = exec.Command("sudo", "-S", "--", "systemctl", "restart", "NetworkManager.service")
		restartCmd.Stdin = strings.NewReader(password)
	}
	restartCmd.Stderr = LogWriter{}
	return restartCmd.Run()
}

// ping returns true if there is internet access of some sort.
// Calls `nc -z` which is lighter weight than actually calling `ping`.
func ping(dst string, echo bool) bool {
	pingCmd := exec.Command("nc", "-zw3", dst, "443")
	pingCmd.Stdin = os.Stdin
	if echo {
		pingCmd.Stdout = LogWriter{}
		pingCmd.Stderr = LogWriter{}
	}
	return pingCmd.Run() == nil
}

type LogWriter struct{}

func (s LogWriter) Write(msg []byte) (int, error) {
	log.Println(string(msg))
	return len(msg), nil
}
