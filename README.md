# Gitignore Generator

A simple CLI tool to generate `.gitignore` files for your projects based on various technology templates from GitHub's gitignore repository.

## Installation

1. Clone this repository
2. Build the binary:
   ```
   go build -o gitignore main.go
   ```
3. (Optional) Add the binary to your PATH:

   **Linux/macOS:**
   ```bash
   # Move the binary to a location in your PATH
   sudo mv gitignore /usr/local/bin/
   
   # Or add the current directory to your PATH in ~/.bashrc or ~/.zshrc
   echo 'export PATH=$PATH:'$(pwd) >> ~/.bashrc
   source ~/.bashrc
   ```

   **Windows:**
   ```cmd
   # Add the directory to your PATH environment variable
   setx PATH "%PATH%;%cd%"
   
   # Or move the binary to a location that's already in your PATH
   # For example:
   # move gitignore.exe C:\Windows\System32\
   ```

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

If the template isn't available locally, the tool will download only that specific template instead of downloading all templates.

### Download all templates

To download all templates at once from GitHub:

```
gitignore download-all
```

You can also use the `--download-all` flag with any command to download all templates:

```
gitignore Go --download-all
```

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

- Efficient template management - only downloads templates as needed
- On-demand downloading of single templates when requested
- Option to download all templates at once when desired
- Works offline after templates are downloaded
- Simple command-line interface with help documentation
- Supports over 200 different technologies and frameworks
- Case-insensitive template matching
- Confirmation before overwriting existing files
- Templates organized by category for easy browsing
- Ability to update or remove templates as needed

## How it Works

When you request a template:
1. The tool checks if it exists in the `.gitignore-cli` directory in your home folder
2. If available locally, it uses the cached version for instant access
3. If not available, it downloads just that specific template from GitHub
4. Templates are stored locally for future use

You can also choose to download all templates at once using the `download-all` command or `--download-all` flag if you prefer to have everything available offline.

## Template Structure

The templates are organized in the following structure:
- Main: Common templates for popular languages and frameworks
- Global: Templates for various editors, tools, and operating systems
- Community: Specialized templates for other tools and frameworks

## License

This project is licensed under the MIT License - see the LICENSE file for details.