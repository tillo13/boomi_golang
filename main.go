package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/joho/godotenv"
)

// Define ANSI escape codes for colors
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

// RequestData structure holds the request information to be sent to Boomi API
type RequestData struct {
	UnixRequestToBoomi    string `json:"unix_request_to_boomi"`
	PayloadRequestToBoomi string `json:"payload_request_to_boomi"`
}

// cleanString function cleans a string by removing all non-numeric characters
func cleanString(str string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	cleanedString := reg.ReplaceAllString(str, "")
	return cleanedString
}

// timer function displays messages in the console with a delay to simulate processing time
func timer(quit chan bool, wg *sync.WaitGroup, start time.Time, done chan bool) {
	defer wg.Done()
	messages := []string{
		color.GreenString("Establishing connection with Boomi API endpoint"),
		color.GreenString("Sending request JSON to Boomi API"),
		color.GreenString("Setting 60-second response threshold for Boomi response"),
		color.GreenString("Waiting for Boomi API response"),
		color.GreenString("Response received: 200 OK"),
		color.GreenString("Parsing and analyzing API response JSON"),
		color.GreenString("Performing data transformation and validation in Boomi"),
		color.GreenString("Preparing final data analysis"),
		color.GreenString("Sending response JSON to Golang program"),
		color.GreenString("Received confirmation: Response JSON successfully parsed"),
		color.GreenString("Validating Boomi JSON format"),
	}

	verbiages := []string{
		color.YellowString("[Engaging thrusters]"),
		color.YellowString("[Engaging hyperdrive engines]"),
		color.YellowString("[Optimizing response modules]"),
		color.YellowString("[Awaiting transmission]"),
		color.YellowString("[Received cosmic signal]"),
		color.YellowString("[Processing space transmission]"),
		color.YellowString("[Analyzing launch coordinates]"),
		color.YellowString("[Generating atmosphere report]"),
		color.YellowString("[Transmitting results to mission control]"),
		color.YellowString("[Receiving confirmation from mission control]"),
		color.YellowString("[Parsing stellar JSON from Boomi]"),
	}

	index := 0
	displayWaitMessage := make(chan bool, 1)

	for {
		select {
		case <-quit:
			// If quit channel is closed and there are remaining messages,
			// print all the remaining messages immediately.
			if index < len(messages) {
				for ; index < len(messages); index++ {
					log.Println(messages[index], verbiages[index]+"...")
				}
			}
			return
		case <-time.After(60 * time.Second):
			log.Println("Boomi Response timer timed out after 60 seconds.")
			return
		case <-time.After(1 * time.Second):
			if index < len(messages) {
				log.Println(messages[index], verbiages[index]+"...")
				index++
			}
			if index == len(messages) {
				displayWaitMessage <- true
			}
		case <-displayWaitMessage:
			if index == len(messages) {
				startWait := time.Now()
				spinner := []string{"-", "\\", "|", "/"}
			waitLoop:
				for {
					select {
					case <-quit: // Stop displaying the wait message when the response is received
						break waitLoop
					case <-done: // Stop displaying the wait message when the done channel is closed
						break waitLoop
					default:
						elapsed := time.Since(startWait).Seconds()
						fmt.Printf("\r%s Reworking a few more things, one moment: %.5fs %s", startWait.Format("2006/01/02 15:04:05"), elapsed, spinner[int(elapsed*2)%len(spinner)])
						time.Sleep(1 * time.Millisecond)
					}
				}
			}
		}
	}
}

// printWithTimestamp function adds a timestamp to a console message
func printWithTimestamp(msg string) {
	currentTime := time.Now()
	timeString := currentTime.Format("2006/01/02 15:04:05")
	fmt.Printf("%s %s\n", timeString, msg)
}

// main function runs the script
func main() {
	printWithTimestamp("Starting system checks")

	var payloads []string // Running record of payloads entered by the user

	reader := bufio.NewReader(os.Stdin) // Move this line to the top of the main function, before the loop

	for {
		screen.Clear()
		screen.MoveTopLeft()

		log.Println(color.GreenString("Program started"), color.YellowString("[Ignition sequence initiated.]"))

		// Read the username and password from the .env file or user input if not available
		err := godotenv.Load()
		var username, password string
		if err != nil {
			fmt.Print("Unable to load .env file. Please enter your username: ")
			username, _ = reader.ReadString('\n')
			username = strings.TrimSuffix(username, "\n")

			fmt.Print("Please enter your password: ")
			password, _ = reader.ReadString('\n')
			password = strings.TrimSuffix(password, "\n")

			log.Println(color.YellowString("Using entered credentials"))
		} else {
			log.Println(color.GreenString("Loaded .env file"), color.YellowString("[Ground control, we are ready for liftoff!]"))
			username = os.Getenv("USERNAME")
			password = os.Getenv("PASSWORD")
		}

		start := time.Now() // Record the start time

		// Measure script initialization time
		scriptInitTime := time.Since(start)

		// Create a time object with the current time
		timestamp := time.Now()

		// Get the Unix timestamp in microseconds
		unixTimestamp := timestamp.UnixNano() / int64(time.Microsecond)

		// Convert the Unix timestamp to a string
		timestampString := strconv.FormatInt(unixTimestamp, 10)

		url := "https://c01-usa-east.integrate.boomi.com/ws/simple/createGeneralListener"

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter the payload to send to Boomi: ")
		payload, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		payload = strings.TrimSuffix(payload, "\n") // This will remove the trailing newline

		log.Println(color.GreenString("Received user input"), color.YellowString("[Launch coordinates received]..."))

		requestData := RequestData{
			UnixRequestToBoomi:    timestampString,
			PayloadRequestToBoomi: payload,
		}

		jsonData, err := json.MarshalIndent(requestData, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		log.Println(color.GreenString("Marshaled request data to JSON"), color.YellowString("[Converting transmission]..."))
		log.Println(color.GreenString("Creating new HTTP request"), color.YellowString("[Creating orbit request]..."))

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		req.SetBasicAuth(username, password)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		quitTimer := make(chan bool)
		var wg sync.WaitGroup
		wg.Add(1)
		doneTimer := make(chan bool)
		go timer(quitTimer, &wg, start, doneTimer)

		retryLimit := 3
		retryWaitTime := 5 * time.Second
		success := false

		for retry := 0; retry < retryLimit; retry++ {
			success = false

			log.Println(color.GreenString("Sending HTTP request to Boomi API"), color.YellowString("[Sending launch status...]"))
			fmt.Println("-----------------------------")
			fmt.Println("Request_to_Boomi:")
			fmt.Println(string(jsonData))
			fmt.Println("----------------------------")

			response, err := client.Do(req)
			if err != nil {
				log.Printf("Error while sending request: %v\n", err)
				log.Printf("Retry - Attempt %d of %d... ", retry+1, retryLimit)
				time.Sleep(retryWaitTime)
				continue
			}

			close(doneTimer) // Close the done channel to stop the timer
			quitTimer <- true
			wg.Wait()

			fmt.Println()
			log.Println("Response Status:", color.BlueString(response.Status), color.YellowString("[We have made contact!]"))

			if response.StatusCode == http.StatusOK {
				log.Println(color.GreenString("Sent HTTP request and received response"), color.YellowString("[Final logics confirmed!]"))
				log.Println(color.GreenString("Reading response body"), color.YellowString("[Validating completed flight log...]"))

				body, err := io.ReadAll(response.Body)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("-----------------------------")
				fmt.Println("Response from Boomi:")
				fmt.Println(string(body))
				fmt.Println("-----------------------------")

				defer response.Body.Close()

				// Get the Boomi received and Golang scriptstart timestamps from the response body
				var responseJSON struct {
					FullResponseFromBoomi string `json:"full_response_from_boomi"`
					IncomingTimestamp     string `json:"incoming_timestamp"`
					BoomiTimestamp        string `json:"boomi_timestamp"`
				}
				err = json.Unmarshal(body, &responseJSON)
				if err != nil {
					log.Fatal(err)
				}

				// Convert the timestamps to integers
				cleanedIncomingTimestamp := cleanString(responseJSON.IncomingTimestamp)
				incomingTimestampMicro, err := strconv.ParseInt(cleanedIncomingTimestamp, 10, 64)
				if err != nil {
					log.Fatal(err)
				}

				cleanedBoomiTimestamp := cleanString(responseJSON.BoomiTimestamp)
				boomiTimestampMicro, err := strconv.ParseInt(cleanedBoomiTimestamp, 10, 64)
				if err != nil {
					log.Fatal(err)
				}

				boomiReceivedTime := time.Unix(0, boomiTimestampMicro*int64(time.Microsecond)).Format("2006-01-02 15:04:05.000000")
				startTime := start.Format("2006-01-02 15:04:05.000000")

				fmt.Println("This Golang script started at:", color.BlueString(startTime))
				fmt.Println("Boomi received it at:", color.BlueString(boomiReceivedTime))
				fmt.Printf("Time taken between Golang creating it and Boomi responding to it: %s\n", color.BlueString(time.Duration(boomiTimestampMicro-incomingTimestampMicro).String()))
				fmt.Printf("Time taken to initialize the script: %s\n", color.BlueString(scriptInitTime.String()))
				scriptProcessingOverhead := time.Since(start) - (time.Duration(boomiTimestampMicro-incomingTimestampMicro) * time.Microsecond)
				fmt.Printf("Script Processing Overhead: %s\n", color.BlueString(scriptProcessingOverhead.String()))
				fmt.Printf("Total execution time: %s\n", color.BlueString(time.Since(start).String()))

				fmt.Println("\n----------------------------")
				fmt.Println(color.YellowString("Current payload entered:"), payload)
				payloads = append(payloads, payload)
				fmt.Println(color.YellowString("Previous payloads:"))
				for _, p := range payloads {
					fmt.Println(p)
				}
				fmt.Println("----------------------------")

				success = true
				break

			} else {
				log.Printf("Retry - Attempt %d of %d... ", retry+1, retryLimit)
				time.Sleep(5 * time.Second)
			}
		}

		if !success {
			log.Println(color.RedString("Exceeded retry limit. Please check your network connection, and try again."))
			break
		}

		fmt.Print(color.GreenString("Again (Y/N): "))
		again, _ := reader.ReadString('\n')
		again = strings.TrimSpace(again)

		if strings.ToUpper(again) != "Y" {
			break
		}
	}
}
