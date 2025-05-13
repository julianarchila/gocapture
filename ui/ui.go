package ui

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/julian/gocapture/internal/analyzer"
	"github.com/julian/gocapture/internal/capture"
	"github.com/julian/gocapture/internal/storage"
	"github.com/julian/gocapture/pkg/models"
)

// UI states
const (
	stateMainMenu = iota
	stateCapturing
	stateFrameList
	stateFrameDetail
	stateSavedCaptures
)

// MainModel is the main UI model
type MainModel struct {
	state          int
	captureEngine  *capture.CaptureEngine
	storageManager *storage.StorageManager
	frameAnalyzer  *analyzer.FrameAnalyzer

	// Captured frames
	frames        []*models.Frame
	selectedFrame int

	// UI components
	mainMenu      *mainMenuModel
	frameList     *frameListModel
	frameDetail   *frameDetailModel
	savedCaptures *savedCapturesModel

	// Error message
	err error
}

// StartUI initializes and starts the UI
func StartUI(captureEngine *capture.CaptureEngine) error {
	// Initialize the storage manager
	storageManager, err := storage.NewStorageManager("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage manager: %v", err)
	}

	// Initialize the frame analyzer
	frameAnalyzer := analyzer.NewFrameAnalyzer()

	// Initialize the main model
	model := &MainModel{
		state:          stateMainMenu,
		captureEngine:  captureEngine,
		storageManager: storageManager,
		frameAnalyzer:  frameAnalyzer,
		frames:         make([]*models.Frame, 0),
		selectedFrame:  0,
	}

	// Initialize UI components
	model.mainMenu = newMainMenuModel()
	model.frameList = newFrameListModel()
	model.frameDetail = newFrameDetailModel()
	model.savedCaptures = newSavedCapturesModel(storageManager)

	// Create and start the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	return nil
}

// Init initializes the model
func (m *MainModel) Init() tea.Cmd {
	// Start with the main menu
	return m.mainMenu.Init()
}

// Update handles user input and updates the model
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Always allow quitting with ctrl+c or q
			return m, tea.Quit
		}
	}

	// Route messages to the appropriate component based on the current state
	switch m.state {
	case stateMainMenu:
		newMainMenu, menuCmd := m.mainMenu.Update(msg)
		m.mainMenu = newMainMenu.(*mainMenuModel)
		cmds = append(cmds, menuCmd)

		// Handle menu selection
		if selectedCmd, ok := msg.(menuSelectedMsg); ok {
			switch selectedCmd.option {
			case "Start Capture":
				m.state = stateCapturing
				cmds = append(cmds, m.startCapturing())
			case "Load Capture":
				m.state = stateSavedCaptures
				cmds = append(cmds, m.savedCaptures.loadSavedCaptures())
			case "Quit":
				return m, tea.Quit
			}
		}

	case stateCapturing:
		// Handle capturing state
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				m.state = stateMainMenu
				cmds = append(cmds, m.stopCapturing())
			case "s":
				// Save the current capture
				metadata := &storage.SaveMetadata{
					Interface:   m.captureEngine.GetInterfaceName(),
					Description: "Captured with GoCapture",
				}
				if err := m.storageManager.SaveFrames(m.frames, metadata); err != nil {
					m.err = err
				} else {
					m.err = fmt.Errorf("Capture saved to %s", metadata.Filename)
				}
			case "enter":
				// Stop capturing and show frame list
				if len(m.frames) > 0 {
					m.state = stateFrameList
					cmds = append(cmds, m.stopCapturing())
					m.frameList.setFrames(m.frames)
				}
			}
		}

		// Check for new frames
		if newFrameMsg, ok := msg.(newFrameMsg); ok {
			m.frames = append(m.frames, newFrameMsg.frame)
			// Analyze the frame
			m.frameAnalyzer.AnalyzeFrame(newFrameMsg.frame)
			cmds = append(cmds, m.checkForMoreFrames())
		}

	case stateFrameList:
		// Update frame list
		newFrameList, frameListCmd := m.frameList.Update(msg)
		m.frameList = newFrameList.(*frameListModel)
		cmds = append(cmds, frameListCmd)

		// Handle key presses in frame list
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				m.state = stateMainMenu
			case "enter":
				if len(m.frames) > 0 {
					m.selectedFrame = m.frameList.cursor
					m.frameDetail.setFrame(m.frames[m.selectedFrame])
					m.state = stateFrameDetail
				}
			case "s":
				// Save the current capture
				metadata := &storage.SaveMetadata{
					Interface:   m.captureEngine.GetInterfaceName(),
					Description: "Captured with GoCapture",
				}
				if err := m.storageManager.SaveFrames(m.frames, metadata); err != nil {
					m.err = err
				} else {
					m.err = fmt.Errorf("Capture saved to %s", metadata.Filename)
				}
			}
		}

	case stateFrameDetail:
		// Update frame detail
		newFrameDetail, frameDetailCmd := m.frameDetail.Update(msg)
		m.frameDetail = newFrameDetail.(*frameDetailModel)
		cmds = append(cmds, frameDetailCmd)

		// Handle key presses in frame detail
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				m.state = stateFrameList
			case "right", "l", "n":
				// Next frame
				if m.selectedFrame < len(m.frames)-1 {
					m.selectedFrame++
					m.frameDetail.setFrame(m.frames[m.selectedFrame])
				}
			case "left", "h", "p":
				// Previous frame
				if m.selectedFrame > 0 {
					m.selectedFrame--
					m.frameDetail.setFrame(m.frames[m.selectedFrame])
				}
			}
		}

	case stateSavedCaptures:
		// Update saved captures list
		newSavedCaptures, savedCapturesCmd := m.savedCaptures.Update(msg)
		m.savedCaptures = newSavedCaptures.(*savedCapturesModel)
		cmds = append(cmds, savedCapturesCmd)

		// Handle saved capture selection
		if loadFramesMsg, ok := msg.(loadFramesMsg); ok {
			m.frames = loadFramesMsg.frames
			if len(m.frames) > 0 {
				m.state = stateFrameList
				m.frameList.setFrames(m.frames)
			} else {
				m.state = stateMainMenu
				m.err = fmt.Errorf("No frames in the selected capture")
			}
		}

		// Handle key presses in saved captures
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				m.state = stateMainMenu
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m *MainModel) View() string {
	var sb strings.Builder

	// Show header
	sb.WriteString("GoCapture - IEEE 802.3/802.11 Frame Analyzer\n")
	sb.WriteString("─────────────────────────────────────────────\n\n")

	// Show appropriate view based on the current state
	switch m.state {
	case stateMainMenu:
		sb.WriteString(m.mainMenu.View())
	case stateCapturing:
		sb.WriteString(fmt.Sprintf("📶 Capturing frames on interface %s\n\n", m.captureEngine.GetInterfaceName()))
		sb.WriteString(fmt.Sprintf("Frames captured: %d\n\n", len(m.frames)))
		sb.WriteString("Press ESC to stop and return to main menu\n")
		sb.WriteString("Press ENTER to stop and view frames\n")
		sb.WriteString("Press S to save the current capture\n")
	case stateFrameList:
		sb.WriteString(m.frameList.View())
	case stateFrameDetail:
		sb.WriteString(m.frameDetail.View())
	case stateSavedCaptures:
		sb.WriteString(m.savedCaptures.View())
	}

	// Show error message if any
	if m.err != nil {
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	}

	// Show help text
	sb.WriteString("\n")
	sb.WriteString("Press q to quit, ESC to go back\n")

	return sb.String()
}

// startCapturing starts capturing frames
func (m *MainModel) startCapturing() tea.Cmd {
	return func() tea.Msg {
		// Start the capture engine
		if err := m.captureEngine.Start(); err != nil {
			m.err = err
			return nil
		}

		// Clear any previous frames
		m.frames = make([]*models.Frame, 0)

		// Return a command to check for frames
		return m.checkForMoreFrames()()
	}
}

// stopCapturing stops capturing frames
func (m *MainModel) stopCapturing() tea.Cmd {
	return func() tea.Msg {
		// Stop the capture engine
		m.captureEngine.Stop()
		return nil
	}
}

// checkForMoreFrames checks for more frames from the capture engine
func (m *MainModel) checkForMoreFrames() tea.Cmd {
	return func() tea.Msg {
		// Get the frame channel
		frameChannel := m.captureEngine.GetFrameChannel()

		// Wait for a frame or channel close
		frame, ok := <-frameChannel
		if !ok {
			// Channel closed, no more frames
			return nil
		}

		// Return the frame
		return newFrameMsg{frame: frame}
	}
}

// newFrameMsg is a message sent when a new frame is captured
type newFrameMsg struct {
	frame *models.Frame
}

// menuSelectedMsg is a message sent when a menu option is selected
type menuSelectedMsg struct {
	option string
}

// loadFramesMsg is a message sent when frames are loaded from storage
type loadFramesMsg struct {
	frames []*models.Frame
}
