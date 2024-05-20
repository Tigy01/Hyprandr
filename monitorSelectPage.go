package main

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type monitorSelectPage struct {
	selection    int
	monitors     map[string]*monitor
	monitorNames []string
}

func (p monitorSelectPage) New(monitors map[string]*monitor) monitorSelectPage {
	monitorNames := make([]string, 0)
	for n := range monitors {
		monitorNames = append(monitorNames, n)
	}
	slices.Sort(monitorNames)
	return monitorSelectPage{
		selection:    0,
		monitors:     monitors,
		monitorNames: monitorNames,
	}
}

func (m monitorSelectPage) Init() tea.Cmd {
	return nil
}

func (m monitorSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowWidth = msg.Width
		windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			m.selection = max(m.selection-1, 0)
			return m, nil
		case "j":
			m.selection = min(m.selection+1, len(m.monitorNames)-1)
			return m, nil
		case "enter":
			nextPage := resolutionSelectPage{}.New(
				m.monitorNames[m.selection],
				m.monitors,
			)
			return nextPage, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func padRight(input string, lineLength int, fill string) string {
	for len(input) < lineLength {
		input += fill
	}
	return input
}

func truncate(input string, lineLength int) string {
	if len(input) > lineLength {
		return input[0:lineLength-2] + "..."
	}
	return input
}

func (m monitorSelectPage) View() string {
	var names string
	var resolutions string
	for i, name := range m.monitorNames {
		var line string
		if i == m.selection {
			line += ">[" + name + "]"
		} else {
			line += "  " + name + " "
		}
		line = padRight(line, 12, " ")
		line = truncate(line, 12)
		names += line + "\n"
		resolutions += m.monitors[name].currentRes + "\n"
	}
	names = lipgloss.NewStyle().
		Border(
			lipgloss.NormalBorder(),
			false,
			true,
			false,
			false,
		).MarginRight(1).
		Render(names)
	return lipgloss.Place(
		windowWidth,
		windowHeight,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				names,
				resolutions,
			),
			lipgloss.PlaceVertical(
				windowHeight-(len(m.monitorNames)+5),
				lipgloss.Bottom,
				lipgloss.NewStyle().Render(getHelp()),
			),
		),
	)
}
