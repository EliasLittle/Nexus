package main

import (
	"fmt"
	"nexus/pkg/logger"
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
	table         table.Model
	client        *nc.NexusClient
	path          string
	children      []string
	err           error
	lastKeyMsg    string
	isLeafNode    bool
	searchInput   textinput.Model
	isSearching   bool
	streamingData bool
}

type rowDataMsg struct {
	rows    []table.Row
	message string
}

type streamDataMsg struct {
	channel <-chan []byte
	rows    []table.Row
	message string
}

type errMsg struct {
	err error
}

type moveDownResponse struct {
	newPath    string
	isLeafNode bool
}

type moveUpResponse struct {
	newPath string
}

func initialModel() model {
	log := logger.GetLogger()

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
	if err != nil {
		log.Error("Failed to create gRPC connection", "error", err)
		os.Exit(1)
	}

	searchInput := textinput.New()
	searchInput.Placeholder = "Search..."
	searchInput.Blur()
	searchInput.CharLimit = 156
	//searchInput.Width = 20

	return model{
		table:         t,
		client:        nc.NewNexusClient(conn),
		path:          "/",
		children:      []string{},
		err:           err,
		lastKeyMsg:    "",
		isLeafNode:    false,
		searchInput:   searchInput,
		isSearching:   false,
		streamingData: false,
	}
}

func (m model) Init() (tea.Model, tea.Cmd) {
	return m, fetchChildrenCmd(m.client, m.path)
}

func fetchChildrenCmd(client *nc.NexusClient, path string) tea.Cmd {
	log := logger.GetLogger()
	return func() tea.Msg {
		children, err := client.GetChildren(path)
		if err != nil {
			log.Error("Failed to fetch children", "error", err)
			return errMsg{err}
		}

		rows := []table.Row{}
		for _, child := range children {
			dataType := "Branch" // default assumption
			// Try to get data type if exists
			if pathType, err := client.GetPathType(child); err == nil {
				dataType = pathType
			}
			rows = append(rows, table.Row{child, dataType})
		}

		return rowDataMsg{rows, "fetchChildren"}
	}
}

// Function to process messages from the channel
func processStream(messageChan <-chan []byte) tea.Cmd {
	log := logger.GetLogger()
	return func() tea.Msg {
		message := <-messageChan
		log.Debug("Stream message received", "message", message)
		return streamDataMsg{
			channel: messageChan,
			rows:    []table.Row{{"Streaming Data", string(message)}},
			message: "streamUpdate",
		}
	}
}

func fetchDataCmd(client *nc.NexusClient, path string) tea.Cmd {
	log := logger.GetLogger()
	return func() tea.Msg {
		log.Info("Fetching data", "path", path)
		data, err := client.Get(path)
		if err != nil {
			log.Error("Failed to fetch data", "path", path, "error", err)
			return errMsg{err}
		}

		if data == nil {
			log.Error("No data found", "path", path)
			return errMsg{fmt.Errorf("no data found at path '%s'", path)}
		}

		var rows []table.Row

		switch v := data.(type) {
		case *pb.IntValue:
			valueStr := fmt.Sprintf("%d", v.Value)
			rows = append(rows, table.Row{path, valueStr})
		case *pb.FloatValue:
			valueStr := fmt.Sprintf("%f", v.Value)
			rows = append(rows, table.Row{path, valueStr})
		case *pb.StringValue:
			valueStr := v.Value
			rows = append(rows, table.Row{path, valueStr})
		case *pb.DatabaseTable:
			valueStr := fmt.Sprintf("DatabaseTable: %s", v.TableName)
			rows = append(rows, table.Row{path, valueStr})
		case *pb.Directory:
			valueStr := fmt.Sprintf("Directory: %s", v.DirectoryPath)
			rows = append(rows, table.Row{path, valueStr})
		case *pb.IndividualFile:
			valueStr := fmt.Sprintf("File (%s): %s", v.FileType, v.FilePath)
			rows = append(rows, table.Row{path, valueStr})
		case *pb.EventStream:
			log.Info("Event Stream: ", "server", v.Server, "topic", v.Topic)

			messageChan, err := nc.GetEventStream(v)
			if err != nil {
				log.Error("Failed to get event stream", "error", err)
				os.Exit(1)
			}

			log.Debug("Stream initialized", "channel", messageChan)

			return streamDataMsg{
				channel: messageChan,
				rows:    []table.Row{{"Streaming Data", "Waiting for messages..."}},
				message: "streamInit",
			}

		default:
			valueStr := fmt.Sprintf("value: %v, has unknown type: %T", v, v)
			rows = append(rows, table.Row{path, valueStr})
		}

		return rowDataMsg{rows, "fetchData"}
	}
}

func filterChildrenCmd(client *nc.NexusClient, path string, searchInput textinput.Model) tea.Cmd {
	log := logger.GetLogger()
	// Fetch all children first
	return func() tea.Msg {
		children, err := client.GetChildren(path)
		if err != nil {
			log.Error("Failed to fetch children", "error", err)
			return errMsg{err}
		}

		lastSegment := strings.Split(searchInput.Value(), "/")[len(strings.Split(searchInput.Value(), "/"))-1]

		// Filter rows based on search input
		filteredRows := []table.Row{}
		for _, child := range children {
			if strings.Contains(strings.ToLower(child), strings.ToLower(lastSegment)) {
				dataType := "Branch" // default assumption
				if pathType, err := client.GetPathType(child); err == nil {
					dataType = pathType
				}
				filteredRows = append(filteredRows, table.Row{child, dataType})
			}
		}

		log.Debug("Filtered children", "count", len(filteredRows))
		return rowDataMsg{rows: filteredRows, message: "filterChildren"}
	}
}

func moveUpCmd(path string) tea.Cmd {
	log := logger.GetLogger()
	return func() tea.Msg {
		log.Debug("Moving up", "path", path)
		if path != "/" {
			// Go up one level
			switch lastSlash := strings.LastIndex(path[:len(path)-1], "/"); lastSlash {
			case -1:
				path = "/"
			case 0:
				path = "/"
			default:
				path = path[:lastSlash+1]
			}

			return moveUpResponse{
				newPath: path,
			}
		}
		return nil
	}
}

func moveDownCmd(client *nc.NexusClient, newPath string, isLeafNode bool) tea.Cmd {
	log := logger.GetLogger()
	// Check if the selected path has children
	return func() tea.Msg {
		children, err := client.GetChildren(newPath)
		if err != nil {
			log.Error("Failed to fetch children", "error", err)
			return errMsg{err}
		}
		if len(children) == 0 {
			if !isLeafNode {
				log.Debug("Leaf node", "path", newPath)
				return moveDownResponse{
					newPath:    newPath,
					isLeafNode: true,
				}
			}
			// If the path is a leaf node, do nothing
			return nil
		}
		return moveDownResponse{
			newPath:    newPath,
			isLeafNode: false,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log := logger.GetLogger()
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case moveDownResponse:
		m.isLeafNode = msg.isLeafNode
		m.path = msg.newPath
		m.searchInput.SetValue(msg.newPath)
		log.Debug("Moving down", "new_path", msg.newPath)
		if msg.isLeafNode {
			cmd = fetchDataCmd(m.client, msg.newPath)
			cmds = append(cmds, cmd)
		} else {
			cmd = fetchChildrenCmd(m.client, msg.newPath)
			cmds = append(cmds, cmd)
		}
	case moveUpResponse:
		m.isLeafNode = false
		m.streamingData = false
		m.path = msg.newPath
		m.searchInput.SetValue(msg.newPath)
		log.Debug("Moving up", "new_path", msg.newPath)
		cmd = fetchChildrenCmd(m.client, msg.newPath)
		cmds = append(cmds, cmd)
	case streamDataMsg:
		log.Debug("Stream data message received", "message", msg.message)
		if msg.message == "streamInit" {
			m.streamingData = true
		}
		if m.streamingData {
			m.table.SetRows(msg.rows)
			log.Debug("Continuing stream", "channel", msg.channel)
			cmds = append(cmds, processStream(msg.channel))
		}
	case tea.KeyMsg:
		m.lastKeyMsg = msg.String()
		switch m.isSearching {
		case true:
			switch m.lastKeyMsg {
			case "ctrl+c":
				cmds = append(cmds, tea.Quit)
			case "esc", "enter":
				m.isSearching = false
				m.searchInput.Blur()
			case "backspace":
				char := m.searchInput.Value()[len(m.searchInput.Value())-1]
				if char == '/' {
					m.searchInput.SetValue(m.searchInput.Value()[:len(m.searchInput.Value())-1])
					// Perform type assertion to model
					cmd = moveUpCmd(m.path)
					cmds = append(cmds, cmd)
				}
				m.searchInput, _ = m.searchInput.Update(msg)

				cmd = filterChildrenCmd(m.client, m.path, m.searchInput)
				cmds = append(cmds, cmd)
			case "/":
				m.searchInput, cmd = m.searchInput.Update(msg)
				cmds = append(cmds, cmd)
				cmd = moveDownCmd(m.client, m.searchInput.Value(), m.isLeafNode)
				cmds = append(cmds, cmd)
			default:
				m.searchInput, cmd = m.searchInput.Update(msg)
				cmds = append(cmds, cmd)
				cmd = filterChildrenCmd(m.client, m.path, m.searchInput)
				cmds = append(cmds, cmd)
			}

		case false:
			switch m.lastKeyMsg {
			case "q", "ctrl+c":
				cmds = append(cmds, tea.Quit)
			case "enter":
				if len(m.table.Rows()) > 0 {
					selected := m.table.SelectedRow()[0]
					newPath := m.path + selected + "/"
					cmd = moveDownCmd(m.client, newPath, m.isLeafNode)
					cmds = append(cmds, cmd)
				}
			case "esc", "backspace":
				cmd = moveUpCmd(m.path)
				cmds = append(cmds, cmd)
			case "/":
				m.searchInput.SetCursor(len(m.searchInput.Value()))
				m.isSearching = true
				m.searchInput.Focus()
			default:
				m.table, cmd = m.table.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case rowDataMsg:
		log.Debug("Children message received", "message", msg.message, "rows", msg.rows)
		m.table.SetRows(msg.rows)

	case errMsg:
		m.err = msg.err

	}

	return m, tea.Batch(cmds...)
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
	log := logger.GetLogger()

	log.Info("Starting Yukon")
	p := tea.NewProgram(initialModel(), tea.WithKeyboardEnhancements())

	if _, err := p.Run(); err != nil {
		log.Error("Error running program", "error", err)
		os.Exit(1)
	}
}
