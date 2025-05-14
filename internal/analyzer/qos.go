package analyzer

import (
	"fmt"

	"github.com/julianarchila/gocapture/pkg/models"
)

// QoSAnalyzer analyzes Quality of Service aspects of network frames
type QoSAnalyzer struct{}

// NewQoSAnalyzer creates a new QoS analyzer
func NewQoSAnalyzer() *QoSAnalyzer {
	return &QoSAnalyzer{}
}

// AnalyzeQoS analyzes the QoS aspects of a frame
func (qa *QoSAnalyzer) AnalyzeQoS(frame *models.Frame) {
	if frame.QoS == nil {
		return
	}

	// Add QoS info to analysis results
	if frame.AnalysisResults == nil {
		frame.AnalysisResults = make(map[string]interface{})
	}

	qosInfo := make(map[string]interface{})

	// Add basic QoS info
	qosInfo["Priority"] = frame.QoS.Priority
	qosInfo["TID"] = frame.QoS.TID
	qosInfo["ACKPolicy"] = frame.QoS.ACKPolicy

	// Add TID description
	qosInfo["TrafficType"] = getTrafficTypeByTID(frame.QoS.TID)

	// Add detailed explanations
	switch frame.QoS.ACKPolicy {
	case "Normal ACK":
		qosInfo["ACKContext"] = "Each frame requires acknowledgment, providing reliable delivery but with overhead."
	case "No ACK":
		qosInfo["ACKContext"] = "Frames are not acknowledged. Used for time-sensitive traffic where retransmissions are not useful."
	case "No Explicit ACK":
		qosInfo["ACKContext"] = "Used in power-save mode. Acknowledgment is not immediate but delayed."
	case "Block ACK":
		qosInfo["ACKContext"] = "Multiple frames are acknowledged with a single Block ACK, improving efficiency."
	}

	// Add performance analysis based on priority
	qosInfo["Explanation"] = getQoSPriorityExplanation(frame.QoS.Priority)

	// Add recommended applications for this priority level
	qosInfo["RecommendedApplications"] = getRecommendedApplications(frame.QoS.Priority)

	// Add any details provided by the parser
	for k, v := range frame.QoS.Details {
		qosInfo[k] = v
	}

	// Add TXOP (Transmission Opportunity) analysis if present
	if frame.QoS.TXOP > 0 {
		txopDuration := frame.QoS.TXOP * 32 // Each unit is 32μs
		qosInfo["TXOPDuration"] = fmt.Sprintf("%d μs", txopDuration)
		qosInfo["TXOPContext"] = "Station is allowed to transmit multiple frames within this time window without contending for the medium again."
	}

	frame.AnalysisResults["QoS"] = qosInfo

	// Update the summary to include QoS info
	if summary, ok := frame.AnalysisResults["Summary"].(string); ok {
		trafficType := getTrafficTypeByTID(frame.QoS.TID)
		frame.AnalysisResults["Summary"] = fmt.Sprintf("%s [QoS: %s, Priority: %d]",
			summary, trafficType, frame.QoS.Priority)
	}
}

// getTrafficTypeByTID returns a human-readable description of a traffic ID
func getTrafficTypeByTID(tid int) string {
	switch tid {
	case 0, 3:
		return "Background"
	case 1, 2:
		return "Best Effort"
	case 4, 5:
		return "Video"
	case 6, 7:
		return "Voice"
	default:
		return fmt.Sprintf("Unknown (%d)", tid)
	}
}

// getQoSPriorityExplanation returns an explanation of what a QoS priority level means
func getQoSPriorityExplanation(priority int) string {
	switch priority {
	case 0:
		return "Lowest priority. Used for bulk transfers and background tasks that do not have strict latency requirements."
	case 1:
		return "Low priority. Used for best effort traffic like email and web browsing."
	case 2:
		return "Low-medium priority. Best effort traffic with slightly higher priority."
	case 3:
		return "Medium priority. Used for applications that require better than best effort but are not sensitive to latency."
	case 4:
		return "Medium-high priority. Used for streaming video that can tolerate some delay."
	case 5:
		return "High priority. Used for video applications with lower tolerance for delay."
	case 6:
		return "Very high priority. Used for voice applications that require low latency and jitter."
	case 7:
		return "Highest priority. Reserved for network control traffic."
	default:
		return "Unknown priority level."
	}
}

// getRecommendedApplications returns examples of applications suitable for a given priority level
func getRecommendedApplications(priority int) []string {
	switch priority {
	case 0:
		return []string{"File downloads", "Print jobs", "Backup operations"}
	case 1, 2:
		return []string{"Web browsing", "Email", "Social media", "Chat applications"}
	case 3:
		return []string{"ERP applications", "Database access", "Interactive applications"}
	case 4, 5:
		return []string{"Video streaming", "Video conferencing", "IPTV"}
	case 6:
		return []string{"VoIP", "Video conferencing audio", "Online gaming"}
	case 7:
		return []string{"Network control protocols", "WLAN management"}
	default:
		return []string{"Unknown"}
	}
}
