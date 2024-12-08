# Makefile for Padel Friends Project

.PHONY: all backend frontend clean

# Default target: build backend and frontend if ui directory exists
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

# Build the frontend using Vite and Vue
frontend:
	@echo "Building frontend with Vite and Vue..."
	cd ui && npm install && npm run build
	@echo "Frontend built successfully."

# Build the frontend using Vite and Vue
frontend_reload:
	@echo "Building frontend with Vite and Vue..."
	cd ui && npm run build
	@echo "Frontend built successfully."


# Clean build artifacts
clean:
	@echo "Cleaning Go backend..."
	go clean
	rm -f backend
	@echo "Go backend cleaned."

	@if [ -d ui ]; then \
		echo "Cleaning frontend..."; \
		cd ui && npm run clean; \
		rm -rf ui/node_modules ui/dist; \
		echo "Frontend cleaned."; \
	else \
		echo "UI directory not found, skipping frontend clean"; \
	fi

# Optional: Run the backend (add more commands as needed)
run: backend frontend_reload
	@echo "Running Go backend..."
	./backend
