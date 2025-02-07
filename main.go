package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
	"github.com/syumai/workers/cloudflare/fetch"
)

// AuthHeaderName is the header name for the secret key.
const AuthHeaderName = "X-Proxy-Secret"

const (
	contentTypeJSON  = "application/json"
	contentTypePlain = "text/plain"
)

var (
	errMissingSubdomain = errorResponse{Message: "Missing subdomain."}
	errUnauthorized     = errorResponse{Message: "Unauthorized"}
)

// errorResponse represents an error message structure.
type errorResponse struct {
	Message string `json:"message"`
}

func main() {
	// Get secret key from environment
	secretKey := cloudflare.Getenv("PROXY_SECRET_KEY")
	if secretKey == "" {
		panic("PROXY_SECRET_KEY environment variable is required")
	}

	workers.Serve(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cloudflare.PassThroughOnException()

		// Check secret key
		if req.Header.Get(AuthHeaderName) != secretKey {
			sendJSONError(w, &errUnauthorized, http.StatusUnauthorized)
			return
		}

		// Get first path segment
		path := strings.TrimPrefix(req.URL.Path, "/")
		i := strings.IndexByte(path, '/')
		var subdomain string
		if i == -1 {
			subdomain = path
		} else {
			subdomain = path[:i]
		}

		// Check if subdomain is provided
		if subdomain == "" {
			sendJSONError(w, &errMissingSubdomain, http.StatusBadRequest)
			return
		}

		// Check if subdomain is valid
		if strings.Contains(subdomain, ".") {
			sendJSONError(w, &errorResponse{Message: "Invalid subdomain format."}, http.StatusBadRequest)
			return
		}

		// Construct the target URL
		var targetURL string
		if i == -1 {
			targetURL = "https://" + subdomain + ".roblox.com/"
		} else {
			targetURL = "https://" + subdomain + ".roblox.com/" + path[i+1:]
		}
		if req.URL.RawQuery != "" {
			targetURL += "?" + req.URL.RawQuery
		}

		// Create fetch client and request
		fc := fetch.NewClient()
		fetchReq, err := fetch.NewRequest(req.Context(), req.Method, targetURL, req.Body)
		if err != nil {
			fmt.Printf("Error creating proxy request: %v\n", err)
			sendJSONError(w, &errorResponse{Message: "Failed to create proxy request."}, http.StatusInternalServerError)
			return
		}

		// Copy all headers from original request except the secret key
		for key, values := range req.Header {
			if !strings.EqualFold(key, AuthHeaderName) {
				fetchReq.Header[key] = values
			}
		}

		// Ensure Content-Type is set for POST/PUT requests
		if (req.Method == "POST" || req.Method == "PUT") && fetchReq.Header.Get("Content-Type") == "" {
			fetchReq.Header.Set("Content-Type", contentTypeJSON)
		}

		// Perform the request
		resp, err := fc.Do(fetchReq, &fetch.RequestInit{
			Redirect: fetch.RedirectModeFollow,
		})
		if err != nil {
			fmt.Printf("Error proxying request to %s: %v\n", targetURL, err)
			sendJSONError(w, &errorResponse{Message: "Failed to proxy request: " + err.Error()}, http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Log error responses asynchronously
		if resp.StatusCode >= 400 {
			cloudflare.WaitUntil(func() {
				fmt.Printf("Error proxying request to %s (status: %d)\n", targetURL, resp.StatusCode)
			})
		}

		// Copy response headers
		for key, values := range resp.Header {
			w.Header()[key] = values
		}

		// Set response status code and copy body
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}))
}

// sendJSONError sends a JSON error response.
func sendJSONError(w http.ResponseWriter, err *errorResponse, status int) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	if encErr := json.NewEncoder(w).Encode(err); encErr != nil {
		fmt.Printf("Error encoding JSON response: %v\n", encErr)
		w.Header().Set("Content-Type", contentTypePlain)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error: Failed to encode JSON response"))
	}
}
