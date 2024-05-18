# Heracles: a Prometheus Middleware for go-chi

_**Prometheus was freed by the hero Heracles.**_

**Heracles** provides a Prometheus middleware for the [go-chi](https://github.com/go-chi/chi) router. It enables the collection of HTTP request metrics such as request counts, latencies, and detailed error metrics. The middleware is highly configurable using functional options.

## Features

- Request Count Metrics: Track the number of HTTP requests partitioned by status code, method, and path.

- Latency Metrics: Measure the latency of HTTP requests.

- Detailed Error Metrics: Track detailed client and server error metrics.

- Custom Labels: Add custom labels to metrics.

- Configurable Buckets: Configure latency histogram buckets.

## Installation

```bash
go get githug.com/hyp3rd/heracles
```

## Usage

```go
package main

import (
    "net/http"
    "time"

    "github.com/go-chi/chi"
    "github.com/hyp3rd/heracles/v1"
)

func main() {
    r := chi.NewRouter()

    // Create a new Prometheus middleware instance
    promMiddleware, err := heracles.NewMiddleware("my_service", heracles.WithRequestsEnabled(true), heracles.WithLatencyEnabled(true))
    if err != nil {
        panic(err)
    }

    // Register the middleware
    r.Use(promMiddleware.Handler)

    // Define your routes
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    http.ListenAndServe(":8080", r)
}
```

## Configuration

The Prometheus middleware can be configured using functional options. The following options are available:

- `WithRequestsEnabled`: Enable request count metrics. Default: `true`.

- `WithLatencyEnabled`: Enable latency metrics. Default: `true`.

- `WithCustomLabels`: Add custom labels to metrics.

- `WithLatencyBuckets`: Configure latency histogram buckets.

### Custom Configuration Usage

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/hyp3rd/heracles/v1"
)

func main() {
    r := chi.NewRouter()

    // Create a new Prometheus middleware instance with custom options
    promMiddleware, err := heracles.NewMiddleware(
        "my_service",
        heracles.WithRequestsEnabled(true),
        heracles.WithLatencyEnabled(true),
        heracles.WithCustomLabels("X-Request-ID"),
        heracles.WithLatencyBuckets([]float64{0.1, 0.5, 1.0, 2.5, 5.0}),
    )
    if err != nil {
        panic(err)
    }

    // Register the middleware
    r.Use(promMiddleware.Handler)

    // Define your routes
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    http.ListenAndServe(":8080", r)
}
```

## Metrics

The middleware exposes the following metrics:

- **chi_requests_total**: The total number of HTTP requests partitioned by status code, method, and path.
- **chi_request_duration_seconds**: The duration of HTTP requests in seconds partitioned by status code, method, and path.
- **chi_detailed_errors_total**: Detailed error counts partitioned by type (client_error, server_error), status code, method, and path.

## Registering Collectors

To register the collectors with the default Prometheus registerer:

```go
promMiddleware.MustRegisterDefault()
```

If you are using a custom Prometheus registerer, you can retrieve the collectors and register them manually:

```go
collectors := promMiddleware.Collectors()
for _, collector := range collectors {
    prometheus.MustRegister(collector)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome. Please open an issue or submit a pull request for any bugs, improvements, or feature requests.

## Author

I'm a surfer, a trader, and a software architect with 15 years of experience designing highly available distributed production environments and developing cloud-native apps in public and private clouds. Just your average bloke. Feel free to connect with me on LinkedIn, but no funny business.

[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/francesco-cosentino/)
