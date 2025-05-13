package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/julian/gocapture/internal/storage"
	"github.com/julian/gocapture/pkg/models"
)

// mainMenuModel represents the main menu UI component
type mainMenuModel struct {
	options []string
	cursor  int
}

// newMainMenuModel creates a new main menu model
func newMainMenuModel() *mainMenuModel {
	return &mainMenuModel{
		options: []string{
			"Start Capture",
			"Load Capture",
			"Quit",
		},
		cursor: 0,
	}
}

// Init initializes the main menu model
func (m *mainMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles updates to the main menu model
func (m *mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			// Return a menu selected message
			return m, func() tea.Msg {
				return menuSelectedMsg{option: m.options[m.cursor]}
			}
		}
	}

	return m, nil
}

// View renders the main menu
func (m *mainMenuModel) View() string {
	var sb strings.Builder

	sb.WriteString("🔍 Main Menu\n\n")

	for i, option := range m.options {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		sb.WriteString(fmt.Sprintf("%s %s\n", cursor, option))
	}

	sb.WriteString("\nUse arrow keys to navigate, Enter to select\n")

	return sb.String()
}

// frameListModel represents the frame list UI component
type frameListModel struct {
	frames   []*models.Frame
	cursor   int
	offset   int
	pageSize int
}

// newFrameListModel creates a new frame list model
func newFrameListModel() *frameListModel {
	return &frameListModel{
		frames:   make([]*models.Frame, 0),
		cursor:   0,
		offset:   0,
		pageSize: 10,
	}
}

// setFrames sets the frames to display in the list
func (m *frameListModel) setFrames(frames []*models.Frame) {
	m.frames = frames
	m.cursor = 0
	m.offset = 0
}

// Init initializes the frame list model
func (m *frameListModel) Init() tea.Cmd {
	return nil
}

// Update handles updates to the frame list model
func (m *frameListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.frames)-1 {
				m.cursor++
				if m.cursor >= m.offset+m.pageSize {
					m.offset = m.cursor - m.pageSize + 1
				}
			}
		case "pgup":
			m.cursor -= m.pageSize
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.offset = m.cursor
		case "pgdown":
			m.cursor += m.pageSize
			if m.cursor >= len(m.frames) {
				m.cursor = len(m.frames) - 1
			}
			if m.cursor >= m.offset+m.pageSize {
				m.offset = m.cursor - m.pageSize + 1
			}
		}
	}

	return m, nil
}

// View renders the frame list
func (m *frameListModel) View() string {
	var sb strings.Builder

	sb.WriteString("📋 Frame List\n\n")

	// Show total frame count
	sb.WriteString(fmt.Sprintf("Total frames: %d\n\n", len(m.frames)))

	// Calculate end index for pagination
	end := m.offset + m.pageSize
	if end > len(m.frames) {
		end = len(m.frames)
	}

	// Show frames
	for i := m.offset; i < end; i++ {
		frame := m.frames[i]

		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		// Get frame summary
		var summary string
		if summaryValue, ok := frame.AnalysisResults["Summary"]; ok {
			summary = fmt.Sprintf("%v", summaryValue)
		} else {
			switch frame.FrameType {
			case models.EthernetFrame:
				summary = "Ethernet Frame"
			case models.WLANManagementFrame:
				summary = "WLAN Management Frame"
			case models.WLANControlFrame:
				summary = "WLAN Control Frame"
			case models.WLANDataFrame:
				summary = "WLAN Data Frame"
			default:
				summary = "Unknown Frame Type"
			}
		}

		// Format the line
		sb.WriteString(fmt.Sprintf("%s #%d [%s] %s → %s (%d bytes)\n",
			cursor,
			frame.ID,
			frame.Timestamp.Format("15:04:05.000"),
			frame.SourceMAC,
			frame.DestinationMAC,
			frame.Length,
		))
		sb.WriteString(fmt.Sprintf("  %s\n", summary))
	}

	// Show pagination info
	if len(m.frames) > m.pageSize {
		sb.WriteString(fmt.Sprintf("\nShowing %d-%d of %d frames\n", m.offset+1, end, len(m.frames)))
	}

	sb.WriteString("\nUse arrow keys to navigate, Enter to view frame details\n")

	return sb.String()
}

// frameDetailModel represents the frame detail UI component
type frameDetailModel struct {
	frame     *models.Frame
	scrollPos int
	viewMode  int // 0 = summary, 1 = details, 2 = raw hex
}

// newFrameDetailModel creates a new frame detail model
func newFrameDetailModel() *frameDetailModel {
	return &frameDetailModel{
		frame:     nil,
		scrollPos: 0,
		viewMode:  0,
	}
}

// setFrame sets the frame to display
func (m *frameDetailModel) setFrame(frame *models.Frame) {
	m.frame = frame
	m.scrollPos = 0
}

// Init initializes the frame detail model
func (m *frameDetailModel) Init() tea.Cmd {
	return nil
}

// Update handles updates to the frame detail model
func (m *frameDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.scrollPos > 0 {
				m.scrollPos--
			}
		case "down", "j":
			m.scrollPos++
		case "tab":
			// Cycle through view modes
			m.viewMode = (m.viewMode + 1) % 3
			m.scrollPos = 0
		}
	}

	return m, nil
}

// View renders the frame detail
func (m *frameDetailModel) View() string {
	var sb strings.Builder

	if m.frame == nil {
		return "No frame selected"
	}

	sb.WriteString("🔍 Frame Details\n\n")

	// Show frame ID and basic info
	sb.WriteString(fmt.Sprintf("Frame #%d - Captured at %s\n",
		m.frame.ID,
		m.frame.Timestamp.Format("2006-01-02 15:04:05.000"),
	))

	sb.WriteString(fmt.Sprintf("Size: %d bytes\n\n", m.frame.Length))

	// Show view mode tabs
	sb.WriteString("[ ")
	if m.viewMode == 0 {
		sb.WriteString("Summary")
	} else {
		sb.WriteString("summary")
	}
	sb.WriteString(" | ")
	if m.viewMode == 1 {
		sb.WriteString("Details")
	} else {
		sb.WriteString("details")
	}
	sb.WriteString(" | ")
	if m.viewMode == 2 {
		sb.WriteString("Hex Dump")
	} else {
		sb.WriteString("hex dump")
	}
	sb.WriteString(" ]\n\n")

	// Show content based on view mode
	switch m.viewMode {
	case 0:
		// Summary view
		m.renderSummaryView(&sb)
	case 1:
		// Details view
		m.renderDetailsView(&sb)
	case 2:
		// Raw hex view
		m.renderHexView(&sb)
	}

	sb.WriteString("\nUse arrow keys to scroll, Tab to change view, ESC to go back to list\n")
	sb.WriteString("Press n/right for next frame, p/left for previous frame\n")

	return sb.String()
}

// renderSummaryView renders the summary view of the frame
func (m *frameDetailModel) renderSummaryView(sb *strings.Builder) {
	// Frame type
	var frameTypeStr string
	switch m.frame.FrameType {
	case models.EthernetFrame:
		frameTypeStr = "Ethernet Frame (IEEE 802.3)"
	case models.WLANManagementFrame:
		frameTypeStr = "WLAN Management Frame (IEEE 802.11)"
	case models.WLANControlFrame:
		frameTypeStr = "WLAN Control Frame (IEEE 802.11)"
	case models.WLANDataFrame:
		frameTypeStr = "WLAN Data Frame (IEEE 802.11)"
	default:
		frameTypeStr = "Unknown Frame Type"
	}

	sb.WriteString(fmt.Sprintf("Frame Type: %s\n", frameTypeStr))
	sb.WriteString(fmt.Sprintf("Source MAC: %s\n", m.frame.SourceMAC))
	sb.WriteString(fmt.Sprintf("Destination MAC: %s\n", m.frame.DestinationMAC))

	// Show summary from analysis results
	if summary, ok := m.frame.AnalysisResults["Summary"]; ok {
		sb.WriteString(fmt.Sprintf("\nSummary: %v\n", summary))
	}

	// Show context from analysis results
	if context, ok := m.frame.AnalysisResults["Context"]; ok {
		sb.WriteString(fmt.Sprintf("\nContext: %v\n", context))
	}

	// Show security info if available
	if security, ok := m.frame.AnalysisResults["Security"].(map[string]interface{}); ok {
		sb.WriteString("\n--- Security Information ---\n")

		if encType, ok := security["EncryptionType"].(string); ok {
			sb.WriteString(fmt.Sprintf("Encryption: %s\n", encType))
		}

		if level, ok := security["SecurityLevel"].(string); ok {
			sb.WriteString(fmt.Sprintf("Security Level: %s\n", level))
		}

		if context, ok := security["Context"].(string); ok {
			sb.WriteString(fmt.Sprintf("Security Context: %s\n", context))
		}

		if warning, ok := security["Warning"].(string); ok {
			sb.WriteString(fmt.Sprintf("Warning: %s\n", warning))
		}
	}

	// Show QoS info if available
	if qos, ok := m.frame.AnalysisResults["QoS"].(map[string]interface{}); ok {
		sb.WriteString("\n--- QoS Information ---\n")

		if priority, ok := qos["Priority"].(int); ok {
			sb.WriteString(fmt.Sprintf("Priority: %d\n", priority))
		}

		if trafficType, ok := qos["TrafficType"].(string); ok {
			sb.WriteString(fmt.Sprintf("Traffic Type: %s\n", trafficType))
		}

		if explanation, ok := qos["Explanation"].(string); ok {
			sb.WriteString(fmt.Sprintf("Explanation: %s\n", explanation))
		}
	}
}

// renderDetailsView renders the detailed view of the frame
func (m *frameDetailModel) renderDetailsView(sb *strings.Builder) {
	// Show frame control info for WLAN frames
	if m.frame.FrameType == models.WLANManagementFrame ||
		m.frame.FrameType == models.WLANControlFrame ||
		m.frame.FrameType == models.WLANDataFrame {

		sb.WriteString("--- Frame Control Field ---\n")

		if fc, ok := m.frame.FrameControl.(map[string]interface{}); ok {
			for k, v := range fc {
				sb.WriteString(fmt.Sprintf("%s: %v\n", k, v))
			}
		}

		sb.WriteString(fmt.Sprintf("\nDuration: %d\n", m.frame.Duration))
		sb.WriteString(fmt.Sprintf("Sequence Control: %d\n", m.frame.SequenceControl))

		sb.WriteString("\n--- Address Fields ---\n")
		sb.WriteString(fmt.Sprintf("Address 1: %s\n", m.frame.Address1))
		sb.WriteString(fmt.Sprintf("Address 2: %s\n", m.frame.Address2))
		sb.WriteString(fmt.Sprintf("Address 3: %s\n", m.frame.Address3))

		if m.frame.Address4 != "" {
			sb.WriteString(fmt.Sprintf("Address 4: %s\n", m.frame.Address4))
		}
	}

	// Show EtherType for Ethernet frames
	if m.frame.FrameType == models.EthernetFrame {
		sb.WriteString(fmt.Sprintf("EtherType: 0x%04X\n", m.frame.EtherType))

		// Show VLAN info if present
		if m.frame.VLANInfo != nil {
			sb.WriteString("\n--- VLAN Information ---\n")
			if vlanInfo, ok := m.frame.VLANInfo.(map[string]interface{}); ok {
				for k, v := range vlanInfo {
					sb.WriteString(fmt.Sprintf("%s: %v\n", k, v))
				}
			}
		}
	}

	// Show security details if available
	if m.frame.Security != nil {
		sb.WriteString("\n--- Security Details ---\n")
		sb.WriteString(fmt.Sprintf("Encryption Type: %s\n", m.frame.Security.EncryptionType))

		if len(m.frame.Security.Details) > 0 {
			for k, v := range m.frame.Security.Details {
				sb.WriteString(fmt.Sprintf("%s: %v\n", k, v))
			}
		}
	}

	// Show QoS details if available
	if m.frame.QoS != nil {
		sb.WriteString("\n--- QoS Details ---\n")
		sb.WriteString(fmt.Sprintf("TID: %d\n", m.frame.QoS.TID))
		sb.WriteString(fmt.Sprintf("Priority: %d\n", m.frame.QoS.Priority))
		sb.WriteString(fmt.Sprintf("ACK Policy: %s\n", m.frame.QoS.ACKPolicy))
		sb.WriteString(fmt.Sprintf("TXOP: %d\n", m.frame.QoS.TXOP))

		if len(m.frame.QoS.Details) > 0 {
			for k, v := range m.frame.QoS.Details {
				sb.WriteString(fmt.Sprintf("%s: %v\n", k, v))
			}
		}
	}

	// Show all analysis results
	sb.WriteString("\n--- All Analysis Results ---\n")
	for k, v := range m.frame.AnalysisResults {
		if k != "Summary" && k != "Context" && k != "Security" && k != "QoS" {
			sb.WriteString(fmt.Sprintf("%s: %v\n", k, v))
		}
	}
}

// renderHexView renders the hex dump view of the frame
func (m *frameDetailModel) renderHexView(sb *strings.Builder) {
	rawData := m.frame.RawData
	if len(rawData) == 0 {
		sb.WriteString("No raw data available")
		return
	}

	// Calculate visible lines based on scroll position
	linesPerPage := 16
	bytesPerLine := 16
	totalLines := (len(rawData) + bytesPerLine - 1) / bytesPerLine

	sb.WriteString(fmt.Sprintf("Showing %d of %d bytes (scroll to see more)\n\n", len(rawData), len(rawData)))
	sb.WriteString("       | 00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F | ASCII\n")
	sb.WriteString("-------+------------------------------------------------+----------------\n")

	// Show hex dump
	for i := m.scrollPos; i < m.scrollPos+linesPerPage && i < totalLines; i++ {
		// Address
		sb.WriteString(fmt.Sprintf(" %04X  | ", i*bytesPerLine))

		// Hex values
		for j := 0; j < bytesPerLine; j++ {
			idx := i*bytesPerLine + j
			if idx < len(rawData) {
				sb.WriteString(fmt.Sprintf("%02X ", rawData[idx]))
			} else {
				sb.WriteString("   ")
			}
		}

		sb.WriteString("| ")

		// ASCII representation
		for j := 0; j < bytesPerLine; j++ {
			idx := i*bytesPerLine + j
			if idx < len(rawData) {
				if rawData[idx] >= 32 && rawData[idx] <= 126 {
					sb.WriteString(string(rawData[idx]))
				} else {
					sb.WriteString(".")
				}
			} else {
				sb.WriteString(" ")
			}
		}

		sb.WriteString("\n")
	}
}

// savedCapturesModel represents the saved captures UI component
type savedCapturesModel struct {
	storageManager *storage.StorageManager
	captures       []*storage.SaveMetadata
	cursor         int
	loading        bool
}

// newSavedCapturesModel creates a new saved captures model
func newSavedCapturesModel(storageManager *storage.StorageManager) *savedCapturesModel {
	return &savedCapturesModel{
		storageManager: storageManager,
		captures:       make([]*storage.SaveMetadata, 0),
		cursor:         0,
		loading:        false,
	}
}

// Init initializes the saved captures model
func (m *savedCapturesModel) Init() tea.Cmd {
	return nil
}

// loadSavedCaptures loads the list of saved captures
func (m *savedCapturesModel) loadSavedCaptures() tea.Cmd {
	return func() tea.Msg {
		m.loading = true

		captures, err := m.storageManager.ListSavedCaptures()
		if err != nil {
			return nil
		}

		m.captures = captures
		m.loading = false

		// Return a command to update the UI
		return tea.Batch(
			func() tea.Msg {
				return capturesLoadedMsg{captures: captures}
			},
		)()
	}
}

// capturesLoadedMsg is a message sent when captures are loaded
type capturesLoadedMsg struct {
	captures []*storage.SaveMetadata
}

// Update handles updates to the saved captures model
func (m *savedCapturesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case capturesLoadedMsg:
		m.captures = msg.captures
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.captures)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.captures) > 0 {
				return m, m.loadCapture(m.captures[m.cursor].Filename)
			}
		}
	}

	return m, nil
}

// loadCapture loads a captured file
func (m *savedCapturesModel) loadCapture(filename string) tea.Cmd {
	return func() tea.Msg {
		frames, _, err := m.storageManager.LoadFrames(filename)
		if err != nil {
			return nil
		}

		return loadFramesMsg{frames: frames}
	}
}

// View renders the saved captures list
func (m *savedCapturesModel) View() string {
	var sb strings.Builder

	sb.WriteString("💾 Saved Captures\n\n")

	if m.loading {
		sb.WriteString("Loading captures...\n")
		return sb.String()
	}

	if len(m.captures) == 0 {
		sb.WriteString("No saved captures found\n")
		return sb.String()
	}

	// Show captures
	for i, capture := range m.captures {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		sb.WriteString(fmt.Sprintf("%s %s - %s\n",
			cursor,
			capture.Filename,
			capture.StartTime.Format("2006-01-02 15:04:05"),
		))

		sb.WriteString(fmt.Sprintf("  Interface: %s, Frames: %d\n",
			capture.Interface,
			capture.FrameCount,
		))

		if capture.Description != "" {
			sb.WriteString(fmt.Sprintf("  Description: %s\n", capture.Description))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("\nUse arrow keys to navigate, Enter to load capture\n")

	return sb.String()
}
