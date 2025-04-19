[日本語版 (Japanese version)](./README.md)

# prompt-bom

<!-- Logo (Replace with actual logo in docs/img/logo.png later) -->
![prompt-bom logo](https://via.placeholder.com/300x100.png?text=prompt-bom)

<!-- Badges (Update with actual CI status and version later) -->
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
<!-- [![Build Status](https://img.shields.io/github/actions/workflow/status/<your-org>/prompt-bom/go.yml?branch=main)](https://github.com/<your-org>/prompt-bom/actions) -->

A Bill of Materials (BOM) CLI tool for managing prompts.

This tool provides a prompt management foundation based on BOM/SBOM concepts, addressing the increasing complexity of generative AI prompts. It implements minimal management features using YAML-based BOM definitions and a Go-based CLI tool, supporting commands like `init`/`validate`/`build`, and considering future scalability, governance, and audit requirements.

## Features

- **Manage Prompts as Components**: Define and manage prompts as individual components.
- **Declarative Definition with YAML**: Describe the BOM structure declaratively in `prompt.bom.yaml`.
- **Simple CLI**: Easy operation with basic commands: `init`, `validate`, `build`.
- **Extensibility**: Designed with future integration with external tools (Dify, LangChain, etc.) in mind.

## Table of Contents

- [Installation](#installation)
- [Quick Start (60-second demo)](#quick-start-60-second-demo)
- [Command List](#command-list)
- [Architecture Overview](#architecture-overview)
- [License](#license)
- [Roadmap](#roadmap)

## Installation

### Go

```bash
go install github.com/kannazuki/prompt-bom/cmd/bom@latest
```
*(Note: The module path may change in the future)*

### Homebrew (Planned)

Installation via Homebrew Tap is planned for the future.

```bash
# brew tap <your-org>/tap
# brew install prompt-bom
```

## Quick Start (60-second demo)

1.  **Generate BOM Template:**

    ```bash
    bom init
    # -> prompt.bom.yaml を生成しました。 (Generated prompt.bom.yaml.)
    ```

2.  **Create Sample Components:**

    ```bash
    mkdir -p examples/components
    echo "This is part A." > examples/components/partA.md
    echo "This is part B." > examples/components/partB.md
    ```

3.  **Edit `prompt.bom.yaml`:**
    Add the following to the `components:` section:

    ```yaml
    components:
      - id: "partA"
        version: "1.0.0"
        hash: "sha256:dummy_hash_a"
        description: "Part A prompt"
        metadata:
          owner: "your-team"
      - id: "partB"
        version: "1.0.0"
        hash: "sha256:dummy_hash_b"
        description: "Part B prompt"
        metadata:
          owner: "your-team"
    ```
    *(Note: `hash` is a dummy value. Auto-generation/validation will be added later)*

4.  **Validate BOM:**

    ```bash
    bom validate prompt.bom.yaml
    # -> OK: スキーマと必須フィールド検証に合格 (OK: Schema and required fields validation passed)
    ```

5.  **Build Prompt:**

    ```bash
    bom build prompt.bom.yaml
    # The following will be printed to standard output:
    # This is part A.
    #
    # This is part B.
    #
    ```

    To output to a file:

    ```bash
    bom build prompt.bom.yaml -o final_prompt.md
    # -> final_prompt.md に結合結果を出力しました。 (Output combined result to final_prompt.md.)
    ```

## Command List

| Command          | Description                                      |
| ---------------- | ------------------------------------------------ |
| `bom init`       | Generate a BOM template YAML (`prompt.bom.yaml`) |
| `bom validate`   | Validate BOM YAML schema and required fields   |
| `bom build`      | Build and output prompts based on the BOM        |

Refer to `docs/usage.md` (To be created) for details.

## Architecture Overview

```plaintext
prompt-bom/
├── cmd/bom/          # CLI entrypoint and command implementations
│   ├── main.go
│   ├── init.go
│   ├── validate.go
│   └── build.go
├── internal/
│   ├── domain/       # Core structs (BOM/Component), business logic
│   │   └── bom.go
│   ├── app/          # (Planned) Application layer
│   └── infra/        # (Planned) File I/O, external integrations, etc.
├── spec/             # Specification files
│   └── prompt.bom.schema.json # JSON Schema for BOM YAML
├── examples/         # Usage examples
│   └── components/   # Sample component files (*.md)
├── docs/             # Documentation
│   ├── testing.md
│   └── img/          # (Logo image)
├── LICENSE           # Apache License 2.0 text
├── go.mod
├── go.sum
└── README.ja.md      # Japanese README
└── README.md         # English README (This file)
```

## License

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0).

## Roadmap

After completing the MVP, the following feature enhancements are planned:

- Dify Integration 
- LangChain Integration 
- Automated Regression Testing 
- Hierarchical BOMs

Refer to `document/Prompt-bom-overview.md` for details. 