package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("urls.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		site := strings.TrimSpace(scanner.Text())
		if site != "" {
			fmt.Println(site)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
