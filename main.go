package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const logTemplate = "{\"target\":\"messaging_outbox::services::outbox_transactional_processor\",\"threadName\":\"main\",\"line_number\":315,\"level\":\"DEBUG\",\"timestamp\":\"%s\",\"span\":{\"name\":\"run_iteration\"},\"spans\":[{\"name\":\"run\"},{\"name\":\"run_relay_iteration\"}],\"fields\":{\"message\":\"%s\"}}\n"

func StringWithCharset(length int) string {
	rand.NewSource(time.Now().UTC().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateLog() {
	fmt.Printf(logTemplate, time.Now().Format(time.RFC3339Nano), StringWithCharset(500))
}

func generateManyLogs(done <-chan bool, wg *sync.WaitGroup) {
For:
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			wg.Done()
			break For
		case <-time.Tick(1 * time.Millisecond):
			generateLog()
		}
	}
}

func main() {
	const generatorsNum = 6

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	doneChans := make([]chan bool, 0)

	wg := sync.WaitGroup{}
	wg.Add(generatorsNum)

	for i := 0; i < generatorsNum; i++ {
		done := make(chan bool, 1)
		doneChans = append(doneChans, done)
		go generateManyLogs(done, &wg)
	}

	<-sigs

	for _, done := range doneChans {
		done <- true
	}

	wg.Wait()
}
