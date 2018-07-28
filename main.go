package main

import (
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"io/ioutil"
	"fmt"
	"os/exec"
	"github.com/kr/pty"
	"sync"
	"io"
)

func main() {
	serverConfig := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	privateKeyBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key (./id_rsa)")
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	serverConfig.AddHostKey(privateKey)

	listener, err := net.Listen("tcp", "0.0.0.0:2222")
	if err != nil {
		log.Fatalf("Failed to listen on 2222 (%s)", err)
	}
	log.Print("Listening on 2222...")

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept on 2222 (%s)", err)
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, serverConfig)
		if err != nil {
			log.Fatalf("Failed to handshake (%s)", err)
		}
		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

		go ssh.DiscardRequests(reqs)
		go handleChannels(chans)
	}
}

func handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go handleChannel(newChannel)
	}
}

func handleChannel(newChannel ssh.NewChannel) {
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("Unknown channel type: %s", t))
		return
	}

	connection, _, err := newChannel.Accept()
	if err != nil {
		log.Fatalf("Could not accept channel (%s)", err)
		return
	}

	bash := exec.Command("bash")

	close := func() {
		connection.Close()
		_, err := bash.Process.Wait()
		if err != nil {
			log.Printf("Failed to exit bash (%s)", err)
		}
		log.Printf("Session closed")
	}

	bashf, err := pty.Start(bash)
	if err != nil {
		log.Printf("Could not start pty (%s)", err)
		close()
		return
	}

	var once sync.Once
	go func() {
		io.Copy(connection, bashf)
		once.Do(close)
	}()
	go func() {
		io.Copy(bashf, connection)
		once.Do(close)
	}()
}