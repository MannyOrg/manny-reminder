package utils

import (
	"fmt"
	"io"
	"log"
)

func SendHttpError(w io.Writer, err error) {
	sendHttpError(w, err.Error())
}
func SendHttpStringError(w io.Writer, err string) {
	sendHttpError(w, err)
}

func sendHttpError(w io.Writer, message string) {
	_, err := fmt.Fprintf(w, message)
	if err != nil {
		log.Default().Print(err.Error())
		return
	}
}
