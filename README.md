# anypool
A pool of anything... I suppose?

## Example
```go
package main

import (
	"fmt"
	"github.com/dropout1337/anypool"
	"net/http"
	"time"
)

func main() {
	factory := func() any {
		return &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	pool := anypool.New(5, factory, anypool.WithReuseLimit(2))
	for {
		conn := pool.Get()

		client := conn.Conn.(*http.Client)
		resp, err := client.Get("http://example.com")
		if err != nil {
			fmt.Println("Error making HTTP request:", err)
		} else {
			fmt.Printf("%v / %v / %v\n", resp.Status, conn.ID, conn.ReuseCount)
			resp.Body.Close()
		}
	}
}

```
