# Frequently Asked Questions

## General

### What is Genesis?
Genesis is a project scaffolding and task runner tool that helps you start and manage projects across different tech stacks with a unified interface. It provides project scaffolding from templates, task running, and template management capabilities.

### Why use Genesis?
Genesis reduces setup time and cognitive load by providing:
- Consistent project structure through templates
- Unified task running interface
- Convention over configuration
- Cross-platform compatibility
- Extensibility through hooks and tasks

### How is Genesis different from other tools?
Genesis combines project scaffolding and task running in a single, language-agnostic tool. Unlike language-specific tools, Genesis works with any tech stack and provides a consistent interface across all your projects.

## Templates

### How do I create a template?
See our [Creating Templates](guides/creating-templates.md) guide for detailed instructions. In short:
1. Create a Git repository
2. Add your project structure
3. Create `template.toml`
4. Add template variables using `.tmpl` files
5. Define hooks if needed

### Can I use private templates?
Yes! You can use private templates by:
1. Using SSH: `genesis new myapp --template git@github.com:org/private-template.git`
2. Using HTTPS with token: `genesis new myapp --template https://token@github.com/org/private-template.git`

### How do I update a template?
Templates are versioned using Git. To update:
1. Make changes in the template repository
2. Commit and tag a new version
3. Users can use the new version with `--version`

### Can I extend existing templates?
Yes! See [Template Inheritance](advanced/index.md#template-inheritance) in our advanced guide for details on extending templates.

## Tasks

### How do I define tasks?
Tasks are defined in `genesis.toml` under the `[tasks]` section. Each task needs a description and command:
```toml
[tasks]
  test = { description = "Run tests", cmd = "go test ./..." }
```

### Can tasks depend on other tasks?
Yes! You can create task dependencies using shell operators:
```toml
[tasks]
  check = { 
    description = "Run all checks",
    cmd = "genesis run lint && genesis run test"
  }
```

### How do I run tasks in parallel?
Use shell background execution and `wait`:
```toml
[tasks]
  parallel-test = {
    cmd = "genesis run test-unit & genesis run test-integration & wait"
  }
```

### Can I use environment variables in tasks?
Yes! Define them in the task configuration:
```toml
[tasks]
  deploy = {
    cmd = "./deploy.sh",
    env = { "ENV" = "production" }
  }
```

## Troubleshooting

### Template not found
- Check the template URL is correct
- Ensure you have access to the repository
- Check your Git credentials
- Try using SSH if HTTPS fails

### Hook execution failed
Common causes:
1. Missing dependencies
2. Incorrect permissions
3. Platform-specific commands
4. Environment variables not set

Solutions:
1. Add dependency checks in pre-hooks
2. Use `chmod` in hooks if needed
3. Write cross-platform hooks
4. Document required environment variables

### Task not found
- Ensure you're in a Genesis project directory
- Check `genesis.toml` exists and is valid
- Run `genesis run list` to see available tasks
- Check for typos in task names

### Template validation failed
Common issues:
1. Missing `template.toml`
2. Invalid TOML syntax
3. Invalid template syntax in `.tmpl` files
4. Missing required fields

Solutions:
1. Create `template.toml`
2. Validate TOML syntax
3. Check template syntax
4. Add required fields

## Best Practices

### Should I commit `genesis.toml`?
Yes! `genesis.toml` should be committed as it defines:
- Project template source
- Common development tasks
- Project-specific configuration

### How should I organize tasks?
Best practices:
1. Use descriptive names
2. Group related tasks with prefixes
3. Add clear descriptions
4. Document prerequisites
5. Keep commands simple

### How do I handle sensitive data?
Never store sensitive data in templates or tasks:
1. Use environment variables
2. Add sensitive files to `.gitignore`
3. Document required secrets
4. Use secure credential storage

### Should I use pre-hooks or post-hooks?
Use:
- Pre-hooks for:
  - Dependency checks
  - Environment setup
  - Validation
- Post-hooks for:
  - Dependency installation
  - Build steps
  - Git initialization 