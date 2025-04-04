package comunica_serial

// This package is hosted at: https://github.com/TsukiGva2/comunica_serial

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

type SerialSender struct {
	port   serial.Port
	dataCh chan string // Channel to send data

	BaudRate int
}

// NewSerialSender initializes and returns a SerialSender instance.
//
// Parameters:
//   - baudRate: The baud rate for the serial communication.
//
// Returns:
//   - sender: A pointer to the initialized SerialSender instance.
//   - err: An error if the initialization fails.
func NewSerialSender(baudRate int) (sender *SerialSender, err error) {

	sender = &SerialSender{
		dataCh:   make(chan string),
		BaudRate: baudRate,
	}

	err = sender.Open()

	if err != nil {
		close(sender.dataCh)
		return
	}

	// Start a goroutine to listen to the channel and send data
	go sender.listenAndSend()

	return
}

func (s *SerialSender) Open() (err error) {

	var portName string
	var newPort serial.Port

	backoff := time.Millisecond * 100 // Initial backoff duration
	maxRetries := 5                   // Maximum number of retries
	retries := 0

	for retries < maxRetries {
		<-time.After(backoff) // Wait for the backoff duration

		log.Println("Attempting to reopen the serial port...")

		portName, err = GetFirstAvailablePortName()

		if err != nil {
			log.Printf("Failed to get available port: %v\n", err)
			retries++
			backoff *= 2 // Exponential backoff

			continue
		}

		mode := &serial.Mode{
			BaudRate: s.BaudRate,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		}

		newPort, err = serial.Open(portName, mode)

		if err != nil {
			log.Printf("Failed to reopen serial port: %v\n", err)
			retries++
			backoff *= 2 // Exponential backoff
			continue
		}

		s.port = newPort

		log.Println("Serial port opened successfully.")

		return
	}

	log.Println("Max retries reached. Giving up on reopening the serial port.")

	return
}

// listenAndSend listens to the data channel and sends data through the serial port.
func (s *SerialSender) listenAndSend() {

	for data := range s.dataCh {
		_, err := s.port.Write(append([]byte(data), '\n'))

		if err != nil {
			log.Printf("Error writing to serial port: %v\n", err)

			s.port.Close()
			s.Open()
		}
	}
}

// SendData sends the provided data through the channel.
//
// Parameters:
//   - data: The string data to send.
//
// Returns:
//   - err: An error if sending data fails.
func (s *SerialSender) SendData(data string) {

	s.dataCh <- data // Send data to the channel
}

// Close closes the serial port and the data channel.
func (s *SerialSender) Close() {

	close(s.dataCh) // Close the channel
	s.port.Close()  // Close the serial port
}

// GetAvailablePorts returns a list of available serial ports.
//
// Returns:
//   - port: A string containing the name of the first available serial port.
//   - err: An error if retrieving the ports fails.
func GetFirstAvailablePortName() (port string, err error) {

	ports, err := serial.GetPortsList()

	if err != nil {
		return
	}

	if len(ports) == 0 {
		err = fmt.Errorf("no serial ports found")
		return
	}

	port = ports[0]

	return
}
