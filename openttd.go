package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	cmd   *exec.Cmd
	stdin io.WriteCloser
)

func startOpenTTD() {
	log.Println("*** Starting openttd server")
	cmd = exec.Command(config.OpenTTD, append([]string{"-D"}, config.Args...)...)
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error obtaining stdout: %s", err.Error())
	}

	stdin, err = cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Error obtaining stdin: %s", err.Error())
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin.Write([]byte(scanner.Text() + "\n"))
		}
	}()

	go func() {
		scanner := bufio.NewScanner(bufio.NewReader(stdout))
		for scanner.Scan() {
			line := strings.Replace(scanner.Text(), "‎", "", -1)
			log.Println(line)
			if strings.HasPrefix(line, "*** ") {
				client.SendMsg(ctx, channelID, line)
			} else if !strings.Contains(line, "[Discord]") {
				if strings.HasPrefix(line, config.ChatPrefix) {
					splited := strings.Split(line, ":")
					user := strings.Replace(splited[0], config.ChatPrefix, "", 1)
					if user == "" {
						user = "Server"
					}
					msg := fmt.Sprintf("**%s**:%s", user, strings.Replace(strings.Join(splited[1:], ":"), "@", "(@)", -1))
					client.SendMsg(ctx, channelID, msg)
				} else if sendStats || (sendStats2 && strings.Contains(line, "Company Name")) {
					statsChannel <- line
				} else if sendClients && strings.Contains(line, "Client") {
					clientsChannel <- line
				} else if cmdOut != 0 {
					cmdOut--
					cmdOutChannel <- line
				}
			}
		}
	}()
	if err := cmd.Start(); nil != err {
		log.Fatalf("Error: %s, %s", cmd.Path, err.Error())
	}

	go func() {
		cmd.Wait()
		log.Println("*** Server stopped, restarting...")
		startOpenTTD()
	}()
}
