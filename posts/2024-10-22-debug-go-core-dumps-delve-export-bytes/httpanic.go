package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(w, "panic recovered")
				proc, err := os.FindProcess(syscall.Getpid())
				if err != nil {
					panic(fmt.Sprintf("could not find own process (pid %d): %v", syscall.Getpid(), err))
				}
				proc.Signal(syscall.SIGABRT)
				// Ensure the stack which triggered the core dump sticks around
				proc.Wait()
			}
		}()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}
		_ = body

		fmt.Fprintf(w, "heya")
		panic("this should result in a core dump")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
