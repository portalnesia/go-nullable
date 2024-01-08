[![Go Reference](https://pkg.go.dev/badge/go.portalnesia.com/nullable.svg)](https://pkg.go.dev/go.portalnesia.com/nullable) ![Go](https://github.com/portalnesia/go-nullable/actions/workflows/nullable_test.yml/badge.svg)

# Nullable

Nullable Data Type for json and database

## Install

```bash
go get go.portalnesia.com/nullable
```

## Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"go.portalnesia.com/nullable"
)

type JsonType struct {
	String nullable.String `json:"string"`
}

func main() {
	dataJson := []byte(`{"string":null}`)

    var data JsonType
    if err := json.Unmarshal(dataJson,&data); err != nil {
        panic(err)
    }
	
    fmt.Println(data.String)
}
```

## Go References
[pkg.go.dev/go.portalnesia.com/nullable](https://pkg.go.dev/go.portalnesia.com/nullable)