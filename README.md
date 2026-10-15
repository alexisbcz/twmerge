# twmerge

Merge Tailwind CSS class lists, like tailwind-merge. Later classes win over conflicting earlier ones.

## Install

```sh
go get github.com/alexisbcz/twmerge
```

## Use

```go
package main

import (
	"fmt"

	"github.com/alexisbcz/twmerge"
)

func main() {
	fmt.Println(twmerge.Merge("rounded-md px-2", "rounded-full")) // px-2 rounded-full
}
```

## License

MIT. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
