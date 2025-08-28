# Makefile

# Define where Go binaries are installed. This keeps our tools project-local.
GOBIN ?= $(shell go env GOPATH)/bin
GOSEC := $(GOBIN)/gosec

# The default command to run if you just type 'make'.
.DEFAULT_GOAL := help

# Tells 'make' that these are command names, not files.
.PHONY: gosec-scan help clean install-dependencies

# This target downloads and tidies the Go module dependencies.
install-dependencies:
	@echo "--- Downloading Go module dependencies... ---"
	@go mod tidy
	@go mod download
	@echo "--- Dependencies are up to date. ---"

# This is the main target for your GitHub Actions workflow.
# It ensures gosec is installed first by depending on the $(GOSEC) target.
gosec-scan: $(GOSEC)
	@echo "--- Running GoSec security scan in $(shell pwd) ---"
	# Runs the scan on the current directory and all subdirectories.
	# The output is saved to a file for the workflow to use.
	# '|| true' is the key part that prevents the workflow from failing if issues are found.
	@$(GOSEC) ./... > gosec_results.txt || true
	@echo "--- GoSec scan complete. Results are in gosec_results.txt ---"

# This target checks if the gosec binary exists and installs it if it doesn't.
# This makes the setup automatic for both CI and local development.
$(GOSEC):
	@echo "--- GoSec not found, installing... ---"
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "--- GoSec installed successfully. ---"

# A simple target to remove generated files.
clean:
	@echo "--- Cleaning up generated files... ---"
	@rm -f gosec_results.txt

# A help command that explains what the Makefile can do.
help:
	@echo "Available commands:"
	@echo "  make install-dependencies - Downloads and tidies Go module dependencies."
	@echo "  make gosec-scan       - Installs gosec if needed and runs the security scan."
	@echo "  make clean            - Removes the generated gosec_results.txt report."
	@echo "  make help             - Shows this help message."
