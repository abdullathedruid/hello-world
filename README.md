# Go + HTMX Structured Project

A well-organized Go + HTMX application following scalable architecture patterns.

## Quick Start

```bash
go build
./hello-world
# Visit http://localhost:8080
```

## Structure

- **`main.go`** - Entry point
- **`handlers/`** - HTTP request handlers
- **`templates/`** - HTML templates (layouts, pages, components)
- **`services/`** - Business logic
- **`routes/`** - Route configuration
- **`config/`** - Application configuration
- **`middleware/`** - HTTP middleware
- **`models/`** - Data structures
- **`static/`** - CSS, JS, images

## Key Features

- ✅ HTMX integration for dynamic interactions
- ✅ Template-based architecture (layouts + pages)
- ✅ Separation of full pages vs HTMX fragments
- ✅ Clean API routes under `/api/*`
- ✅ Farcaster MiniApp SDK integration
- ✅ Structured logging with Logrus

## Development Workflow

### Testing & Coverage
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench
```

### Code Quality
```bash
# Format code (REQUIRED before creating PRs)
make fmt

# Check for issues
make vet

# Run linter if available
make lint

# Run all checks
make check
```

### Before Creating a Pull Request
**ALWAYS run these commands before creating a PR:**
```bash
make fmt        # Format code
make vet        # Check for issues  
make test       # Ensure tests pass
```

## For Developers

See `ARCHITECTURE.md` for detailed architectural guidelines and LLM instructions.