package display

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ewan-valentine/reqstat/internal/analyzer"
	"github.com/ewan-valentine/reqstat/internal/client"
)

var (
	// Color palette - Catppuccin Mocha inspired
	rosewater = lipgloss.Color("#f5e0dc")
	flamingo  = lipgloss.Color("#f2cdcd")
	pink      = lipgloss.Color("#f5c2e7")
	mauve     = lipgloss.Color("#cba6f7")
	red       = lipgloss.Color("#f38ba8")
	maroon    = lipgloss.Color("#eba0ac")
	peach     = lipgloss.Color("#fab387")
	yellow    = lipgloss.Color("#f9e2af")
	green     = lipgloss.Color("#a6e3a1")
	teal      = lipgloss.Color("#94e2d5")
	sky       = lipgloss.Color("#89dceb")
	sapphire  = lipgloss.Color("#74c7ec")
	blue      = lipgloss.Color("#89b4fa")
	lavender  = lipgloss.Color("#b4befe")
	text      = lipgloss.Color("#cdd6f4")
	subtext0  = lipgloss.Color("#a6adc8")
	surface0  = lipgloss.Color("#313244")
	base      = lipgloss.Color("#1e1e2e")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(mauve).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(blue).
			PaddingLeft(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(surface0)

	labelStyle = lipgloss.NewStyle().
			Foreground(subtext0).
			Width(18)

	valueStyle = lipgloss.NewStyle().
			Foreground(text)

	successStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	urlStyle = lipgloss.NewStyle().
			Foreground(sapphire).
			Italic(true)

	headerKeyStyle = lipgloss.NewStyle().
			Foreground(teal)

	headerValStyle = lipgloss.NewStyle().
			Foreground(text)

	jsonKeyStyle = lipgloss.NewStyle().
			Foreground(mauve)

	jsonTypeStyle = lipgloss.NewStyle().
			Foreground(peach)

	jsonExampleStyle = lipgloss.NewStyle().
				Foreground(subtext0).
				Italic(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface0).
			Padding(1, 2)

	dividerStyle = lipgloss.NewStyle().
			Foreground(surface0)
)

type Options struct {
	ShowBody    bool
	PrettyJSON  bool
	MaxBodySize int
}

func Render(result *client.Result, opts Options) {
	fmt.Println()

	// Title with URL
	fmt.Println(titleStyle.Render("⚡ reqstat"))
	fmt.Println(urlStyle.Render("   GET " + result.URL))
	fmt.Println()

	// Status section
	renderStatus(result)
	fmt.Println()

	// Timing section
	renderTiming(result)
	fmt.Println()

	// Size section
	renderSize(result)
	fmt.Println()

	// Headers section
	renderHeaders(result)
	fmt.Println()

	// JSON Shape section
	if result.IsJSON() {
		renderJSONShape(result)
		fmt.Println()
	}

	// Body (optional)
	if opts.ShowBody {
		renderBody(result, opts)
		fmt.Println()
	}
}

func renderStatus(result *client.Result) {
	fmt.Println(sectionStyle.Render("STATUS"))
	fmt.Println()

	statusText := fmt.Sprintf("%d %s", result.StatusCode, statusText(result.StatusCode))
	var styled string
	switch {
	case result.IsSuccess():
		styled = successStyle.Render("● " + statusText)
	case result.IsRedirect():
		styled = warningStyle.Render("● " + statusText)
	case result.IsClientError():
		styled = warningStyle.Render("● " + statusText)
	default:
		styled = errorStyle.Render("● " + statusText)
	}

	fmt.Printf("   %s\n", styled)
}

func renderTiming(result *client.Result) {
	fmt.Println(sectionStyle.Render("TIMING"))
	fmt.Println()

	timings := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"Total", formatDuration(result.Duration), green},
		{"DNS Lookup", formatDuration(result.DNSLookup), teal},
		{"TCP Connect", formatDuration(result.TCPConnection), sky},
		{"TLS Handshake", formatDuration(result.TLSHandshake), sapphire},
		{"Server Response", formatDuration(result.ServerResponse), blue},
	}

	for _, t := range timings {
		label := labelStyle.Render(t.label)
		value := lipgloss.NewStyle().Foreground(t.color).Render(t.value)
		fmt.Printf("   %s %s\n", label, value)
	}

	// Visual timing bar
	fmt.Println()
	renderTimingBar(result)
}

func renderTimingBar(result *client.Result) {
	total := result.Duration.Milliseconds()
	if total == 0 {
		return
	}

	barWidth := 50
	dns := int(float64(result.DNSLookup.Milliseconds()) / float64(total) * float64(barWidth))
	tcp := int(float64(result.TCPConnection.Milliseconds()) / float64(total) * float64(barWidth))
	tls := int(float64(result.TLSHandshake.Milliseconds()) / float64(total) * float64(barWidth))
	server := int(float64(result.ServerResponse.Milliseconds()) / float64(total) * float64(barWidth))
	rest := barWidth - dns - tcp - tls - server

	if rest < 0 {
		rest = 0
	}

	bar := ""
	bar += lipgloss.NewStyle().Foreground(teal).Render(strings.Repeat("█", dns))
	bar += lipgloss.NewStyle().Foreground(sky).Render(strings.Repeat("█", tcp))
	bar += lipgloss.NewStyle().Foreground(sapphire).Render(strings.Repeat("█", tls))
	bar += lipgloss.NewStyle().Foreground(blue).Render(strings.Repeat("█", server))
	bar += lipgloss.NewStyle().Foreground(lavender).Render(strings.Repeat("█", rest))

	fmt.Printf("   %s\n", bar)
	fmt.Printf("   %s\n",
		lipgloss.NewStyle().Foreground(subtext0).Render(
			"DNS │ TCP │ TLS │ Server │ Transfer"))
}

func renderSize(result *client.Result) {
	fmt.Println(sectionStyle.Render("SIZE"))
	fmt.Println()

	label := labelStyle.Render("Content Length")
	value := valueStyle.Render(formatBytes(result.ContentLength))
	fmt.Printf("   %s %s\n", label, value)

	label = labelStyle.Render("Content Type")
	value = valueStyle.Render(result.ContentType)
	fmt.Printf("   %s %s\n", label, value)
}

func renderHeaders(result *client.Result) {
	fmt.Println(sectionStyle.Render("HEADERS"))
	fmt.Println()

	keys := make([]string, 0, len(result.Headers))
	for k := range result.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		values := result.Headers[k]
		key := headerKeyStyle.Render(k)
		val := headerValStyle.Render(strings.Join(values, ", "))
		fmt.Printf("   %s: %s\n", key, val)
	}
}

func renderJSONShape(result *client.Result) {
	fmt.Println(sectionStyle.Render("JSON SHAPE"))
	fmt.Println()

	schema, err := analyzer.AnalyzeJSON(result.Body)
	if err != nil {
		fmt.Printf("   %s\n", errorStyle.Render("Failed to parse JSON: "+err.Error()))
		return
	}

	// Summary
	fmt.Printf("   %s\n", lipgloss.NewStyle().Foreground(subtext0).Render(schema.Summary()))
	fmt.Println()

	// Render the shape
	lines := strings.Split(schema.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		styled := styleLine(line)
		fmt.Printf("   %s\n", styled)
	}
}

func styleLine(line string) string {
	// Style JSON shape output
	if strings.Contains(line, ": {") || strings.HasPrefix(strings.TrimSpace(line), "{") {
		return lipgloss.NewStyle().Foreground(text).Render(line)
	}
	if strings.HasPrefix(strings.TrimSpace(line), "}") {
		return lipgloss.NewStyle().Foreground(text).Render(line)
	}

	// Handle field lines like "  name: string // e.g. "value""
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		key := parts[0]
		rest := parts[1]

		// Check for example comment
		if idx := strings.Index(rest, "//"); idx != -1 {
			typeStr := rest[:idx]
			example := rest[idx:]
			return jsonKeyStyle.Render(key) + ":" +
				jsonTypeStyle.Render(typeStr) +
				jsonExampleStyle.Render(example)
		}

		return jsonKeyStyle.Render(key) + ":" + jsonTypeStyle.Render(rest)
	}

	return line
}

func renderBody(result *client.Result, opts Options) {
	fmt.Println(sectionStyle.Render("BODY"))
	fmt.Println()

	body := string(result.Body)

	if opts.PrettyJSON && result.IsJSON() {
		var parsed any
		if err := json.Unmarshal(result.Body, &parsed); err == nil {
			pretty, _ := json.MarshalIndent(parsed, "", "  ")
			body = string(pretty)
		}
	}

	if len(body) > opts.MaxBodySize {
		body = body[:opts.MaxBodySize]
		body += lipgloss.NewStyle().Foreground(subtext0).Render("\n   ... (truncated)")
	}

	// Indent body
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		fmt.Printf("   %s\n", line)
	}
}

func formatDuration(d interface{ Milliseconds() int64 }) string {
	ms := d.Milliseconds()
	if ms == 0 {
		return "-"
	}
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.2fs", float64(ms)/1000)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func statusText(code int) string {
	texts := map[int]string{
		200: "OK",
		201: "Created",
		204: "No Content",
		301: "Moved Permanently",
		302: "Found",
		304: "Not Modified",
		400: "Bad Request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		429: "Too Many Requests",
		500: "Internal Server Error",
		502: "Bad Gateway",
		503: "Service Unavailable",
		504: "Gateway Timeout",
	}
	if t, ok := texts[code]; ok {
		return t
	}
	return ""
}
