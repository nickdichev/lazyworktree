package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen types for modal dialogs
type screenType int

const (
	screenNone screenType = iota
	screenConfirm
	screenInput
	screenHelp
	screenTrust
	screenWelcome
	screenCommit
)

// ConfirmScreen represents a confirmation dialog
type ConfirmScreen struct {
	message string
	result  chan bool
}

// InputScreen represents an input dialog
type InputScreen struct {
	prompt      string
	placeholder string
	value       string
	input       textinput.Model
	result      chan string
}

// HelpScreen represents a help screen
type HelpScreen struct {
	viewport viewport.Model
}

// TrustScreen represents a trust confirmation screen
type TrustScreen struct {
	filePath    string
	commands    []string
	viewport    viewport.Model
	result      chan string
}

// WelcomeScreen represents a welcome screen
type WelcomeScreen struct {
	currentDir  string
	worktreeDir string
	result      chan bool
}

// CommitScreen represents a commit detail screen
type CommitScreen struct {
	header      string
	diff        string
	useDelta    bool
	viewport    viewport.Model
	headerShown bool
}

// NewConfirmScreen creates a new confirmation screen
func NewConfirmScreen(message string) *ConfirmScreen {
	return &ConfirmScreen{
		message: message,
		result:  make(chan bool, 1),
	}
}

// Init initializes the confirm screen
func (s *ConfirmScreen) Init() tea.Cmd {
	return nil
}

// Update handles updates for the confirm screen
func (s *ConfirmScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			s.result <- true
			return s, tea.Quit
		case "n", "N", "esc", "q":
			s.result <- false
			return s, tea.Quit
		}
	}
	return s, nil
}

// View renders the confirm screen
func (s *ConfirmScreen) View() string {
	width := 60
	height := 11

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(width).
		Height(height)

	messageStyle := lipgloss.NewStyle().
		Width(width - 4).
		Height(height - 6).
		Align(lipgloss.Center, lipgloss.Center)

	buttonStyle := lipgloss.NewStyle().
		Width((width - 6) / 2).
		Align(lipgloss.Center).
		Padding(0, 1)

	confirmButton := buttonStyle.
		Foreground(lipgloss.Color("9")).
		Render("[Confirm]")

	cancelButton := buttonStyle.
		Foreground(lipgloss.Color("4")).
		Render("[Cancel]")

	content := fmt.Sprintf("%s\n\n%s  %s",
		messageStyle.Render(s.message),
		confirmButton,
		cancelButton,
	)

	return boxStyle.Render(content)
}

// NewInputScreen creates a new input screen
func NewInputScreen(prompt, placeholder, value string) *InputScreen {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return &InputScreen{
		prompt:      prompt,
		placeholder: placeholder,
		value:       value,
		input:       ti,
		result:      make(chan string, 1),
	}
}

// Init initializes the input screen
func (s *InputScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles updates for the input screen
func (s *InputScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			value := s.input.Value()
			s.result <- value
			return s, tea.Quit
		case "esc":
			s.result <- ""
			return s, tea.Quit
		}
	}

	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

// View renders the input screen
func (s *InputScreen) View() string {
	width := 60

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(width)

	promptStyle := lipgloss.NewStyle().
		MarginBottom(1)

	content := fmt.Sprintf("%s\n\n%s",
		promptStyle.Render(s.prompt),
		s.input.View(),
	)

	return boxStyle.Render(content)
}

// NewHelpScreen creates a new help screen
func NewHelpScreen() *HelpScreen {
	helpText := `# Git Worktree Status Help

**Navigation**
- j / Down: Move cursor down
- k / Up: Move cursor up
- 1: Focus Worktree pane
- 2: Focus Info/Diff pane
- 3: Focus Log pane
- Enter: Jump to selected worktree (exit and cd)
- Tab: Cycle focus (table → status → log)
- j / k in Recent Log: Move between commits
- Enter in Recent Log: Open commit details and diff

**Actions**
- c: Create new worktree
- d: Refresh diff in the status pane (auto-shown when dirty; uses delta if available)
- D: Delete selected worktree
- f: Fetch all remotes
- p: Fetch PR status from GitHub
- r: Refresh list
- s: Sort (toggle Name/Last Active)
- /: Filter worktrees
- g: Open LazyGit
- ?: Show this help

**Status Indicators**
- ✔ Clean: No local changes
- ✎ Dirty: Uncommitted changes
- ↑N: Ahead of remote by N commits
- ↓N: Behind remote by N commits

**Performance Note**
PR data is not fetched by default for speed.
Press p to fetch PR information from GitHub.`

	vp := viewport.New(80, 30)
	vp.SetContent(helpText)

	return &HelpScreen{
		viewport: vp,
	}
}

// Init initializes the help screen
func (s *HelpScreen) Init() tea.Cmd {
	return nil
}

// Update handles updates for the help screen
func (s *HelpScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return s, tea.Quit
		}
	}
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

// View renders the help screen
func (s *HelpScreen) View() string {
	width := 80
	height := 30

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		Padding(1, 2).
		Width(width).
		Height(height)

	return boxStyle.Render(s.viewport.View())
}

// NewTrustScreen creates a new trust screen
func NewTrustScreen(filePath string, commands []string) *TrustScreen {
	commandsText := strings.Join(commands, "\n")
	question := fmt.Sprintf("The repository config '%s' defines the following commands.\nThis file has changed or hasn't been trusted yet.\nDo you trust these commands to run?", filePath)

	content := fmt.Sprintf("%s\n\n%s", question, commandsText)

	vp := viewport.New(70, 20)
	vp.SetContent(content)

	return &TrustScreen{
		filePath: filePath,
		commands: commands,
		viewport: vp,
		result:   make(chan string, 1),
	}
}

// Init initializes the trust screen
func (s *TrustScreen) Init() tea.Cmd {
	return nil
}

// Update handles updates for the trust screen
func (s *TrustScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t", "T":
			s.result <- "trust"
			return s, tea.Quit
		case "b", "B":
			s.result <- "block"
			return s, tea.Quit
		case "esc", "c", "C":
			s.result <- "cancel"
			return s, tea.Quit
		}
	}
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

// View renders the trust screen
func (s *TrustScreen) View() string {
	width := 70
	height := 25

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(width).
		Height(height)

	buttonStyle := lipgloss.NewStyle().
		Width(20).
		Align(lipgloss.Center).
		Padding(0, 1).
		Margin(0, 1)

	trustButton := buttonStyle.
		Foreground(lipgloss.Color("2")).
		Render("[Trust & Run]")

	blockButton := buttonStyle.
		Foreground(lipgloss.Color("3")).
		Render("[Block (Skip)]")

	cancelButton := buttonStyle.
		Foreground(lipgloss.Color("1")).
		Render("[Cancel Operation]")

	content := fmt.Sprintf("%s\n\n%s  %s  %s",
		s.viewport.View(),
		trustButton,
		blockButton,
		cancelButton,
	)

	return boxStyle.Render(content)
}

// NewWelcomeScreen creates a new welcome screen
func NewWelcomeScreen(currentDir, worktreeDir string) *WelcomeScreen {
	return &WelcomeScreen{
		currentDir:  currentDir,
		worktreeDir: worktreeDir,
		result:       make(chan bool, 1),
	}
}

// Init initializes the welcome screen
func (s *WelcomeScreen) Init() tea.Cmd {
	return nil
}

// Update handles updates for the welcome screen
func (s *WelcomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r", "R":
			s.result <- true
			return s, tea.Quit
		case "q", "Q", "esc":
			s.result <- false
			return s, tea.Quit
		}
	}
	return s, nil
}

// View renders the welcome screen
func (s *WelcomeScreen) View() string {
	width := 70
	height := 20

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		Padding(1, 2).
		Width(width).
		Height(height)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")).
		Align(lipgloss.Center).
		MarginBottom(1)

	messageStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		MarginBottom(2)

	buttonStyle := lipgloss.NewStyle().
		Width(20).
		Align(lipgloss.Center).
		Padding(0, 1).
		Margin(0, 1)

	quitButton := buttonStyle.
		Foreground(lipgloss.Color("1")).
		Render("[Quit]")

	retryButton := buttonStyle.
		Foreground(lipgloss.Color("4")).
		Render("[Retry]")

	message := fmt.Sprintf("No worktrees found.\n\nCurrent Directory: %s\nWorktree Root: %s\n\nPlease ensure you are in a git repository or the configured worktree root.\nYou may need to initialize a repository or configure 'worktree_dir' in config.",
		s.currentDir,
		s.worktreeDir,
	)

	content := fmt.Sprintf("%s\n%s\n\n%s  %s",
		titleStyle.Render("Welcome to LazyWorktree"),
		messageStyle.Render(message),
		quitButton,
		retryButton,
	)

	return boxStyle.Render(content)
}

// NewCommitScreen creates a new commit detail screen
func NewCommitScreen(header, diff string, useDelta bool) *CommitScreen {
	content := fmt.Sprintf("%s\n\n%s", header, diff)
	vp := viewport.New(95, 95)
	vp.SetContent(content)

	return &CommitScreen{
		header:      header,
		diff:        diff,
		useDelta:    useDelta,
		viewport:    vp,
		headerShown: true,
	}
}

// Init initializes the commit screen
func (s *CommitScreen) Init() tea.Cmd {
	return nil
}

// Update handles updates for the commit screen
func (s *CommitScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return s, tea.Quit
		case "j", "down":
			s.viewport.LineDown(1)
			return s, nil
		case "k", "up":
			s.viewport.LineUp(1)
			return s, nil
		case "ctrl+d", " ":
			s.viewport.HalfViewDown()
			return s, nil
		case "ctrl+u":
			s.viewport.HalfViewUp()
			return s, nil
		case "g":
			s.viewport.GotoTop()
			return s, nil
		case "G":
			s.viewport.GotoBottom()
			return s, nil
		}
	}
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

// View renders the commit screen
func (s *CommitScreen) View() string {
	width := 95
	height := 95

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(width).
		Height(height)

	return boxStyle.Render(s.viewport.View())
}

// Key bindings for screens
type keyMap struct {
	Confirm key.Binding
	Cancel  key.Binding
	Quit    key.Binding
	Scroll  key.Binding
}

var defaultKeyMap = keyMap{
	Confirm: key.NewBinding(
		key.WithKeys("enter", "y"),
		key.WithHelp("enter/y", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc", "q", "n"),
		key.WithHelp("esc/q/n", "cancel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
	Scroll: key.NewBinding(
		key.WithKeys("j", "k", "ctrl+d", "ctrl+u", "g", "G"),
		key.WithHelp("j/k/ctrl+d/ctrl+u/g/G", "scroll"),
	),
}
