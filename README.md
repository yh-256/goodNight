# ZenWatch

ZenWatch is a command-line tool for analyzing Git repositories.

## Usage

ZenWatch is a command-line tool that currently supports one main command: `analyze`.

### `analyze`

This command analyzes a Git repository.

**Synopsis:**

```shell
zenwatch analyze <repository-url> [flags]
```

**Arguments:**

*   `<repository-url>`: The URL of the Git repository to analyze.

**Flags:**

*   `--out <output-file>`: Specifies the path to save the output Markdown report. Defaults to `reports/latest.md`.

**Example:**

```shell
zenwatch analyze https://github.com/example/project.git --out project_report.md
```

## Building from Source

To build ZenWatch from source, you need to have Go installed on your system.

1.  **Clone the repository (if you haven't already):**
    ```shell
    git clone <repository-url> # Replace <repository-url> with the actual URL
    cd zenwatch
    ```
2.  **Build the project:**
    ```shell
    go build ./cmd/zenwatch
    ```
    This will create a `zenwatch` executable in the current directory.

## Running Tests

To run the tests for ZenWatch, navigate to the root directory of the project and use the following Go command:

```shell
go test ./...
```
This command will discover and run all test files in the project.

## Contributing

Contributions to ZenWatch are welcome! If you find any issues or have suggestions for improvements, please feel free to:

1.  **Open an issue:** Describe the bug or feature request in detail.
2.  **Fork the repository:** Create your own fork of the project.
3.  **Create a new branch:** Make your changes in a dedicated branch.
4.  **Commit your changes:** Write clear and concise commit messages.
5.  **Push your changes:** Push your branch to your fork.
6.  **Submit a pull request:** Explain the changes you've made and why they should be merged.

Please ensure your code adheres to the existing style and that all tests pass before submitting a pull request.

## License

This project is licensed under the terms of the LICENSE file.
