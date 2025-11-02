package datepicker

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	GoToStart,
	GoToEnd,
	Down,
	Up,
	Select,
	Help,
	Quit,
	Validate key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.GoToStart, k.GoToEnd},
		{k.Help, k.Quit, k.Validate, k.Select},
	}
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		GoToStart: key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "start")),
		GoToEnd:   key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "end")),
		Down:      key.NewBinding(key.WithKeys("j", "down", "ctrl+n"), key.WithHelp("j", "down")),
		Up:        key.NewBinding(key.WithKeys("k", "up", "ctrl+p"), key.WithHelp("k", "up")),
		Select:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "select")),
		Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "ctrl+c")),
		Validate:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "validate")),
	}
}

type Styles struct {
	Cursor,
	Date,
	Selected,
	Deltatime lipgloss.Style
}

func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {

	return Styles{
		Cursor:    r.NewStyle().Foreground(lipgloss.Color("212")),
		Date:      r.NewStyle(),
		Selected:  r.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),
		Deltatime: r.NewStyle().Foreground(lipgloss.Color("240")),
	}
}

// Model represents the state of our date picker
type model struct {
	dates    []time.Time         // List of dates to display
	cursor   int                 // Current cursor position
	selected map[string]struct{} // Map to keep track of selected dates
	offset   int                 // Offset for infinite scrolling

	help       help.Model
	Keys       KeyMap
	Styles     Styles
	TimeLayout string
}

// GetSelected returns selected dates sorted.
// No selected if program is quit and not validated.
func (m model) GetSelected() ([]time.Time, error) {
	dates := make([]time.Time, 0, len(m.selected))
	for dateStr := range m.selected {
		date, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", dateStr, err)
		}

		dates = append(dates, date)
	}
	slices.SortFunc(dates, func(a, b time.Time) int {
		return a.Compare(b)
	})

	return dates, nil
}

// Initial model setup
func New() model {
	// Generate a list of dates starting from today
	now := time.Now()
	dates := make([]time.Time, 7) // Let's display the next 30 days
	for i := range dates {
		dates[i] = now.AddDate(0, 0, i)
	}

	return model{
		dates:    dates,
		cursor:   0,
		selected: make(map[string]struct{}),

		help:       help.New(),
		Keys:       DefaultKeyMap(),
		Styles:     DefaultStyles(),
		TimeLayout: time.DateOnly,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.GoToStart):
		case key.Matches(msg, m.Keys.GoToEnd):
		case key.Matches(msg, m.Keys.Up):
			if m.cursor > 0 {
				m.cursor--
			} else {
				// If we're at the top, scroll up by adjusting the offset
				m.offset--
				m.updateDates()
			}
		case key.Matches(msg, m.Keys.Down):
			if m.cursor < len(m.dates)-1 {
				m.cursor++
			} else {
				// If we're at the bottom, scroll down by adjusting the offset
				m.offset++
				m.updateDates()
			}
		case key.Matches(msg, m.Keys.Select):
			currDate := m.dates[m.cursor]
			_, ok := m.selected[currDate.Format(time.DateOnly)]
			if ok {
				delete(m.selected, currDate.Format(time.DateOnly))
			} else {
				m.selected[currDate.Format(time.DateOnly)] = struct{}{}
			}
		case key.Matches(msg, m.Keys.Validate):
			return m, tea.Quit
		case key.Matches(msg, m.Keys.Quit):
			clear(m.selected)
			return m, tea.Quit
		case key.Matches(msg, m.Keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

// Update the dates based on the current offset
func (m *model) updateDates() {
	now := time.Now()
	for i := range m.dates {
		m.dates[i] = now.AddDate(0, 0, i+m.offset)
	}
}

func (m model) View() string {

	var s strings.Builder
	for i, date := range m.dates {
		cursor := " " // Default: no cursor
		if m.cursor == i {
			cursor = "▶︎" // Cursor is on this line
		}
		dateRender := m.Styles.Date.Render
		if _, ok := m.selected[date.Format(time.DateOnly)]; ok {
			dateRender = m.Styles.Selected.Render
		}
		since := time.Until(date).Round(time.Hour * 24)

		fmt.Fprintf(&s, "%s %s %s\n",
			m.Styles.Deltatime.Render(fmt.Sprintf("%2.0f", since.Hours()/24)),
			m.Styles.Cursor.Render(cursor),
			dateRender(date.Format(m.TimeLayout)),
		)
	}
	s.WriteRune('\n')
	s.WriteString(m.help.View(m.Keys))
	return s.String()
}
