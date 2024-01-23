package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
)

func HandleConnection(session ssh.Session) {
	log.Printf("Received a connection from %s\n", session.RemoteAddr())

	defer session.Close()

	_, _, enabled := session.Pty()

	if !enabled {
		session.Write([]byte("Your terminal must be PTY-compatible in order to use our service.\n"))
		session.CloseWrite()
		return
	}

	session.Write([]byte("\x1B[?1000h"))

	screen := tcell.NewSimulationScreen("")

	if err := screen.Init(); err != nil {
		log.Println(err)

		return
	}

	defer screen.Fini()

	go func() {
		for {
			data := make([]byte, 1)

			if _, err := session.Read(data); err != nil {
				log.Println(err)

				return
			}

			if data[0] == 0x1B {
				if _, err := session.Read(data); err != nil {
					log.Println(err)

					return
				}

				if data[0] == 0x5B {
					if _, err := session.Read(data); err != nil {
						log.Println(err)

						return
					}

					if data[0] == 0x4D {
						if _, err := session.Read(data); err != nil {
							log.Println(err)

							return
						}

						switch data[0] {
						case 0x43:
							{
								// Mouse Move

								coords := make([]byte, 2)

								if _, err := session.Read(coords); err != nil {
									log.Println(err)

									return
								}

								log.Printf("Mouse Move: (%d, %d)\n", coords[0], coords[1])

								break
							}
						case 0x20:
							{
								// Left Mouse Down

								coords := make([]byte, 2)

								if _, err := session.Read(coords); err != nil {
									log.Println(err)

									return
								}

								log.Printf("Left Down: (%d, %d)\n", coords[0], coords[1])

								break
							}
						case 0x21:
							{
								// Middle Mouse Down

								coords := make([]byte, 2)

								if _, err := session.Read(coords); err != nil {
									log.Println(err)

									return
								}

								log.Printf("Middle Down: (%d, %d)\n", coords[0], coords[1])

								break
							}
						case 0x22:
							{
								// Right Mouse Down

								coords := make([]byte, 2)

								if _, err := session.Read(coords); err != nil {
									log.Println(err)

									return
								}

								log.Printf("Right Down: (%d, %d)\n", coords[0], coords[1])

								break
							}
						case 0x23:
							{
								// Mouse Up

								coords := make([]byte, 2)

								if _, err := session.Read(coords); err != nil {
									log.Println(err)

									return
								}

								log.Printf("Up: (%d, %d)\n", coords[1], coords[0])

								session.Write([]byte(fmt.Sprintf("\x1B[%d;%dH", coords[1]-32, coords[0]-32)))

								break
							}
						case 0x40:
							{
								// Mouse Up

								coords := make([]byte, 2)

								if _, err := session.Read(coords); err != nil {
									log.Println(err)

									return
								}

								log.Printf("Drag To: (%d, %d)\n", coords[0], coords[1])

								break
							}
						}
					}
				}
			} else {
				session.Write(data)
			}
		}
	}()

	s := make(chan bool, 1)
	<-s
}
