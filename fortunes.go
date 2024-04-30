package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Function to read fortunes from the file and return a slice of strings
func readFortunesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fortunes []string             // Slice to store the fortunes
	scanner := bufio.NewScanner(file) // Scanner to read the file line by line
	scanner.Split(bufio.ScanLines)    // Split the file into lines

	var fortuneBuilder strings.Builder // Builder to build the fortunes

	for scanner.Scan() { // Iterate over the lines of the file
		line := scanner.Text() // Get the current line
		if line == "%%" {      // Check if the line is the delimiter
			fortunes = append(fortunes, fortuneBuilder.String())
			fortuneBuilder.Reset()
		} else { // If the line is not the delimiter, add it to the current fortune
			fortuneBuilder.WriteString(line)
			fortuneBuilder.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil { // Check for any errors during the scan
		return nil, err
	}

	return fortunes, nil // Return the slice of fortunes
}

// Function to pick a random fortune from the provided slice of fortunes
func pickRandomFortune(fortunes []string) string {
	rand.Seed(time.Now().UnixNano())        // Seed the random number generator
	randomIndex := rand.Intn(len(fortunes)) // Generate a random index
	return fortunes[randomIndex]            // Return the fortune at the random index
}

func fortune(fortunes []string, requestChan chan bool, responseChan chan string) {
	for { // Infinite loop to keep the goroutine running
		<-requestChan // Wait for a message on the channel
		selectedFortune := pickRandomFortune(fortunes)
		responseChan <- selectedFortune // Send the selected fortune back via the response channel
	}
}

func main() {
	fortuneFile := "fortunes.txt"
	fortunes, err := readFortunesFromFile(fortuneFile)
	if err != nil { // Check for errors while reading the fortunes
		fmt.Printf("Error reading fortunes from file: %v\n", err)
		return
	}

	requestChan := make(chan bool)
	responseChan := make(chan string)

	// Start the fortune goroutine
	go fortune(fortunes, requestChan, responseChan)

	// Main loop to interact with the user
	fmt.Println("Welcome to the Fortune Teller!")
	for {
		var input string
		fmt.Print("Would you like another fortune? (yes/no): ")
		fmt.Scanln(&input)

		switch strings.ToLower(input) {
		case "yes": // If the user wants another fortune, send a message down the channel to request a fortune
			requestChan <- true // Send a message down the channel to request a fortune
			fortune := <-responseChan
			fmt.Println(fortune)
		case "no": // If the user does not want another fortune, exit the program
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid input. Please answer with 'yes' or 'no'.")
		}
	}
}
