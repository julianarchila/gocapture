package storage

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/julian/gocapture/pkg/models"
)

// StorageManager handles saving and loading of captured frames
type StorageManager struct {
	outputDir string
}

// SaveMetadata contains metadata about a saved capture
type SaveMetadata struct {
	Filename    string    `json:"filename"`
	Interface   string    `json:"interface"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	FrameCount  int       `json:"frame_count"`
	Description string    `json:"description"`
}

// NewStorageManager creates a new storage manager
func NewStorageManager(outputDir string) (*StorageManager, error) {
	// Create output directory if it doesn't exist
	if outputDir == "" {
		var homeDir string
		var err error

		// Check if we're running as root (via sudo)
		if os.Geteuid() == 0 {
			// Try to get the original user's home directory
			if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
				// Use the original user's home directory
				homeDir = filepath.Join("/home", sudoUser)
			} else {
				// Fallback to current user's home directory
				homeDir, err = os.UserHomeDir()
				if err != nil {
					return nil, fmt.Errorf("could not determine user home directory: %v", err)
				}
			}
		} else {
			// Not running as root, use current user's home directory
			homeDir, err = os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("could not determine user home directory: %v", err)
			}
		}

		outputDir = filepath.Join(homeDir, ".gocapture", "captures")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %v", err)
	}

	return &StorageManager{
		outputDir: outputDir,
	}, nil
}

// SaveFrames saves a slice of frames to a file
func (sm *StorageManager) SaveFrames(frames []*models.Frame, metadata *SaveMetadata) error {
	if len(frames) == 0 {
		return fmt.Errorf("no frames to save")
	}

	// Generate a filename if not provided
	filename := metadata.Filename
	if filename == "" {
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("capture_%s.gcap", timestamp)
	}

	// Ensure the filename has the correct extension
	if filepath.Ext(filename) != ".gcap" {
		filename += ".gcap"
	}

	// Create the output file
	filePath := filepath.Join(sm.outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Create a gob encoder for serializing the frames
	encoder := gob.NewEncoder(file)

	// Register any complex types that might be in the frame
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	// First write the metadata
	if metadata.StartTime.IsZero() && len(frames) > 0 {
		metadata.StartTime = frames[0].Timestamp
	}
	if metadata.EndTime.IsZero() && len(frames) > 0 {
		metadata.EndTime = frames[len(frames)-1].Timestamp
	}
	if metadata.FrameCount == 0 {
		metadata.FrameCount = len(frames)
	}
	metadata.Filename = filename

	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to encode metadata: %v", err)
	}

	// Create a serializable copy of the frames without OriginalPacket
	serializableFrames := make([]*models.Frame, len(frames))
	for i, frame := range frames {
		// Create a copy of the frame without OriginalPacket
		serializableFrame := *frame
		serializableFrame.OriginalPacket = nil
		serializableFrames[i] = &serializableFrame
	}

	// Then write the frames
	if err := encoder.Encode(serializableFrames); err != nil {
		return fmt.Errorf("failed to encode frames: %v", err)
	}

	return nil
}

// LoadFrames loads frames from a file
func (sm *StorageManager) LoadFrames(filename string) ([]*models.Frame, *SaveMetadata, error) {
	// Check if the file exists
	filePath := filepath.Join(sm.outputDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("file %s does not exist", filename)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a gob decoder for deserializing the frames
	decoder := gob.NewDecoder(file)

	// Register any complex types that might be in the frame
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	// First read the metadata
	var metadata SaveMetadata
	if err := decoder.Decode(&metadata); err != nil {
		return nil, nil, fmt.Errorf("failed to decode metadata: %v", err)
	}

	// Then read the frames
	var frames []*models.Frame
	if err := decoder.Decode(&frames); err != nil {
		return nil, nil, fmt.Errorf("failed to decode frames: %v", err)
	}

	return frames, &metadata, nil
}

// ListSavedCaptures returns a list of all saved captures
func (sm *StorageManager) ListSavedCaptures() ([]*SaveMetadata, error) {
	// Get all files in the output directory
	files, err := os.ReadDir(sm.outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read output directory: %v", err)
	}

	// Filter for .gcap files and load their metadata
	var captures []*SaveMetadata
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".gcap" {
			continue
		}

		// Open the file
		filePath := filepath.Join(sm.outputDir, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			continue // Skip files we can't open
		}

		// Create a gob decoder for deserializing the metadata
		decoder := gob.NewDecoder(file)

		// Register any complex types that might be in the frame
		gob.Register(map[string]interface{}{})
		gob.Register([]interface{}{})

		// Read the metadata
		var metadata SaveMetadata
		if err := decoder.Decode(&metadata); err != nil {
			file.Close()
			continue // Skip files with invalid metadata
		}

		file.Close()
		captures = append(captures, &metadata)
	}

	return captures, nil
}

// DeleteCapture deletes a saved capture
func (sm *StorageManager) DeleteCapture(filename string) error {
	// Check if the file exists
	filePath := filepath.Join(sm.outputDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filename)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
