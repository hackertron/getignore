# Gitignore Generator

A simple CLI tool to generate `.gitignore` files for your projects based on various technology templates.

## Installation

1. Clone this repository
2. Build the binary:
   ```
   go build -o gitignore main.go
   ```
3. (Optional) Add the binary to your PATH

## Usage

### Generate a gitignore file

To generate a `.gitignore` file for a specific framework or technology:

```
gitignore <framework-name>
```

For example:
```
gitignore Go
```

This will create a `.gitignore` file in the current directory with Go-specific ignore patterns.

### List available templates

To see a list of all available templates:

```
gitignore list
```

### Specify output file

By default, the tool creates a file named `.gitignore` in the current directory. You can specify a different output file as the second argument:

```
gitignore Python my-python-gitignore
```

## Features

- Works completely offline with bundled templates
- Simple command-line interface
- Supports over 200 different technologies and frameworks
- Case-insensitive template matching
- Confirmation before overwriting existing files

## Template Structure

The templates are organized in the following structure:
- Root directory: Common templates for popular languages and frameworks
- `Global/`: Templates for various editors, tools, and operating systems
- `community/`: Specialized templates for other tools and frameworks

## License

This project is licensed under the MIT License - see the LICENSE file for details.
