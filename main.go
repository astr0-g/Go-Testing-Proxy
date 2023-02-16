package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// Set up the proxy URL and credentials

	// Read the list of proxies from the file
	proxies, err := readProxies("proxy.txt")
	if err != nil {
		fmt.Println("Failed to read proxies:", err)
		return
	}

	// Loop through the list of proxies and test each one
	for _, proxy := range proxies {
		// Construct the full proxy URL
		proxyURL, err := url.Parse(fmt.Sprintf("%s%s", proxy))
		if err != nil {
			fmt.Println("Failed to parse proxy URL:", err)
			return
		}

		// Create a new HTTP client with the proxy settings
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}

		// Create a new HTTP request
		req, err := http.NewRequest("GET", "https://www.example.com/", nil)
		if err != nil {
			fmt.Println("Failed to create HTTP request:", err)
			return
		}

		// Send the HTTP request using the proxy
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Proxy %s failed: %s\n", proxy, err)
			continue
		}

		// Check the HTTP response status code
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Proxy %s is working!\n", proxy)
		} else {
			fmt.Printf("Proxy %s failed with status code %d\n", proxy, resp.StatusCode)
		}

		// Close the HTTP response body
		resp.Body.Close()
	}
}

// Reads a list of proxies from a text file
func readProxies(filename string) ([]string, error) {
	// Open the text file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Initialize a slice to hold the proxy addresses
	proxies := []string{}

	// Read each line of the file
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			// Split the line into four parts using the colon separator
			parts := strings.Split(line, ":")
			if len(parts) == 4 {
				// Construct the proxy address string and add it to the slice
				proxy := fmt.Sprintf("%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1])
				proxies = append(proxies, proxy)
			}
		}
	}

	// Check if there was an error reading the file
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return the list of proxies
	return proxies, nil
}
