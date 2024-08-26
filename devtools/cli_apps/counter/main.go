package counter

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	lines := flag.Bool("l", false, "count lines")
	bytes := flag.Bool("b", false, "count bytes")
	flag.Parse()
	fmt.Println(count(os.Stdin, *lines, *bytes))
}

func count(r io.Reader, countLines bool, bytes bool) int {
	scanner := bufio.NewScanner(r)

	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	if bytes {
		scanner.Split(bufio.ScanBytes)
	}

	wc := 0
	for scanner.Scan() {
		wc++
	}

	return wc
}
