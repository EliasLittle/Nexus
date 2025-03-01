package main

import (
	"fmt"
	"nexus/pkg/logger"
	"os"
	"strconv"
	"strings"

	nc "nexus/pkg/client"
	pb "nexus/pkg/proto"

	"flag"

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
	initPath      string
	children      []*pb.ChildInfo
	rows          []table.Row
	err           error
	lastKeyMsg    string
	isLeafNode    bool
	searchInput   textinput.Model
	isSearching   bool
	streamingData bool
	showPopup     bool
	popupType     string // "add" or "delete"
	popupInput    textinput.Model
	// New fields for type selection and form
	selectedType    int
	showTypeSelect  bool
	showForm        bool
	formInputs      []textinput.Model
	formInputCursor int
	formLabels      []string
	currentFormType string
}

// Add constants for data types
const (
	typeValue = iota
	typeEventStream
	typeFile
	typeDirectory
	typeDBTable
)

var dataTypes = []string{
	"Value",
	"Event Stream",
	"File",
	"Directory",
	"Database Table",
}

var formFields = map[string][]string{
	"Value": {
		"Value",
	},
	"Event Stream": {
		"Server",
		"Topic",
	},
	"File": {
		"File Path",
		"File Type",
	},
	"Directory": {
		"Directory Path",
		"File Type",
	},
	"Database Table": {
		"DB Type",
		"Host",
		"Port",
		"DB Name",
		"Table Name",
	},
}

type rowDataMsg struct {
	rows    []table.Row
	message string
}

type streamDataMsg struct {
	channel <-chan []byte
	row     table.Row
	label   string
	rowNum  int
	message string
}

type errMsg struct {
	err error
}

type moveDownResponse struct {
	newPath    string
	children   []*pb.ChildInfo
	isLeafNode bool
}

type moveUpResponse struct {
	newPath string
}

type addPathResponse struct {
	success bool
	message string
}

type deletePathResponse struct {
	success bool
	message string
}

func initialModel(initialPath, host string, port int) model {
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

	conn, err := nc.CreateGRPCConnection(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Error("Failed to create gRPC connection", "error", err)
		os.Exit(1)
	}

	searchInput := textinput.New()
	searchInput.Placeholder = "Search..."
	searchInput.Blur()
	searchInput.CharLimit = 156

	// Initialize form inputs for the largest possible form
	maxFormFields := 0
	for _, fields := range formFields {
		if len(fields) > maxFormFields {
			maxFormFields = len(fields)
		}
	}

	formInputs := make([]textinput.Model, maxFormFields)
	for i := range formInputs {
		input := textinput.New()
		input.CharLimit = 156
		input.Blur()
		formInputs[i] = input
	}

	return model{
		table:           t,
		client:          nc.NewNexusClient(conn),
		path:            "/",
		initPath:        initialPath,
		children:        nil,
		err:             err,
		lastKeyMsg:      "",
		isLeafNode:      false,
		searchInput:     searchInput,
		isSearching:     false,
		streamingData:   false,
		showPopup:       false,
		popupType:       "",
		popupInput:      textinput.New(),
		selectedType:    0,
		showTypeSelect:  false,
		showForm:        false,
		formInputs:      formInputs,
		formInputCursor: 0,
		formLabels:      nil,
		currentFormType: "",
	}
}

func (m model) Init() (tea.Model, tea.Cmd) {
	if m.initPath != "/" {
		return m, moveDownCmd(m.client, m.initPath, m.isLeafNode)
	}
	return m, fetchRowsCmd(m.client, m.path)
}

func fetchRowsCmd(client *nc.NexusClient, path string) tea.Cmd {
	log := logger.GetLogger()
	cmds := []tea.Cmd{}
	children, err := client.GetChildren(path)
	if err != nil {
		log.Error("Failed to fetch children", "error", err)
		cmds = append(cmds, func() tea.Msg {
			return errMsg{err}
		})
	}

	rows := []table.Row{}
	for index, child := range children {
		log.Info("Fetching data", "path", path+child.Name)
		data, dataType, err := client.GetFull(path + child.Name)
		if err != nil {
			log.Error("Failed to fetch data", "path", path, "error", err)
			cmds = append(cmds, func() tea.Msg {
				return errMsg{err}
			})
		}

		/*
			if data == nil {
				log.Error("No data found", "path", path)
				cmds = append(cmds, func() tea.Msg {
					return errMsg{fmt.Errorf("no data found at path '%s'", path)}
				})
			} */

		log.Debug("Matching data type", "type", dataType)
		switch v := data.(type) {
		case nil:
			if dataType == "InternalNode" {
				log.Debug("Internal node", "path", path+child.Name)
				rows = append(rows, table.Row{child.Name, "⮧"})
			} else {
				log.Error("No data found", "path", path+child.Name)
				rows = append(rows, table.Row{child.Name, "No data found"})
			}
		case *pb.IntValue:
			valueStr := fmt.Sprintf("%d", v.Value)
			rows = append(rows, table.Row{child.Name, valueStr})
		case *pb.FloatValue:
			valueStr := fmt.Sprintf("%f", v.Value)
			rows = append(rows, table.Row{child.Name, valueStr})
		case *pb.StringValue:
			valueStr := v.Value
			rows = append(rows, table.Row{child.Name, valueStr})
		case *pb.DatabaseTable:
			valueStr := fmt.Sprintf("DatabaseTable: %s", v.TableName)
			rows = append(rows, table.Row{child.Name, valueStr})
		case *pb.Directory:
			valueStr := fmt.Sprintf("Directory: %s", v.DirectoryPath)
			rows = append(rows, table.Row{child.Name, valueStr})
		case *pb.IndividualFile:
			valueStr := fmt.Sprintf("File (%s): %s", v.FileType, v.FilePath)
			rows = append(rows, table.Row{child.Name, valueStr})
		case *pb.EventStream:
			log.Info("Event Stream: ", "server", v.Server, "topic", v.Topic)

			messageChan, err := nc.GetEventStream(v)
			if err != nil {
				log.Error("Failed to get event stream", "error", err)
				os.Exit(1)
			}

			log.Debug("Stream initialized", "channel", messageChan)

			rows = append(rows, table.Row{child.Name, "Waiting for messages..."})

			cmds = append(cmds, processStream(messageChan, child.Name, index, "streamInit"))

		default:
			log.Debug("Unknown data type", "type", dataType)
			valueStr := fmt.Sprintf("value: %v, has unknown type: %T", v, v)
			rows = append(rows, table.Row{child.Name, valueStr})
		}

		//return rowDataMsg{rows, "fetchData"}
	}
	cmds = append(cmds, func() tea.Msg {
		return rowDataMsg{rows, "fetchData"}
	})
	return tea.Batch(cmds...)
}

/*
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
			rows = append(rows, table.Row{child.Name, child.Type})
		}

		sort.Slice(rows, func(i, j int) bool {
			return rows[i][0] < rows[j][0]
		})
		return rowDataMsg{rows, "fetchChildren"}
	}
}

func fetchDataCmd(client *nc.NexusClient, path string) tea.Cmd {
	log := logger.GetLogger()
	return func() tea.Msg {
		log.Info("Fetching data", "path", path)
		data, _, err := client.GetFull(path)
		if err != nil {
			log.Error("Failed to fetch data", "path", path, "error", err)
			return errMsg{err}
		}

		if data == nil {
			log.Error("fetchDataCmd:No data found", "path", path)
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
				row:     table.Row{"Waiting for messages..."},
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
			if strings.Contains(strings.ToLower(child.Name), strings.ToLower(lastSegment)) {
				filteredRows = append(filteredRows, table.Row{child.Name, child.Type})
			}
		}

		log.Debug("Filtered children", "count", len(filteredRows))
		sort.Slice(filteredRows, func(i, j int) bool {
			return filteredRows[i][0] < filteredRows[j][0]
		})
		return rowDataMsg{rows: filteredRows, message: "filterChildren"}
	}
}
*/

func filterRows(rows []table.Row, searchInput textinput.Model) []table.Row {
	lastSegment := strings.Split(searchInput.Value(), "/")[len(strings.Split(searchInput.Value(), "/"))-1]

	filteredRows := []table.Row{}
	for _, row := range rows {
		if strings.Contains(strings.ToLower(row[0]), strings.ToLower(lastSegment)) {
			filteredRows = append(filteredRows, row)
		}
	}
	return filteredRows
}

// Function to process messages from the channel
func processStream(messageChan <-chan []byte, label string, rowNum int, note string) tea.Cmd {
	log := logger.GetLogger()
	return func() tea.Msg {
		message := <-messageChan
		log.Debug("Stream message received", "message", message)
		return streamDataMsg{
			channel: messageChan,
			row:     table.Row{label, string(message)},
			label:   label,
			rowNum:  rowNum,
			message: note,
		}
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
					children:   children,
					isLeafNode: true,
				}
			}
			// If the path is a leaf node, do nothing
			return nil
		}
		return moveDownResponse{
			newPath:    newPath,
			children:   children,
			isLeafNode: false,
		}
	}
}

func addPathCmd(client *nc.NexusClient, path string, dataType string, popupInput string, formInputs []textinput.Model) tea.Cmd {
	log := logger.GetLogger()
	return func() tea.Msg {
		var err error
		switch dataType {
		case "Value":
			err = client.PublishValue(path+popupInput, formInputs[0].Value())
			if err != nil {
				log.Error("Failed to add value", "error", err)
				return addPathResponse{success: false, message: err.Error()}
			}
			return addPathResponse{success: true, message: "Value added"}
		case "Event Stream":
			es := nc.CreateEventStream(formInputs[1].Value())
			es.Server = formInputs[0].Value()
			err = client.PublishEventStream(path+popupInput, es)
		case "File":
			file := nc.CreateIndividualFile(formInputs[0].Value())
			err = client.PublishIndividualFile(path+popupInput, file)
		case "Directory":
			dir, err := nc.CreateDirectory(formInputs[0].Value())
			if err == nil {
				err = client.PublishDirectory(path+popupInput, dir)
			}
		case "Database Table":
			port, _ := strconv.ParseInt(formInputs[2].Value(), 10, 32)
			table := nc.CreateDatabaseTable(
				formInputs[0].Value(), // DB Type
				formInputs[1].Value(), // Host
				int32(port),           // Port
				formInputs[3].Value(), // DB Name
				formInputs[4].Value(), // Table Name
			)
			err = client.PublishDatabaseTable(path+popupInput, table)
		}

		if err != nil {
			log.Error("Failed to add path", "error", err)
			return addPathResponse{success: false, message: err.Error()}
		}
		return addPathResponse{success: true, message: dataType + " added"}
	}
}

func deletePathCmd(client *nc.NexusClient, path string) tea.Cmd {
	log := logger.GetLogger()
	return func() tea.Msg {
		err := client.Delete(path)
		if err != nil {
			log.Error("Failed to delete path", "error", err)
			return deletePathResponse{
				success: false,
				message: fmt.Sprintf("Failed to delete path: %v", err),
			}
		}
		return deletePathResponse{
			success: true,
			message: "Path deleted successfully",
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log := logger.GetLogger()
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case moveDownResponse:
		m.isLeafNode = msg.isLeafNode
		m.path = msg.newPath
		m.searchInput.SetValue(msg.newPath)
		m.children = msg.children
		log.Debug("Moving down", "new_path", msg.newPath)
		cmds = append(cmds, fetchRowsCmd(m.client, msg.newPath))
	case moveUpResponse:
		m.isLeafNode = false
		m.streamingData = false
		m.path = msg.newPath
		m.children, _ = m.client.GetChildren(msg.newPath)
		m.searchInput.SetValue(msg.newPath)
		log.Debug("Moving up", "new_path", msg.newPath)
		cmd = fetchRowsCmd(m.client, msg.newPath)
		cmds = append(cmds, cmd)
	case addPathResponse:
		m.showPopup = false
		m.showTypeSelect = false
		m.showForm = false
		if msg.success {
			cmd = fetchRowsCmd(m.client, m.path)
			cmds = append(cmds, cmd)
		} else {
			m.err = fmt.Errorf(msg.message)
		}
	case deletePathResponse:
		m.showPopup = false
		if msg.success {
			cmd = fetchRowsCmd(m.client, m.path)
			cmds = append(cmds, cmd)
		} else {
			m.err = fmt.Errorf(msg.message)
		}
	case streamDataMsg:
		log.Debug("Stream data message received", "message", msg.message)
		if msg.message == "streamInit" {
			m.streamingData = true
		}
		if m.streamingData {
			rows := m.table.Rows()
			rows[msg.rowNum] = msg.row
			m.table.SetRows(rows)
			log.Debug("Continuing stream", "channel", msg.channel)
			cmds = append(cmds, processStream(msg.channel, msg.label, msg.rowNum, "streamUpdate"))
		}
	case tea.KeyMsg:
		m.lastKeyMsg = msg.String()

		if m.showPopup {
			if m.showTypeSelect {
				switch msg.String() {
				case "esc":
					m.showPopup = false
					m.showTypeSelect = false
				case "enter":
					m.showTypeSelect = false
					m.showForm = true
					m.currentFormType = dataTypes[m.selectedType]
					m.formLabels = formFields[m.currentFormType]
					// Reset and focus first input
					for i := range m.formInputs {
						m.formInputs[i].Reset()
						m.formInputs[i].Blur()
					}
					m.formInputCursor = 0
					m.formInputs[0].Focus()
				case "up":
					m.selectedType--
					if m.selectedType < 0 {
						m.selectedType = len(dataTypes) - 1
					}
				case "down":
					m.selectedType = (m.selectedType + 1) % len(dataTypes)
				}
			} else if m.showForm {
				switch msg.String() {
				case "esc":
					m.showPopup = false
					m.showForm = false
				case "enter":
					if m.formInputCursor == len(m.formLabels)-1 {
						// Submit form
						cmd = addPathCmd(m.client, m.path, m.currentFormType, m.popupInput.Value(), m.formInputs)
						cmds = append(cmds, cmd)
					} else {
						// Move to next input
						m.formInputs[m.formInputCursor].Blur()
						m.formInputCursor++
						m.formInputs[m.formInputCursor].Focus()
					}
				case "tab":
					// Move to next input
					m.formInputs[m.formInputCursor].Blur()
					m.formInputCursor = (m.formInputCursor + 1) % len(m.formLabels)
					m.formInputs[m.formInputCursor].Focus()
				case "shift+tab":
					// Move to previous input
					m.formInputs[m.formInputCursor].Blur()
					m.formInputCursor--
					if m.formInputCursor < 0 {
						m.formInputCursor = len(m.formLabels) - 1
					}
					m.formInputs[m.formInputCursor].Focus()
				default:
					// Update the focused input
					m.formInputs[m.formInputCursor], cmd = m.formInputs[m.formInputCursor].Update(msg)
					cmds = append(cmds, cmd)
				}
			} else {
				switch msg.String() {
				case "esc":
					m.showPopup = false
					m.popupInput.Reset()
				case "enter":
					if m.popupType == "add" {
						m.showTypeSelect = true
					} else if m.popupType == "delete" {
						if len(m.table.Rows()) > 0 {
							selected := m.table.SelectedRow()[0]
							cmd = deletePathCmd(m.client, m.path+selected)
							cmds = append(cmds, cmd)
						}
					}
				default:
					if m.popupType == "add" {
						m.popupInput, cmd = m.popupInput.Update(msg)
						cmds = append(cmds, cmd)
					}
				}
			}
			return m, tea.Batch(cmds...)
		}

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
				m.table.SetRows(filterRows(m.rows, m.searchInput))
				//cmd = filterChildrenCmd(m.client, m.path, m.searchInput)
				//cmds = append(cmds, cmd)
			case "/":
				if m.searchInput.Value() != "" { // Check if searchInput is empty
					char := m.searchInput.Value()[len(m.searchInput.Value())-1]
					if char != '/' {
						m.searchInput, cmd = m.searchInput.Update(msg)
						cmds = append(cmds, cmd)
						cmd = moveDownCmd(m.client, m.searchInput.Value(), m.isLeafNode)
						cmds = append(cmds, cmd)
					}
				}
			default:
				m.searchInput, cmd = m.searchInput.Update(msg)
				cmds = append(cmds, cmd)
				m.table.SetRows(filterRows(m.rows, m.searchInput))
				//cmd = filterChildrenCmd(m.client, m.path, m.searchInput)
				//cmds = append(cmds, cmd)
			}

		case false:
			switch m.lastKeyMsg {
			case "q", "ctrl+c":
				cmds = append(cmds, tea.Quit)
			case "enter":
				if len(m.rows) > 0 {
					selected := m.table.SelectedRow()[0]
					newPath := m.path + selected + "/"
					cmd = moveDownCmd(m.client, newPath, m.isLeafNode)
					cmds = append(cmds, cmd)
				} else {
					log.Debug("No rows selected", "path", m.path)
				}
			case "esc", "backspace":
				cmd = moveUpCmd(m.path)
				cmds = append(cmds, cmd)
			case "/":
				m.searchInput.SetCursor(len(m.searchInput.Value()))
				m.isSearching = true
				m.searchInput.Focus()
			case "a":
				m.showPopup = true
				m.popupType = "add"
				m.popupInput.Focus()
				m.popupInput.Placeholder = "Enter path name..."
			case "d":
				m.showPopup = true
				m.popupType = "delete"
			default:
				log.Debug("Default table key press", "key", m.lastKeyMsg)
				m.table, cmd = m.table.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case rowDataMsg:
		log.Debug("Row data message received", "message", msg.message, "rows", msg.rows)
		m.rows = msg.rows
		m.table.SetRows(m.rows)
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

	mainView := baseStyle.Render(
		fmt.Sprintf("Path: %s\n\n%s\n\n%s\n\nPress q to quit, / to search, enter to navigate, backspace/esc to go up, a to add, d to delete",
			m.path,
			searchBar,
			m.table.View(),
		))

	if m.showPopup {
		popup := lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(50)

		var popupContent string
		if m.showTypeSelect {
			// Build type selection menu
			var typeList strings.Builder
			for i, t := range dataTypes {
				if i == m.selectedType {
					typeList.WriteString("> " + t + "\n")
				} else {
					typeList.WriteString("  " + t + "\n")
				}
			}
			popupContent = fmt.Sprintf("Select Data Type\n\n%s\n\nPress enter to select, esc to cancel", typeList.String())
		} else if m.showForm {
			// Build form view
			var form strings.Builder
			form.WriteString(fmt.Sprintf("Add New %s\n\n", m.currentFormType))
			form.WriteString(fmt.Sprintf("Path: %s\n\n", m.popupInput.Value()))

			// Add form fields
			for i, label := range m.formLabels {
				if i == m.formInputCursor {
					form.WriteString(fmt.Sprintf("> %s: %s\n", label, m.formInputs[i].View()))
				} else {
					form.WriteString(fmt.Sprintf("  %s: %s\n", label, m.formInputs[i].View()))
				}
			}

			form.WriteString("\nPress enter to confirm, tab to move between fields, esc to cancel")
			popupContent = form.String()
		} else if m.popupType == "add" {
			popupContent = fmt.Sprintf("Enter Path Name\n\n%s\n\nPress enter to continue, esc to cancel", m.popupInput.View())
		} else if m.popupType == "delete" {
			if len(m.table.Rows()) > 0 {
				selected := m.table.SelectedRow()[0]
				popupContent = fmt.Sprintf("Delete Path\n\nAre you sure you want to delete '%s'?\n\nPress enter to confirm, esc to cancel", selected)
			} else {
				popupContent = "No path selected to delete\n\nPress esc to cancel"
			}
		}

		return lipgloss.JoinVertical(lipgloss.Center,
			mainView,
			"\n",
			popup.Render(popupContent),
		)
	}

	return mainView
}

func main() {
	log := logger.GetLogger()

	// Add flags for the initial path, host, and port
	var initialPath, host string
	var port int
	flag.StringVar(&initialPath, "path", "/", "Initial path to start from")
	flag.StringVar(&host, "host", "localhost", "Host of the server")
	flag.IntVar(&port, "port", 50051, "Port of the server")
	flag.Parse()

	log.Info("Starting Yukon", "host", host, "port", port)
	// Use host and port as needed in your application logic
	p := tea.NewProgram(initialModel(initialPath, host, port), tea.WithKeyboardEnhancements())

	if _, err := p.Run(); err != nil {
		log.Error("Error running program", "error", err)
		os.Exit(1)
	}
}
