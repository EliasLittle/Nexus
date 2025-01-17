package main

import (
	"fmt"
	"os"
	"strings"

	nc "nexus/pkg/client"
	pb "nexus/pkg/proto"

	"github.com/charmbracelet/bubbles/v2/table"
	"github.com/charmbracelet/bubbles/v2/textinput"
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
	searchInput textinput.Model
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

	searchInput := textinput.New()
	searchInput.Placeholder = "Search..."
	searchInput.Blur()
	searchInput.CharLimit = 156
	//searchInput.Width = 20

	return model{
		table:       t,
		client:      nc.NewNexusClient(conn),
		path:        "/",
		children:    []string{},
		err:         err,
		lastKeyMsg:  "",
		isLeafNode:  false,
		searchInput: searchInput,
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
		dataType := "Branch" // default assumption
		// Try to get data type if exists
		if pathType, err := m.client.GetPathType(child); err == nil {
			dataType = pathType
		}
		rows = append(rows, table.Row{child, dataType})
	}

	return childrenMsg{rows}
}

func (m model) fetchData() tea.Msg {
	data, err := m.client.Get(m.path)
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
	switch v := data.(type) {
	case *pb.IntValue:
		valueStr = fmt.Sprintf("%d", v.Value)
	case *pb.FloatValue:
		valueStr = fmt.Sprintf("%f", v.Value)
	case *pb.StringValue:
		valueStr = v.Value
	case *pb.DatabaseTable:
		valueStr = fmt.Sprintf("DatabaseTable: %s", v.TableName)
	case *pb.Directory:
		valueStr = fmt.Sprintf("Directory: %s", v.DirectoryPath)
	case *pb.IndividualFile:
		valueStr = fmt.Sprintf("File (%s): %s", v.FileType, v.FilePath)
	case *pb.EventStream:
		valueStr = fmt.Sprintf("EventStream: %s@%s", v.Topic, v.Server)
	default:
		valueStr = fmt.Sprintf("value: %v, has unknown type: %T", v, v)
	}

	return childrenMsg{rows: []table.Row{{"value", valueStr}}}
}

type childrenMsg struct {
	rows []table.Row
}

type errMsg struct {
	err error
}

func (m *model) filterChildren() tea.Msg {
	// Fetch all children first
	children, err := m.client.GetChildren(m.path)
	if err != nil {
		fmt.Printf("Error fetching children: %v\n", err)
		return errMsg{err}
	}

	lastSegment := strings.Split(m.searchInput.Value(), "/")[len(strings.Split(m.searchInput.Value(), "/"))-1]

	// Filter rows based on search input
	filteredRows := []table.Row{}
	for _, child := range children {
		if strings.Contains(strings.ToLower(child), strings.ToLower(lastSegment)) {
			dataType := "Branch" // default assumption
			if pathType, err := m.client.GetPathType(child); err == nil {
				dataType = pathType
			}
			filteredRows = append(filteredRows, table.Row{child, dataType})
		}
	}

	return childrenMsg{rows: filteredRows}
}

func (m model) moveUp() (tea.Model, tea.Cmd) {
	if m.path != "/" {
		// Go up one level
		switch lastSlash := strings.LastIndex(m.path[:len(m.path)-1], "/"); lastSlash {
		case -1:
			m.path = "/"
		case 0:
			m.path = "/"
		default:
			m.path = m.path[:lastSlash+1]
		}

		m.isLeafNode = false
		m.searchInput.SetValue(m.path)
		return m, m.fetchChildren
	}
	return m, nil
}

func (m model) moveDown(newPath string) (tea.Model, tea.Cmd) {
	// Check if the selected path has children
	if children, err := m.client.GetChildren(newPath); err == nil {
		if len(children) == 0 {
			if !m.isLeafNode {
				m.isLeafNode = true
				m.path = newPath
				m.searchInput.SetValue(newPath)
				return m, m.fetchData
			}
			// If the path is a leaf node, do nothing
			return m, nil
		}
		m.isLeafNode = false
		m.path = newPath
		m.searchInput.SetValue(newPath)
		return m, m.fetchChildren
	}
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.lastKeyMsg = msg.String()
		switch m.isSearching {
		case true:
			switch m.lastKeyMsg {
			case "ctrl+c":
				return m, tea.Quit
			case "esc", "enter":
				m.isSearching = false
				m.searchInput.Blur()
				return m, nil
			case "backspace":
				char := m.searchInput.Value()[len(m.searchInput.Value())-1]
				if char == '/' {
					m.searchInput.SetValue(m.searchInput.Value()[:len(m.searchInput.Value())-1])
					return m.moveUp()
				}
				m.searchInput, _ = m.searchInput.Update(msg)
				return m, m.filterChildren
			case "/":
				m.searchInput, _ = m.searchInput.Update(msg)
				return m.moveDown(m.searchInput.Value())
			default:
				m.searchInput, _ = m.searchInput.Update(msg)
				return m, m.filterChildren
			}

		case false:
			switch m.lastKeyMsg {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				if len(m.table.Rows()) > 0 {
					selected := m.table.SelectedRow()[0]
					newPath := m.path + selected + "/"
					return m.moveDown(newPath)
				}
			case "esc", "backspace":
				return m.moveUp()
			case "/":
				m.isSearching = true
				m.searchInput.Focus()
				return m, nil
			default:
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}
		}

	case childrenMsg:
		m.table.SetRows(msg.rows)
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	default:
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
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
			Render(fmt.Sprintf(" Search: %s ", m.searchInput.View()))
	} else {
		searchBar = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Bold(false).
			Render(fmt.Sprintf(" Search: %s ", m.searchInput.View()))
	}

	return baseStyle.Render(
		//fmt.Sprintf("Current path: %s\n\n isSearching: %t\n\n Last Key Pressed: %s\n\n%s\n\nPress q to quit, enter to navigate, backspace/esc to go up",
		fmt.Sprintf("Path: %s\n\n%s\n\n%s\n\nPress q to quit, / to search, enter to navigate, backspace/esc to go up",
			m.path,
			searchBar,
			m.table.View(),
		))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithKeyboardEnhancements())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
