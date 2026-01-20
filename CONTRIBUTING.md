# Contributing to Copy Image Tool 🤝

First off, thank you for considering contributing to Copy Image Tool! It's people like you that make Copy Image Tool such a great tool.

There are many ways to contribute, from writing code to filing issues on GitHub.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
    - [Reporting Bugs](#reporting-bugs)
    - [Suggesting Enhancements](#suggesting-enhancements)
    - [Pull Requests](#pull-requests)
- [Styleguides](#styleguides)
    - [Go Styleguide](#go-styleguide)
    - [Frontend Styleguide](#frontend-styleguide)
- [Development Workflow](#development-workflow)
    - [Setup](#setup)
    - [Running Tests](#running-tests)
    - [Linting](#linting)

## 📜 Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## ❓ How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report. Following these guidelines helps maintainers and the community understand your report, reproduce the behavior, and find related reports.

Beyond reporting bugs, you can also contribute by **fixing bugs**! Check out the [issue tracker](https://github.com/hoangtran1411/copy-image/issues) for bugs labeled "bug" or "good first issue".

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion, including completely new features and minor improvements to existing functionality.

- **Check for existing suggestions**: Before creating a new one, please check the [issue tracker](https://github.com/hoangtran1411/copy-image/issues) to see if someone has already suggested it.
- **Use a clear and descriptive title**: For the issue to identify the suggestion.
- **Provide a step-by-step description of the suggested enhancement**: In as many details as possible.

### Pull Requests

1.  **Fork the repo** and create your branch from `main`.
2.  **Ensure your code follows the styleguide**.
3.  **Add tests** for the new functionality or bug fix.
4.  **Run the linter and tests** to ensure everything is working as expected.
5.  **Submit a PR** with a clear description of the changes.

## 🎨 Styleguides

### Go Styleguide

- All Go code should be formatted using `gofmt` or `goimports`.
- Follow the guidelines in [Effective Go](https://golang.org/doc/effective_go.html).
- Use descriptive variable and function names.
- Document all exported functions, variables, and types.

### Frontend Styleguide

- Use Vanilla CSS for styling (unless Tailwind is requested).
- Maintain a consistent theme (Dark Mode by default).
- Ensure the UI is responsive and accessible.

## 🛠️ Development Workflow

### Setup

```bash
# Clone the repository
git clone https://github.com/hoangtran1411/copy-image.git
cd copy-image

# Install dependencies
go mod download
```

### Running Tests

We use the standard Go testing tool. Please make sure all tests pass before submitting a PR.

```bash
# Run all tests
make test

# Run tests with coverage
make coverage
```

### Linting

We use `golangci-lint` to maintain code quality. Please fix any linting errors before submitting a PR.

```bash
# Run the linter
make lint
```

## 💎 Premium Aesthetics

If you are working on the UI, please focus on "Premium" aesthetics:
- Use curated color palettes.
- Implement smooth transitions and micro-animations.
- Follow a consistent design language.

## 📧 Contact

If you have any questions, feel free to open an issue or reach out to the maintainer.

Happy coding! 🚀
