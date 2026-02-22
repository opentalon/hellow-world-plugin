package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/opentalon/opentalon/pkg/plugin"
)

const fixedFragment = "In which year was the first Hello, World! printed in Go?"

func TestCapabilities(t *testing.T) {
	h := &helloWorldHandler{fixedFragment: fixedFragment}
	caps := h.Capabilities()
	if caps.Name != "hello-world" {
		t.Errorf("Name = %q, want hello-world", caps.Name)
	}
	if len(caps.Actions) != 1 {
		t.Fatalf("len(Actions) = %d, want 1", len(caps.Actions))
	}
	if caps.Actions[0].Name != "prepare" {
		t.Errorf("Actions[0].Name = %q, want prepare", caps.Actions[0].Name)
	}
	if len(caps.Actions[0].Parameters) != 1 || caps.Actions[0].Parameters[0].Name != "text" {
		t.Errorf("Parameters = %v", caps.Actions[0].Parameters)
	}
}

func TestExecute_UnknownAction(t *testing.T) {
	h := &helloWorldHandler{}
	resp := h.Execute(plugin.Request{ID: "id1", Action: "unknown", Args: map[string]string{"text": "hello"}})
	if resp.Error == "" {
		t.Error("expected Error set for unknown action")
	}
	if !strings.Contains(resp.Error, "unknown") {
		t.Errorf("Error = %q", resp.Error)
	}
}

func TestExecute_NoHello_ReturnsGuard(t *testing.T) {
	h := &helloWorldHandler{}
	resp := h.Execute(plugin.Request{ID: "id1", Action: "prepare", Args: map[string]string{"text": "foo bar"}})
	if resp.Error != "" {
		t.Fatalf("unexpected Error: %s", resp.Error)
	}
	var guard struct {
		SendToLLM bool   `json:"send_to_llm"`
		Message   string `json:"message"`
	}
	if err := json.Unmarshal([]byte(resp.Content), &guard); err != nil {
		t.Fatalf("Content is not JSON: %v", err)
	}
	if guard.SendToLLM {
		t.Error("send_to_llm should be false")
	}
	if guard.Message == "" {
		t.Error("message should be set")
	}
}

func TestExecute_EmptyText_ReturnsGuard(t *testing.T) {
	h := &helloWorldHandler{}
	resp := h.Execute(plugin.Request{ID: "id1", Action: "prepare", Args: map[string]string{"text": ""}})
	if resp.Error != "" {
		t.Fatalf("unexpected Error: %s", resp.Error)
	}
	var guard struct {
		SendToLLM bool `json:"send_to_llm"`
	}
	if err := json.Unmarshal([]byte(resp.Content), &guard); err != nil {
		t.Fatalf("Content is not JSON: %v", err)
	}
	if guard.SendToLLM {
		t.Error("send_to_llm should be false for empty text")
	}
}

func TestExecute_Hello_AppendsWorldAndFragment(t *testing.T) {
	h := &helloWorldHandler{fixedFragment: fixedFragment}
	resp := h.Execute(plugin.Request{ID: "id1", Action: "prepare", Args: map[string]string{"text": "hello"}})
	if resp.Error != "" {
		t.Fatalf("unexpected Error: %s", resp.Error)
	}
	if resp.CallID != "id1" {
		t.Errorf("CallID = %q, want id1", resp.CallID)
	}
	parts := strings.SplitN(resp.Content, "\n\n", 2)
	if len(parts) != 2 {
		t.Fatalf("expected \"transformed\\n\\nfragment\", got %q", resp.Content)
	}
	if parts[0] != "hello world" {
		t.Errorf("transformed = %q, want hello world", parts[0])
	}
	if parts[1] != fixedFragment {
		t.Errorf("fragment = %q, want %q", parts[1], fixedFragment)
	}
}

func TestExecute_HelloWorld_NoDoubleWorld(t *testing.T) {
	h := &helloWorldHandler{fixedFragment: fixedFragment}
	resp := h.Execute(plugin.Request{ID: "id2", Action: "prepare", Args: map[string]string{"text": "hello world"}})
	if resp.Error != "" {
		t.Fatalf("unexpected Error: %s", resp.Error)
	}
	parts := strings.SplitN(resp.Content, "\n\n", 2)
	if len(parts) != 2 {
		t.Fatalf("expected \"transformed\\n\\nfragment\", got %q", resp.Content)
	}
	if parts[0] != "hello world" {
		t.Errorf("transformed = %q, want hello world (no double world)", parts[0])
	}
}

func TestExecute_HelloMixedCase_KeepsCaseAndAppendsWorld(t *testing.T) {
	h := &helloWorldHandler{fixedFragment: fixedFragment}
	resp := h.Execute(plugin.Request{ID: "id3", Action: "prepare", Args: map[string]string{"text": "  HELLO  "}})
	if resp.Error != "" {
		t.Fatalf("unexpected Error: %s", resp.Error)
	}
	parts := strings.SplitN(resp.Content, "\n\n", 2)
	if len(parts) != 2 {
		t.Fatalf("expected \"transformed\\n\\nfragment\", got %q", resp.Content)
	}
	// TrimSpace applied, then " world" appended; original casing kept.
	if parts[0] != "HELLO world" {
		t.Errorf("transformed = %q, want HELLO world", parts[0])
	}
}
