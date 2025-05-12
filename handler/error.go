package handler
import (
	"fmt"
	"log"
	"net/http"
)

// ErrorHandler manages HTTP errors by sending appropriate status codes and messages
func ErrorHandler(w http.ResponseWriter, r *http.Request, statusCode int) {
	// Define custom error messages based on status codes
	var message string

	switch statusCode {
	case http.StatusBadRequest:
		message = "Bad request: The server cannot process the request due to client error"
	case http.StatusUnauthorized:
		message = "Unauthorized: Authentication is required and has failed or not been provided"
	case http.StatusForbidden:
		message = "Forbidden: You don't have permission to access this resource"
	case http.StatusNotFound:
		message = "Not Found: The requested resource could not be found"
	case http.StatusMethodNotAllowed:
		message = "Method Not Allowed: The request method is not supported for this resource"
	case http.StatusInternalServerError:
		message = "Internal Server Error: Something went wrong on the server"
	default:
		message = "An error occurred"
	}

	// Log the error with request details
	log.Printf("Error %d (%s): %s %s from %s",
		statusCode,
		http.StatusText(statusCode),
		r.Method,
		r.URL.Path,
		r.RemoteAddr)

	// Set the content type and status code
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	// Create a simple HTML error page
	errorHTML := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>Error %d - %s</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    margin: 40px;
                    line-height: 1.6;
                }
                .error-container {
                    max-width: 600px;
                    margin: 0 auto;
                    padding: 20px;
                    border: 1px solid #ddd;
                    border-radius: 5px;
                    background-color: #f8f8f8;
                }
                h1 {
                    color: #d9534f;
                }
                .back-link {
                    margin-top: 20px;
                }
            </style>
        </head>
        <body>
            <div class="error-container">
                <h1>Error %d - %s</h1>
                <p>%s</p>
                <div class="back-link">
                    <a href="/">Return to Home Page</a>
                </div>
            </div>
        </body>
        </html>
    `, statusCode, http.StatusText(statusCode), statusCode, http.StatusText(statusCode), message)

	// Write the error page to the response
	fmt.Fprint(w, errorHTML)
}


