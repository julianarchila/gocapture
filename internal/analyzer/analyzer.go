package analyzer

import (
	"fmt"

	"github.com/julianarchila/gocapture/pkg/models"
)

// FrameAnalyzer analyzes network frames and provides insights
type FrameAnalyzer struct {
	securityAnalyzer *SecurityAnalyzer
	qosAnalyzer      *QoSAnalyzer
}

// NewFrameAnalyzer creates a new frame analyzer
func NewFrameAnalyzer() *FrameAnalyzer {
	return &FrameAnalyzer{
		securityAnalyzer: NewSecurityAnalyzer(),
		qosAnalyzer:      NewQoSAnalyzer(),
	}
}

// AnalyzeFrame performs analysis on a frame to provide insights
func (fa *FrameAnalyzer) AnalyzeFrame(frame *models.Frame) {
	// Initialize analysis results if needed
	if frame.AnalysisResults == nil {
		frame.AnalysisResults = make(map[string]interface{})
	}

	// Add basic frame type information
	switch frame.FrameType {
	case models.EthernetFrame:
		fa.analyzeEthernetFrame(frame)
	case models.WLANManagementFrame:
		fa.analyzeWLANManagementFrame(frame)
	case models.WLANControlFrame:
		fa.analyzeWLANControlFrame(frame)
	case models.WLANDataFrame:
		fa.analyzeWLANDataFrame(frame)
	}

	// Analyze security if present
	if frame.Security != nil {
		fa.securityAnalyzer.AnalyzeSecurity(frame)
	}

	// Analyze QoS if present
	if frame.QoS != nil {
		fa.qosAnalyzer.AnalyzeQoS(frame)
	}
}

// analyzeEthernetFrame provides analysis for Ethernet frames
func (fa *FrameAnalyzer) analyzeEthernetFrame(frame *models.Frame) {
	// Add Ethernet-specific analysis
	etherTypeDescription := getEtherTypeDescription(frame.EtherType)

	frame.AnalysisResults["Summary"] = fmt.Sprintf("Ethernet frame: %s", etherTypeDescription)
	frame.AnalysisResults["Details"] = map[string]interface{}{
		"EtherType":      frame.EtherType,
		"EtherTypeDesc":  etherTypeDescription,
		"SourceMAC":      frame.SourceMAC,
		"DestinationMAC": frame.DestinationMAC,
	}

	// Add VLAN info if present
	if frame.VLANInfo != nil {
		frame.AnalysisResults["VLANInfo"] = frame.VLANInfo
	}
}

// analyzeWLANManagementFrame provides analysis for WLAN management frames
func (fa *FrameAnalyzer) analyzeWLANManagementFrame(frame *models.Frame) {
	// Extract frame control info
	frameControl, ok := frame.FrameControl.(map[string]interface{})
	if !ok {
		return
	}

	subtype, _ := frameControl["Subtype"].(uint16)

	// Get management frame type
	var frameTypeStr string
	if managementInfo, ok := frame.AnalysisResults["ManagementInfo"].(map[string]interface{}); ok {
		if typeStr, ok := managementInfo["Type"].(string); ok {
			frameTypeStr = typeStr
		}
	}

	// Add analysis based on frame subtype
	frame.AnalysisResults["Summary"] = fmt.Sprintf("WLAN Management Frame: %s", frameTypeStr)

	// Add more detailed analysis based on the specific management frame type
	switch subtype {
	case 0: // Association Request
		frame.AnalysisResults["Context"] = "A station is requesting association with an access point"
	case 1: // Association Response
		frame.AnalysisResults["Context"] = "An access point is responding to an association request"
	case 4: // Probe Request
		frame.AnalysisResults["Context"] = "A station is actively scanning for access points"
	case 5: // Probe Response
		frame.AnalysisResults["Context"] = "An access point is responding to a probe request"
	case 8: // Beacon
		frame.AnalysisResults["Context"] = "An access point is broadcasting its presence and capabilities"
	case 10: // Disassociation
		frame.AnalysisResults["Context"] = "A station or AP is terminating an association"
	case 11: // Authentication
		frame.AnalysisResults["Context"] = "A station is attempting to authenticate with an access point"
	case 12: // Deauthentication
		frame.AnalysisResults["Context"] = "A station or AP is terminating authentication"
	}
}

// analyzeWLANControlFrame provides analysis for WLAN control frames
func (fa *FrameAnalyzer) analyzeWLANControlFrame(frame *models.Frame) {
	// Extract frame control info
	frameControl, ok := frame.FrameControl.(map[string]interface{})
	if !ok {
		return
	}

	subtype, _ := frameControl["Subtype"].(uint16)

	// Add control frame specific analysis
	var subtypeStr string
	switch subtype {
	case 8: // Block ACK Request
		subtypeStr = "Block ACK Request"
		frame.AnalysisResults["Context"] = "Station is requesting block acknowledgment for multiple frames"
	case 9: // Block ACK
		subtypeStr = "Block ACK"
		frame.AnalysisResults["Context"] = "Station is acknowledging receipt of multiple frames"
	case 10: // PS-Poll
		subtypeStr = "PS-Poll"
		frame.AnalysisResults["Context"] = "Power-save station is requesting buffered frames from access point"
	case 11: // RTS
		subtypeStr = "RTS (Request to Send)"
		frame.AnalysisResults["Context"] = "Station is initiating RTS/CTS mechanism to reserve medium"
	case 12: // CTS
		subtypeStr = "CTS (Clear to Send)"
		frame.AnalysisResults["Context"] = "Station is responding to RTS, clearing sender to transmit"
	case 13: // ACK
		subtypeStr = "ACK"
		frame.AnalysisResults["Context"] = "Station is acknowledging receipt of a frame"
	case 14: // CF-End
		subtypeStr = "CF-End"
		frame.AnalysisResults["Context"] = "Access point is indicating end of contention-free period"
	case 15: // CF-End + CF-ACK
		subtypeStr = "CF-End + CF-ACK"
		frame.AnalysisResults["Context"] = "Access point is indicating end of contention-free period with ACK"
	default:
		subtypeStr = fmt.Sprintf("Unknown control frame (%d)", subtype)
	}

	frame.AnalysisResults["Summary"] = fmt.Sprintf("WLAN Control Frame: %s", subtypeStr)
}

// analyzeWLANDataFrame provides analysis for WLAN data frames
func (fa *FrameAnalyzer) analyzeWLANDataFrame(frame *models.Frame) {
	// Extract frame control info
	frameControl, ok := frame.FrameControl.(map[string]interface{})
	if !ok {
		return
	}

	subtype, _ := frameControl["Subtype"].(uint16)
	toDS, _ := frameControl["ToDS"].(bool)
	fromDS, _ := frameControl["FromDS"].(bool)

	// Determine direction
	var direction string
	if toDS && !fromDS {
		direction = "Station to Distribution System"
	} else if !toDS && fromDS {
		direction = "Distribution System to Station"
	} else if toDS && fromDS {
		direction = "Distribution System to Distribution System"
	} else {
		direction = "Station to Station (Ad-Hoc)"
	}

	// Add data frame specific analysis
	var subtypeStr string
	switch subtype {
	case 0: // Data
		subtypeStr = "Data"
	case 4: // Null function (no data)
		subtypeStr = "Null function (no data)"
		frame.AnalysisResults["Context"] = "Station is informing AP of its power state without sending data"
	case 8: // QoS Data
		subtypeStr = "QoS Data"
		frame.AnalysisResults["Context"] = "Data frame with QoS prioritization"
	case 12: // QoS Null function (no data)
		subtypeStr = "QoS Null function (no data)"
		frame.AnalysisResults["Context"] = "QoS station is informing AP of its power state without sending data"
	default:
		subtypeStr = fmt.Sprintf("Unknown data frame (%d)", subtype)
	}

	frame.AnalysisResults["Summary"] = fmt.Sprintf("WLAN Data Frame: %s, Direction: %s", subtypeStr, direction)
	frame.AnalysisResults["Direction"] = direction
}

// getEtherTypeDescription returns a description of the Ethertype
func getEtherTypeDescription(etherType uint16) string {
	switch etherType {
	case 0x0800:
		return "IPv4"
	case 0x0806:
		return "ARP"
	case 0x8100:
		return "VLAN-tagged frame (IEEE 802.1Q)"
	case 0x86DD:
		return "IPv6"
	case 0x8863:
		return "PPPoE Discovery"
	case 0x8864:
		return "PPPoE Session"
	case 0x888E:
		return "802.1X Authentication"
	case 0x8035:
		return "RARP"
	case 0x8847:
		return "MPLS unicast"
	case 0x8848:
		return "MPLS multicast"
	default:
		return fmt.Sprintf("Unknown (0x%04X)", etherType)
	}
}
