// Copyright Oula Kuuva 2019
// Logging stuff courtesy of https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html
package main

import (
	"io"
	"log"
	"os"
)

type logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

func (l logger) checkFatal(err error, message string, returnCode int) {
	if err != nil {
		l.error.Println(message)
		l.error.Println(err)
		os.Exit(returnCode)
	}
}

func newLogger(
	debugHandle io.Writer,
	stdHandle io.Writer,
	errorHandle io.Writer) *logger {

	var l logger
	l.debug = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.info = log.New(stdHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.warning = log.New(stdHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return &l
}
