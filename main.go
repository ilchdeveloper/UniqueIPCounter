package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/axiomhq/hyperloglog" // Importing the HyperLogLog library for estimating unique elements
)

func main() {
	// Path to the file containing the list of IP addresses.
	// Replace "ip_addresses" with the actual path of your IP file.
	filePath := "ip_addresses"
	// Call the function to count unique IP addresses and handle potential errors.
	uniqueCount, err := countUniqueIPs(filePath)
	if err != nil {
		log.Fatalf("Error counting unique IPs: %v", err)
	}
	// Print the estimated number of unique IP addresses to the console.
	fmt.Printf("Estimated number of unique IP addresses: %d\n", uniqueCount)
}

// countUniqueIPs reads a file and estimates the count of unique IP addresses.
func countUniqueIPs(filePath string) (uint64, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		// Return an error if the file cannot be opened.
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	// Ensure the file is closed when the function exits.
	defer file.Close()

	// Create a new HyperLogLog counter with a precision of 14
	hll := hyperloglog.New14()
	if hll == nil {
		return 0, fmt.Errorf("failed to create HyperLogLog: %w", err)
	}

	// Use a scanner to read the file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		// Add each IP address to the HyperLogLog counter
		hll.Insert([]byte(ip))
	}

	// Check if there were any errors while reading the file.
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}

	// HyperLogLog uses its internal data to provide an approximate count.
	return hll.Estimate(), nil
}
