package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

const (
	port            = ":6666"
	framesPerBuffer = 400
)

type Data struct {
	Angle int `json:"angle"`
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		// block until func receive a signal
		sig := <-sigCh
		fmt.Println("got signal", sig)
		cancel()
	}()

	angleCh := make(chan int)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		_, err := fmt.Fprintln(os.Stderr, "udp server is running on", port)
		if err != nil {
			return fmt.Errorf("print log failed -> %w", err)
		}
		conn, err := net.ListenPacket("udp", port)
		if err != nil {
			return fmt.Errorf("listen packet failed -> %w", err)
		}
		defer conn.Close()
		buf := make([]byte, 1024)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("angle receiver closed")
				return nil
			default:
				time.Sleep(5 * time.Millisecond)
				if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
					log.Println(err)
				}
				l, _, err := conn.ReadFrom(buf)
				if err != nil {
					log.Println("Connect ERROR : ", err)
					continue
				}
				var d Data
				if err := json.Unmarshal(buf[:l], &d); err != nil {
					return fmt.Errorf("json unmarshal failed -> %w", err)
				}
				angleCh <- d.Angle
			}
		}
	})

	err := portaudio.Initialize()
	defer portaudio.Terminate()
	if err != nil {
		err = fmt.Errorf("initialize failed -> %w", err)
		log.Fatalln(err)
	}

	h, err := portaudio.DefaultHostApi()
	if err != nil {
		log.Fatalln(fmt.Errorf("get DefaultHostAPI failed -> %w", err))
	}
	paParam := portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice)
	paParam.Output.Channels = 2
	paParam.SampleRate = 48000
	fmt.Printf("output parameters: %+v", paParam.Output)
	bufOut := make([]int16)



	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

func processAudio( out []float32) {
	for i := range out {
		out[i] = .7 * e.buffer[e.i]
		e.buffer[e.i] = in[i]
		e.i = (e.i + 1) % len(e.buffer)
	}
}

func readAllSLTF()  {
	
}

// func read_socket
// チャネルを使って角度の変更を読み取る
