// Copyright Oula Kuuva 2019
// Logging stuff courtesy of https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html
package main

import (
	"log"
	"io"
)

var (
    debug   *log.Logger
    info    *log.Logger
    warning *log.Logger
    error   *log.Logger
)

func initLoggers(
    debugHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    debug = log.New(debugHandle,
        "DEBUG: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}