This article is built independently from the main Infer Origins website.

## Markdown rendering

The blog builder uses the same Goldmark Markdown engine as the main portal.

### Code example

```go
package main

import "fmt"

func main() {
    fmt.Println("Infer Origins")
}
```

### Shell example

```bash
go run ./cmd/blogbuild
```

## Why separate the blog?

Publishing a new article should not require rebuilding products, services, case studies, or the rest of the static website.
