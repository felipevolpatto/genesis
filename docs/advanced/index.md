# Advanced Usage

This section covers advanced features and usage patterns for Genesis.

## Template Development

### Template Inheritance

You can create templates that extend other templates:

1. Clone the base template
2. Add your customizations
3. Update `template.toml`:
   ```toml
   version = "1.0"
   
   [vars]
   # Include base template variables
   name = { prompt = "Project name:", default = "my-project" }
   # Add your custom variables
   database = { prompt = "Database type:", default = "postgres" }
   
   [hooks]
   pre = [
     # Run base template pre-hooks
     "git clone https://github.com/base/template base",
     "cp -r base/* .",
     "rm -rf base"
   ]
   post = [
     # Run your custom hooks
     "./setup-db.sh",
     # Run base template post-hooks
     "go mod tidy"
   ]
   ```

### Dynamic Templates

Use Go template features for advanced customization:

{% raw %}
```go
// config.go.tmpl
package config

type Config struct {
  {{- if eq .database "postgres" }}
  PostgresConfig
  {{- else if eq .database "mysql" }}
  MySQLConfig
  {{- end }}
}

{{- if eq .database "postgres" }}
type PostgresConfig struct {
  Host string
  Port int
}
{{- else if eq .database "mysql" }}
type MySQLConfig struct {
  Socket string
}
{{- end }}
```
{% endraw %}

### Conditional Files

Use hooks to conditionally include files:

```toml
[hooks]
post = [
  '''
  if [[ "${DATABASE}" == "postgres" ]]; then
    cp postgres.yml docker-compose.yml
  else
    cp mysql.yml docker-compose.yml
  fi
  '''
]
```

## Task Automation

### Task Dependencies

Create tasks that depend on other tasks:

```toml
[tasks]
  lint = { description = "Run linters", cmd = "golangci-lint run" }
  test = { description = "Run tests", cmd = "go test ./..." }
  
  # Combine tasks using shell operators
  check = {
    description = "Run all checks",
    cmd = "genesis run lint && genesis run test"
  }
```

### Parallel Tasks

Run tasks in parallel:

```toml
[tasks]
  # Use & to run in background, wait for all
  parallel-test = {
    description = "Run tests in parallel",
    cmd = """
      genesis run test-unit & \
      genesis run test-integration & \
      wait
    """
  }
```

### Environment-specific Tasks

Configure tasks for different environments:

```toml
[tasks]
  # Development database
  db-dev = {
    description = "Start development database",
    cmd = "docker-compose up db",
    env = { "POSTGRES_DB" = "myapp_dev" }
  }

  # Test database
  db-test = {
    description = "Start test database",
    cmd = "docker-compose up db",
    env = { "POSTGRES_DB" = "myapp_test" }
  }
```

## Git Integration

### Private Templates

Access private templates:

1. Using SSH:
   ```bash
   genesis new myapp --template git@github.com:org/private-template.git
   ```

2. Using HTTPS with token:
   ```bash
   genesis new myapp --template https://token@github.com/org/private-template.git
   ```

### Template Versioning

Use specific template versions:

```bash
# Use a tag
genesis new myapp --template https://github.com/org/template.git --version v1.0.0

# Use a branch
genesis new myapp --template https://github.com/org/template.git --version feature/new-stack

# Use a commit hash
genesis new myapp --template https://github.com/org/template.git --version 4eeee43
```

## Hook Scripts

### Error Handling

Add error handling to hooks:

```toml
[hooks]
pre = [
  '''
  set -e  # Exit on error
  
  echo "Checking dependencies..."
  if ! command -v docker &> /dev/null; then
    echo "Error: docker is required but not installed"
    exit 1
  fi
  
  echo "Running setup..."
  ./setup.sh || {
    echo "Error: setup failed"
    exit 1
  }
  '''
]
```

### Cross-platform Hooks

Write hooks that work on multiple platforms:

```toml
[hooks]
post = [
  '''
  # Detect OS
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    brew install dependency
  elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    apt-get install dependency
  elif [[ "$OSTYPE" == "msys" ]]; then
    # Windows
    choco install dependency
  else
    echo "Unsupported OS: $OSTYPE"
    exit 1
  fi
  '''
]
```

## Custom Template Functions

Genesis uses Go's template engine, which provides many built-in functions:

{% raw %}
```go
// template.go.tmpl
package main

import "time"

// Using built-in functions
const Version = "{{ .version | printf "v%s" }}"
const BuildTime = "{{ now | date "2006-01-02" }}"

// Conditional logic
{{- if eq .env "production" }}
const Debug = false
{{- else }}
const Debug = true
{{- end }}

// Looping
var Dependencies = []string{
{{- range .dependencies }}
  "{{ . }}",
{{- end }}
}
```
{% endraw %}

## Security Considerations

1. **Template Verification**:
   - Always verify template sources
   - Review hooks before running
   - Use trusted templates

2. **Hook Security**:
   - Avoid running untrusted scripts
   - Review commands before execution
   - Use environment variables for secrets

3. **Private Data**:
   - Don't commit sensitive data
   - Use environment variables
   - Add sensitive files to `.gitignore` 