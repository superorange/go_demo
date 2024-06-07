package proxy

import (
	"bufio"
	"fmt"
	"os"
)

func Read() {

	file, err := os.OpenFile("empty.txt", os.O_RDWR, 0755)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)

	scanner.Buffer(make([]byte, 3), 1024)

	for scanner.Scan() {

		fmt.Printf("%s END\n", scanner.Text())
	}

}
