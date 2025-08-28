package root

import (
	"nectar/components/shared"
	"nectar/types"
	"nectar/utils"
	"strings"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConnectionFormModel struct {
	connection     types.Connection
	inputs         []textinput.Model
	filePicker     shared.FilePickerModel
	focused        int
	editing        bool
	showFilePicker bool
	selectedColor  int
}

func NewConnectionForm() ConnectionFormModel {
	// Initialize text inputs for all form fields
	inputs := make([]textinput.Model, 5)

	// Host input
	inputs[utils.InputHost] = textinput.New()
	inputs[utils.InputHost].Placeholder = "localhost"
	inputs[utils.InputHost].CharLimit = 255
	inputs[utils.InputHost].Width = 40

	// Port input
	inputs[utils.InputPort] = textinput.New()
	inputs[utils.InputPort].Placeholder = utils.GetDefaultPort(types.PostgreSQL)
	inputs[utils.InputPort].CharLimit = 6
	inputs[utils.InputPort].Width = 10

	// User input
	inputs[utils.InputUser] = textinput.New()
	inputs[utils.InputUser].Placeholder = "username"
	inputs[utils.InputUser].CharLimit = 100
	inputs[utils.InputUser].Width = 40

	// Password input
	inputs[utils.InputPassword] = textinput.New()
	inputs[utils.InputPassword].Placeholder = "password"
	inputs[utils.InputPassword].EchoMode = textinput.EchoPassword
	inputs[utils.InputPassword].EchoCharacter = '•'
	inputs[utils.InputPassword].CharLimit = 100
	inputs[utils.InputPassword].Width = 40

	// Connection Name input
	inputs[utils.InputConnectionName] = textinput.New()
	inputs[utils.InputConnectionName].Placeholder = "Connection Name"
	inputs[utils.InputConnectionName].CharLimit = 100
	inputs[utils.InputConnectionName].Width = 40

	filePicker := shared.NewFilePicker()

	return ConnectionFormModel{
		connection: types.Connection{
			Type:      types.PostgreSQL,
			EnableSSL: false,
		},
		inputs:        inputs,
		filePicker:    filePicker,
		focused:       0,
		selectedColor: 0,
	}
}

func (m ConnectionFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ConnectionFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showFilePicker {
			return m.handleFilePickerKeys(msg)
		}
		return m.handleFormKeys(msg)
	}

	if m.showFilePicker {
		var cmd tea.Cmd
		m.filePicker, cmd = m.filePicker.Update(msg)
		return m, cmd
	}

	return m.updateInputs(msg)
}

// Handle file picker navigation and disable sidebar keys
func (m ConnectionFormModel) handleFilePickerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.showFilePicker = false
		return m, nil
	case "up", "down", "j", "k":
		// Pass these keys to the file picker when it's active
		var cmd tea.Cmd
		m.filePicker, cmd = m.filePicker.Update(msg)

		// Check if a file was selected
		if m.filePicker.SelectedFile() != "" {
			m.connection.DatabaseFile = m.filePicker.SelectedFile()
			m.showFilePicker = false
		}

		return m, cmd
	}

	// Pass all other keys to the file picker
	var cmd tea.Cmd
	m.filePicker, cmd = m.filePicker.Update(msg)

	if m.filePicker.SelectedFile() != "" {
		m.connection.DatabaseFile = m.filePicker.SelectedFile()
		m.showFilePicker = false
	}

	return m, cmd
}

// Handle form navigation and input
func (m ConnectionFormModel) handleFormKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editing {
		switch msg.String() {
		case "enter", "esc", "tab", "shift+tab":
			m.editing = false
			m.blurAllInputs()

			if msg.String() == "tab" {
				m.nextField()
			} else if msg.String() == "shift+tab" {
				m.prevField()
			}
			return m, nil
		default:
			// Handle text input in the focused field
			for i, input := range m.inputs {
				if input.Focused() {
					var cmd tea.Cmd
					m.inputs[i], cmd = input.Update(msg)
					return m, cmd
				}
			}
		}
	}

	// Handle navigation and actions when not editing
	switch msg.String() {
	case "tab":
		m.nextField()
	case "shift+tab":
		m.prevField()
	case "enter":
		return m.handleEnterKey()
	case "left":
		return m.handleLeftKey()
	case "right":
		return m.handleRightKey()
	}
	return m, nil
}

// Utility methods for form navigation and state management
func (m ConnectionFormModel) getInputIndexFromFocus() int {
	return utils.GetInputIndex(m.connection.Type, m.focused)
}

func (m ConnectionFormModel) handleEnterKey() (tea.Model, tea.Cmd) {
	// Handle database file selection for SQLite
	if m.focused == utils.SQLiteFieldDatabaseFile && m.connection.Type == types.SQLite {
		m.showFilePicker = true
		return m, m.filePicker.Init()
	}

	// Handle SSL toggle for non-SQLite databases
	if m.focused == utils.FieldSSL && m.connection.Type != types.SQLite {
		m.connection.EnableSSL = !m.connection.EnableSSL
		return m, nil
	}

	// Handle text input editing
	inputIndex := m.getInputIndexFromFocus()
	if inputIndex >= 0 {
		if m.editing {
			m.editing = false
			m.inputs[inputIndex].Blur()
		} else {
			m.editing = true
			m.inputs[inputIndex].Focus()
		}
	}

	return m, nil
}

func (m ConnectionFormModel) handleLeftKey() (tea.Model, tea.Cmd) {
	// Handle connection type navigation
	if m.focused == utils.FieldConnectionType {
		m.connection.Type = utils.PrevConnectionType(m.connection.Type)
		m.updatePortPlaceholder()
	} else if m.isColorField() {
		// Handle color selection
		if m.selectedColor > 0 {
			m.selectedColor--
		}
	}
	return m, nil
}

func (m ConnectionFormModel) handleRightKey() (tea.Model, tea.Cmd) {
	// Handle connection type navigation
	if m.focused == utils.FieldConnectionType {
		m.connection.Type = utils.NextConnectionType(m.connection.Type)
		m.updatePortPlaceholder()
	} else if m.isColorField() {
		// Handle color selection
		if m.selectedColor < len(shared.ConnectionColors)-1 {
			m.selectedColor++
		}
	}
	return m, nil
}

// Helper method to check if current field is the color field
func (m ConnectionFormModel) isColorField() bool {
	if m.connection.Type == types.SQLite {
		return m.focused == utils.SQLiteFieldColor
	}
	return m.focused == utils.FieldColor
}

// Update port placeholder based on connection type
func (m *ConnectionFormModel) updatePortPlaceholder() {
	port := utils.GetDefaultPort(m.connection.Type)
	if port != "" {
		m.inputs[utils.InputPort].Placeholder = port
	}
}

// Navigation methods
func (m *ConnectionFormModel) nextField() {
	m.blurAllInputs()
	maxFields := utils.GetFieldCount(m.connection.Type)
	m.focused = (m.focused + 1) % maxFields
}

func (m *ConnectionFormModel) prevField() {
	m.blurAllInputs()
	m.focused--
	if m.focused < 0 {
		maxFields := utils.GetFieldCount(m.connection.Type)
		m.focused = maxFields - 1
	}
}

func (m *ConnectionFormModel) blurAllInputs() {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

func (m ConnectionFormModel) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

// View rendering methods
func (m ConnectionFormModel) View() string {
	if m.showFilePicker {
		return m.renderFilePicker()
	}
	return m.renderForm()
}

func (m ConnectionFormModel) renderFilePicker() string {
	return m.filePicker.View()
}

func (m ConnectionFormModel) getEditingText() string {
	if m.editing {
		return "press Enter to finish editing"
	}
	return "press Enter to edit"
}

func (m ConnectionFormModel) renderForm() string {
	var content strings.Builder

	// Define consistent styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Mauve().Hex,
			Dark:  catppuccin.Mocha.Mauve().Hex,
		}).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Text().Hex,
			Dark:  catppuccin.Mocha.Text().Hex,
		})

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Mauve().Hex,
			Dark:  catppuccin.Mocha.Mauve().Hex,
		})

	// Form title
	content.WriteString(titleStyle.Render("New Connection") + "\n\n")

	// Connection Type selector
	m.renderConnectionTypeField(&content, focusedStyle, labelStyle)

	// Database-specific fields
	if m.connection.Type == types.SQLite {
		m.renderSQLiteFields(&content, focusedStyle, labelStyle)
	} else {
		m.renderDatabaseServerFields(&content, focusedStyle, labelStyle)
	}

	// Separator
	content.WriteString("─────────────────────────────────────────\n\n")

	// Connection saving fields (common to all database types)
	m.renderConnectionSavingFields(&content, focusedStyle, labelStyle)

	// Apply form styling with fixed width and border
	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Mauve().Hex,
			Dark:  catppuccin.Mocha.Mauve().Hex,
		}).
		Padding(2, 4).
		Width(60)

	return formStyle.Render(content.String())
}

// Render the connection type selector field
func (m ConnectionFormModel) renderConnectionTypeField(content *strings.Builder, focusedStyle, labelStyle lipgloss.Style) {
	connTypeLabel := "Connection Type: " + m.connection.Type.String()
	if m.focused == utils.FieldConnectionType {
		connTypeLabel = focusedStyle.Render("> " + connTypeLabel + " (use ← → to change)")
	} else {
		connTypeLabel = labelStyle.Render("  " + connTypeLabel)
	}
	content.WriteString(connTypeLabel + "\n\n")
}

// Render SQLite-specific fields (database file selection)
func (m ConnectionFormModel) renderSQLiteFields(content *strings.Builder, focusedStyle, labelStyle lipgloss.Style) {
	// Database File field
	fileLabel := "Database File: "
	if m.connection.DatabaseFile != "" {
		fileLabel += m.connection.DatabaseFile
	} else {
		fileLabel += "No file selected"
	}

	if m.focused == utils.SQLiteFieldDatabaseFile {
		fileLabel = focusedStyle.Render("> " + fileLabel + " (press Enter to browse)")
	} else {
		fileLabel = labelStyle.Render("  " + fileLabel)
	}
	content.WriteString(fileLabel + "\n\n")
}

// Render database server fields (PostgreSQL/MySQL)
func (m ConnectionFormModel) renderDatabaseServerFields(content *strings.Builder, focusedStyle, labelStyle lipgloss.Style) {
	// Host field
	m.renderInputField(content, "Host", utils.FieldHost, utils.InputHost, focusedStyle, labelStyle)

	// Port field
	m.renderInputField(content, "Port", utils.FieldPort, utils.InputPort, focusedStyle, labelStyle)

	// SSL toggle field
	m.renderSSLField(content, focusedStyle, labelStyle)

	// User field
	m.renderInputField(content, "User", utils.FieldUser, utils.InputUser, focusedStyle, labelStyle)

	// Password field
	m.renderInputField(content, "Password", utils.FieldPassword, utils.InputPassword, focusedStyle, labelStyle)
}

// Helper to render a standard input field with label
func (m ConnectionFormModel) renderInputField(content *strings.Builder, fieldName string, fieldIndex, inputIndex int, focusedStyle, labelStyle lipgloss.Style) {
	label := fieldName + ":"
	if m.focused == fieldIndex {
		label = focusedStyle.Render("> " + label + " (" + m.getEditingText() + ")")
	} else {
		label = labelStyle.Render("  " + label)
	}
	content.WriteString(label + "\n")
	content.WriteString("  " + m.inputs[inputIndex].View() + "\n\n")
}

// Render SSL toggle field for database servers
func (m ConnectionFormModel) renderSSLField(content *strings.Builder, focusedStyle, labelStyle lipgloss.Style) {
	sslLabel := "Enable SSL: "
	if m.connection.EnableSSL {
		sslLabel += "✓ Yes"
	} else {
		sslLabel += "✗ No"
	}

	if m.focused == utils.FieldSSL {
		sslLabel = focusedStyle.Render("> " + sslLabel + " (press Enter to toggle)")
	} else {
		sslLabel = labelStyle.Render("  " + sslLabel)
	}
	content.WriteString(sslLabel + "\n\n")
}

// Render connection saving fields (name and color)
func (m ConnectionFormModel) renderConnectionSavingFields(content *strings.Builder, focusedStyle, labelStyle lipgloss.Style) {
	// Connection Name field - field index depends on database type
	nameFieldIndex := utils.SQLiteFieldConnectionName
	if m.connection.Type != types.SQLite {
		nameFieldIndex = utils.FieldConnectionName
	}

	saveNameLabel := "Save Connection:"
	if m.focused == nameFieldIndex {
		saveNameLabel = focusedStyle.Render("> " + saveNameLabel + " (" + m.getEditingText() + ")")
	} else {
		saveNameLabel = labelStyle.Render("  " + saveNameLabel)
	}
	content.WriteString(saveNameLabel + "\n")
	content.WriteString("  " + m.inputs[utils.InputConnectionName].View() + "\n\n")

	// Color selection field - field index depends on database type
	colorFieldIndex := utils.SQLiteFieldColor
	if m.connection.Type != types.SQLite {
		colorFieldIndex = utils.FieldColor
	}

	colorLabel := "Color: " + shared.ConnectionColors[m.selectedColor].Name
	if m.focused == colorFieldIndex {
		colorLabel = focusedStyle.Render("> " + colorLabel + " (use ← → to change)")
	} else {
		colorLabel = labelStyle.Render("  " + colorLabel)
	}

	// Color preview box
	colorPreview := lipgloss.NewStyle().
		Background(shared.ConnectionColors[m.selectedColor].Color).
		Render("  ")
	content.WriteString(colorLabel + " " + colorPreview + "\n")
}
