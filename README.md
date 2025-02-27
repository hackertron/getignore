# Gitignore Generator

A simple CLI tool to generate `.gitignore` files for your projects based on various technology templates from GitHub's gitignore repository.

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

### Update templates

To force update templates from GitHub:

```
gitignore update
```

### Remove templates

To remove all locally stored templates:

```
gitignore clean
```

### Get help

To see all available commands and usage information:

```
gitignore help
```

You can also use `-h` or `--help` flags.

### Specify output file

By default, the tool creates a file named `.gitignore` in the current directory. You can specify a different output file as the second argument:

```
gitignore Python my-python-gitignore
```

## Features

- Works offline after initial template download
- Auto-downloads templates from GitHub when first used
- Simple command-line interface with help documentation
- Supports over 200 different technologies and frameworks
- Case-insensitive template matching
- Confirmation before overwriting existing files
- Templates organized by category for easy browsing
- Ability to update or remove templates as needed

## How it Works

The first time you run the tool, it will:
1. Create a `.gitignore-cli` directory in your home folder
2. Download all templates from GitHub's gitignore repository
3. Store them locally for fast, offline use

All subsequent uses will be instant and offline, using the cached templates.

## Template Structure

The templates are organized in the following structure:
- Main: Common templates for popular languages and frameworks
- Global: Templates for various editors, tools, and operating systems
- Community: Specialized templates for other tools and frameworks

## License

This project is licensed under the MIT License - see the LICENSE file for details.
