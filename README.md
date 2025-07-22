# Genesis

```
  ____                      _     
 / ___| ___ _ __   ___  ___(_)___ 
| |  _ / _ \ '_ \ / _ \/ __| / __|
| |_| |  __/ | | |  __/\__ \ \__ \
 \____|\___|_| |_|\___||___/_|___/
```

> Begin any project, unified.

Genesis is a powerful project scaffolding and task runner tool that helps you start and manage projects across different tech stacks with a unified interface.

## Features

- 🚀 **Project Scaffolding**: Initialize new projects from predefined or custom templates
- 🔄 **Task Runner**: Execute common development tasks defined in a project-local configuration file
- 🎯 **Convention over Configuration**: Sensible defaults with full customization options
- 🌐 **Language Agnostic**: Works with any programming language or framework
- 🔌 **Extensible**: Easy to create and share custom templates
- ⚡ **Fast & Portable**: Single binary with no external runtime dependencies

## Installation

### Using Go Install

```bash
go install github.com/felipevolpatto/genesis@latest
```

### From GitHub Releases

1. Visit the [Releases](https://github.com/felipevolpatto/genesis/releases) page
2. Download the binary for your platform
3. Add it to your PATH

## Quick Start

### Create a New Project

```bash
# Create a new project using a template
genesis new my-project --template https://github.com/example/template-go-cli

# Skip prompts and use defaults
genesis new my-project -t https://github.com/example/template-go-cli -y

# Use a specific template version
genesis new my-project -t https://github.com/example/template-go-cli -v v1.0.0
```

### Run Tasks

```bash
# List available tasks
genesis run list

# Run a specific task
genesis run test
genesis run build
```

### Manage Templates

```bash
# List available templates
genesis template list

# Validate a template
genesis template validate ./my-template
```

## Configuration Files

### template.toml

This file defines the variables and hooks for a template:

```toml
# The version of the template spec
version = "1.0"

# Variables that the user will be prompted to enter
[vars]
  description = { prompt = "Enter a short project description:", default = "A new Genesis project." }
  author = { prompt = "Enter the author's name:", default = "Your Name" }

# Hooks to run before or after scaffolding
[hooks]
  post = [
    "go mod tidy",
    "git init",
    "git add .",
    "git commit -m 'Initial commit from Genesis'",
  ]
```

### genesis.toml

This file defines the project's tasks:

```toml
# The version of the genesis config spec
version = "1.0"

[project]
  template_url = "https://github.com/user/go-cli-template"
  template_version = "v1.0.2"

[tasks]
  test = { description = "Run all unit tests.", cmd = "go test ./..." }
  build = { description = "Build the application binary.", cmd = "go build -o myapp main.go" }
  lint = { description = "Lint the source code.", cmd = "golangci-lint run" }
```

## Creating Templates

A template is a Git repository containing:

1. A `template.toml` file defining variables and hooks
2. Project files and directories
3. Template files (ending in `.tmpl`) that will be processed with variables

Example template structure:

```
my-template/
├── .gitignore.tmpl
├── main.go.tmpl
├── README.md.tmpl
└── template.toml
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/felipevolpatto/genesis.git
   cd genesis
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

4. Build the binary:
   ```bash
   go build -o genesis
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

### V1 (Current)

- ✅ Project scaffolding from Git templates
- ✅ Task running with genesis.toml
- ✅ Basic template validation
- ✅ Interactive variable prompts

### V2 and Beyond

- 🔄 Interactive mode with TUI
- 🔌 Plugin system for custom commands
- 📦 Template discovery and registry
- 🖥️ Optional GUI wrapper
- 🔄 Auto-updating mechanism 