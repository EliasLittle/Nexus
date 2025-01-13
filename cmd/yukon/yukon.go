package main

import (
	"fmt"
	"os"
	"strings"

	nc "nexus/pkg/client"
	pb "nexus/pkg/proto"

	"github.com/charmbracelet/bubbles/v2/table"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table       table.Model
	client      *nc.NexusClient
	path        string
	children    []string
	err         error
	lastKeyMsg  string
	isLeafNode  bool
	searchInput string
	isSearching bool
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
		table:       t,
		client:      nc.NewNexusClient(conn),
		path:        "/",
		children:    []string{},
		err:         err,
		lastKeyMsg:  "",
		isLeafNode:  false,
		searchInput: "",
		isSearching: false,
	}
}

func (m model) Init() (tea.Model, tea.Cmd) {
	return m, m.fetchChildren
}

func (m model) fetchChildren() tea.Msg {
	children, err := m.client.GetChildren(m.path)
	if err != nil {
		fmt.Printf("Error fetching children: %v\n", err)
		return errMsg{err}
	}

	rows := []table.Row{}
	for _, child := range children {
		dataType := "directory" // default assumption
		// Try to get data type if exists
		if pathType, err := m.client.GetPathType(child); err == nil {
			dataType = pathType
		}
		rows = append(rows, table.Row{child, dataType})
	}

	// Filter rows based on search input if search bar is active
	if m.isSearching {
		filteredRows := []table.Row{}
		for _, row := range rows {
			if strings.Contains(strings.ToLower(row[0]), strings.ToLower(m.searchInput)) {
				filteredRows = append(filteredRows, row)
			}
		}
		rows = filteredRows
	}

	return childrenMsg{rows}
}

func (m model) fetchData() tea.Msg {
	data, err := m.client.GetValue(m.path)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error fetching data from path '%s': %v\n", m.path, err)
		return errMsg{err}
	}

	if data == nil {
		// Handle case where data is nil
		fmt.Printf("No data found at path '%s'\n", m.path)
		return errMsg{fmt.Errorf("no data found at path '%s'", m.path)}
	}

	var valueStr string
	switch v := data.Value.(type) {
	case *pb.DirectValue_IntValue:
		valueStr = fmt.Sprintf("%d", v.IntValue.Value)
	case *pb.DirectValue_FloatValue:
		valueStr = fmt.Sprintf("%f", v.FloatValue.Value)
	case *pb.DirectValue_StringValue:
		valueStr = v.StringValue.Value
	default:
		valueStr = "unknown type"
	}

	return childrenMsg{rows: []table.Row{{"value", valueStr}}}
}

type childrenMsg struct {
	rows []table.Row
}

type errMsg struct {
	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg_type := msg.(type) {
	case tea.KeyMsg:
		m.lastKeyMsg = msg_type.String()
		switch m.isSearching {
		case true:
			switch m.lastKeyMsg {
			case "esc":
				m.isSearching = false
				return m, nil
			case "backspace":
				if len(m.searchInput) > 0 {
					m.searchInput = m.searchInput[:len(m.searchInput)-1]
				}
				return m.filterChildren()
			case "/":
				m.searchInput += msg_type.String()

				m.path = m.searchInput
				return m, m.fetchChildren
			case "ctrl+c":
				return m, tea.Quit
			default:
				// Update search input with typed characters
				m.searchInput += msg_type.String()
				// Filter children based on search input
				return m.filterChildren()
			}
		case false:
			switch m.lastKeyMsg {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				if len(m.table.Rows()) > 0 {
					selected := m.table.SelectedRow()[0]
					newPath := m.path + selected + "/"
					// Check if the selected path has children
					if children, err := m.client.GetChildren(newPath); err == nil {
						if len(children) == 0 {
							if !m.isLeafNode {
								m.isLeafNode = true
								m.path = newPath
								return m, m.fetchData
							}
							return m, nil
						}
						m.isLeafNode = false
						m.path = newPath
						return m, m.fetchChildren
					}
				}
			case "esc", "backspace":
				return m.moveUp()
			case "/":
				m.isSearching = true // Enable search mode
				m.searchInput = m.path
				return m, nil
			default:
				return m, nil
			}
		}

	case childrenMsg:
		m.table.SetRows(msg_type.rows)
		return m, nil

	case errMsg:
		m.err = msg_type.err
		return m, nil

	default:
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) filterChildren() (tea.Model, tea.Cmd) {
	// Fetch all children first
	children := m.children

	lastSegment := strings.Split(m.searchInput, "/")[len(strings.Split(m.searchInput, "/"))-1]

	// Filter rows based on search input
	filteredRows := []table.Row{}
	for _, child := range children {
		if strings.Contains(strings.ToLower(child), strings.ToLower(lastSegment)) {
			dataType := "directory" // default assumption
			if pathType, err := m.client.GetPathType(child); err == nil {
				dataType = pathType
			}
			filteredRows = append(filteredRows, table.Row{child, dataType})
		}
	}

	m.table.SetRows(filteredRows)
	return m, nil
}

func (m model) moveUp() (tea.Model, tea.Cmd) {
	if m.path != "/" {
		// Go up one level
		lastSlash := strings.LastIndex(m.path[:len(m.path)-1], "/")
		if lastSlash == 0 {
			m.path = "/"
		} else {
			m.path = m.path[:lastSlash+1]
		}
		m.isLeafNode = false
		return m, m.fetchChildren
	}
	return m, nil
}

func (m model) moveDown() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	searchBar := ""
	if m.isSearching {
		searchBar = lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("240")).
			Bold(true).
			Render(fmt.Sprintf(" Search: %s ", m.searchInput))
	} else {
		searchBar = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Bold(false).
			Render(fmt.Sprintf(" Search: %s ", m.searchInput))
	}

	return baseStyle.Render(
		fmt.Sprintf("Current path: %s\n\n isSearching: %t\n\n Last Key Pressed: %s\n\n%s\n\nPress q to quit, enter to navigate, backspace/esc to go up",
			m.path,
			m.isSearching,
			m.lastKeyMsg,
			searchBar+m.table.View(),
		))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithKeyboardEnhancements())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
