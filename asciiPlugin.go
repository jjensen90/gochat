package main

import "log"

type AsciiPlugin struct {
	Channel chan []byte
}

func (p *AsciiPlugin) Update() {
	for {
		select {
		case message := <-p.Channel:
			switch {
			case string(message) ==  "/ascii":
				log.Printf("%s", message)
			default:
				log.Print("don't care")
			}
		}
	}
}

func (p *AsciiPlugin) GetChannel() chan []byte {
	return p.Channel
}