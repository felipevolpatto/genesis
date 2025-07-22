# Template Configuration Reference

The `template.toml` file is the heart of a Genesis template. It defines variables, hooks, and other settings that control how the template works.

## File Structure

```toml
version = "1.0"  # Required

[vars]
  # Variables to collect from the user
  name = { prompt = "Project name:", default = "my-project" }
  description = { prompt = "Description:", default = "" }
  author = { 
    prompt = "Author name:", 
    default = "", 
    regex = "^[a-zA-Z ]+$"  # Optional validation
  }

[hooks]
  # Commands to run before scaffolding
  pre = [
    "echo 'Setting up project...'"
  ]
  # Commands to run after scaffolding
  post = [
    "go mod tidy",
    "git init"
  ]
```

## Version Field

The `version` field is required and should be set to `"1.0"`. This helps Genesis ensure compatibility with future template formats.

## Variables Section

The `vars` section defines variables that will be:
1. Collected from the user during project creation
2. Available in template files
3. Available in hook commands as environment variables

Each variable has the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `prompt` | string | Yes | The prompt shown to the user |
| `default` | string | No | Default value if user provides no input |
| `regex` | string | No | Regular expression for input validation |

Example:
```toml
[vars]
  # Basic variable with prompt and default
  name = { prompt = "Project name:", default = "my-project" }
  
  # Variable with validation
  email = { 
    prompt = "Email:", 
    default = "", 
    regex = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  }
```

## Hooks Section

The `hooks` section defines commands to run before (`pre`) and after (`post`) the template is applied.

### Pre-hooks

Pre-hooks run before any files are created. Use them for:
- Validating prerequisites
- Setting up the environment
- Downloading dependencies

```toml
[hooks]
  pre = [
    "command check-deps",
    "echo 'Setting up...'"
  ]
```

### Post-hooks

Post-hooks run after all files are created. Use them for:
- Initializing Git repository
- Installing dependencies
- Running formatters or linters
- Any other project setup

```toml
[hooks]
  post = [
    "go mod tidy",
    "git init",
    "pre-commit install"
  ]
```

### Hook Environment Variables

All variables defined in the `vars` section are available to hooks as environment variables:
- Variable names are converted to uppercase
- Special characters are replaced with underscores

Example:
```toml
[vars]
  name = { prompt = "Name:", default = "project" }
  go_version = { prompt = "Go version:", default = "1.21" }

[hooks]
  post = [
    # NAME="project" GO_VERSION="1.21" will be set
    "echo \"Creating $NAME with Go $GO_VERSION\""
  ]
```

## Using Variables in Templates

Variables can be used in template files (files ending in `.tmpl`) using Go's template syntax:

```go
// main.go.tmpl
package main

// {{ .name }} - Created by {{ .author }}
func main() {
    println("Welcome to {{ .name }}!")
}
```

## Template File Processing

- Files ending in `.tmpl` are processed using Go's template engine
- Other files are copied as-is
- Files and directories starting with `.` are ignored by default

## Best Practices

1. **Variable Names**:
   - Use descriptive names
   - Stick to lowercase letters, numbers, and underscores
   - Avoid special characters

2. **Prompts**:
   - Make them clear and specific
   - Include units or format if applicable
   - Use consistent punctuation

3. **Validation**:
   - Use regex validation for critical fields
   - Keep regex patterns simple and maintainable
   - Document expected formats

4. **Hooks**:
   - Keep commands idempotent
   - Handle errors gracefully
   - Document prerequisites

5. **Templates**:
   - Use consistent variable naming
   - Add comments explaining template logic
   - Test with different variable values 