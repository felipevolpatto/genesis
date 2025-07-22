# Contributing to Genesis

First off, thank you for considering contributing to Genesis!

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* Use a clear and descriptive title
* Describe the exact steps which reproduce the problem
* Provide specific examples to demonstrate the steps
* Describe the behavior you observed after following the steps
* Explain which behavior you expected to see instead and why
* Include any error messages or logs

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* Use a clear and descriptive title
* Provide a step-by-step description of the suggested enhancement
* Provide specific examples to demonstrate the steps
* Describe the current behavior and explain which behavior you expected to see instead
* Explain why this enhancement would be useful

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the Go coding style
* Include appropriate test coverage
* Document new code based on the Documentation Styleguide

## Development Process

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Write your code
4. Write or update tests
5. Run the test suite
6. Push your branch and submit a pull request

### Setting Up the Development Environment

1. Install Go 1.21 or later
2. Clone your fork:
   ```bash
   git clone git@github.com:your-username/genesis.git
   ```
3. Add the main repository as a remote:
   ```bash
   git remote add upstream https://github.com/felipevolpatto/genesis.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Code Style

* Follow the standard Go formatting rules (use `gofmt`)
* Use meaningful variable and function names
* Write comments for non-obvious code sections
* Keep functions focused and small
* Use proper error handling

### Documentation Styleguide

* Use Markdown
* Reference functions and variables with backticks: \`myFunction()\`
* Include code examples where appropriate
* Keep line length to a maximum of 80 characters
* Use proper spelling and grammar

## Project Structure

```
genesis/
├── cmd/                # Command definitions
│   ├── root.go        # Root command
│   ├── new.go         # Project creation
│   ├── run.go         # Task runner
│   └── template.go    # Template management
├── internal/
│   ├── config/        # Configuration parsing
│   ├── runner/        # Task execution
│   ├── scaffolder/    # Project scaffolding
│   └── tui/          # Terminal UI helpers
├── templates/         # Built-in templates
└── main.go           # Entry point
```

## Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line
* Consider starting the commit message with an applicable emoji:
    * 🎨 `:art:` when improving the format/structure of the code
    * 🐎 `:racehorse:` when improving performance
    * 🚱 `:non-potable_water:` when plugging memory leaks
    * 📝 `:memo:` when writing docs
    * 🐛 `:bug:` when fixing a bug
    * 🔥 `:fire:` when removing code or files
    * 💚 `:green_heart:` when fixing the CI build
    * ✅ `:white_check_mark:` when adding tests
    * 🔒 `:lock:` when dealing with security
    * ⬆️ `:arrow_up:` when upgrading dependencies
    * ⬇️ `:arrow_down:` when downgrading dependencies

## Additional Notes

### Issue and Pull Request Labels

* `bug` - Something isn't working
* `enhancement` - New feature or request
* `documentation` - Improvements or additions to documentation
* `good first issue` - Good for newcomers
* `help wanted` - Extra attention is needed
* `question` - Further information is requested

## Recognition

Contributors who submit a PR that gets merged will be added to the Contributors list in the README. 