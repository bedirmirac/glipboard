package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/bedirmirac/glipboard/storage"
)

func StartTUI(s *storage.Storage) {
	items, err := s.Fetch()
	if err != nil {
		fmt.Printf("there's been an error during fetching from database: %v", err)
		os.Exit(1)
	}
	p := tea.NewProgram(initialModel(s, items))
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
