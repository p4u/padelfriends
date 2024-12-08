# Makefile for Padel Friends Project

.PHONY: all backend frontend frontend_reload clean run check_npm

# Check if npm is available
HAVE_NPM := $(shell command -v npm >/dev/null 2>&1 && echo yes || echo no)

ifeq ($(HAVE_NPM),no)
NPM_SETUP = curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash && \
	export NVM_DIR="$(HOME)/.nvm" && [ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh" && nvm install node
RUN_NPM = export NVM_DIR="$(HOME)/.nvm"; [ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh" && npm
else
NPM_SETUP = echo "npm is already installed ($(shell which npm)), skipping nvm installation."
RUN_NPM = npm
endif

# Default target: build backend and then frontend if UI is present
all: backend
	@if [ -d ui ]; then \
		$(MAKE) frontend; \
	else \
		echo "UI directory not found, skipping frontend build"; \
	fi

# Build the Go backend
backend:
	@echo "Building Go backend..."
	go build -o backend .
	@echo "Go backend built successfully."

# Ensure npm is available (install if not)
check_npm:
	@$(NPM_SETUP)

# Build the frontend using Vite and Vue
frontend: check_npm
	@echo "Building frontend with Vite and Vue..."
	cd ui && $(RUN_NPM) install && $(RUN_NPM) run build
	@echo "Frontend built successfully."

frontend_reload: check_npm
	@echo "Building frontend with Vite and Vue (reload)..."
	cd ui && $(RUN_NPM) run build
	@echo "Frontend built successfully."

# Clean build artifacts
clean:
	@echo "Cleaning Go backend..."
	go clean
	rm -f backend
	@echo "Go backend cleaned."

	@if [ -d ui ]; then \
		echo "Cleaning frontend..."; \
		cd ui && $(RUN_NPM) run clean; \
		rm -rf ui/node_modules ui/dist; \
		echo "Frontend cleaned."; \
	else \
		echo "UI directory not found, skipping frontend clean"; \
	fi

# Optional: Run the backend
run: backend frontend_reload
	@echo "Running Go backend..."
	./backend

execute:
	@echo "Running Go backend..."
	./backend
