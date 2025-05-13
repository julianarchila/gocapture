# GoCapture: IEEE 802.3/802.11 Frame Capture and Analysis Tool

GoCapture is a comprehensive LAN/WLAN frame capture and analysis tool written in Go. It provides detailed inspection of network frames from both IEEE 802.3 (Ethernet) and IEEE 802.11 (WLAN) networks, with a particular focus on analyzing security fields and Quality of Service (QoS) parameters in wireless frames.

## Features

- **Frame Capture**: Capture raw network frames from both wired (IEEE 802.3) and wireless (IEEE 802.11) interfaces
- **Promiscuous Mode**: Capture all network traffic passing through the interface
- **BPF Filters**: Apply Berkeley Packet Filter expressions to filter captures
- **Frame Parsing**: Detailed parsing of all frame fields
- **Security Analysis**: Identification and analysis of security mechanisms (WEP, WPA, WPA2, WPA3)
- **QoS Analysis**: Analysis of QoS parameters and traffic prioritization
- **Terminal UI**: Intuitive terminal-based user interface using Bubble Tea
- **Save/Load**: Save captures to disk and load them later for analysis

## Installation

### Prerequisites

- Go 1.18 or later
- libpcap development package
  - Debian/Ubuntu: `sudo apt-get install libpcap-dev`
  - CentOS/RHEL: `sudo yum install libpcap-devel`
  - macOS: `brew install libpcap`

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/USERNAME/gocapture.git
   cd gocapture
   ```

2. Build the project:
   ```
   go build -o gocapture ./cmd/gocapture
   ```

3. Install to your system (optional):
   ```
   go install ./cmd/gocapture
   ```

## Usage

### Starting the Application

```
./gocapture
```

Without any arguments, GoCapture will list the available network interfaces.

### Capture Options

```
./gocapture -interface eth0 -promiscuous=true -filter "port 80"
```

Options:
- `-interface`: Network interface to capture from
- `-promiscuous`: Enable promiscuous mode (default: true)
- `-filter`: BPF filter expression to apply

### Interface Navigation

- Use arrow keys to navigate menus and lists
- Press Enter to select options
- Press Esc to go back
- Press Tab to cycle through views in frame details
- Press q to quit the application

## Understanding Frame Types

### Ethernet Frames (IEEE 802.3)

Ethernet frames include:
- Source and destination MAC addresses
- EtherType field indicating the payload protocol
- Optional VLAN tagging information

### WLAN Frames (IEEE 802.11)

WLAN frames are divided into three main types:

1. **Management Frames**: Used for station association, authentication, and beacon transmission
2. **Control Frames**: Used for medium access control (ACK, RTS, CTS, etc.)
3. **Data Frames**: Carry actual network data, possibly with QoS information

Each frame type has specific fields and purposes explained within the application.

## Security Analysis

GoCapture identifies encryption types used in wireless networks:
- WEP (Wired Equivalent Privacy) - Legacy, insecure
- WPA (Wi-Fi Protected Access) using TKIP
- WPA2 using CCMP (AES)
- WPA3 using GCMP

The tool provides context about the security implications of different encryption methods.

## QoS Analysis

For frames containing QoS information, GoCapture analyzes:
- Traffic prioritization levels (0-7)
- Traffic types (Background, Best Effort, Video, Voice)
- TXOP (Transmission Opportunity) allocation
- ACK policies

## License

[MIT License](LICENSE)

## Acknowledgments

- The Go programming language team
- The gopacket and libpcap developers
- The Bubble Tea framework developers 