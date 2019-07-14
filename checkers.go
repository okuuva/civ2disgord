package main

import "net/http"

func checkResponses(responses []*http.Response, logger *logger) bool {
	success := true
	for _, response := range responses {
		logger.debug.Println(response)
		logger.debug.Println(response.Request)
		url := response.Request.URL.String()
		if response.StatusCode != 204 {
			logger.error.Printf("Failed to send message to %s!", url)
			success = false
		} else {
			logger.info.Printf("Successfully sent message to %s", url)
		}
	}
	return success
}

func checkErrors(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
