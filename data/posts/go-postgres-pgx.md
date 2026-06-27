---
title: PostgreSQL Connection Pooling in Go with pgx
date: 2024-05-10
excerpt: The five pgxpool settings that actually matter in production — and why the defaults will hurt you at scale.
---

# PostgreSQL Connection Pooling in Go with `pgx`

If you're building a Go service that talks to PostgreSQL, `pgx` is the library to reach for. But most tutorials show the minimum viable setup. In production, the defaults will quietly bite you. Here's what I've learned.

## The naive setup

```go
conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
```

This creates a **single connection**. Fine for a toy project; catastrophic under any real load. Each goroutine that needs the database blocks every other one.

## Use pgxpool

```go
import "github.com/jackc/pgx/v5/pgxpool"

pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
```

`pgxpool` manages a connection pool automatically. But the defaults are conservative. Here's what to tune.

## The five settings that matter

```go
config, _ := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))

// 1. Max connections — default is 4. Too low for any concurrent service.
config.MaxConns = 20

// 2. Min idle connections — keeps warm connections ready.
config.MinConns = 5

// 3. Max connection lifetime — recycle connections before the DB closes them.
config.MaxConnLifetime = 30 * time.Minute

// 4. Max idle time — close connections that have been idle too long.
config.MaxConnIdleTime = 5 * time.Minute

// 5. Connect timeout — fail fast rather than queue forever.
config.ConnConfig.ConnectTimeout = 5 * time.Second

pool, err := pgxpool.NewWithConfig(context.Background(), config)
```

## Always pass context

```go
// Good — respects cancellation / timeouts
row := pool.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id)

// Bad — can hang indefinitely
row := pool.QueryRow(context.Background(), ...)
```

Wire your HTTP request context all the way down. A slow query that the client already gave up on shouldn't keep holding a connection.

## Check your pool health

```go
stats := pool.Stat()
log.Printf(
    "pool: total=%d idle=%d acquired=%d",
    stats.TotalConns(),
    stats.IdleConns(),
    stats.AcquiredConns(),
)
```

Expose this on a `/metrics` or `/health` endpoint and alert when `AcquiredConns` stays near `MaxConns` — that's a sign you need to raise the limit or optimise your queries.

## Rule of thumb for MaxConns

A good starting point: `MaxConns = num_cpu_cores * 2 + effective_spindle_count`. For a modern cloud database with SSDs, that's roughly `num_cores * 3`. Start there, load test, and adjust.

Next up: **structured logging in Go with `slog`** — the stdlib package that finally makes JSON logs ergonomic.
