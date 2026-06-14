# Meshbrow Go SDK

The official Go SDK for [Meshbrow](https://meshbrow.dev) — Managed Browser Fleet for AI Agents.

## Installation

```bash
go get github.com/meshbrow-dev/meshbrow-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/meshbrow-dev/meshbrow-go"
)

func main() {
    client := meshbrow.NewClient("your-api-key")
    ctx := context.Background()

    // Launch a stealth browser session
    session, err := client.CreateSession(ctx, &meshbrow.CreateSessionParams{
        Stealth:      "max",
        ProxyType:    "residential",
        ProxyCountry: "US",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.DestroySession(ctx, session.ID, false)

    // Navigate and extract
    client.Navigate(ctx, session.ID, "https://example.com", "load")
    content, _ := client.Extract(ctx, session.ID, "", 5000)
    fmt.Println(content.Text)

    // Screenshot
    screenshot, _ := client.TakeScreenshot(ctx, session.ID, "", false)
    fmt.Printf("Screenshot: %d bytes (base64)\n", len(screenshot.Data))
}
```

## Fleet Operations

```go
fleet, err := client.CreateFleet(ctx, &meshbrow.CreateFleetParams{
    Count:        10,
    ProxyType:    "residential",
    ProxyCountry: "US",
})
if err != nil {
    log.Fatal(err)
}
defer client.DestroyFleet(ctx, fleet.ID)

for _, session := range fleet.Sessions {
    client.Navigate(ctx, session.ID, "https://example.com", "load")
}
```

## License

MIT
