package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket/pcap"
	"github.com/julian/gocapture/internal/capture"
	"github.com/julian/gocapture/ui"
)

func main() {
	// Parse command line arguments
	interfaceName := flag.String("interface", "", "Network interface to capture from")
	promiscuous := flag.Bool("promiscuous", true, "Enable promiscuous mode")
	filter := flag.String("filter", "", "BPF filter expression")
	flag.Parse()

	// List available interfaces if none specified
	if *interfaceName == "" {
		interfaces, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatalf("Could not find network interfaces: %v", err)
		}

		fmt.Println("Available network interfaces:")
		for _, iface := range interfaces {
			fmt.Printf("- %s: %s\n", iface.Name, iface.Description)
			for _, addr := range iface.Addresses {
				fmt.Printf("  %s\n", addr.IP)
			}
		}
		os.Exit(0)
	}

	// Initialize the capture engine
	captureEngine, err := capture.NewCaptureEngine(*interfaceName, *promiscuous, *filter)
	if err != nil {
		log.Fatalf("Failed to initialize capture engine: %v", err)
	}

	// Start the UI
	if err := ui.StartUI(captureEngine); err != nil {
		log.Fatalf("UI error: %v", err)
	}
} 