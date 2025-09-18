# Project settings
BINARY := evalevm
PKG := ./...

# Go build flags
GOFLAGS := -trimpath
LDFLAGS := -s -w

# Default target
all: build

# Build optimized binary
build: clean
	@echo "Building optimized binary..."
	mkdir -p dist
	GOFLAGS="$(GOFLAGS)" go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf dist
	rm -rf debug_*
	rm -rf result_*
	rm -rf cfg_*

# Run the app (after build)
run: build
	@echo "Running $(BINARY)..."
	./dist/$(BINARY)

# Run tests
test:
	@echo "Running tests..."
	go test -v $(PKG)

setup:
	chmod +rwx ./docker/measure.sh
	chmod +rwx ./docker/helper.sh

	cp docker/measure.sh tools/byte-inspector/measure.sh
	cp docker/helper.sh tools/byte-inspector/helper.sh

	cp docker/measure.sh tools/conkas/measure.sh
	cp docker/helper.sh tools/conkas/helper.sh

	cp docker/measure.sh tools/ethersolve_creator/measure.sh
	cp docker/helper.sh tools/ethersolve_creator/helper.sh

	cp docker/measure.sh tools/ethersolve_runtime/measure.sh
	cp docker/helper.sh tools/ethersolve_runtime/helper.sh

	cp docker/measure.sh tools/evm_cfg_builder/measure.sh
	cp docker/helper.sh tools/evm_cfg_builder/helper.sh

	cp docker/measure.sh tools/evm-lisa/measure.sh
	cp docker/helper.sh tools/evm-lisa/helper.sh

	cp docker/measure.sh tools/honeybadger/measure.sh
	cp docker/helper.sh tools/honeybadger/helper.sh

	cp docker/measure.sh tools/paper/measure.sh
	cp docker/helper.sh tools/paper/helper.sh

	cp docker/measure.sh tools/rattle/measure.sh
	cp docker/helper.sh tools/rattle/helper.sh

	cp docker/measure.sh tools/securify/measure.sh
	cp docker/helper.sh tools/securify/helper.sh

	cp docker/measure.sh tools/evmole/measure.sh
	cp docker/helper.sh tools/evmole/helper.sh

	cp docker/measure.sh tools/evm-cfg/measure.sh
	cp docker/helper.sh tools/evm-cfg/helper.sh

	cp docker/measure.sh tools/vandal/measure.sh
	cp docker/helper.sh tools/vandal/helper.sh

generate-images: build setup
	./dist/$(BINARY) analyzer build --tools ./tools

scan-csv: build
	./dist/$(BINARY) scan csv --dataset /Users/sergio/Desktop/phd/dataset_eth_mainnet/unique_eth_contracts_order_by_size_2015_pruning.csv

scan-minimal: build
	./dist/$(BINARY) scan evm --bytecode 0x606060405260e060020a60003504631f50eadc8114601a575b005b6004356060908152602080822060243514825290f3

scan-sample: build
	 ./dist/$(BINARY) scan evm --bytecode 0x60606040526000543411601157600080fd5b60015473ffffffffffffffffffffffffffffffffffffffff1663d7bb99ba6247b760346040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381858988f15050505050500000a165627a7a72305820d43c494964b14aa56c5eee3c28da8f5049cd2c382d21ee39a116116b5c1253db0029

scan-dataset: build
	# /Users/sergio/Desktop/phd/github/EtherSolve_ICPC2021_ReplicationPackage/Benchmark/Bytecode-dataset-1000-contracts
	 ./dist/$(BINARY) scan dir --dataset /Users/sergio/Desktop/phd/github/smartbugs/samples/0.4.x

scan-dataset-1000: build
	 ./dist/$(BINARY) scan dir --dataset /Users/sergio/Desktop/phd/github/EtherSolve_ICPC2021_ReplicationPackage/Benchmark/Bytecode-dataset-1000-contracts
