package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/theb0imanuu/salam/internal/models"
	"github.com/theb0imanuu/salam/internal/monitor"
)

type tickMsg time.Time

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Help   key.Binding
	Reload key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Reload: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reload"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Reload}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Reload, k.Help, k.Quit},
	}
}

type Model struct {
	Stats    *models.HealthData
	Monitor  *monitor.Monitor
	LastTick time.Time
	Err      error
	Help     help.Model
	KeyMap   keyMap
}

func NewModel(m *monitor.Monitor) Model {
	return Model{
		Monitor: m,
		Help:    help.New(),
		KeyMap:  keys,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.tick())
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.Reload):
			m.fetchStats()
			return m, nil
		}

	case tickMsg:
		m.fetchStats()
		m.LastTick = time.Time(msg)
		return m, m.tick()

	case tea.WindowSizeMsg:
		m.Help.Width = msg.Width
	}

	return m, nil
}

func (m *Model) fetchStats() {
	stats, err := m.Monitor.Collect()
	if err != nil {
		m.Err = err
		return
	}
	m.Stats = stats
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error: %v", m.Err)
	}

	if m.Stats == nil {
		return mainStyle.Render(titleStyle.Render("🕊️ SALAM - Real-time Infrastructure Monitoring") + "\n\nInitializing dashboard...")
	}

	header := titleStyle.Render("🕊️ SALAM - Real-time Infrastructure Monitoring")

	// CPU & Memory
	cpuBox := renderBox("CPU Usage", renderProgressBar(30, m.Stats.CPU.Usage)+
		fmt.Sprintf("\n\nCores: %d\nLoad: %.2f", m.Stats.CPU.Cores, m.Stats.CPU.LoadAverage))

	memBox := renderBox("Memory Usage", renderProgressBar(30, m.Stats.Memory.Usage)+
		fmt.Sprintf("\n\nUsed: %s\nFree: %s", m.Stats.Memory.Used, m.Stats.Memory.Free))

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, cpuBox, memBox)

	// Disk
	diskContent := ""
	for _, d := range m.Stats.Disk {
		diskContent += fmt.Sprintf("%-10s (%s)\n%s\n\n", d.Path, d.FSType, renderProgressBar(30, d.Usage))
	}
	diskBox := renderBox("Storage", strings.TrimSpace(diskContent))

	// Network
	netContent := ""
	for i, n := range m.Stats.Network {
		if i > 3 {
			break
		} // Limit to 4 interfaces
		netContent += fmt.Sprintf("%-8s ⬆️ %s  ⬇️ %s\n", n.Name, models.FormatBytes(n.BytesSent), models.FormatBytes(n.BytesRecv))
	}
	netBox := renderBox("Network", strings.TrimSpace(netContent))

	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, diskBox, netBox)

	info := metricLabelStyle.Render(fmt.Sprintf("\nLast update: %s", m.LastTick.Format("15:04:05")))
	helpView := "\n\n" + m.Help.View(m.KeyMap)

	return mainStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			topRow,
			bottomRow,
			info,
			helpView,
		),
	)
}

func StartDashboard(m *monitor.Monitor) error {
	p := tea.NewProgram(NewModel(m))
	_, err := p.Run()
	return err
}
