package main

import (
	"fmt"
	"os"
	"strings"

	nc "nexus/pkg/client"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/grpc"
	// Importing the nexus-client package using a relative path
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table    table.Model
	conn     *grpc.ClientConn
	path     string
	children []string
	err      error
}

func initialModel() model {
	columns := []table.Column{
		{Title: "Path", Width: 40},
		{Title: "Type", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s)

	conn, err := nc.CreateGRPCConnection(nc.DefaultConnection)

	return model{
		table:    t,
		conn:     conn,
		path:     "/",
		children: []string{},
		err:      err,
	}
}

func (m model) Init() tea.Cmd {
	return m.fetchChildren
}

func (m model) fetchChildren() tea.Msg {
	children, err := nc.GetChildren(m.conn, m.path)
	if err != nil {
		return errMsg{err}
	}

	rows := []table.Row{}
	for _, child := range children {
		dataType := "directory" // default assumption
		// Try to get data type if exists
		if data, err := nc.Get(m.conn, child); err == nil && data != nil {
			dataType = "data" // you might want to be more specific based on the data
		}
		rows = append(rows, table.Row{child, dataType})
	}

	return childrenMsg{rows}
}

type childrenMsg struct {
	rows []table.Row
}

type errMsg struct {
	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(m.table.Rows()) > 0 {
				selected := m.table.SelectedRow()[0]
				if m.path == "/" {
					m.path = m.path + selected
				} else {
					m.path = m.path + "/" + selected
				}
				return m, m.fetchChildren
			}
		case "backspace", "esc":
			if m.path != "/" {
				// Go up one level
				lastSlash := strings.LastIndex(m.path[:len(m.path)-1], "/")
				if lastSlash == 0 {
					m.path = "/"
				} else {
					m.path = m.path[:lastSlash]
				}
				return m, m.fetchChildren
			}
		}

	case childrenMsg:
		m.table.SetRows(msg.rows)
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	return baseStyle.Render(
		fmt.Sprintf("Current path: %s\n\n%s\n\nPress q to quit, enter to navigate, backspace/esc to go up",
			m.path,
			m.table.View(),
		))
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
