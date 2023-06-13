package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type RequestData struct {
	UnixRequestToBoomi    string `json:"unix_request_to_boomi"`
	PayloadRequestToBoomi string `json:"payload_request_to_boomi"`
}

func sendHTTPRequest(url, username, password string, jsonData []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func cleanString(str string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	cleanedString := reg.ReplaceAllString(str, "")
	return cleanedString
}

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

func printWithTimestamp(msg string) {
	currentTime := time.Now()
	timeString := currentTime.Format("2006/01/02 15:04:05")
	fmt.Printf("%s %s\n", timeString, msg)
}

func main() {
	printWithTimestamp("Starting system checks")

	var payloads []string // Running record of payloads entered by the user

	reader := bufio.NewReader(os.Stdin)

	for {
		screen.Clear()
		screen.MoveTopLeft()

		log.Println(color.GreenString("Program started"), color.YellowString("[Ignition sequence initiated.]"))

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
		scriptInitTime := time.Since(start)

		timestamp := time.Now()
		unixTimestamp := timestamp.UnixNano() / int64(time.Microsecond)
		timestampString := strconv.FormatInt(unixTimestamp, 10)

		url := "https://c01-usa-east.integrate.boomi.com/ws/simple/createGeneralListener"
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter the payload to send to Boomi: ")
		payload, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		payload = strings.TrimSuffix(payload, "\n")

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

		quitTimer := make(chan bool)
		var wg sync.WaitGroup
		wg.Add(1)
		doneTimer := make(chan bool)
		go timer(quitTimer, &wg, start, doneTimer)

		retryLimit := 3
		retryWaitTime := 5 * time.Second
		var response *http.Response
		var success bool

		for retry := 0; retry < retryLimit; retry++ {
			success = false
			response, err = sendHTTPRequest(url, username, password, jsonData)

			if err == nil && response.StatusCode == http.StatusOK {
				success = true
				break
			} else {
				log.Printf("Retry - Attempt %d of %d... ", retry+1, retryLimit)
				time.Sleep(retryWaitTime)
			}
		}

		if success {
			close(doneTimer) // Close the done channel to stop the timer
			quitTimer <- true
			wg.Wait()

			fmt.Println()
			log.Println("Response Status:", color.BlueString(response.Status), color.YellowString("[We have made contact!]"))

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			var responseJSON struct {
				FullResponseFromBoomi string `json:"full_response_from_boomi"`
				IncomingTimestamp     string `json:"incoming_timestamp"`
				BoomiTimestamp        string `json:"boomi_timestamp"`
			}
			err = json.Unmarshal(body, &responseJSON)
			if err != nil {
				log.Fatal(err)
			}

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

			fmt.Print(color.GreenString("Again (Y/N): "))
			again, _ := reader.ReadString('\n')
			again = strings.TrimSpace(again)

			if strings.ToUpper(again) != "Y" {
				break
			}
		} else {
			log.Println(color.RedString("Exceeded retry limit. Please check your network connection, and try again."))
			break
		}
	}
}
