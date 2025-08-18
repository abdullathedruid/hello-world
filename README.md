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

## For Developers

See `ARCHITECTURE.md` for detailed architectural guidelines and LLM instructions.