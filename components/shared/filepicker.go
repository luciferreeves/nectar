package shared

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	MaxVisibleFiles = 10
)

type FilePickerModel struct {
	currentDir   string
	files        []fs.DirEntry
	selected     int
	selectedFile string
	err          error
	scrollOffset int
}

func NewFilePicker() FilePickerModel {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	fp := FilePickerModel{
		currentDir:   homeDir,
		selected:     0,
		scrollOffset: 0,
	}
	fp.loadDirectory()
	return fp
}

func (m *FilePickerModel) loadDirectory() {
	files, err := os.ReadDir(m.currentDir)
	if err != nil {
		m.err = err
		m.files = nil
		return
	}

	var filteredFiles []fs.DirEntry

	// Always add parent directory option unless we're at root
	parentDir := filepath.Dir(m.currentDir)
	if parentDir != m.currentDir {
		// Create a synthetic directory entry for ".."
		filteredFiles = append(filteredFiles, &parentDirEntry{})
	}

	// Add directories and SQLite files
	for _, file := range files {
		if file.IsDir() {
			filteredFiles = append(filteredFiles, file)
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".db") ||
			strings.HasSuffix(strings.ToLower(file.Name()), ".sqlite") ||
			strings.HasSuffix(strings.ToLower(file.Name()), ".sqlite3") {
			filteredFiles = append(filteredFiles, file)
		}
	}

	// Sort: parent directory first, then directories, then files
	sort.Slice(filteredFiles, func(i, j int) bool {
		nameI := filteredFiles[i].Name()
		nameJ := filteredFiles[j].Name()

		// Parent directory always comes first
		if nameI == ".." {
			return true
		}
		if nameJ == ".." {
			return false
		}

		// Then sort by type (directories first) and name
		if filteredFiles[i].IsDir() != filteredFiles[j].IsDir() {
			return filteredFiles[i].IsDir()
		}
		return nameI < nameJ
	})

	m.files = filteredFiles
	m.selected = 0
	m.scrollOffset = 0
	m.err = nil
}

// parentDirEntry implements fs.DirEntry for the ".." parent directory
type parentDirEntry struct{}

func (p *parentDirEntry) Name() string               { return ".." }
func (p *parentDirEntry) IsDir() bool                { return true }
func (p *parentDirEntry) Type() fs.FileMode          { return fs.ModeDir }
func (p *parentDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func (m FilePickerModel) Init() tea.Cmd {
	return nil
}

func (m FilePickerModel) Update(msg tea.Msg) (FilePickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
				// Adjust scroll offset to keep selection visible
				if m.selected < m.scrollOffset {
					m.scrollOffset = m.selected
				}
			}
		case "down", "j":
			if m.selected < len(m.files)-1 {
				m.selected++
				// Adjust scroll offset to keep selection visible
				if m.selected >= m.scrollOffset+MaxVisibleFiles {
					m.scrollOffset = m.selected - MaxVisibleFiles + 1
				}
			}
		case "enter":
			if len(m.files) == 0 {
				break
			}
			selectedFile := m.files[m.selected]
			if selectedFile.IsDir() {
				if selectedFile.Name() == ".." {
					// Navigate to parent directory
					m.currentDir = filepath.Dir(m.currentDir)
				} else {
					// Navigate to selected directory
					m.currentDir = filepath.Join(m.currentDir, selectedFile.Name())
				}
				m.loadDirectory()
			} else {
				// Select the SQLite file
				m.selectedFile = filepath.Join(m.currentDir, selectedFile.Name())
			}
		}
	}
	return m, nil
}

func (m FilePickerModel) View() string {
	// Fixed height container to prevent UI pushing
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(80).
		Height(MaxVisibleFiles + 8) // Fixed height: header + files + help + padding

	if m.err != nil {
		errorContent := "Select SQLite Database File\n\nError: " + m.err.Error()
		return containerStyle.Render(errorContent)
	}

	var content strings.Builder

	// Header section
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	content.WriteString(headerStyle.Render("Select SQLite Database File"))
	content.WriteString("\n")
	content.WriteString("Current: " + filepath.Base(m.currentDir))
	content.WriteString("\n\n")

	// Files section with fixed height
	if len(m.files) == 0 {
		content.WriteString("No SQLite files found in this directory")
		// Add padding to maintain consistent height
		for i := 0; i < MaxVisibleFiles-1; i++ {
			content.WriteString("\n")
		}
	} else {
		visibleStart := m.scrollOffset
		visibleEnd := min(visibleStart+MaxVisibleFiles, len(m.files))

		// Display visible files
		for i := visibleStart; i < visibleEnd; i++ {
			file := m.files[i]
			var line string

			if file.Name() == ".." {
				line = "📁 .. (parent directory)"
			} else if file.IsDir() {
				line = "📁 " + file.Name() + "/"
			} else {
				line = "📄 " + file.Name()
			}

			if i == m.selected {
				selectedStyle := lipgloss.NewStyle().
					Background(lipgloss.Color("12")).
					Foreground(lipgloss.Color("0"))
				line = selectedStyle.Render("> " + line)
			} else {
				line = "  " + line
			}
			content.WriteString(line + "\n")
		}

		// Fill remaining lines to maintain consistent height
		remainingLines := MaxVisibleFiles - (visibleEnd - visibleStart)
		for range remainingLines {
			content.WriteString("\n")
		}
	}

	// Footer with navigation help
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	help := helpStyle.Render("↑/↓: navigate, Enter: select/open, Esc: cancel")
	content.WriteString("\n" + help)

	return containerStyle.Render(content.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m FilePickerModel) SelectedFile() string {
	return m.selectedFile
}
