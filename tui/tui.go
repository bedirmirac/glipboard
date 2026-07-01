package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
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
		case "enter":
			_, ok := m.selected[m.cursor]
			if ok {
				/*copy selected logic */
				textToCopy := m.choices[m.cursor].Context
				err := clipboard.WriteAll(textToCopy)
				if err != nil {
					m.message = fmt.Sprintf("The content couldn't be copied: %v", err)
				}
				m.message = "The content is copied."
				return m, tea.Quit
			} else {
				if len(m.selected) > 0 {
					m.message = "You can only select one item"
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "d":
			_, ok := m.selected[m.cursor]
			if ok {
				/*delete function*/
				idToDelete := m.choices[m.cursor].ID
				err := m.db.Delete(idToDelete)
				if err != nil {
					m.message = fmt.Sprintf("The content couldn't be deleted: %v", err)
					return m, tea.Quit
				}
				m.choices = append(m.choices[:m.cursor], m.choices[m.cursor+1:]...)
				m.message = "The content is deleted."
				delete(m.selected, m.cursor)
				m.cursor = 0
				return m, nil
			} else {
				if len(m.selected) > 0 {
					m.message = "You can only select one item"
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "ctrl+r":
			err := m.db.DeleteAll()
			if err != nil {
				m.message = fmt.Sprintf("The contents could't be deleted: %v", err)
				return m, tea.Quit
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
	s += "\nPress 'q' to quit, 'enter' twice to copy, 'd' twice to delete one, 'ctrl+r' to delete all \n If you pressed 'enter' or 'd' once and select item unintentialy, you can deselect it by pressing 'space'."

	return tea.NewView(s)
}
