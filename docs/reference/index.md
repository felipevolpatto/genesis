# Reference Documentation

This section contains detailed reference documentation for Genesis configuration files and features.

## Configuration Files

- [Template Configuration](template-config.md) - Detailed reference for `template.toml` files
- [Project Configuration](project-config.md) - Detailed reference for `genesis.toml` files

## Command Line Interface

### Global Flags

- `--help` - Show help for any command
- `--version` - Show Genesis version

### Commands

#### `new`
Create a new project from a template:
```bash
genesis new [project-name] --template [url] [--version version] [--yes]
```

Flags:
- `--template` - Git URL of the template repository
- `--version` - Specific version of the template (commit hash, tag, or branch)
- `--yes` - Skip prompts and use default values

#### `run`
Run a task defined in `genesis.toml`:
```bash
genesis run [task-name]
genesis run list  # List available tasks
```

#### `template`
Manage and validate templates:
```bash
genesis template list  # List available templates
genesis template validate [path]  # Validate a template
``` 