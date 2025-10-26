# Go Code Standards and Best Practices

## 🎯 Quality Goal: A+ Go Report Card

This project is committed to **idiomatic Go practices** and achieving an **A+ grade on [Go Report Card](https://goreportcard.com/)**.

---

## 📊 Go Report Card Criteria

We target 100% compliance with all Go Report Card checks:

### 1. gofmt ✅
- **Requirement**: All code must be formatted with `gofmt`
- **How**: Run `make fmt` before every commit
- **CI**: Enforced in CI/CD pipeline
- **Target**: 100% compliance

### 2. go vet ✅
- **Requirement**: No issues reported by `go vet`
- **How**: Run `make vet` before every commit
- **CI**: Enforced in CI/CD pipeline
- **Target**: 0 issues

### 3. golint / staticcheck ✅
- **Requirement**: No lint warnings (using golangci-lint)
- **How**: Run `make lint` before every commit
- **CI**: Enforced with golangci-lint in CI/CD
- **Target**: 0 warnings

### 4. gocyclo ✅
- **Requirement**: Cyclomatic complexity < 15 per function
- **How**: Keep functions small and focused
- **Refactor**: Break complex functions into smaller ones
- **Target**: < 10 average complexity

### 5. ineffassign ✅
- **Requirement**: No ineffectual assignments
- **How**: Remove unused variable assignments
- **CI**: Checked by golangci-lint
- **Target**: 0 issues

### 6. misspell ✅
- **Requirement**: No spelling errors in comments/strings
- **How**: Use proper English spelling
- **CI**: Checked by golangci-lint
- **Target**: 0 issues

---

## 🏛️ Idiomatic Go Principles

### Package Design

**✅ DO**:
```go
// Package gcs provides a client for the Globus Connect Server Manager API.
//
// The GCS Manager API runs on individual GCS v5 endpoints (not centralized)
// and is accessed via HTTPS at the endpoint's FQDN.
package gcs

// Clear, focused package with single responsibility
```

**❌ DON'T**:
```go
package utils  // Too generic
package helpers // Too vague
package misc   // Anti-pattern
```

### Naming Conventions

**✅ DO**:
```go
// Use MixedCaps, not underscores
type EndpointInfo struct { }
func GetEndpoint() { }

// Abbreviations should be consistent
var apiClient *APIClient  // API, not Api
var jsonData []byte       // JSON, not Json
var urlPath string        // URL, not Url
var gcsClient *GCSClient  // GCS, not Gcs

// Interface names: one-method interfaces end in -er
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**❌ DON'T**:
```go
type endpoint_info struct { }  // No underscores
func get_endpoint() { }        // No underscores
var ApiClient *APIClient       // Wrong caps
```

### Error Handling

**✅ DO**:
```go
// Return errors, don't panic
func GetEndpoint(id string) (*Endpoint, error) {
    if id == "" {
        return nil, fmt.Errorf("endpoint ID cannot be empty")
    }

    endpoint, err := client.fetchEndpoint(id)
    if err != nil {
        return nil, fmt.Errorf("fetch endpoint: %w", err)
    }

    return endpoint, nil
}

// Use error wrapping with %w for context
// Use error types for domain errors
type ValidationError struct {
    Field string
    Value interface{}
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %v", e.Field, e.Value)
}
```

**❌ DON'T**:
```go
func GetEndpoint(id string) *Endpoint {
    if id == "" {
        panic("empty id")  // Don't panic in library code
    }

    endpoint, err := client.fetchEndpoint(id)
    if err != nil {
        log.Fatal(err)  // Don't call os.Exit or log.Fatal in libraries
    }

    return endpoint
}
```

### Context Usage

**✅ DO**:
```go
// Accept context.Context as first parameter
func (c *Client) GetEndpoint(ctx context.Context, id string) (*Endpoint, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", c.buildURL("/endpoint/"+id), nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.httpClient.Do(req)
    // ...
}

// Use context for cancellation and timeouts
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

endpoint, err := client.GetEndpoint(ctx, "abc123")
```

**❌ DON'T**:
```go
// Don't ignore context or put it last
func (c *Client) GetEndpoint(id string) (*Endpoint, error) {
    // No context support
}

func (c *Client) GetEndpoint(id string, ctx context.Context) {
    // Wrong parameter order
}
```

### Struct Design

**✅ DO**:
```go
// Use struct embedding for composition
type BaseClient struct {
    httpClient *http.Client
    baseURL    string
}

type GCSClient struct {
    BaseClient  // Embedded
    token      *oauth2.Token
}

// Constructor pattern with options
type ClientOption func(*Client)

func WithTimeout(d time.Duration) ClientOption {
    return func(c *Client) {
        c.httpClient.Timeout = d
    }
}

func NewClient(opts ...ClientOption) *Client {
    c := &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }

    for _, opt := range opts {
        opt(c)
    }

    return c
}

// Usage
client := gcs.NewClient(
    gcs.WithTimeout(60*time.Second),
)
```

### Interfaces

**✅ DO**:
```go
// Accept interfaces, return structs
func ProcessData(r io.Reader) (*Result, error) {
    // Flexible - accepts any Reader
}

// Small interfaces (1-3 methods)
type Authenticator interface {
    Authenticate(ctx context.Context) (*Token, error)
}

// Define interfaces at usage point, not implementation
// (consumer-driven interfaces)
```

**❌ DON'T**:
```go
// Don't create large interfaces
type Manager interface {
    Create(...) error
    Read(...) error
    Update(...) error
    Delete(...) error
    List(...) error
    Search(...) error
    Export(...) error
    Import(...) error
    // Too many methods!
}
```

### Concurrency

**✅ DO**:
```go
// Use channels for communication
func worker(jobs <-chan Job, results chan<- Result) {
    for j := range jobs {
        results <- process(j)
    }
}

// Use sync.WaitGroup for coordination
var wg sync.WaitGroup
for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        worker(jobs, results)
    }()
}
wg.Wait()

// Use sync.Once for initialization
var instance *Client
var once sync.Once

func GetInstance() *Client {
    once.Do(func() {
        instance = newClient()
    })
    return instance
}
```

### Testing

**✅ DO**:
```go
// Table-driven tests
func TestGetEndpoint(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *Endpoint
        wantErr bool
    }{
        {
            name:  "valid endpoint",
            input: "abc123",
            want:  &Endpoint{ID: "abc123"},
        },
        {
            name:    "empty ID",
            input:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := GetEndpoint(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetEndpoint() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetEndpoint() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Use testify for assertions (optional but recommended)
import "github.com/stretchr/testify/assert"

func TestEndpointCreate(t *testing.T) {
    endpoint, err := CreateEndpoint("test")
    assert.NoError(t, err)
    assert.NotNil(t, endpoint)
    assert.Equal(t, "test", endpoint.Name)
}
```

### Documentation

**✅ DO**:
```go
// Package documentation at package clause
// Package gcs provides a client for the Globus Connect Server Manager API.
package gcs

// Exported functions need doc comments starting with function name
// GetEndpoint retrieves endpoint information by ID.
//
// The endpoint ID must be a valid UUID. Returns an error if the endpoint
// does not exist or the API request fails.
func GetEndpoint(ctx context.Context, id string) (*Endpoint, error) {
    // ...
}

// Exported types need doc comments
// Endpoint represents a Globus Connect Server v5 endpoint.
type Endpoint struct {
    // ID is the unique identifier for this endpoint.
    ID string `json:"id"`

    // DisplayName is the human-readable name shown in Globus Web App.
    DisplayName string `json:"display_name"`
}
```

**❌ DON'T**:
```go
// retrieves endpoint (lowercase start)
func GetEndpoint(id string) {
}

// No doc comment at all
type Endpoint struct {
}
```

---

## 🔧 Development Workflow

### Before Every Commit

```bash
# 1. Format code
make fmt

# 2. Run vet
make vet

# 3. Run linter
make lint

# 4. Run tests
make test

# Or run all checks at once:
make verify
```

### golangci-lint Configuration

We use a comprehensive `.golangci.yml` configuration:

```yaml
linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - misspell
    - gocyclo
    - gocritic
    - gocognit
    - godot
    - gofumpt
    - revive

linters-settings:
  gocyclo:
    min-complexity: 15
  gocognit:
    min-complexity: 20
```

---

## 📚 Reference Resources

### Official Go Documentation
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)

### Style Guides
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)

### Go Report Card
- [Go Report Card](https://goreportcard.com/)
- [goreportcard on GitHub](https://github.com/gojp/goreportcard)

---

## ✅ Pre-Release Checklist

Before each release, verify:

- [ ] Go Report Card grade: **A+**
- [ ] Test coverage: **> 80%**
- [ ] `make verify` passes with 0 issues
- [ ] All exported symbols documented
- [ ] Examples provided for main use cases
- [ ] CHANGELOG.md updated
- [ ] Version bumped in appropriate files

---

## 🎓 Learning from globus-go-sdk

We follow the patterns established in `globus-go-sdk`:

```go
// Service-based organization
type Service struct {
    client *Client
}

// Context-first parameters
func (s *Service) Get(ctx context.Context, id string) (*Resource, error)

// Error wrapping with context
return nil, fmt.Errorf("get resource %s: %w", id, err)

// Options pattern for complex parameters
type GetOptions struct {
    Include []string
    Fields  []string
}

func (s *Service) Get(ctx context.Context, id string, opts *GetOptions) (*Resource, error)
```

---

**Last Updated**: October 26, 2025
**Review Cadence**: Every phase completion
**Owner**: Project Lead
