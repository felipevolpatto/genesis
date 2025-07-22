# Getting Started with Genesis

This guide will help you get started with Genesis, a powerful project scaffolding tool.

## Installation

You can install Genesis using Go's package manager:

```bash
go install github.com/felipevolpatto/genesis@latest
```

## Basic Usage

### Creating a New Project

To create a new project from a template:

```bash
genesis new my-project --template https://github.com/example/template
```

This command will:
1. Clone the template repository
2. Process any variables defined in `template.toml`
3. Apply the template to create your new project
4. Run any post-hooks defined in the template

### Using Project Tasks

Once you have a project created with Genesis, you can use the task runner:

```bash
# List available tasks
genesis run list

# Run a specific task
genesis run test
```

## Next Steps

- Learn about [creating templates](../guides/creating-templates.md)
- Explore [template configuration](../reference/template-config.md)
- See [project configuration](../reference/project-config.md) 