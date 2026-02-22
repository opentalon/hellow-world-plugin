# Hello World Plugin

[![CI](https://github.com/opentalon/hellow-world-plugin/actions/workflows/ci.yml/badge.svg)](https://github.com/opentalon/hellow-world-plugin/actions/workflows/ci.yml)

Standalone OpenTalon plugin intended as a **content preparer**: run before the first LLM call to transform or guard user input. When the user says "hello", it turns the input into "hello world" plus a random question; otherwise it returns a guard message and no LLM call.

Depends on [github.com/opentalon/opentalon](https://github.com/opentalon/opentalon) `pkg/plugin` SDK only.

## Spec

| Item | Value |
|------|--------|
| **Plugin ID** | `hello-world` |
| **Action** | `prepare` |

**Parameters**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `text` | string | yes | User message or text to transform |

**Return (content preparer use)**

- **If input does not contain "hello"** (case-insensitive, trimmed): return JSON so the orchestrator skips the LLM and sends a message to the user:
  ```json
  {"send_to_llm": false, "message": "Plugin only accepts send hello to LLM. All another knows human brain."}
  ```
- **If input contains "hello"**: return a **single plain string** (not JSON):  
  `"<transformed>" + "\n\n" + "<random prompt fragment>"`  
  - Transformed: input with `" world"` appended only if the trimmed text does not already end with `"world"` (so one "world", not double).
  - Random fragment: one of the fixed questions (e.g. "In which year was the first Hello, World! printed in Ruby?").

**Environment**

- `HELLO_WORLD_PROMPT_FRAGMENT`: if set, use this string instead of a random fragment (e.g. for tests).

## Build and test

From this repo (standalone; no dependency on the OpenTalon repo path):

```bash
make build   # → binary named "hello-world-plugin" (or set BINARY_NAME)
make test    # go test ./...
make lint    # golangci-lint run
```

## Config (in OpenTalon config.yaml)

Register the plugin from GitHub and use it as a content preparer:

```yaml
plugins:
  hello-world:
    enabled: true
    github: "opentalon/hello-world-plugin"
    ref: "master"
    config: {}

orchestrator:
  content_preparers:
    - plugin: hello-world
      action: prepare
      arg_key: text
```
