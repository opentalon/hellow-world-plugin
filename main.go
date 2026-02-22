package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"

	"github.com/opentalon/opentalon/pkg/plugin"
)

// Random prompt fragments (questions for the LLM). One is picked per prepare() call. Returned as the user message to the LLM.
var promptFragments = []string{
	"In which language was the first Hello, World! program printed?",
	"In which year was the first Hello, World! program ever printed?",
	"In which year was the first Hello, World! printed in C?",
	"In which year was the first Hello, World! printed in Java?",
	"In which year was the first Hello, World! printed in Python?",
	"In which year was the first Hello, World! printed in Ruby?",
	"In which year was the first Hello, World! printed in Go?",
	"In which year was the first Hello, World! printed in JavaScript?",
	"In which year was the first Hello, World! printed in Rust?",
	"In which year was the first Hello, World! printed in PHP?",
	"What is the most printed programming language for Hello, World!?",
}

type helloWorldHandler struct {
	// If set (e.g. via HELLO_WORLD_PROMPT_FRAGMENT), used instead of random pick.
	fixedFragment string
}

func (h *helloWorldHandler) Capabilities() plugin.CapabilitiesMsg {
	return plugin.CapabilitiesMsg{
		Name:        "hello-world",
		Description: "When the user types 'hello', adds 'world' (so 'hello' becomes 'hello world') and attaches a random prompt fragment about Hello, World! history (first language, which year for C/Java/Ruby/Go etc.).",
		Actions: []plugin.ActionMsg{
			{
				Name:        "prepare",
				Description: "If user text contains 'hello', add ' world' to it and return a random prompt fragment (e.g. which year Java/Ruby/Go first printed Hello, World!).",
				Parameters: []plugin.ParameterMsg{
					{Name: "text", Description: "User message or text to transform", Type: "string", Required: true},
				},
			},
		},
	}
}

func (h *helloWorldHandler) pickFragment() string {
	if h.fixedFragment != "" {
		return h.fixedFragment
	}
	return promptFragments[rand.Intn(len(promptFragments))]
}

func (h *helloWorldHandler) Execute(req plugin.Request) plugin.Response {
	if req.Action != "prepare" {
		return plugin.Response{Error: "unknown action: " + req.Action}
	}

	text := req.Args["text"]
	// Guard example: only "hello" is allowed to reach the LLM; everything else gets this message.
	if !strings.Contains(strings.ToLower(strings.TrimSpace(text)), "hello") {
		guard := map[string]interface{}{
			"send_to_llm": false,
			"message":     "Plugin only accepts send hello to LLM. All another knows human brain.",
		}
		body, _ := json.Marshal(guard)
		return plugin.Response{
			CallID:  req.ID,
			Content: string(body),
		}
	}
	// When user types "hello", add " world" only if not already there (avoid "hello world world").
	trimmed := strings.TrimSpace(text)
	lower := strings.ToLower(trimmed)
	if !strings.HasSuffix(lower, "world") {
		trimmed += " world"
	}
	transformed := trimmed
	question := h.pickFragment()
	content := transformed + "\n\n" + question
	return plugin.Response{
		CallID:  req.ID,
		Content: content,
	}
}

func main() {
	// Optional: force one fragment for all calls (e.g. in tests or for a specific campaign).
	fixed := os.Getenv("HELLO_WORLD_PROMPT_FRAGMENT")

	handler := &helloWorldHandler{fixedFragment: fixed}
	if err := plugin.Serve(handler); err != nil {
		os.Exit(1)
	}
}
