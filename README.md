# Ditto

[![Go Version](https://img.shields.io/badge/go-1.22-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/arthvm/ditto)](https://goreportcard.com/report/github.com/arthvm/ditto)

Ditto is a command-line tool that uses AI to generate git commit messages from your staged changes, helping you create clean and conventional commits with ease.

## Overview

Writing conventional commit messages can be tedious. Ditto simplifies this process by analyzing your staged changes (`git diff --staged`) and generating a concise, conventional commit message for you. It uses Large Language Models (LLMs) to generate the messages and supports multiple providers.

## Features

-   Generate conventional commit messages from staged files.
-   Supports multiple LLM providers (Gemini, Ollama).
-   Customizable prompts to give the model more context.
-   Simple and easy-to-use command-line interface.

## Installation

You can install `ditto` using `go install`:

```sh
go install github.com/arthvm/ditto@latest
```

Alternatively, you can check the [Releases](https://github.com/arthvm/ditto/releases) page for pre-compiled binaries.

## Usage

1.  Stage your changes using `git add`.
2.  Run the `ditto commit` command:

```sh
ditto commit
```

This will generate a commit message using the default provider (Gemini) and commit the changes.

### Options

-   `--provider`: Select the LLM provider to use.

```sh
ditto commit --provider ollama
```

-   `--prompt`: Provide additional context to the model for a more accurate commit message.

```sh
ditto commit --prompt "This change is part of the new authentication feature."
```

## Configuration

### Gemini

To use the Gemini provider, you need to have the [Gemini CLI](https://github.com/google/generative-ai-go) installed and configured on your machine. Please follow their official instructions to set up your API key.

### Ollama

To use the Ollama provider, you need to have an Ollama server running. By default, `ditto` will try to connect to `http://localhost:11434`. You can configure a different host by setting the `OLLAMA_HOST` environment variable.

```sh
export OLLAMA_HOST="http://your-ollama-host:11434"
```

## Supported Providers

-   `gemini` (default)
-   `gemini-flash`
-   `gemini-flash-lite`
-   `ollama`

## Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature`).
3.  Make your changes.
4.  Commit your changes (`ditto commit`).
5.  Push to the branch (`git push origin feature/your-feature`).
6.  Open a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

**Note:** The `LICENSE` file in this repository is currently empty. It is recommended to add the full MIT License text to it.
