# Tracker Package - Core Components Guide

## Package Overview

The `tracker` package contains the core Snowplow event tracking implementation. It handles event creation, validation, enrichment, and transport to Snowplow collectors with automatic batching and retry mechanisms.

## Component Architecture

```
Tracker (orchestrator)
    ├── Emitter (transport layer)
    │   ├── Storage (persistence)
    │   ├── HTTP Client (network)
    │   └── Send Channel (concurrency)
    ├── Subject (user context)
    └── Events (type definitions)
        ├── PageViewEvent
        ├── StructuredEvent
        ├── SelfDescribingEvent
        ├── ScreenViewEvent
        ├── TimingEvent
        └── EcommTransactionEvent
```

## Core Component Patterns

### Tracker Initialization
```go
// ✅ Correct: Builder pattern with validation
t := InitTracker(
    RequireEmitter(emitter),    // Required
    OptionNamespace("ns"),       // Optional
)

// ❌ Wrong: Direct struct creation
t := &Tracker{Emitter: emitter}
```

### Emitter Configuration
```go
// ✅ Correct: Functional options with defaults
e := InitEmitter(
    RequireCollectorUri("collector.com"),
    RequireStorage(storage),
    OptionRequestType("POST"),  // GET or POST
    OptionSendLimit(500),        // Batch size
)

// ❌ Wrong: Missing required parameters
e := InitEmitter()  // Will panic
```

### Event Creation Pattern
```go
// ✅ Correct: Struct initialization then Init()
event := PageViewEvent{
    PageUrl: common.NewString("url"),
}
event.Init()  // Validates and sets defaults

// ❌ Wrong: Tracking without initialization
tracker.TrackPageView(PageViewEvent{})
```

## Event Type Patterns

### Required vs Optional Fields
- **Required**: Validated in `Init()`, panic if missing
- **Optional**: Use nil to omit from payload
- **Auto-generated**: Timestamp, EventId if nil

### Event Enrichment Flow
1. Event struct created with user data
2. `Init()` validates and adds defaults
3. `SetSubjectIfNil()` adds tracker subject
4. `Get()` builds final payload
5. Emitter adds protocol fields

## Emitter Transport Patterns

### Batch Processing
```go
// GET: Individual requests up to ByteLimitGet
// POST: Batched requests up to SendLimit items

// ✅ Correct: Configure for high throughput
OptionSendLimit(100),        // Batch 100 events
OptionByteLimitPost(50000),  // 50KB POST limit

// ❌ Wrong: Limits too low for production
OptionSendLimit(1),  // Inefficient
```

### Retry Mechanism
```go
// Automatic retry with exponential backoff
// Failed events returned to storage
// Callback notified of success/failure

// ✅ Correct: Set callback for monitoring
OptionCallback(func(success, failure []CallbackResult) {
    log.Printf("Sent: %d, Failed: %d", len(success), len(failure))
})
```

### Storage Integration
```go
// Events persisted before sending
// Retrieved in batches for transmission
// Deleted only after successful send

// ✅ Correct: Choose appropriate storage
RequireStorage(*memory.Init())     // Fast, volatile
RequireStorage(*sqlite3.Init(db))  // Persistent

// ❌ Wrong: No storage
RequireStorage(nil)  // Will panic
```

## Subject Management Patterns

### Subject Hierarchy
```go
// Priority: Event > Tracker > None
event.Subject = eventSubject       // Highest priority
tracker.Subject = trackerSubject   // Default
// No subject = anonymous tracking

// ✅ Correct: Event-specific subject
event := PageViewEvent{
    Subject: InitSubject(),  // Override tracker subject
}
```

### ID Cookie Parsing
```go
// ✅ Correct: Parse Snowplow ID cookie
subject.FromIdCookie("duid.sessionId.sessionCount.ts")

// Sets: DomainUserId, SessionId, SessionIndex
```

## Self-Describing JSON Pattern

### Custom Event Creation
```go
// ✅ Correct: Iglu schema with data
sdj := InitSelfDescribingJson(
    "iglu:com.acme/event/jsonschema/1-0-0",
    map[string]interface{}{"key": "value"},
)

// ❌ Wrong: Invalid schema format
sdj := InitSelfDescribingJson("schema", data)
```

### Context Addition
```go
// ✅ Correct: Multiple contexts as array
event.Contexts = []SelfDescribingJson{
    *InitSelfDescribingJson(schema1, data1),
    *InitSelfDescribingJson(schema2, data2),
}

// ❌ Wrong: Single context not in array
event.Contexts = InitSelfDescribingJson(s, d)
```

## Constants & Protocol

### Event Types
- `EVENT_PAGE_VIEW = "pv"`
- `EVENT_STRUCTURED = "se"`
- `EVENT_UNSTRUCTURED = "ue"`

### Protocol Versions
- `POST_PROTOCOL_VERSION = "tp2"`
- `TRACKER_VERSION = "golang-3.3.0"`

### Schema Versions
- Payload: `1-0-4`
- Contexts: `1-0-1`
- Unstruct: `1-0-0`

## Testing Patterns

### HTTP Mocking
```go
// ✅ Correct: Use httpmock for tests
httpmock.Activate()
defer httpmock.DeactivateAndReset()
httpmock.RegisterResponder("POST", url, responder)

// ❌ Wrong: Real HTTP in tests
http.Post(url, "application/json", body)
```

### Event Verification
```go
// ✅ Correct: Assert payload contents
assert.Equal("pv", payload.Get()["e"])
assert.NotNil(payload.Get()["eid"])
assert.NotNil(payload.Get()["dtm"])

// ❌ Wrong: Only check success
assert.True(tracker.TrackPageView(event))
```

## Error Handling Patterns

### Validation Panics
```go
// ✅ Correct: Validate before tracking
if url == "" {
    return errors.New("URL required")
}
event := PageViewEvent{PageUrl: &url}

// ❌ Wrong: Let tracker panic
event := PageViewEvent{}  // Missing required
tracker.TrackPageView(event)  // Panic!
```

### Network Failures
```go
// ✅ Correct: Events persisted, will retry
// Storage ensures no data loss
// Callback monitors failures

// ❌ Wrong: Ignore send failures
tracker.TrackPageView(event)
// Assume success without verification
```

## Performance Considerations

### Batching Strategy
- **GET**: No batching, immediate send
- **POST**: Batch up to SendLimit
- Trade-off: Latency vs Throughput

### Storage Choice
- **Memory**: Fast, good for high volume
- **SQLite**: Persistent, survives restarts
- Consider: Data criticality vs Performance

### Concurrent Sending
- Emitter uses goroutines for sending
- SendChannel controls concurrency
- Storage handles thread safety

## Common Integration Patterns

### Web Server Integration
```go
// ✅ Correct: Shared tracker, per-request events
var tracker *Tracker  // Initialize once

func handler(w http.ResponseWriter, r *http.Request) {
    event := PageViewEvent{
        PageUrl: common.NewString(r.URL.String()),
    }
    tracker.TrackPageView(event)
}
```

### Graceful Shutdown
```go
// ✅ Correct: Ensure all events sent
signal.Notify(stop, os.Interrupt)
<-stop
tracker.Emitter.Stop()
tracker.BlockingFlush(10, 100)
```

### Error Recovery
```go
// ✅ Correct: Recover from panics
defer func() {
    if r := recover(); r != nil {
        log.Printf("Tracking error: %v", r)
    }
}()
tracker.TrackPageView(event)
```

## Quick Reference

### Component Initialization Order
1. Storage → 2. Emitter → 3. Subject → 4. Tracker

### Event Tracking Flow
1. Create event struct → 2. Set fields → 3. Track method → 4. Validation → 5. Enrichment → 6. Storage → 7. Batch → 8. Send

### Required Parameters
- **Tracker**: Emitter
- **Emitter**: CollectorUri, Storage
- **Events**: Vary by type (check struct tags)