package services

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/revel/revel"
	"golang.org/x/crypto/ssh"
)

//Payload TODO remove from export
type Payload struct {
	Str string
}

func initGitSshd() {

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	keyFile, _ := revel.Config.String("ssh.key")
	privateBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatal("Failed to load private key " + keyFile)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	config.AddHostKey(private)

	sshPort, _ := revel.Config.String("ssh.port")

	listener, err := net.Listen("tcp", "0.0.0.0:"+sshPort)
	if err != nil {
		log.Fatal("Failed to listen on 2022")
	}

	log.Print("Listening on 2022...")
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			log.Printf("Failed to handshake (%s)", err)
			continue
		}

		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		go handleRequests(reqs)
		go handleGitCommands(chans)

	}
}

func handleRequests(reqs <-chan *ssh.Request) {
	for req := range reqs {
		log.Printf("recieved out-of-band request: %+v", req)
	}
}

func handleGitCommands(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		channel, requets, _ := newChannel.Accept()
		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Printf("declining %s request...\n", req.Type)
				switch req.Type {
				case "exec":

					payload := Payload{}
					ssh.Unmarshal(req.Payload, &payload)

					parts := strings.Fields(payload.Str)
					//fmt.Printf("run command: %s \n", payload.Str)
					command := string(parts[0])
					cmd := exec.Command("")

					repoRoot, _ := revel.Config.String("repo.root")

					repo := repoRoot + strings.Replace(string(parts[1]), "'", "", -1)[1:]

					git, _ := revel.Config.String("git.cmd")

					switch command {
					case "git-upload-pack":
						cmd = exec.Command(git, "upload-pack", repo)
					case "git-receive-pack":
						cmd = exec.Command(git, "receive-pack", repo)
					case "git-upload-archive":
						cmd = exec.Command(git, "upload-archive", repo)
					}
					out, _ := cmd.StdoutPipe()
					input, _ := cmd.StdinPipe()
					cmd.Stderr = os.Stderr
					go io.Copy(channel, out)
					go io.Copy(input, channel)
					cmd.Start()
					cmd.Wait()
					channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
					channel.Close()
				case "env":
					payload := Payload{}
					ssh.Unmarshal(req.Payload, &payload)
					fmt.Printf("run command: %s \n", payload.Str)
				default:
					channel.Close()
				}
			}
		}(requets)
	}
}
