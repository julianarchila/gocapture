package parser

import (
	"encoding/binary"
	"net"

	"github.com/google/gopacket/layers"
	"github.com/julianarchila/gocapture/pkg/models"
)

// FrameParser handles parsing of different types of network frames
type FrameParser struct {
	ethernetParser *EthernetParser
	wlanParser     *WLANParser
}

// NewFrameParser creates a new frame parser
func NewFrameParser() *FrameParser {
	return &FrameParser{
		ethernetParser: NewEthernetParser(),
		wlanParser:     NewWLANParser(),
	}
}

// ParseFrame identifies the frame type and parses it accordingly
func (fp *FrameParser) ParseFrame(frame *models.Frame) {
	// Try to determine the frame type based on the packet
	packet := frame.OriginalPacket

	// Check if it's an Ethernet frame (802.3)
	if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
		frame.FrameType = models.EthernetFrame
		fp.ethernetParser.Parse(frame)
		return
	}

	// Check if it's a wireless frame (802.11)
	// This requires the interface to be in monitor mode to capture 802.11 headers
	// Try to parse the raw data as a WLAN frame
	if len(frame.RawData) >= 2 {
		// Check frame control field to identify if it's a WLAN frame
		frameControl := binary.LittleEndian.Uint16(frame.RawData[0:2])
		frameType := (frameControl >> 2) & 0x3

		switch frameType {
		case 0: // Management frame
			frame.FrameType = models.WLANManagementFrame
			fp.wlanParser.Parse(frame)
			return
		case 1: // Control frame
			frame.FrameType = models.WLANControlFrame
			fp.wlanParser.Parse(frame)
			return
		case 2: // Data frame
			frame.FrameType = models.WLANDataFrame
			fp.wlanParser.Parse(frame)
			return
		}
	}

	// If we reach here, we couldn't identify the frame type
	// We'll leave it to the Ethernet parser as a fallback
	frame.FrameType = models.EthernetFrame
	fp.ethernetParser.Parse(frame)
}

// EthernetParser parses IEEE 802.3 Ethernet frames
type EthernetParser struct{}

// NewEthernetParser creates a new Ethernet parser
func NewEthernetParser() *EthernetParser {
	return &EthernetParser{}
}

// Parse parses an Ethernet frame
func (ep *EthernetParser) Parse(frame *models.Frame) {
	packet := frame.OriginalPacket

	// Parse Ethernet layer
	if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
		ethernet, _ := ethernetLayer.(*layers.Ethernet)
		frame.SourceMAC = ethernet.SrcMAC.String()
		frame.DestinationMAC = ethernet.DstMAC.String()
		frame.EtherType = uint16(ethernet.EthernetType)

		// Check for VLAN tagging
		if vlanLayer := packet.Layer(layers.LayerTypeDot1Q); vlanLayer != nil {
			vlan, _ := vlanLayer.(*layers.Dot1Q)
			frame.VLANInfo = map[string]interface{}{
				"Priority": vlan.Priority,
				"VID":      vlan.VLANIdentifier,
			}
		}
	}

	frame.Parsed = true
}

// WLANParser parses IEEE 802.11 WLAN frames
type WLANParser struct{}

// NewWLANParser creates a new WLAN parser
func NewWLANParser() *WLANParser {
	return &WLANParser{}
}

// Parse parses a WLAN frame
func (wp *WLANParser) Parse(frame *models.Frame) {
	// We need to manually parse the raw data for 802.11 frames
	// since gopacket may not have full support for all 802.11 frame types

	// Parse frame control field
	if len(frame.RawData) < 24 { // Minimum size for a valid 802.11 frame
		return
	}

	frameControl := binary.LittleEndian.Uint16(frame.RawData[0:2])
	duration := binary.LittleEndian.Uint16(frame.RawData[2:4])

	// Extract frame type and subtype
	frameType := (frameControl >> 2) & 0x3
	frameSubtype := (frameControl >> 4) & 0xF

	// Extract flags
	toDS := (frameControl >> 8) & 0x1
	fromDS := (frameControl >> 9) & 0x1
	moreFragments := (frameControl >> 10) & 0x1
	retry := (frameControl >> 11) & 0x1
	powerManagement := (frameControl >> 12) & 0x1
	moreData := (frameControl >> 13) & 0x1
	protected := (frameControl >> 14) & 0x1
	order := (frameControl >> 15) & 0x1

	// Store frame control info
	frame.FrameControl = map[string]interface{}{
		"Type":            frameType,
		"Subtype":         frameSubtype,
		"ToDS":            toDS == 1,
		"FromDS":          fromDS == 1,
		"MoreFragments":   moreFragments == 1,
		"Retry":           retry == 1,
		"PowerManagement": powerManagement == 1,
		"MoreData":        moreData == 1,
		"Protected":       protected == 1,
		"Order":           order == 1,
	}

	frame.Duration = duration

	// Parse addresses based on frame type
	// Address fields depend on the ToDS and FromDS flags

	// Address 1 is always present (DA/RA)
	frame.Address1 = net.HardwareAddr(frame.RawData[4:10]).String()

	// Address 2 is always present (SA/TA)
	frame.Address2 = net.HardwareAddr(frame.RawData[10:16]).String()

	// Address 3 is always present (varies based on ToDS/FromDS)
	frame.Address3 = net.HardwareAddr(frame.RawData[16:22]).String()

	// Get sequence control field
	frame.SequenceControl = binary.LittleEndian.Uint16(frame.RawData[22:24])

	// Address 4 is only present if both ToDS and FromDS are set
	offset := 24
	if toDS == 1 && fromDS == 1 && len(frame.RawData) >= 30 {
		frame.Address4 = net.HardwareAddr(frame.RawData[24:30]).String()
		offset = 30
	}

	// Set SourceMAC and DestinationMAC based on DS flags
	switch {
	case toDS == 0 && fromDS == 0:
		// Ad hoc or IBSS
		frame.DestinationMAC = frame.Address1 // DA
		frame.SourceMAC = frame.Address2      // SA
		// Address3 is BSSID
	case toDS == 1 && fromDS == 0:
		// Infrastructure, to AP
		frame.DestinationMAC = frame.Address3 // DA (final destination)
		frame.SourceMAC = frame.Address2      // SA
		// Address1 is BSSID
	case toDS == 0 && fromDS == 1:
		// Infrastructure, from AP
		frame.DestinationMAC = frame.Address1 // DA
		frame.SourceMAC = frame.Address3      // SA (original source)
		// Address2 is BSSID
	case toDS == 1 && fromDS == 1:
		// WDS or mesh
		frame.DestinationMAC = frame.Address3 // DA
		frame.SourceMAC = frame.Address4      // SA
		// Address1 is RA (receiver address)
		// Address2 is TA (transmitter address)
	}

	// Parse QoS info for QoS data frames
	if frame.FrameType == models.WLANDataFrame && (frameSubtype == 8 || frameSubtype == 9 || frameSubtype == 10 || frameSubtype == 11) {
		if offset+2 <= len(frame.RawData) {
			qosControl := binary.LittleEndian.Uint16(frame.RawData[offset : offset+2])
			tid := qosControl & 0xF
			eosp := (qosControl >> 4) & 0x1
			ackPolicy := (qosControl >> 5) & 0x3
			amsdu := (qosControl >> 7) & 0x1
			txop := (qosControl >> 8) & 0xFF

			frame.QoS = &models.QoSInfo{
				Priority:  int(tid),
				TID:       int(tid),
				TXOP:      int(txop),
				ACKPolicy: getACKPolicyString(ackPolicy),
				Details: map[string]interface{}{
					"EOSP":  eosp == 1,
					"AMSDU": amsdu == 1,
				},
			}

			offset += 2
		}
	}

	// Check for security info based on Protected flag
	if protected == 1 {
		frame.Security = &models.SecurityInfo{
			EncryptionType: "Unknown",
			Details:        make(map[string]interface{}),
		}

		// Try to identify encryption type
		// This is a simplified approach and may need to be refined
		if offset+4 <= len(frame.RawData) {
			// Check for WEP
			if len(frame.RawData) >= offset+4 && len(frame.RawData) <= offset+12 {
				frame.Security.EncryptionType = "WEP"
				frame.Security.Details["IV"] = frame.RawData[offset : offset+3]
				frame.Security.Details["KeyID"] = frame.RawData[offset+3] >> 6
			} else if len(frame.RawData) >= offset+8 {
				// Check TKIP/CCMP/GCMP
				if offset+12 <= len(frame.RawData) {
					// Check for CCMP
					if (frame.RawData[offset+3] & 0x20) == 0 {
						frame.Security.EncryptionType = "CCMP (WPA2)"
						frame.Security.Details["PN"] = frame.RawData[offset : offset+6]
					} else {
						// TKIP
						frame.Security.EncryptionType = "TKIP (WPA)"
						frame.Security.Details["IV"] = frame.RawData[offset : offset+4]
						frame.Security.Details["ExtIV"] = frame.RawData[offset+4 : offset+8]
					}
				}
			}
		}
	}

	// Handle management frames special parsing
	if frame.FrameType == models.WLANManagementFrame {
		wp.parseManagementFrame(frame, frameSubtype, offset)
	}

	frame.Parsed = true
}

// parseManagementFrame parses management frame details
func (wp *WLANParser) parseManagementFrame(frame *models.Frame, subtype uint16, offset int) {
	// Add management frame specific details
	managementInfo := make(map[string]interface{})

	switch subtype {
	case 0: // Association Request
		managementInfo["Type"] = "Association Request"
	case 1: // Association Response
		managementInfo["Type"] = "Association Response"
	case 2: // Reassociation Request
		managementInfo["Type"] = "Reassociation Request"
	case 3: // Reassociation Response
		managementInfo["Type"] = "Reassociation Response"
	case 4: // Probe Request
		managementInfo["Type"] = "Probe Request"
	case 5: // Probe Response
		managementInfo["Type"] = "Probe Response"
	case 8: // Beacon
		managementInfo["Type"] = "Beacon"
		// Parse beacon specific fields
		if offset+12 <= len(frame.RawData) {
			// Parse timestamp
			timestamp := binary.LittleEndian.Uint64(frame.RawData[offset : offset+8])
			offset += 8

			// Parse beacon interval
			beaconInterval := binary.LittleEndian.Uint16(frame.RawData[offset : offset+2])
			offset += 2

			// Parse capability info
			capabilityInfo := binary.LittleEndian.Uint16(frame.RawData[offset : offset+2])
			offset += 2

			managementInfo["Timestamp"] = timestamp
			managementInfo["BeaconInterval"] = beaconInterval
			managementInfo["CapabilityInfo"] = map[string]bool{
				"ESS":               (capabilityInfo & 0x0001) != 0,
				"IBSS":              (capabilityInfo & 0x0002) != 0,
				"CF-Pollable":       (capabilityInfo & 0x0004) != 0,
				"CF-Poll-Request":   (capabilityInfo & 0x0008) != 0,
				"Privacy":           (capabilityInfo & 0x0010) != 0,
				"ShortPreamble":     (capabilityInfo & 0x0020) != 0,
				"PBCC":              (capabilityInfo & 0x0040) != 0,
				"ChannelAgility":    (capabilityInfo & 0x0080) != 0,
				"SpectrumMgmt":      (capabilityInfo & 0x0100) != 0,
				"QoS":               (capabilityInfo & 0x0200) != 0,
				"ShortSlotTime":     (capabilityInfo & 0x0400) != 0,
				"APSD":              (capabilityInfo & 0x0800) != 0,
				"RadioMeasurement":  (capabilityInfo & 0x1000) != 0,
				"DSSS-OFDM":         (capabilityInfo & 0x2000) != 0,
				"DelayedBlockAck":   (capabilityInfo & 0x4000) != 0,
				"ImmediateBlockAck": (capabilityInfo & 0x8000) != 0,
			}
		}
	case 9: // ATIM
		managementInfo["Type"] = "ATIM"
	case 10: // Disassociation
		managementInfo["Type"] = "Disassociation"
	case 11: // Authentication
		managementInfo["Type"] = "Authentication"
	case 12: // Deauthentication
		managementInfo["Type"] = "Deauthentication"
	case 13: // Action
		managementInfo["Type"] = "Action"
	default:
		managementInfo["Type"] = "Unknown"
	}

	if frame.AnalysisResults == nil {
		frame.AnalysisResults = make(map[string]interface{})
	}
	frame.AnalysisResults["ManagementInfo"] = managementInfo
}

// getACKPolicyString returns a string representation of the ACK policy
func getACKPolicyString(policy uint16) string {
	switch policy {
	case 0:
		return "Normal ACK"
	case 1:
		return "No ACK"
	case 2:
		return "No Explicit ACK"
	case 3:
		return "Block ACK"
	default:
		return "Unknown"
	}
}
