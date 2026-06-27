---
title: Building Fast Websites with Go + Vanilla JavaScript
date: 2026-06-26
excerpt: Everyone reaches for React, Vue, or Next.js by default. My experience after shipping several production applications? You don't always need a frontend framework.
---

# Building Fast Websites with Go + Vanilla JavaScript

Everyone reaches for React, Vue, or Next.js by default. My experience after shipping several production applications? You don't always need a frontend framework.

For portfolios, dashboards, admin panels, SaaS products, and countless business websites, Go on the backend with Vanilla JavaScript on the frontend gets you surprisingly far. Less tooling, fewer dependencies, and a codebase that's easy to understand months later.

Here's the setup I've settled on.

## Project layout

```text
myapp/
├── main.go
├── handler/
│   ├── api.go
│   └── page.go
├── store/
│   └── postgres.go
├── templates/
│   └── index.html
└── static/
    ├── css/
    ├── js/
    └── images/
```

Simple, predictable, and easy to navigate. Every folder has one responsibility, so finding code never feels like a scavenger hunt.

## The frontend pattern

```javascript
async function loadProjects() {
    const res = await fetch("/api/projects");
    const projects = await res.json()

    renderProjects(projects)
}
```

The backend owns the data. JavaScript owns the interactions. Keeping those responsibilities separate makes features easier to build, debug, and extend without introducing unnecessary complexity.

## Let Go do the heavy lifting

```go
fs := http.FileServer(http.Dir("./static"))

http.Handle("/static/",
    http.StripPrefix("/static/", fs),
)

http.HandleFunc("/", homeHandler)
```

Go doesn't just power your API—it can serve your HTML, CSS, JavaScript, and assets too. One binary, one deployment, and no extra runtime just to deliver static files.

## APIs that stay consistent

```go
func respondJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

Every endpoint should respond the same way. Consistent JSON responses make your frontend smaller because every request follows the same predictable pattern.

## What about React?

React is an excellent tool, and I use it when applications become highly interactive or the frontend grows into its own product. But for many websites and CRUD applications, Go with Vanilla JavaScript is faster to build, easier to deploy, and much simpler to maintain.

## Closing thought

Modern web development doesn't have to mean dozens of build tools and thousands of dependencies. A Go backend paired with Vanilla JavaScript delivers fast websites, clean architecture, and a developer experience that stays enjoyable as the project grows. Sometimes the simplest stack is also the most productive.