# GoCapture - User Documentation

## Overview

GoCapture is a comprehensive network frame capture and analysis tool for both IEEE 802.3 (Ethernet) and IEEE 802.11 (WLAN) networks. It provides deep insights into network traffic, security mechanisms, and Quality of Service (QoS) parameters, making it valuable for network administrators, security analysts, and students learning about network protocols.

## Installation

### Prerequisites

- Go 1.18 or later
- libpcap development package
  - Debian/Ubuntu: `sudo apt-get install libpcap-dev`
  - CentOS/RHEL: `sudo yum install libpcap-devel`
  - macOS: `brew install libpcap`

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/USERNAME/gocapture.git
   cd gocapture
   ```

2. Build the project:
   ```bash
   go build -o gocapture ./cmd/gocapture
   ```

3. Install to your system (optional):
   ```bash
   go install ./cmd/gocapture
   ```

## Basic Usage

### View Available Network Interfaces

To view all available network interfaces on your system:

```bash
./gocapture
```

### Start Packet Capture

To start capturing packets on a specific interface:

```bash
sudo ./gocapture -interface eth0
```

Note: Administrator/root privileges are required for capturing packets on most interfaces.

### Command-Line Options

- `-interface`: Network interface to capture from (e.g., eth0, wlan0)
- `-promiscuous`: Enable promiscuous mode (default: true)
- `-filter`: BPF filter expression (e.g., "port 80" to capture only HTTP traffic)

Example with filter:
```bash
sudo ./gocapture -interface eth0 -filter "port 53"
```

## Application Architecture

GoCapture follows a modular architecture with clear separation of concerns:

### Components

1. **Capture Engine** (`internal/capture`)
   - Interfaces with network hardware using the libpcap library
   - Handles frame capture in promiscuous mode
   - Provides a channel-based API for consuming captured frames

2. **Parser Module** (`internal/parser`)
   - Decodes raw frames into structured data
   - Implements specific parsers for Ethernet and WLAN frames
   - Extracts header fields, addresses, and protocol-specific information

3. **Analyzer Engine** (`internal/analyzer`)
   - Interprets parsed frame data
   - Provides security analysis for encryption methods
   - Analyzes QoS parameters and traffic prioritization
   - Gives context and recommendations based on observed frames

4. **Storage Module** (`internal/storage`)
   - Serializes captured frames to disk
   - Loads previously saved captures
   - Manages capture metadata

5. **User Interface** (`ui/`)
   - Terminal-based UI built with Bubble Tea
   - Multiple views: main menu, capture, frame list, frame details
   - Interactive navigation and frame inspection

### Data Flow

1. The Capture Engine captures raw frames from the network interface
2. Raw frames are sent to the Parser Module for decoding
3. Parsed frames are passed to the Analyzer Engine for interpretation
4. The UI displays the analyzed frames and allows for interaction
5. The Storage Module can save/load captures at any point

## Frame Types

### Ethernet Frames (IEEE 802.3)

Ethernet frames are the foundation of wired networks and include:

- MAC header with source and destination addresses
- EtherType field indicating the payload protocol (e.g., IPv4, IPv6, ARP)
- Optional VLAN tagging for network segmentation
- Payload data
- Frame Check Sequence (FCS) for error detection

### WLAN Frames (IEEE 802.11)

WLAN frames are more complex and are categorized into three types:

1. **Management Frames**
   - Establish and maintain communications
   - Examples: Beacons, Association Requests/Responses, Authentication frames
   - Used for network discovery and connection management

2. **Control Frames**
   - Assist in the delivery of data frames
   - Examples: Acknowledgments (ACK), Request to Send (RTS), Clear to Send (CTS)
   - Help manage access to the shared wireless medium

3. **Data Frames**
   - Carry the actual network data
   - Can include QoS parameters for traffic prioritization
   - May be protected by various encryption methods

## Security Analysis

GoCapture identifies and analyzes encryption methods used in wireless networks:

1. **WEP (Wired Equivalent Privacy)**
   - Legacy encryption with serious security flaws
   - Uses RC4 cipher with static keys
   - Vulnerable to statistical attacks

2. **WPA (Wi-Fi Protected Access)**
   - Uses TKIP (Temporal Key Integrity Protocol)
   - Stronger than WEP, but still has vulnerabilities
   - Designed as a transitional solution

3. **WPA2**
   - Uses CCMP based on the AES algorithm
   - Currently the most widely deployed Wi-Fi security standard
   - Vulnerable to KRACK attacks (if not patched)

4. **WPA3**
   - Latest security standard for Wi-Fi
   - Uses SAE (Simultaneous Authentication of Equals)
   - Provides forward secrecy and protection against offline dictionary attacks

## QoS Analysis

For frames with Quality of Service information, GoCapture analyzes:

- **Traffic Categories**: Background, Best Effort, Video, Voice
- **Priority Levels**: 0 (lowest) to 7 (highest)
- **TXOP Allocation**: Transmission opportunity durations
- **ACK Policies**: How frames are acknowledged

## Troubleshooting

### Permission Issues

If you encounter "Operation not permitted" errors:

1. Ensure you're running with administrator/root privileges:
   ```bash
   sudo ./gocapture -interface eth0
   ```

2. Verify the interface exists and is up:
   ```bash
   ip link show
   ```

### No Frames Captured

If the application runs but doesn't capture any frames:

1. Check if your interface is in monitor mode (for WLAN captures)
2. Try a different interface
3. Ensure there is actual network traffic on the interface
4. Try using the loopback interface (`lo`) and generate some local traffic

### Building Errors

If you encounter build errors related to pcap:

1. Ensure libpcap development packages are installed
2. Verify Go modules are properly initialized
3. Run `go mod tidy` to resolve dependencies

## Advanced Usage

### Capturing Wireless Frames

To capture raw 802.11 frames, your wireless interface must be in monitor mode:

1. Put your interface in monitor mode (may vary by OS and driver)
2. Start GoCapture with the monitor mode interface

### Using BPF Filters

Berkeley Packet Filter (BPF) expressions allow for precise capture filtering:

- `port 80 or port 443`: Capture HTTP and HTTPS traffic
- `host 192.168.1.1`: Capture traffic to/from a specific host
- `icmp`: Capture only ICMP packets
- `not port 22`: Exclude SSH traffic

## Contributing

Contributions to GoCapture are welcome! Please see our contributing guidelines for more information.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 