package comunica_serial

// This package is hosted at: https://github.com/TsukiGva2/comunica_serial

import (
	"fmt"

	"go.bug.st/serial"
)

type SerialSender struct {
	port   serial.Port
	dataCh chan string // Channel to send data
}

// NewSerialSender initializes and returns a SerialSender instance.
//
// Parameters:
//   - portName: The name of the serial port to open.
//   - baudRate: The baud rate for the serial communication.
//
// Returns:
//   - sender: A pointer to the initialized SerialSender instance.
//   - err: An error if the initialization fails.
func NewSerialSender(portName string, baudRate int) (sender *SerialSender, err error) {

	mode := &serial.Mode{
		BaudRate: baudRate,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	var port serial.Port
	port, err = serial.Open(portName, mode)

	if err != nil {
		return
	}

	sender = &SerialSender{
		port:   port,
		dataCh: make(chan string), // Initialize the channel
	}

	// Start a goroutine to listen to the channel and send data
	go sender.listenAndSend()

	return
}

// listenAndSend listens to the data channel and sends data through the serial port.
func (s *SerialSender) listenAndSend() {

	for data := range s.dataCh {
		s.port.Write([]byte(data))
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
//   - ports: A slice of strings representing the available serial ports.
//   - err: An error if retrieving the ports fails.
func GetAvailablePorts() (ports []string, err error) {

	ports, err = serial.GetPortsList()

	if err != nil {
		return
	}

	if len(ports) == 0 {
		err = fmt.Errorf("no serial ports found")
		return
	}

	return
}
