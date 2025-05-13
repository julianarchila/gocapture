# GoCapture - User Interface Controls

This document details the cursor navigation rules and keyboard controls for the GoCapture terminal user interface.

## Global Controls

These controls work throughout the application:

| Key       | Action                            |
|-----------|-----------------------------------|
| `q`       | Quit application                  |
| `Ctrl+C`  | Quit application                  |
| `Esc`     | Go back to previous screen        |

## Main Menu Navigation

The main menu is the starting point of the application:

| Key       | Action                            |
|-----------|-----------------------------------|
| `↑` / `k` | Move cursor up                    |
| `↓` / `j` | Move cursor down                  |
| `Enter`   | Select highlighted option         |

Available menu options:
- **Start Capture**: Begin capturing frames on the specified interface
- **Load Capture**: Browse and load previously saved captures
- **Quit**: Exit the application

## Capture Screen Controls

When actively capturing frames:

| Key       | Action                                 |
|-----------|----------------------------------------|
| `Esc`     | Stop capturing and return to main menu |
| `Enter`   | Stop capturing and view captured frames|
| `s`       | Save the current capture               |

## Frame List Navigation

When viewing the list of captured frames:

| Key       | Action                               |
|-----------|--------------------------------------|
| `↑` / `k` | Move cursor up                       |
| `↓` / `j` | Move cursor down                     |
| `PgUp`    | Move cursor up one page              |
| `PgDn`    | Move cursor down one page            |
| `Enter`   | View detailed information for the selected frame |
| `s`       | Save the current list of frames      |
| `Esc`     | Return to main menu                  |

The frame list view shows:
- Frame ID
- Timestamp
- Source and destination MAC addresses
- Frame length
- Summary of frame type and content

## Frame Detail View Controls

When viewing detailed information about a specific frame:

| Key         | Action                                   |
|-------------|------------------------------------------|
| `Tab`       | Cycle through view modes (Summary, Details, Hex Dump) |
| `↑` / `k`   | Scroll content up                        |
| `↓` / `j`   | Scroll content down                      |
| `→` / `l` / `n` | View next frame in sequence          |
| `←` / `h` / `p` | View previous frame in sequence      |
| `Esc`       | Return to frame list                     |

### View Modes

1. **Summary View**
   - Shows a concise overview of the frame
   - Displays frame type, addresses, and important fields
   - Provides analysis context and security/QoS information

2. **Details View**
   - Shows all frame fields and their values
   - Displays raw header information
   - Shows complete analysis results

3. **Hex Dump View**
   - Shows raw binary data in hexadecimal format
   - Displays both hex values and ASCII representation
   - Shows byte offsets for easy reference

## Saved Captures Screen

When browsing saved captures:

| Key       | Action                               |
|-----------|--------------------------------------|
| `↑` / `k` | Move cursor up                       |
| `↓` / `j` | Move cursor down                     |
| `Enter`   | Load selected capture                |
| `Esc`     | Return to main menu                  |

## Cursor Rules and Behavior

The cursor in GoCapture follows these consistent rules:

1. **Visibility**: The cursor is always visible as a `>` character at the start of the selected item

2. **Wrapping**: Cursors do not wrap around from bottom to top or top to bottom

3. **Pagination**: When a list exceeds the visible area:
   - The view automatically scrolls to keep the cursor visible
   - Page Up/Down keys move the cursor by a full page
   - The current page position is indicated (e.g., "Showing 1-10 of 50 frames")

4. **Selection**: The currently selected item is always highlighted with the cursor

5. **Bounds Checking**: 
   - Pressing up when at the top of a list has no effect
   - Pressing down when at the bottom of a list has no effect
   - Next/previous frame navigation stops at the beginning or end of the frame list

## Special Input Handling

### Filter Expressions

When applying BPF filter expressions via command-line arguments:

- Expressions must be properly quoted if they contain spaces
- Complex expressions can use boolean operators (AND, OR, NOT)
- Examples:
  ```
  -filter "host 192.168.1.1 and port 80"
  -filter "not port 22"
  -filter "ether host 00:11:22:33:44:55"
  ```

### Interface Selection

Interface names are case-sensitive and must match exactly as listed when running GoCapture without arguments.

## Accessibility Considerations

- Vim-style navigation (`h`, `j`, `k`, `l`) is provided as an alternative to arrow keys
- High-contrast cursor indicator (`>`) makes selection clearly visible
- Summary information is consistently formatted for screen readers 