# folderA/Makefile

# Define where Go binaries are installed to avoid polluting the system path.
# This ensures gosec is installed in a predictable location within the project's context.
GOBIN ?= $(shell go env GOPATH)/bin
GOSEC := $(GOBIN)/gosec

# Default target for when 'make' is run without arguments.
.DEFAULT_GOAL := help

# Phony targets are rules that don't represent actual files.
.PHONY: gosec-scan help

# This is the primary target for our CI workflow.
# It depends on the $(GOSEC) target, which ensures gosec is installed before running.
gosec-scan: $(GOSEC)
	@echo "--- Running GoSec security scan in $(shell pwd) ---"
	# Run gosec on all packages within the current directory.
	# Output is redirected to a file for the workflow to use later.
	# '|| true' is critical: it makes this step always succeed, even if vulnerabilities are found.
	@$(GOSEC) ./... > gosec_results.txt || true
	@echo "--- GoSec scan complete. Results are in gosec_results.txt ---"

# This target acts as a dependency check. 'make' will only run the command
# if the file specified by $(GOSEC) does not exist.
$(GOSEC):
	@echo "--- GoSec not found, installing... ---"
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "--- GoSec installed successfully. ---"

# A simple help command to explain available targets.
help:
	@echo "Available targets:"
	@echo "  gosec-scan  - Run the GoSec security scan and generate a report."
	@echo "  help        - Show this help message."
