// Package ui implements the interactive TUI wizard for Tube Trend Buddy.
// Built on Bubble Tea + Lip Gloss + Bubbles. The TUI shells out to the
// current binary (os.Args[0]) with the right args, so it never duplicates
// the cobra command logic.
package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// Run launches the TUI. Blocks until the user quits.
func Run() error {
	m := newModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

// -- feature catalog --

type argSpec struct {
	key         string
	label       string
	required    bool
	defaultVal  string
	placeholder string
}

type feature struct {
	key  string
	name string
	desc string
	args []argSpec
	flag bool // true if --no-youtube / --deep / etc. are useful; placeholder
}

var features = []feature{
	{
		key:  "title",
		name: "title",
		desc: "CTR-optimized YouTube titles",
		args: []argSpec{
			{key: "niche", label: "Niche", required: true, placeholder: "AI tools, resep Nusantara..."},
			{key: "audience", label: "Audience", defaultVal: "general"},
			{key: "count", label: "Count", defaultVal: "10"},
		},
	},
	{
		key:  "tags",
		name: "tags",
		desc: "SEO tags (broad + mid + long tail)",
		args: []argSpec{
			{key: "title", label: "Video title", required: true, placeholder: "How to Learn AI in 30 Days"},
			{key: "description", label: "Description (opt.)"},
			{key: "count", label: "Count", defaultVal: "30"},
		},
	},
	{
		key:  "trend",
		name: "trend",
		desc: "Rising niches + content angles (live YouTube data)",
		args: []argSpec{
			{key: "region", label: "Region", defaultVal: "US", placeholder: "ID, US, JP..."},
			{key: "category", label: "Category", defaultVal: "general"},
			{key: "period", label: "Period", defaultVal: "7d", placeholder: "1d | 7d | 30d | 90d"},
			{key: "count", label: "Count", defaultVal: "10"},
		},
	},
	{
		key:  "niche",
		name: "niche",
		desc: "Audit a channel + find adjacent niches",
		args: []argSpec{
			{key: "channel", label: "Channel handle", required: true, placeholder: "@mkbhd"},
			{key: "deep", label: "Deep analysis", defaultVal: "false"},
		},
	},
	{
		key:  "describe",
		name: "describe",
		desc: "Description + hashtags + CTA",
		args: []argSpec{
			{key: "title", label: "Video title", required: true},
			{key: "description", label: "Summary (opt.)"},
			{key: "timestamps", label: "Timestamps (opt.)", placeholder: "00:00 intro, 02:30 basics"},
			{key: "cta", label: "Include CTA", defaultVal: "true"},
		},
	},
	{
		key:  "thumbnail",
		name: "thumbnail",
		desc: "Thumbnail concept (overlay + image-gen prompt)",
		args: []argSpec{
			{key: "title", label: "Video title", required: true},
			{key: "mood", label: "Mood", defaultVal: "curious", placeholder: "shocked | curious | excited | serious | funny"},
			{key: "face", label: "Faces (0-3)", defaultVal: "1"},
			{key: "colors", label: "Color palette", placeholder: "yellow,black,red"},
		},
	},
	{
		key:  "monetize",
		name: "monetize",
		desc: "Monetization plan (ads, sponsorships, products)",
		args: []argSpec{
			{key: "niche", label: "Niche", required: true},
			{key: "subs", label: "Subscribers", defaultVal: "1000"},
			{key: "region", label: "Region", defaultVal: "US"},
			{key: "focus", label: "Focus", defaultVal: "all", placeholder: "ads | sponsorships | products | all"},
		},
	},
	{
		key:  "calendar",
		name: "calendar",
		desc: "Full-month content calendar",
		args: []argSpec{
			{key: "month", label: "Month (e.g. 7 or 2026-07)", defaultVal: "next month"},
			{key: "niche", label: "Niche", required: true},
			{key: "frequency", label: "Frequency", defaultVal: "2/week", placeholder: "1/week | 2/week | 3/week | daily"},
		},
	},
}

// -- model --

type viewState int

const (
	viewMenu viewState = iota
	viewForm
	viewLoading
	viewResult
)

type model struct {
	state   viewState
	width   int
	height  int
	err     error

	// menu
	list list.Model

	// form
	formForm  []textinput.Model
	formLabel []string
	formFeat  *feature
	formFocus int

	// loading
	spinner  spinner.Model
	loading  string

	// result
	viewport viewport.Model
	result   string
}

func newModel() *model {
	items := make([]list.Item, len(features))
	for i, f := range features {
		items[i] = featureItem{feat: f, idx: i}
	}
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.SetHeight(2)
	delegate.SetSpacing(0)
	l := list.New(items, delegate, 0, 0)
	l.Title = "Tube Trend Buddy - 8 features"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return &model{
		state:   viewMenu,
		list:    l,
		spinner: s,
	}
}

type featureItem struct {
	feat feature
	idx  int
}

func (f featureItem) Title() string       { return fmt.Sprintf("  %d. %s", f.idx+1, f.feat.name) }
func (f featureItem) Description() string { return f.feat.desc }
func (f featureItem) FilterValue() string { return f.feat.name }

func (m *model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-2)
		if m.state == viewResult {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 2
		}
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case viewMenu:
			return m.updateMenu(msg)
		case viewForm:
			return m.updateForm(msg)
		case viewLoading:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil
		case viewResult:
			return m.updateResult(msg)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case runResultMsg:
		m.state = viewResult
		m.result = string(msg)
		m.viewport = viewport.New(m.width, m.height-2)
		m.viewport.SetContent(m.result)
		return m, nil

	case runErrorMsg:
		m.state = viewResult
		m.result = "ERROR:\n\n" + string(msg)
		m.viewport = viewport.New(m.width, m.height-2)
		m.viewport.SetContent(m.result)
		return m, nil
	}

	return m, nil
}

func (m *model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		sel, ok := m.list.SelectedItem().(featureItem)
		if !ok {
			return m, nil
		}
		m.formFeat = &sel.feat
		m.buildForm()
		m.state = viewForm
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) buildForm() {
	m.formForm = make([]textinput.Model, len(m.formFeat.args))
	m.formLabel = make([]string, len(m.formFeat.args))
	for i, a := range m.formFeat.args {
		ti := textinput.New()
		ti.CharLimit = 200
		ti.Placeholder = a.placeholder
		if a.defaultVal != "" {
			ti.SetValue(a.defaultVal)
		}
		ti.Width = 40
		if i == 0 {
			ti.Focus()
		}
		m.formForm[i] = ti
		m.formLabel[i] = a.label
		if a.required {
			m.formLabel[i] += " *"
		}
	}
	m.formFocus = 0
}

func (m *model) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = viewMenu
		return m, nil
	case "tab", "down":
		m.formFocus = (m.formFocus + 1) % len(m.formForm)
		for i, ti := range m.formForm {
			if i == m.formFocus {
				ti.Focus()
			} else {
				ti.Blur()
			}
			m.formForm[i] = ti
		}
		return m, nil
	case "shift+tab", "up":
		m.formFocus = (m.formFocus - 1 + len(m.formForm)) % len(m.formForm)
		for i, ti := range m.formForm {
			if i == m.formFocus {
				ti.Focus()
			} else {
				ti.Blur()
			}
			m.formForm[i] = ti
		}
		return m, nil
	case "enter":
		// validate required
		for i, a := range m.formFeat.args {
			if a.required && strings.TrimSpace(m.formForm[i].Value()) == "" {
				m.err = fmt.Errorf("field %q is required", a.label)
				return m, nil
			}
		}
		m.state = viewLoading
		m.loading = "Calling LLM..."
		featName := m.formFeat.name
		args := map[string]string{}
		for i, a := range m.formFeat.args {
			args[a.key] = m.formForm[i].Value()
		}
		return m, runFeatureCmd(featName, args)
	}
	var cmd tea.Cmd
	m.formForm[m.formFocus], cmd = m.formForm[m.formFocus].Update(msg)
	return m, cmd
}

func (m *model) updateResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "r", "esc", "enter":
		m.state = viewMenu
		return m, nil
	}
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	switch m.state {
	case viewMenu:
		return m.list.View() + "\n" + menuFooter(m.width)
	case viewForm:
		return m.viewForm()
	case viewLoading:
		return fmt.Sprintf("\n  %s %s\n\n  (ctrl+c to cancel)\n", m.spinner.View(), m.loading)
	case viewResult:
		return m.viewport.View() + "\n" + resultFooter(m.width)
	}
	return ""
}

func (m *model) viewForm() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("  %s", m.formFeat.name)))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  " + m.formFeat.desc))
	b.WriteString("\n\n")
	for i, lbl := range m.formLabel {
		if i == m.formFocus {
			b.WriteString(focusedLabelStyle.Render("  > " + lbl + ":"))
		} else {
			b.WriteString(labelStyle.Render("    " + lbl + ":"))
		}
		b.WriteString("\n")
		b.WriteString("      " + m.formForm[i].View() + "\n\n")
	}
	if m.err != nil {
		b.WriteString(errStyle.Render("  " + m.err.Error()))
		b.WriteString("\n")
	}
	b.WriteString(dimStyle.Render("  tab/arrows to move - enter to run - esc to back"))
	b.WriteString("\n")
	return b.String()
}

func menuFooter(w int) string {
	return dimStyle.Render(strings.Repeat(" ", 0) + "  arrows/enter to pick - q to quit")
}

func resultFooter(w int) string {
	return dimStyle.Render("  r/esc to back to menu - q to quit")
}

// -- run feature (shells out to self) --

type runResultMsg string
type runErrorMsg string

func runFeatureCmd(feat string, args map[string]string) tea.Cmd {
	return func() tea.Msg {
		out, err := runFeature(feat, args)
		if err != nil {
			return runErrorMsg(err.Error())
		}
		return runResultMsg(out)
	}
}

func runFeature(feat string, args map[string]string) (string, error) {
	argv := []string{feat}
	// global flags first (provider, model, lang, count, json)
	for _, a := range []string{"provider", "model", "lang", "count", "json"} {
		if v, ok := args[a]; ok && v != "" {
			argv = append(argv, "--"+a, v)
		}
	}
	// feature-specific args
	for k, v := range args {
		if k == "provider" || k == "model" || k == "lang" || k == "count" || k == "json" {
			continue
		}
		if v == "" {
			continue
		}
		argv = append(argv, "--"+k, v)
	}
	cmd := exec.Command(os.Args[0], argv...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		// include stderr tail in error
		return "", fmt.Errorf("%w: %s", err, trim(string(out), 500))
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func trim(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// -- styles --

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			MarginTop(1).
			MarginLeft(2)

	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	focusedLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205"))

	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

	errStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
)
