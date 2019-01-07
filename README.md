# tmx
Parses json formatted Tiled maps.

### Example
```go
package main

import (
  "fmt"
  "github.com/drakbar/tmx"
)

func main() {
  
  m, err := tmx.LoadTileMap("path/to/mapfile.json")
  
  if err != nil {
    panic(err)
  }

  fmt.Printf("%+v\n", m)
}
```
