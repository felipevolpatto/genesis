# Core Concepts

This section explains the fundamental concepts behind Genesis and how they work together.

## Project Scaffolding

Project scaffolding is the process of creating a new project from a template. In Genesis, this involves:

1. **Templates**: Git repositories containing:
   - Base project structure
   - Template files (`.tmpl`)
   - Configuration (`template.toml`)
   - Hooks for customization

2. **Variables**: Dynamic values that:
   - Are collected from the user
   - Customize the template
   - Can be validated
   - Are available in hooks

3. **Processing**: The scaffolding process:
   - Clones the template
   - Collects variables
   - Runs pre-hooks
   - Processes template files
   - Creates project structure
   - Runs post-hooks

## Task Running

Task running provides a unified way to execute common development commands:

1. **Tasks**: Defined in `genesis.toml` as:
   - Named commands
   - With descriptions
   - Optional environment variables
   - Optional working directories

2. **Execution**: Tasks are run via:
   - `genesis run <task-name>`
   - In the project directory
   - With specified environment
   - In the defined working directory

## Template Management

Templates in Genesis are:

1. **Git-based**: 
   - Stored in Git repositories
   - Versioned with tags/branches
   - Shareable and reusable

2. **Validated**:
   - Must contain `template.toml`
   - Must have valid template syntax
   - Can be checked with `genesis template validate`

3. **Discoverable**:
   - Listed via `genesis template list`
   - Can be official or community templates
   - Can be private or public

## Convention over Configuration

Genesis follows these conventions:

1. **File Processing**:
   - Files ending in `.tmpl` are processed
   - Other files are copied as-is
   - Hidden files are ignored by default

2. **Configuration**:
   - `template.toml` for template settings
   - `genesis.toml` for project settings
   - Version field for compatibility

3. **Directory Structure**:
   - Project root contains `genesis.toml`
   - Templates define their own structure
   - Working directory for tasks

## Extensibility

Genesis is extensible through:

1. **Templates**:
   - Custom variables
   - Pre/post hooks
   - Custom file structures

2. **Tasks**:
   - Custom commands
   - Environment variables
   - Working directories

3. **Integration**:
   - Works with any Git repository
   - Supports any build tools
   - Platform independent 