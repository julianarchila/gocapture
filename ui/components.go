package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/julianarchila/gocapture/internal/storage"
	"github.com/julianarchila/gocapture/pkg/models"
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
			"Iniciar Captura",
			"Cargar Captura",
			"Salir",
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

	sb.WriteString("🔍 Menú Principal\n\n")

	for i, option := range m.options {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		sb.WriteString(fmt.Sprintf("%s %s\n", cursor, option))
	}

	sb.WriteString("\nUse las teclas de flecha para navegar, Enter para seleccionar\n")

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

	sb.WriteString("📋 Lista de Tramas\n\n")

	// Show total frame count
	sb.WriteString(fmt.Sprintf("Total de tramas: %d\n\n", len(m.frames)))

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
				summary = "Trama Ethernet"
			case models.WLANManagementFrame:
				summary = "Trama de Gestión WLAN"
			case models.WLANControlFrame:
				summary = "Trama de Control WLAN"
			case models.WLANDataFrame:
				summary = "Trama de Datos WLAN"
			default:
				summary = "Tipo de Trama Desconocido"
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
		sb.WriteString(fmt.Sprintf("\nMostrando %d-%d de %d tramas\n", m.offset+1, end, len(m.frames)))
	}

	sb.WriteString("\nUse las teclas de flecha para navegar, Enter para ver detalles de la trama\n")

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

	sb.WriteString("🔍 Detalles de Trama\n\n")

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

	sb.WriteString("\nUse Tab para cambiar modos de vista, flechas para desplazar\n")
	sb.WriteString("Press n/right for next frame, p/left for previous frame\n")

	return sb.String()
}

// renderSummaryView renders the summary view of the frame
func (m *frameDetailModel) renderSummaryView(sb *strings.Builder) {
	frame := m.frame

	sb.WriteString(fmt.Sprintf("ID: %d\n", frame.ID))
	sb.WriteString(fmt.Sprintf("Tiempo: %s\n", frame.Timestamp.Format("2006-01-02 15:04:05.000")))
	sb.WriteString(fmt.Sprintf("Longitud: %d bytes\n\n", frame.Length))

	sb.WriteString("Direcciones:\n")
	sb.WriteString(fmt.Sprintf("  Origen: %s\n", frame.SourceMAC))
	sb.WriteString(fmt.Sprintf("  Destino: %s\n\n", frame.DestinationMAC))

	sb.WriteString("Tipo de Trama: ")
	switch frame.FrameType {
	case models.EthernetFrame:
		sb.WriteString("Ethernet\n")
	case models.WLANManagementFrame:
		sb.WriteString("Gestión WLAN\n")
	case models.WLANControlFrame:
		sb.WriteString("Control WLAN\n")
	case models.WLANDataFrame:
		sb.WriteString("Datos WLAN\n")
	default:
		sb.WriteString("Desconocido\n")
	}

	// Show analysis results
	if len(frame.AnalysisResults) > 0 {
		sb.WriteString("\nAnálisis:\n")
		for key, value := range frame.AnalysisResults {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}
}

// renderDetailsView renders the detailed view of the frame
func (m *frameDetailModel) renderDetailsView(sb *strings.Builder) {
	frame := m.frame

	sb.WriteString("Información Detallada de la Trama\n\n")

	// Basic information
	sb.WriteString("Información Básica:\n")
	sb.WriteString(fmt.Sprintf("  ID: %d\n", frame.ID))
	sb.WriteString(fmt.Sprintf("  Tiempo: %s\n", frame.Timestamp.Format("2006-01-02 15:04:05.000")))
	sb.WriteString(fmt.Sprintf("  Longitud: %d bytes\n", frame.Length))

	// MAC addresses
	sb.WriteString("\nDirecciones MAC:\n")
	sb.WriteString(fmt.Sprintf("  Origen: %s\n", frame.SourceMAC))
	sb.WriteString(fmt.Sprintf("  Destino: %s\n", frame.DestinationMAC))

	// Frame type specific information
	sb.WriteString("\nInformación Específica del Tipo:\n")
	switch frame.FrameType {
	case models.EthernetFrame:
		sb.WriteString("  Tipo: Ethernet\n")
		if frame.EtherType != 0 {
			sb.WriteString(fmt.Sprintf("  EtherType: 0x%04x\n", frame.EtherType))
		}
	case models.WLANManagementFrame:
		sb.WriteString("  Tipo: Gestión WLAN\n")
		if frame.FrameControl != nil {
			sb.WriteString(fmt.Sprintf("  Control de Trama: %v\n", frame.FrameControl))
		}
	case models.WLANControlFrame:
		sb.WriteString("  Tipo: Control WLAN\n")
		if frame.FrameControl != nil {
			sb.WriteString(fmt.Sprintf("  Control de Trama: %v\n", frame.FrameControl))
		}
	case models.WLANDataFrame:
		sb.WriteString("  Tipo: Datos WLAN\n")
		if frame.FrameControl != nil {
			sb.WriteString(fmt.Sprintf("  Control de Trama: %v\n", frame.FrameControl))
		}
	}

	// Analysis results
	if len(frame.AnalysisResults) > 0 {
		sb.WriteString("\nResultados del Análisis:\n")
		for key, value := range frame.AnalysisResults {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}
}

// renderHexView renders the hex dump view of the frame
func (m *frameDetailModel) renderHexView(sb *strings.Builder) {
	frame := m.frame

	sb.WriteString("Vista Hexadecimal de la Trama\n\n")

	// Format the hex dump
	for i := 0; i < len(frame.RawData); i += 16 {
		// Print offset
		sb.WriteString(fmt.Sprintf("%08x  ", i))

		// Print hex values
		for j := 0; j < 16; j++ {
			if i+j < len(frame.RawData) {
				sb.WriteString(fmt.Sprintf("%02x ", frame.RawData[i+j]))
			} else {
				sb.WriteString("   ")
			}
			if j == 7 {
				sb.WriteString(" ")
			}
		}

		// Print ASCII representation
		sb.WriteString(" |")
		for j := 0; j < 16; j++ {
			if i+j < len(frame.RawData) {
				b := frame.RawData[i+j]
				if b >= 32 && b <= 126 {
					sb.WriteString(fmt.Sprintf("%c", b))
				} else {
					sb.WriteString(".")
				}
			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("|\n")
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

// View renders the saved captures screen
func (m *savedCapturesModel) View() string {
	var sb strings.Builder

	sb.WriteString("💾 Capturas Guardadas\n\n")

	if m.loading {
		sb.WriteString("Cargando capturas...\n")
		return sb.String()
	}

	if len(m.captures) == 0 {
		sb.WriteString("No hay capturas guardadas\n")
		return sb.String()
	}

	// Show captures
	for i, capture := range m.captures {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		sb.WriteString(fmt.Sprintf("%s %s\n", cursor, capture.Filename))
		sb.WriteString(fmt.Sprintf("  Tiempo: %s\n", capture.StartTime.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("  Tramas: %d\n", capture.FrameCount))
	}

	sb.WriteString("\nUse las teclas de flecha para navegar, Enter para cargar\n")

	return sb.String()
}
