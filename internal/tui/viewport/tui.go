package viewport

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Constants
const viewportHeight = 4
const outerMargin = 3

const maxLogLines = 25 // 🔥 Track max number of lines

// Message types
type outputMsg string
type doneMsg struct {
	exitCode int
	err      error
}
type tickMsg time.Duration

// Run a function with a spinner and viewport
func RunWithSpinner(title string, fn func(chan string, *[]string) (int, error)) (*Model, error) {
	m := newModel(title, fn)
	p := tea.NewProgram(m)

	updatedModel, err := p.Run()       // ✅ Get the final updated model
	finalModel := updatedModel.(Model) // ✅ Type assertion to model

	return &finalModel, err // ✅ Return the correct model instance
}

// Execute a command, stream its output, and return exit code
func RunCommand(outputChan chan string, logLines *[]string, cmd *exec.Cmd) (int, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 1, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 1, err
	}

	if err := cmd.Start(); err != nil {
		return 1, err
	}

	var wg sync.WaitGroup
	output := make(chan string, 100)

	wg.Add(2)
	go func() {
		defer wg.Done()
		streamOutput(stdout, output, logLines)
	}()
	go func() {
		defer wg.Done()
		streamOutput(stderr, output, logLines)
	}()

	// 🔥 Dedicated goroutine to forward all output safely
	// Send all output to the UI
	go func() {
		for line := range output {
			outputChan <- line
		}
	}()

	go func() {
		wg.Wait()
		close(output)
	}()

	err = cmd.Wait()
	// ✅ Check if process was terminated by a signal
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// 🔥 Use syscall.WaitStatus to get termination signal
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() { // ✅ Check if the process was killed by a signal
					sig := status.Signal()
					return -int(sig), fmt.Errorf("terminated by signal: %v", sig)
				}
			}
		}
	}

	// ✅ Ensure ProcessState is valid before calling ExitCode()
	if cmd.ProcessState == nil {
		return 1, fmt.Errorf("process did not start or was force-killed")
	}

	return cmd.ProcessState.ExitCode(), err
}

func streamOutput(r io.Reader, output chan string, logLines *[]string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// 🔥 Push new line
		*logLines = append(*logLines, line)

		// 🔥 Pop/Shift: Remove oldest if exceeding maxLogLines
		if len(*logLines) > maxLogLines {
			*logLines = (*logLines)[1:] // Remove the first (oldest) element
		}

		output <- line // Send to viewport
	}
}

// Wait for the next output message
func waitForOutput(outputChan chan string) tea.Cmd {
	return func() tea.Msg {
		return outputMsg(<-outputChan) // 🔥 Blocks until a message is received
	}
}

// Bubble Tea model
type Model struct {
	title    string
	spinner  spinner.Model
	fn       func(chan string, *[]string) (int, error)
	sub      chan string
	viewport viewport.Model
	LogLines *[]string
	Start    time.Time
	width    int
	ExitCode int
	done     bool
}

// Create a new model
func newModel(title string, fn func(chan string, *[]string) (int, error)) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	vp := viewport.New(80, viewportHeight)
	vp.SetContent("")

	// ✅ Allocate logLines as a pointer
	logLines := &[]string{}

	return Model{
		title:    title,
		spinner:  s,
		fn:       fn,
		sub:      make(chan string),
		LogLines: logLines,
		viewport: vp,
		Start:    time.Now(),
		ExitCode: -1,
		width:    80,
	}
}

// Initialize the Bubble Tea program
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			code, err := m.fn(m.sub, m.LogLines)
			m.ExitCode = code
			return doneMsg{exitCode: code, err: err}
		},
		waitForOutput(m.sub),
		tickCmd(m.Start),
	)
}

// Tick function to update duration timer
func tickCmd(start time.Time) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(time.Since(start))
	})
}

// Update the UI
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - (2 * outerMargin)
		m.viewport.Width = m.width - 4
		m.viewport.Height = viewportHeight
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" || msg.String() == "esc" {
			m.done = true
			m.ExitCode = -1 // ✅ Set special exit code for user abort
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case outputMsg:
		m.appendOutput(string(msg))
		cmds = append(cmds, waitForOutput(m.sub))
		return m, tea.Batch(cmds...) // 🔥 Explicitly update the viewport
	case doneMsg:
		m.done = true
		m.ExitCode = msg.exitCode
		return m, tea.Quit
	case tickMsg:
		cmds = append(cmds, tickCmd(m.Start))
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// Append new output to the viewport
func (m *Model) appendOutput(s string) {
	// Pop/Shift: Keep only last `maxLogLines` lines
	if len(*m.LogLines) > maxLogLines {
		*m.LogLines = (*m.LogLines)[1:]
	}

	// Update viewport content

	// Manually wrap lines for viewport display
	var wrappedLines []string
	for _, line := range *m.LogLines {
		wrappedLines = append(wrappedLines, wrapText(line, m.viewport.Width-4)...) // Adjust for borders
	}

	// Set viewport content
	m.viewport.SetContent(strings.Join(wrappedLines, "\n"))
	m.viewport.GotoBottom()
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text} // Prevent division by zero
	}

	var lines []string
	for len(text) > width {
		lines = append(lines, text[:width]) // ✅ Store a chunk
		text = text[width:]                 // ✅ Remove chunk from original text
	}
	lines = append(lines, text) // ✅ Add the last piece
	return lines
}

// Render the UI
func (m Model) View() string {
	if m.done {
		return ""
	}

	// Compute elapsed time
	elapsed := time.Since(m.Start).Round(time.Second)
	timer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(fmt.Sprintf("(%s)", elapsed))

	// Right-align the timer
	timerStyled := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Right).
		Render(timer)

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Render(m.spinner.View() + m.title + timerStyled)

	// Viewport box
	viewportBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingLeft(1).
		MarginLeft(2).
		Width(m.width - (2 * outerMargin)).
		Render(m.viewport.View())

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		MarginLeft(4).
		Render("press `q` to quit")

	return lipgloss.NewStyle().
		Render(header + "\n" + viewportBox + "\n" + footer)
}
