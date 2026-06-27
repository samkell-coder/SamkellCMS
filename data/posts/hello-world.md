---
title: Hello World — Why I Started Writing
date: 2024-03-15
excerpt: Every engineer I admire writes. Here's why I'm finally joining them — and what I plan to cover.
---

# Hello World — Why I Started Writing

Every engineer I admire writes. Not just commits and PRs, but actual prose — blog posts, essays, explanations of hard-won lessons. I've been meaning to start for years. Today I finally am.

## What you'll find here

This blog is where I'll document:

- **Deep dives** into Go internals and patterns I've found useful in production
- **Postmortems** — things I broke and what I learned
- **Tool breakdowns** — honest takes on the infrastructure I use daily
- **Side projects** — the chaotic, fun stuff that doesn't make it onto a résumé

## Why Go?

I fell in love with Go around 2019, coming from a Node.js background. What hooked me wasn't the syntax (it's famously boring, in the best way) — it was **predictability**. Go programs do what you expect. The garbage collector is well-behaved, goroutines are cheap, and the standard library is extraordinary.

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world.")
}
```

That's it. No config, no boilerplate, no bundler. A single binary, cross-compiled to any target. For backend systems, that clarity is a superpower.

## A promise

I'll publish at least twice a month. No fluff, no filler — just things I'd have wanted to read two years ago.

See you in the next one.
