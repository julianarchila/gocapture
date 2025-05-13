package analyzer

import (
	"fmt"

	"github.com/julian/gocapture/pkg/models"
)

// SecurityAnalyzer analyzes security aspects of network frames
type SecurityAnalyzer struct{}

// NewSecurityAnalyzer creates a new security analyzer
func NewSecurityAnalyzer() *SecurityAnalyzer {
	return &SecurityAnalyzer{}
}

// AnalyzeSecurity analyzes the security aspects of a frame
func (sa *SecurityAnalyzer) AnalyzeSecurity(frame *models.Frame) {
	if frame.Security == nil {
		return
	}

	// Add security info to analysis results
	if frame.AnalysisResults == nil {
		frame.AnalysisResults = make(map[string]interface{})
	}

	securityInfo := make(map[string]interface{})
	securityInfo["EncryptionType"] = frame.Security.EncryptionType

	// Add different analysis based on encryption type
	switch frame.Security.EncryptionType {
	case "WEP":
		securityInfo["SecurityLevel"] = "Low"
		securityInfo["Vulnerabilities"] = []string{
			"RC4 encryption easily broken",
			"IV reuse vulnerability",
			"Susceptible to packet forgery",
		}
		securityInfo["Recommendation"] = "WEP is deprecated and insecure. Upgrade to WPA2 or WPA3."
		securityInfo["Context"] = "WEP uses a 24-bit initialization vector (IV) with RC4 stream cipher, which is vulnerable to statistical attacks."

	case "TKIP (WPA)":
		securityInfo["SecurityLevel"] = "Medium"
		securityInfo["Vulnerabilities"] = []string{
			"Michael MIC attack",
			"TKIP uses RC4 which has weaknesses",
		}
		securityInfo["Recommendation"] = "TKIP/WPA is outdated. Upgrade to WPA2 or WPA3."
		securityInfo["Context"] = "TKIP was designed as a stopgap replacement for WEP, improving security while maintaining hardware compatibility."

	case "CCMP (WPA2)":
		securityInfo["SecurityLevel"] = "High"
		securityInfo["Vulnerabilities"] = []string{
			"KRACK attack (fixed in most modern devices)",
		}
		securityInfo["Recommendation"] = "Ensure device firmware is updated to patch against KRACK vulnerabilities."
		securityInfo["Context"] = "CCMP uses AES and provides strong encryption considered secure for most applications."

	case "GCMP (WPA3)":
		securityInfo["SecurityLevel"] = "Very High"
		securityInfo["Vulnerabilities"] = []string{
			"Potential side-channel attacks (theoretical)",
		}
		securityInfo["Recommendation"] = "Currently the most secure option for WiFi."
		securityInfo["Context"] = "WPA3 provides stronger encryption, protection against brute force attacks, and forward secrecy."

	default:
		if frame.FrameType == models.WLANManagementFrame || frame.FrameType == models.WLANControlFrame || frame.FrameType == models.WLANDataFrame {
			// Add info about unencrypted WLAN frames
			frameControl, ok := frame.FrameControl.(map[string]interface{})
			if ok {
				protected, _ := frameControl["Protected"].(bool)
				if !protected {
					securityInfo["SecurityLevel"] = "None"
					securityInfo["Warning"] = "Unencrypted wireless frame"
					
					// Add context based on frame type
					if frame.FrameType == models.WLANManagementFrame {
						securityInfo["Context"] = "Management frames are typically unencrypted in older standards, but newer standards support Protected Management Frames (PMF)."
					} else if frame.FrameType == models.WLANDataFrame {
						securityInfo["Context"] = "Unencrypted data frames expose payload contents to anyone monitoring the channel."
						securityInfo["Recommendation"] = "Configure wireless network to use encryption."
					}
				}
			}
		}
	}

	// Add any details provided by the parser
	for k, v := range frame.Security.Details {
		securityInfo[k] = v
	}

	frame.AnalysisResults["Security"] = securityInfo
	
	// Update the summary to include security info
	if summary, ok := frame.AnalysisResults["Summary"].(string); ok {
		frame.AnalysisResults["Summary"] = fmt.Sprintf("%s [%s]", summary, frame.Security.EncryptionType)
	}
} 