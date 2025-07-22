# Genesis Documentation

Welcome to the Genesis documentation! Genesis is a powerful project scaffolding and task runner tool that helps you start and manage projects across different tech stacks with a unified interface.

## Quick Links

- [Getting Started](getting-started/index.md) - Install Genesis and create your first project
- [Core Concepts](concepts/index.md) - Learn about the fundamental concepts
- [Guides](guides/index.md) - Step-by-step guides for common tasks
- [Advanced Usage](advanced/index.md) - Advanced features and patterns
- [Reference](reference/index.md) - Detailed technical reference
- [FAQ](faq.md) - Frequently asked questions

## Features

- **Project Scaffolding**: Create new projects from templates with `genesis new`
- **Task Running**: Execute common development tasks with `genesis run`
- **Template Management**: Create, validate, and share templates
- **Convention over Configuration**: Sensible defaults with flexibility when needed
- **Cross-platform**: Works on Linux, macOS, and Windows
- **Extensible**: Customize with hooks and tasks

## Installation

```bash
go install github.com/felipevolpatto/genesis@latest
```

## Quick Start

Create a new project:
```bash
genesis new myapp --template https://github.com/example/template
```

List available tasks:
```bash
genesis run list
```

Run a task:
```bash
genesis run test
```

## Contributing

We welcome contributions! See our [Contributing Guide](https://github.com/felipevolpatto/genesis/blob/main/CONTRIBUTING.md) for details.

## License

Genesis is open source software licensed under the [MIT License](https://github.com/felipevolpatto/genesis/blob/main/LICENSE). 