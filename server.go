package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
)

type Flusher struct {
	w *bufio.Writer
}

func NewFlusher(w io.Writer) *Flusher {
	return &Flusher{
		w: bufio.NewWriter(w),
	}
}

func (f *Flusher) Write(p []byte) (n int, err error) {
	n, err = f.w.Write(p)
	if err != nil {
		return
	}
	err = f.w.Flush()
	if err != nil {
		return
	}
	return
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port cmd", os.Args[0])
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
	defer tcpListener.Close()

	fmt.Printf("Listening TCP connection on %v\n", tcpAddr)

	for {
		tcpConn, err := tcpListener.Accept()
		if err != nil {
			continue
		}
		go handleClient(tcpConn)
	}
}

func handleClient(tcpConn net.Conn) {
	defer tcpConn.Close()

	command := os.Args[2:]
	cmd := exec.Command(command[0], command[1:]...)

	// stdin, err := cmd.StdinPipe()
	// if err != nil {
	//     fmt.Fprintf(os.Stderr, "Error: %v", err)
	//     return
	// }
	// defer stdin.Close()

	cmd.Stdin = tcpConn

	cmd.Stdout = NewFlusher(tcpConn)
	cmd.Stderr = NewFlusher(tcpConn)

	fmt.Println(command)

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return
	}
}
