package stack_component_select

import (
	tea "github.com/charmbracelet/bubbletea"
	mouseZone "github.com/lrstanley/bubblezone"
)

// Execute starts the TUI app and returns the selected items from the views
func Execute(commands []string, stacksComponentsMap map[string]any, componentsStacksMap map[string]any) (*App, error) {
	mouseZone.NewGlobal()
	mouseZone.SetEnabled(true)
	app := NewApp(commands, stacksComponentsMap, componentsStacksMap)
	p := tea.NewProgram(app, tea.WithMouseCellMotion())

	_, err := p.Run()
	if err != nil {
		return nil, err
	}

	return app, nil
}
