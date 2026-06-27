---
title: Building a REST API in Go — The Right Way
date: 2026-04-02
excerpt: Skip the framework debates. Here's a production-ready pattern for Go HTTP services using only the standard library and two well-chosen dependencies.
---

# Building a REST API in Go — The Right Way

There's a lot of noise about which Go web framework to use. Gin vs Echo vs Chi vs stdlib. My take after building several production services: **start with `net/http`, reach for a router when routing becomes the problem**.

Here's the pattern I've landed on.

## Project layout

```
myapi/
├── main.go
├── handler/
│   ├── user.go
│   └── health.go
├── store/
│   └── postgres.go
└── model/
    └── user.go
```

Flat and obvious. No `pkg/`, no `internal/` maze until you need it.

## The handler pattern

```go
type UserHandler struct {
    store store.UserStore
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    user, err := h.store.GetByID(r.Context(), id)
    if err != nil {
        respondError(w, http.StatusNotFound, "user not found")
        return
    }
    respondJSON(w, http.StatusOK, user)
}
```

Dependencies are explicit (injected via struct fields), not hidden in globals. This makes testing trivial — swap the store for a mock and you're done.

## Error handling that doesn't lie

```go
func respondError(w http.ResponseWriter, code int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
```

Never send a `200 OK` with `{"error": "something broke"}` in the body. HTTP status codes exist for a reason.

## Middleware the stdlib way

```go
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

Chain them: `http.ListenAndServe(":8080", logging(router))`.

## What about Gin?

Gin is excellent and I use it for larger projects (binding validation alone saves hours). But for a simple CRUD API, the stdlib + Chi keeps your binary lean and your dependency tree clean.

## Closing thought

The best Go code reads like a spec. If you can't explain what each function does in one sentence, the function is doing too much. Go's verbosity isn't a bug — it's the language telling you to slow down and be explicit.

Next post: **connection pooling in PostgreSQL with `pgx`** — the settings that actually matter.