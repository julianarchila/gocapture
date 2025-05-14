package capture

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/julianarchila/gocapture/internal/parser"
	"github.com/julianarchila/gocapture/pkg/models"
)

// CaptureEngine handles network frame capture
type CaptureEngine struct {
	handle        *pcap.Handle
	interfaceName string
	promiscuous   bool
	filter        string
	isRunning     bool
	frameChannel  chan *models.Frame
	stopChannel   chan struct{}
	frameCounter  int64
	mutex         sync.Mutex
	frameParser   *parser.FrameParser
}

// NewCaptureEngine creates a new capture engine
func NewCaptureEngine(interfaceName string, promiscuous bool, filter string) (*CaptureEngine, error) {
	// Check if interface exists
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("error finding devices: %v", err)
	}

	var found bool
	for _, device := range devices {
		if device.Name == interfaceName {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("interface %s not found", interfaceName)
	}

	frameParser := parser.NewFrameParser()

	return &CaptureEngine{
		interfaceName: interfaceName,
		promiscuous:   promiscuous,
		filter:        filter,
		frameChannel:  make(chan *models.Frame, 1000), // Buffer for 1000 frames
		stopChannel:   make(chan struct{}),
		frameParser:   frameParser,
	}, nil
}

// Start begins the capture process
func (ce *CaptureEngine) Start() error {
	ce.mutex.Lock()
	defer ce.mutex.Unlock()

	if ce.isRunning {
		return fmt.Errorf("capture already running")
	}

	// Reinitialize channels for new capture session
	ce.frameChannel = make(chan *models.Frame, 1000) // Buffer for 1000 frames
	ce.stopChannel = make(chan struct{})

	// Open the device for capturing
	// snaplen: 65535 - maximum capture size
	// promiscuous: true - capture all packets, not just those addressed to this interface
	// timeout: 500ms - read timeout
	var err error
	ce.handle, err = pcap.OpenLive(ce.interfaceName, 65535, ce.promiscuous, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("error opening interface %s: %v", ce.interfaceName, err)
	}

	// Set BPF filter if specified
	if ce.filter != "" {
		if err := ce.handle.SetBPFFilter(ce.filter); err != nil {
			ce.handle.Close()
			return fmt.Errorf("error setting BPF filter: %v", err)
		}
	}

	ce.isRunning = true

	// Start packet processing in a goroutine
	go ce.captureFrames()

	return nil
}

// Stop stops the capture process
func (ce *CaptureEngine) Stop() {
	ce.mutex.Lock()
	defer ce.mutex.Unlock()

	if !ce.isRunning {
		return
	}

	close(ce.stopChannel)
	if ce.handle != nil {
		ce.handle.Close()
	}
	ce.isRunning = false
}

// GetFrameChannel returns the channel for receiving captured frames
func (ce *CaptureEngine) GetFrameChannel() <-chan *models.Frame {
	return ce.frameChannel
}

// IsRunning returns whether the capture engine is currently running
func (ce *CaptureEngine) IsRunning() bool {
	ce.mutex.Lock()
	defer ce.mutex.Unlock()
	return ce.isRunning
}

// GetInterfaceName returns the name of the interface being captured
func (ce *CaptureEngine) GetInterfaceName() string {
	return ce.interfaceName
}

// captureFrames is the main capture loop
func (ce *CaptureEngine) captureFrames() {
	packetSource := gopacket.NewPacketSource(ce.handle, ce.handle.LinkType())
	packetChannel := packetSource.Packets()

	for {
		select {
		case <-ce.stopChannel:
			close(ce.frameChannel)
			return
		case packet, ok := <-packetChannel:
			if !ok {
				close(ce.frameChannel)
				return
			}

			// Process the packet and convert it to a Frame
			ce.mutex.Lock()
			ce.frameCounter++
			frameID := ce.frameCounter
			ce.mutex.Unlock()

			// Create the base frame
			frame := &models.Frame{
				ID:             frameID,
				Timestamp:      time.Now(),
				RawData:        packet.Data(),
				Length:         len(packet.Data()),
				OriginalPacket: packet,
			}

			// Parse the frame based on its type
			ce.frameParser.ParseFrame(frame)

			// Send the frame to the channel
			ce.frameChannel <- frame
		}
	}
}
