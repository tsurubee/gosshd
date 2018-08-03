package main

import (
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"log"
	"bytes"
	"errors"
	"io"
	"strings"
)

// TODO The linkage information between the user name and the host name is acquired from the database.
func findUpstreamByUsername(username string) (string, error) {
	if username == "tsurubee" {
		return "host-tsurubee", nil
	} else if username == "bob" {
		return "host-bob", nil
	}
	return "", errors.New(username + "'s host is not found!")
}

func main() {
	ssh.Handle(func(sess ssh.Session) {
		username := sess.User()
		upstream, err := findUpstreamByUsername(sess.User())
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("Connecting for %s by %s\n", upstream, username)

		config := &gossh.ClientConfig{
			User: username,
			Auth: []gossh.AuthMethod{
				gossh.Password("test"),
			},
			HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		}

		clientConn, err := gossh.Dial("tcp", upstream + ":22", config)
		if err != nil {
			panic(err)
		}
		defer clientConn.Close()

		usess, err := clientConn.NewSession()
		if err != nil {
			panic(err)
		}
		defer usess.Close()

		var b bytes.Buffer
		usess.Stdout = &b
		if err := usess.Run("hostname"); err != nil {
			log.Fatal("Failed to run: " + err.Error())
		}
		r := strings.NewReader(b.String())
		io.Copy(sess, r)
	})

	log.Println("Starting ssh server on port 2222")
	ssh.ListenAndServe(":2222", nil)
}