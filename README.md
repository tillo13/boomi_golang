# Boomi API Interface

This project is a command line interface for sending payloads to the Boomi API. It's written in Go and includes automatic payload handling, request and response timing, and script initialization timing.

## Table of Contents

- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

To run this project locally, you will need to clone the repository, install the dependencies, and compile the Go script.

## Prerequisites

- Go version 1.16 or higher
- An internet connection
- Access to the Boomi API

## Installation

1. Clone the repository to your local machine.
2. Navigate to the directory where you cloned the repository.
3. Run `go mod tidy` to fetch the dependencies.
4. Compile the Go script with `go build main.go`.

## Usage

1. Run the compiled script: `./main`.
2. Follow the prompts to provide your Boomi username and password.
3. Enter the payload to send to Boomi.
4. The program will send the payload to the Boomi API and print out the response. You can repeat this process as many times as you like.

## Contributing

Feel free to fork this project and submit pull requests. All contributions are welcome.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
