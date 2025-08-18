# Go + HTMX Project Architecture

This document describes the standardized project structure for Go + HTMX applications. **LLMs working on this codebase MUST follow these conventions.**

## Directory Structure

```
project/
├── main.go                 # Entry point - minimal, delegates to routes
├── go.mod                  # Go module dependencies
├── config/
│   └── config.go          # Application configuration (env vars, settings)
├── handlers/              # HTTP handlers by feature/domain
│   ├── home.go           # Full page handlers
│   ├── debug.go          # Debug/admin pages
│   └── api.go            # HTMX fragment handlers
├── middleware/            # HTTP middleware
│   ├── auth.go           # Authentication middleware
│   └── logging.go        # Request logging
├── models/               # Data structures and types
│   └── page.go          # Page data models
├── services/             # Business logic services
│   └── click_service.go  # Domain logic, state management
├── templates/            # HTML templates
│   ├── layouts/
│   │   └── base.html     # Base layout with HTMX script
│   ├── pages/            # Full page templates
│   │   ├── home.html     # Page content (uses layouts)
│   │   └── debug.html    # Debug page
│   └── components/       # Reusable HTMX fragments
│       └── navbar.html   # Shared components
├── static/              # Static assets (CSS, JS, images)
│   └── css/
│       └── style.css    # Compiled CSS
└── routes/              # Route definitions and setup
    └── routes.go        # Central route configuration
```

## Core Principles for LLMs

### 1. **Separation of Concerns**
- **main.go**: Entry point only - configure logging, load config, setup routes, start server
- **handlers/**: HTTP request handling - parse requests, call services, render responses
- **services/**: Business logic - pure functions, state management, data processing
- **models/**: Data structures - types, structs, interfaces
- **templates/**: HTML rendering - layouts, pages, components

### 2. **HTMX Patterns**

#### Full Pages vs Fragments
```go
// Full page handler (handlers/home.go)
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, _ := template.ParseFiles(
        "templates/layouts/base.html",
        "templates/pages/home.html",
    )
    tmpl.Execute(w, data)
}

// Fragment handler (handlers/api.go)
func TimeFragmentHandler(w http.ResponseWriter, r *http.Request) {
    // Return just HTML fragment, no layout
    tmpl := `<div>{{.Time}}</div>`
    t, _ := template.New("time").Parse(tmpl)
    t.Execute(w, data)
}
```

#### Route Organization
```go
// routes/routes.go
func SetupRoutes() *mux.Router {
    // Full pages
    r.HandleFunc("/", handlers.HomePage)
    r.HandleFunc("/dashboard", handlers.DashboardPage)
    
    // HTMX fragments under /api prefix
    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/users", handlers.UsersFragment)
    api.HandleFunc("/notifications", handlers.NotificationsFragment)
}
```

### 3. **Template Structure**

#### Base Layout
```html
<!-- templates/layouts/base.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
    {{template "content" .}}
</body>
</html>
```

#### Page Templates
```html
<!-- templates/pages/home.html -->
{{define "content"}}
<div class="container">
    <h1>{{.Title}}</h1>
    <div hx-get="/api/data" hx-target="#content">
        <!-- HTMX interactions -->
    </div>
</div>
{{end}}
```

### 4. **Import Conventions**

Use module-based imports for local packages:
```go
import (
    "net/http"
    "github.com/gorilla/mux"        // External packages
    "hello-world/config"            // Local packages (use module name)
    "hello-world/handlers"
    "hello-world/services"
)
```

## LLM Instructions

### When Adding New Features:

1. **New Page**: 
   - Create handler in `handlers/[feature].go`
   - Create template in `templates/pages/[feature].html`
   - Add route in `routes/routes.go`

2. **New HTMX Fragment**:
   - Add handler function to `handlers/api.go`
   - Return HTML fragment only (no layout)
   - Use `/api/*` route prefix

3. **New Business Logic**:
   - Create service in `services/[domain]_service.go`
   - Keep handlers thin - delegate to services
   - Use dependency injection pattern

4. **Shared Components**:
   - Create in `templates/components/[name].html`
   - Use `{{template "component-name" .}}` to include

### File Modification Rules:

- **ALWAYS** use existing project structure
- **NEVER** put HTML directly in handler strings (use templates)
- **ALWAYS** separate full pages from HTMX fragments
- **NEVER** create new directories without justification
- **ALWAYS** follow the `/api/*` convention for HTMX endpoints

### Common Patterns:

```go
// Smart handler - detects HTMX vs full page
func SmartHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") == "true" {
        // Return fragment
        renderFragment(w, "components/table.html", data)
    } else {
        // Return full page
        renderPage(w, "pages/dashboard.html", data)
    }
}
```

## Testing

- Test handlers by mocking services
- Test services independently 
- Test templates by checking rendered output
- Use table-driven tests for multiple scenarios

This architecture scales from simple demos to production applications while maintaining clear separation of concerns and HTMX best practices.