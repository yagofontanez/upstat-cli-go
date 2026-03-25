package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

const (
	BASE_URL         = "https://upstat-backend.onrender.com/api/v1"
	CONFIG_FILE_NAME = ".upstat"
	REFRESH_INTERVAL = 30 * time.Second
)

type Lang string

const (
	PT Lang = "pt"
	EN Lang = "en"
)

var I18N = map[Lang]map[string]string{
	PT: {
		"title":            "▲ UpStat CLI",
		"connecting":       "Conectando...",
		"noMonitors":       "Nenhum monitor encontrado.",
		"monitors":         "monitores",
		"online":           "online",
		"offline":          "offline",
		"updating":         "Atualizando a cada",
		"seconds":          "s",
		"quit":             "Ctrl+C para sair",
		"invalidKey":       "API key inválida. Rode: upstat logout",
		"fetchError":       "Erro ao buscar monitores",
		"keyRemoved":       "API key removida.",
		"keySaved":         "Configurações salvas em ~/.upstat",
		"askKey":           "Cole sua API key (Settings → API Keys no painel):",
		"askKeyValidation": "A key deve começar com ups_",
		"goodbye":          "Até mais! 👋",
	},
	EN: {
		"title":            "▲ UpStat CLI",
		"connecting":       "Connecting...",
		"noMonitors":       "No monitors found.",
		"monitors":         "monitors",
		"online":           "online",
		"offline":          "offline",
		"updating":         "Updating every",
		"seconds":          "s",
		"quit":             "Ctrl+C to quit",
		"invalidKey":       "Invalid API key. Run: upstat logout",
		"fetchError":       "Failed to fetch monitors",
		"keyRemoved":       "API key removed.",
		"keySaved":         "Config saved at ~/.upstat",
		"askKey":           "Paste your API key (Settings → API Keys in the dashboard):",
		"askKeyValidation": "Key must start with ups_",
		"goodbye":          "Goodbye! 👋",
	},
}

func tr(lang Lang, key string) string {
	if val, ok := I18N[lang][key]; ok {
		return val
	}
	return key
}

type Config struct {
	ApiKey string `json:"apiKey"`
	Lang   Lang   `json:"lang"`
}

type Monitor struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	URL              string  `json:"url"`
	Status           string  `json:"status"`
	LatencyMs        *int    `json:"latency_ms"`
	UptimePercentage float64 `json:"uptime_percentage"`
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, CONFIG_FILE_NAME)
}

func saveConfig(cfg Config) error {
	data, _ := json.Marshal(cfg)
	return ioutil.WriteFile(getConfigPath(), data, 0644)
}

func loadConfig() (*Config, error) {
	data, err := ioutil.ReadFile(getConfigPath())
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.ApiKey == "" || cfg.Lang == "" {
		return nil, fmt.Errorf("invalid config")
	}
	return &cfg, nil
}

func clearConfig() error {
	return os.Remove(getConfigPath())
}

func fetchMonitors(apiKey string) ([]Monitor, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", BASE_URL+"/monitors", nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("unauthorized")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var monitors []Monitor
	if err := json.Unmarshal(body, &monitors); err != nil {
		return nil, err
	}
	return monitors, nil
}

func renderMonitors(monitors []Monitor, lang Lang) {
	now := time.Now().Format("15:04:05")
	fmt.Print("\033[H\033[2J")

	color.New(color.FgHiCyan).Add(color.Bold).Printf("%s — %s\n", tr(lang, "title"), now)
	fmt.Println(strings.Repeat("─", 60))

	if len(monitors) == 0 {
		color.New(color.FgHiBlack).Printf("\n  %s\n\n", tr(lang, "noMonitors"))
		return
	}

	up := 0
	for _, m := range monitors {
		isUp := m.Status == "up"
		dot := color.HiGreenString("●")
		name := color.HiWhiteString(m.Name)
		uptime := color.HiGreenString("%.1f%%", m.UptimePercentage)
		latency := "—"
		if m.LatencyMs != nil {
			latency = fmt.Sprintf("%dms", *m.LatencyMs)
		}
		if !isUp {
			dot = color.RedString("●")
			name = color.RedString(m.Name)
			uptime = color.RedString("%.1f%%", m.UptimePercentage)
		} else {
			up++
		}
		fmt.Printf("  %s  %-28s %-8s %s\n", dot, name, latency, uptime)
	}

	down := len(monitors) - up
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("  %d %s  ·  %d %s", len(monitors), tr(lang, "monitors"), up, tr(lang, "online"))
	if down > 0 {
		fmt.Printf("  ·  %d %s", down, tr(lang, "offline"))
	}
	fmt.Println()
	fmt.Printf("\n  %s %d%s — %s\n\n", tr(lang, "updating"), int(REFRESH_INTERVAL.Seconds()), tr(lang, "seconds"), tr(lang, "quit"))
}

func startWatch(cfg Config) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + tr(cfg.Lang, "connecting")
	s.Start()

	run := func() {
		monitors, err := fetchMonitors(cfg.ApiKey)
		s.Stop()
		if err != nil {
			if err.Error() == "unauthorized" {
				color.Red("\n  %s\n", tr(cfg.Lang, "invalidKey"))
				os.Exit(1)
			}
			color.Red("\n  %s: %s\n", tr(cfg.Lang, "fetchError"), err.Error())
			return
		}
		renderMonitors(monitors, cfg.Lang)
	}

	run()
	ticker := time.NewTicker(REFRESH_INTERVAL)
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ticker.C:
				run()
			case <-done:
				ticker.Stop()
				color.HiBlack("\n  %s\n\n", tr(cfg.Lang, "goodbye"))
				os.Exit(0)
			}
		}
	}()

	select {}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "logout" {
		cfg, _ := loadConfig()
		lang := PT
		if cfg != nil {
			lang = cfg.Lang
		}
		clearConfig()
		color.HiGreen("  %s\n", tr(lang, "keyRemoved"))
		os.Exit(0)
	}

	cfg, err := loadConfig()
	if err != nil {
		langChoice := ""
		promptLang := &survey.Select{
			Message: "Language / Idioma:",
			Options: []string{"Português", "English"},
		}
		survey.AskOne(promptLang, &langChoice)

		var lang Lang
		if langChoice == "Português" {
			lang = PT
		} else {
			lang = EN
		}

		apiKey := ""
		promptKey := &survey.Password{
			Message: tr(lang, "askKey"),
		}
		survey.AskOne(promptKey, &apiKey, survey.WithValidator(func(ans interface{}) error {
			if s, ok := ans.(string); ok && strings.HasPrefix(s, "ups_") {
				return nil
			}
			return fmt.Errorf(tr(lang, "askKeyValidation"))
		}))

		cfg = &Config{ApiKey: apiKey, Lang: lang}
		saveConfig(*cfg)
		color.HiBlack("  %s\n\n", tr(cfg.Lang, "keySaved"))
	}

	startWatch(*cfg)
}
