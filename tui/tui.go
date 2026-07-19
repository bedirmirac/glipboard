package tui

import (
	"bytes"
	"fmt"
	"net/http"

	tea "charm.land/bubbletea/v2"

	"github.com/bedirmirac/glipboard/storage"
)

type model struct {
	choices  []storage.Clipboard
	cursor   int
	selected map[int]struct{}
	db       *storage.Storage
	message  string
}

func initialModel(db *storage.Storage, items []storage.Clipboard) model {
	return model{
		choices:  items,
		selected: make(map[int]struct{}),
		db:       db,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "right", "l":
			if m.cursor+10 < len(m.choices)-1 {
				m.cursor += 10
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "left", "h":
			if m.cursor-10 > 0 {
				m.cursor -= 10
			} else {
				m.cursor = 0
			}
		case "enter":
			_, ok := m.selected[m.cursor]
			if ok {
				/*copy selected logic */
				choice := m.choices[m.cursor]
				var req *http.Request
				var err error

				if choice.Type == "text" {
					req, err = http.NewRequest("POST", "http://127.0.0.1:49321/request?type=text", bytes.NewBufferString(choice.Context))
				} else if choice.Type == "image" {
					req, err = http.NewRequest("POST", "http://127.0.0.1:49321/request?type=image", bytes.NewBufferString(choice.FilePath))
				} else {
					m.message = "Error: Unsupported data type."
					return m, nil
				}

				if err != nil {
					m.message = "Error: Could not create HTTP request."
					return m, nil
				}

				/* Send the request to the background Daemon */
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					m.message = fmt.Sprintf("ERROR: Could not connect to daemon: %v", err)
					return m, nil
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					m.message = fmt.Sprintf("ERROR: Daemon returned status code %d", resp.StatusCode)
					return m, nil
				}

				return m, tea.Quit
			}
		case "d":
			_, ok := m.selected[m.cursor]
			if ok {
				/*delete function*/
				choice := m.choices[m.cursor]
				var req, req1 *http.Request
				var reqs []*http.Request
				var err error
				req, err = http.NewRequest("POST", "http://127.0.0.1:49321/request?type=delete", bytes.NewBufferString(choice.Hash))
				if err == nil {
					reqs = append(reqs, req)
				}
				if choice.Type == "image" {
					req1, err = http.NewRequest("POST", "http://127.0.0.1:49321/request?type=deleteImagePath", bytes.NewBufferString(choice.FilePath))
					if err == nil {
						reqs = append(reqs, req1)
					}
				}

				if err != nil {
					m.message = "Error: Could not create HTTP request."
					return m, nil
				}

				/* Send the request to the background Daemon */
				for _, r := range reqs {
					client := &http.Client{}
					r, err := client.Do(r)
					if err != nil {
						m.message = fmt.Sprintf("ERROR: Could not connect to daemon: %v", err)
						return m, nil
					}
					defer r.Body.Close()
					if r.StatusCode != http.StatusOK {
						m.message = fmt.Sprintf("ERROR: Daemon returned status code %d", r.StatusCode)
						return m, nil
					}
				}

				return m, tea.Quit
			}

		case "ctrl+r":
			var req *http.Request
			var err error
			req, err = http.NewRequest("POST", "http://127.0.0.1:49321/request?type=deleteAll", nil)
			if err != nil {
				m.message = "Error: Could not create HTTP request."
				return m, nil
			}

			/* Send the request to the background Daemon */
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				m.message = fmt.Sprintf("ERROR: Could not connect to daemon: %v", err)
				return m, nil
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				m.message = fmt.Sprintf("ERROR: Daemon returned status code %d", resp.StatusCode)
				return m, nil
			}

			m.message = "All of the contents is deleted!"
			return m, tea.Quit
		case "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				if len(m.selected) > 0 {
					m.message = "You can only select one item"
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	// The header
	s := "--------GLIPBOARD--------"

	pageSize := 10

	startingPoint := (m.cursor / pageSize) * pageSize
	endPoint := startingPoint + pageSize

	if endPoint > len(m.choices) {
		endPoint = len(m.choices)
	}

	itemsOnPage := m.choices[startingPoint:endPoint]

	for i, choice := range itemsOnPage {
		realIndex := startingPoint + i

		cursor := " "
		if m.cursor == realIndex {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		if choice.Context == "" {
			choice.Context = "Image->" + choice.FilePath
		}

		choice.Context = maxCharOfString(choice.Context)
		s += fmt.Sprintf("\n%v [%v] %v\n", cursor, checked, choice.Context)
	}

	if m.message != "" {
		s += fmt.Sprintf("\n---\n%v\n---\n", m.message)
	}
	totalPage := (len(m.choices) + pageSize - 1) / pageSize
	if totalPage == 0 {
		totalPage = 1
	}
	currentPage := (m.cursor / pageSize) + 1
	pageInfo := fmt.Sprintf("\n--- Page %v / %v ---\n", currentPage, totalPage)
	s += pageInfo
	// The footer
	s += "\nPress 'q' to quit, 'enter' twice to copy, 'd' twice to delete one, \n 'ctrl+r' to delete all \n If you pressed 'enter' or 'd' once and select item unintentialy, \n you can deselect it by pressing 'space'."

	return tea.NewView(s)
}
