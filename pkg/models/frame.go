package models

import (
	"time"

	"github.com/google/gopacket"
)

// FrameType represents the type of network frame
type FrameType int

const (
	// Ethernet frame types
	EthernetFrame FrameType = iota
	
	// 802.11 frame types
	WLANManagementFrame
	WLANControlFrame
	WLANDataFrame
)

// SecurityInfo contains information about security mechanisms
type SecurityInfo struct {
	EncryptionType string // WEP, WPA, WPA2, WPA3, None
	Details        map[string]interface{}
}

// QoSInfo contains Quality of Service information
type QoSInfo struct {
	Priority    int
	TID         int // Traffic ID
	TXOP        int // Transmission Opportunity
	ACKPolicy   string
	QueueSize   int
	Details     map[string]interface{}
}

// Frame represents a network frame with all its information
type Frame struct {
	ID              int64
	Timestamp       time.Time
	FrameType       FrameType
	RawData         []byte
	Length          int
	SourceMAC       string
	DestinationMAC  string
	
	// For 802.3 Ethernet frames
	EtherType       uint16
	VLANInfo        interface{}
	
	// For 802.11 WLAN frames
	FrameControl    interface{}
	Duration        uint16
	SequenceControl uint16
	Address1        string // Usually destination
	Address2        string // Usually source
	Address3        string // Usually BSSID
	Address4        string // Used in ad-hoc mode
	
	// Security and QoS info
	Security        *SecurityInfo
	QoS             *QoSInfo
	
	// Analysis results
	Parsed          bool
	AnalysisResults map[string]interface{}
	
	// Original packet for further analysis
	OriginalPacket  gopacket.Packet
}

// FrameFilter defines criteria for filtering captured frames
type FrameFilter struct {
	FrameTypes      []FrameType
	SourceMAC       string
	DestinationMAC  string
	BSSID           string
	EncryptionTypes []string
	ContainsBytes   []byte
} 