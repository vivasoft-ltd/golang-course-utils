package main

import "github.com/vivasoft-golang-course/utils/logger"

func main() {
	logger.SetFileLogger("app.log")
	logger.Error("This is an error message")
	// This is a placeholder for the main function.
	// The actual implementation would depend on the specific requirements of the application.
}
