package main

// Package is called aw
import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"

	aw "github.com/deanishe/awgo"

	"com.grybkov.alfred-ynab/src/config"
	"com.grybkov.alfred-ynab/src/ynab"
)
import "golang.design/x/clipboard"

// Workflow is the main API
var wf *aw.Workflow
var cacheTTL time.Duration

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
	cacheTTL = config.GetCacheTTL()

}

type Arg struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Value      string `json:"value"`
	BrowserURL string `json:"browser_url"`
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

// Your workflow starts here
func run() {
	query := ""

	// Use wf.Args so magic actions are handled
	args := wf.Args()
	if len(args) > 0 {
		query = args[0]
	}

	// wf.NewItem(cfg.Get("YNAB_BUDGET"))

	// Disable UIDs so Alfred respects our sort order. Without this,
	// it may bump read/unpublished books to the top of results, but
	// we want to force them to always be below unread books.
	wf.Configure(aw.SuppressUIDs(true))

	if query == "" {
		apiClient := ynab.MakeYNABClient(config.GetYnabAccessToken())

		handlers := []ItemHandler{
			&TransactionsHandler{},
			&AccountsHandler{},
			&CategoriesHandler{},
		}
		var wg sync.WaitGroup
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Preload(&wg, apiClient, cacheTTL)
		}
		wg.Wait()

		for _, handler := range handlers {
			handler.AddItems(wf)
		}
	} else if query == "browse" {
		arg := Arg{}
		argJson := args[1]
		log.Println(argJson)
		_ = json.Unmarshal([]byte(argJson), &arg)
		openbrowser(arg.BrowserURL)
	} else if query == "clipboard" {
		err := clipboard.Init()
		if err != nil {
			panic(err)
		}
		arg := Arg{}
		argJson := args[1]
		log.Println(argJson)
		_ = json.Unmarshal([]byte(argJson), &arg)
		clipboard.Write(clipboard.FmtText, []byte(arg.Value))
	} else {
		wf.WarnEmpty("No matching items", "Try a different query?")
	}
	// Add a "Script Filter" result
	// Send results to Alfred
	wf.SendFeedback()
}

func main() {
	// Wrap your entry point with Run() to catch and log panics and
	// show an error in Alfred instead of silently dying
	wf.Run(run)
}
