# CONTRIBUTING.md

# Contributing to Jou

Thank you for your interest in contributing to jou! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.23.0 or later
- Git
- A terminal that supports ANSI escape codes

### Setting Up Your Development Environment

1. **Fork the repository**

   ```bash
   # Go to GitHub and fork the repository
   # Then clone your fork
   git clone https://github.com/YOUR_USERNAME/jou.git
   cd jou
   ```

2. **Add the upstream remote**

   ```bash
   git remote add upstream https://github.com/cheersmas/jou.git
   ```

3. **Install dependencies**

   ```bash
   go mod download
   ```

4. **Build and test**
   ```bash
   go build -o jou main.go
   ./jou
   ```

## 🔄 Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b bugfix/your-bugfix-name
```

### 2. Make Your Changes

- Write clean, readable code
- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Test the application manually
go run main.go
```

### 4. Commit Your Changes

```bash
git add .
git commit -m "feat: add new feature description"
```

**Commit Message Format:**

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `style:` for formatting changes
- `refactor:` for code refactoring
- `test:` for adding tests
- `chore:` for maintenance tasks

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## 📝 Code Style Guidelines

### Go Code Style

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Project Structure

- Follow the existing clean architecture pattern
- Keep domain logic separate from UI logic
- Use interfaces for dependency injection
- Maintain separation between layers (app, services, repositories, domains)

### Testing

- Write unit tests for new functionality
- Test edge cases and error conditions
- Aim for good test coverage
- Use table-driven tests where appropriate

## �� Reporting Issues

### Before Creating an Issue

1. Search existing issues to avoid duplicates
2. Check if the issue is already fixed in the latest version
3. Gather relevant information (OS, Go version, error messages)

### Issue Template

Use the provided issue templates:

- Bug Report
- Feature Request
- Documentation Issue
- Question

## 🚀 Feature Development

### Before Starting

1. Check if a similar feature is already requested
2. Discuss major features in an issue first
3. Consider the impact on existing functionality
4. Think about backward compatibility

### Implementation Guidelines

1. **Start Small**: Implement the core functionality first
2. **Add Tests**: Write tests as you develop
3. **Update Documentation**: Keep README and code comments up to date
4. **Consider UX**: Think about how users will interact with your feature
5. **Performance**: Consider the performance impact of your changes

## �� UI/UX Guidelines

### Terminal Interface

- Keep the interface clean and intuitive
- Use consistent keyboard shortcuts
- Provide helpful feedback to users
- Handle edge cases gracefully
- Ensure the interface works in different terminal sizes

### Styling

- Use the existing style system in `app/styles/`
- Maintain consistency with the current design
- Consider accessibility (color contrast, etc.)
- Test in different terminal themes

## 📋 Pull Request Guidelines

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] Tests pass locally
- [ ] Documentation is updated
- [ ] No merge conflicts
- [ ] Commit messages are clear and descriptive

### PR Description

Include:

- Description of changes
- Motivation for the change
- How to test the changes
- Screenshots (if UI changes)
- Breaking changes (if any)

### Review Process

1. Automated checks must pass
2. Code review by maintainers
3. Testing by maintainers
4. Approval and merge

## 🏷️ Labels and Milestones

### Issue Labels

- `bug`: Something isn't working
- `enhancement`: New feature or request
- `documentation`: Improvements or additions to documentation
- `question`: Further information is requested
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention is needed
- `needs-triage`: Needs maintainer attention

### Pull Request Labels

- `ready for review`: Ready for maintainer review
- `work in progress`: Still being developed
- `needs testing`: Requires testing
- `breaking change`: Contains breaking changes

## 🎯 Good First Issues

Look for issues labeled with `good first issue` if you're new to the project. These are typically:

- Small bug fixes
- Documentation improvements
- Simple feature additions
- Code cleanup tasks

## 💬 Communication

- Use GitHub issues for bug reports and feature requests
- Use GitHub discussions for general questions and ideas
- Be respectful and constructive in all interactions
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md)

## 🏆 Recognition

Contributors will be recognized in:

- README contributors section
- Release notes
- GitHub contributors page

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## ❓ Questions?

If you have questions about contributing:

1. Check existing issues and discussions
2. Create a new issue with the "question" label
3. Join our community discussions

Thank you for contributing to jou! 🎉
