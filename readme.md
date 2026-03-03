# EvalEVM 🦁

EvalEVM is a modular framework for evaluating and comparing Ethereum Virtual Machine (EVM) static analysis tools using Docker containers [1](#0-0) . It provides a unified interface to run multiple analysis tools on smart contract bytecode and standardizes their output for easy comparison [2](#0-1) .

---

## What EvalEVM Does

EvalEVM orchestrates multiple EVM analysis tools through a plugin architecture. The core `Comparator` struct manages a list of analyzers and executes them in parallel using Docker containers [3](#0-2) . Currently integrated tools include EVMLisa, Vandal, Rattle, EVM-CFG, EVMole, Paper, and others [4](#0-3) .

Each tool is implemented as a Go struct that satisfies the `datatype.Analyzer` interface, with methods for creating Docker tasks and parsing output [5](#0-4) . The framework handles task execution through a `WorkerPool` that manages Docker container lifecycle, timeouts, and result collection [6](#0-5) .

## Key Features

- **Docker Isolation**: Each tool runs in an ephemeral container with embedded Dockerfiles [7](#0-6) 
- **Parallel Execution**: Utilizes multiple CPU cores for batch processing [8](#0-7) 
- **Unified Metrics**: Standardizes output formats including CFG graphs, vulnerability detection, and performance metrics [9](#0-8) 
- **Result Caching**: Automatically skips redundant scans by caching results [10](#0-9) 
- **CSV Export**: Supports streaming results to CSV for large-scale analysis [11](#0-10) 

## Beneficial Use Cases

### Security Research and Academia
Researchers can compare the effectiveness of different static analysis techniques across multiple tools. The framework provides standardized metrics like CFG node/edge counts, vulnerability detection rates, and performance measurements [12](#0-11) .

### Smart Contract Auditing
Security firms can evaluate which tools best suit their audit workflows by running them against real contract bytecode and comparing results. Tools like Vandal and EVMLisa specialize in vulnerability detection [13](#0-12) [14](#0-13) .

### Tool Development
Developers creating new EVM analysis tools can benchmark their implementations against established tools using the same interface and evaluation criteria [15](#0-14) .

### Large-scale Analysis
The framework supports scanning entire datasets of contracts from CSV files or directories, making it suitable for blockchain-wide security studies [16](#0-15) [17](#0-16) .

## Notes

The project includes both active and deprecated tools, with some like EVM-CFG-Builder marked as deprecated but retained for reference [18](#0-17) . Each tool has specific capabilities - some focus on CFG generation while others detect vulnerabilities or extract function signatures [19](#0-18) .

## Documentation 📚

Comprehensive documentation is available in: https://deepwiki.com/zerjioang/evalevm/1-overview

## License

MIT
