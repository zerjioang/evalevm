# EvalEVM 🦁
[![Go Report Card](https://goreportcard.com/badge/github.com/zerjioang/evalevm)](https://goreportcard.com/report/github.com/zerjioang/evalevm)
[![GoDoc](https://godoc.org/github.com/zerjioang/evalevm?status.svg)](https://godoc.org/github.com/zerjioang/evalevm)
[![Build Status](https://travis-ci.org/zerjioang/evalevm.svg?branch=master)](https://travis-ci.org/zerjioang/evalevm)
![GitHub](https://img.shields.io/github/license/zerjioang/evalevm)

> A modular framework for evaluating and comparing EVM static analysis tools using Docker.

## Documentation 📚

Comprehensive documentation is available in the [`docs/`](./docs) directory:

-   [**Introduction**](./docs/introduction.md): Overview and scientific context.
-   [**Usage Guide**](./docs/usage.md): Installation, CLI commands, and configuration.
-   [**Tools**](./docs/tools.md): List of supported analyzers and how to add new ones.
-   [**Architecture**](./docs/architecture.md): System design and components.
-   [**Glossary**](./docs/glossary.md): Terminology.

## Quick Start 🚀

1.  **Build**: `make build`
2.  **Build Analyzer Images**: `./dist/evalevm analyzer build`
3.  **Scan**: `./dist/evalevm scan evm -b 6080604052...`

## Features ✨

-   **Docker Isolation**: Runs tools in ephemeral containers.
-   **Parallel Execution**: Utilizes all CPU cores for batch processing.
-   **Unified Metrics**: Standardizes output for easy comparison.
-   **Result Caching**: Skips redundant scans automatically.

## License

MIT
