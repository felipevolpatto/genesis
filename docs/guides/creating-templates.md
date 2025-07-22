# Creating Templates

This guide explains how to create templates for Genesis.

## Template Structure

A Genesis template is a Git repository with the following structure:

```
my-template/
├── template.toml         # Template configuration
├── main.go.tmpl         # Template files (*.tmpl)
├── go.mod.tmpl
└── README.md.tmpl
```

## Template Configuration

The `template.toml` file defines:
- Variables to collect from the user
- Pre and post-hooks to run
- Other template settings

Example `template.toml`:

```toml
version = "1.0"

[vars]
  name = { prompt = "Project name:", default = "my-project" }
  description = { prompt = "Project description:", default = "" }
  author = { prompt = "Author name:", default = "" }

[hooks]
  pre = [
    "echo 'Setting up project...'"
  ]
  post = [
    "go mod tidy",
    "git init"
  ]
```

## Template Files

Template files use Go's `text/template` syntax:

```go
// main.go.tmpl
package main

func main() {
    println("Welcome to {{ .name }}!")
}
```

## Variables

Variables defined in `template.toml` can be used in:
- Template files (with `{{ .varname }}`)
- Hook commands (with environment variables)

## Testing Templates

Use the `genesis template validate` command to test your template:

```bash
genesis template validate path/to/template
```

## Publishing Templates

1. Push your template to a Git repository
2. Share the repository URL with users
3. Users can use it with:
   ```bash
   genesis new my-project --template https://github.com/user/template
   ``` 