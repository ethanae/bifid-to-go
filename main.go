package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// PolybiusSquare is an alias for a 5x5 matrix
type polybiusSquare [5][5]string

func main() {
	var configFilePath string = ""
	flag.StringVar(
		&configFilePath,
		"pb",
		"",
		"Specify the file of your 5x5 polybius square",
	)

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("bifid -polybius ./<some_file_path>.json")
	}
	flag.Parse()

	configFile, fileErr := ioutil.ReadFile(configFilePath)
	if fileErr != nil {
		panic("Unable to read file " + configFilePath + ".\n" + fileErr.Error())
	}

	var polybius polybiusSquare
	err := json.Unmarshal(configFile, &polybius)
	if err != nil {
		panic("Unable to parse configuration file" + err.Error())
	}

	var output string
	for {
		fmt.Println("------------")
		fmt.Print("input  > ")
		input := readPlainText()
		operation := string(input[0])

		if operation == "+" {
			fmt.Println("Encrypting message...")
			output = encrypt(input[1:], polybius)
		} else if operation == "-" {
			fmt.Println("Decrypting message...")
			output = decrypt(input[1:], polybius)
		} else {
			fmt.Println("Missing operation. The very first character must be either '+' for encryption or '-' for decryption.")
			fmt.Println("Example: +HELLO or -YUOYO")
			continue
		}

		fmt.Printf("output > %s\n", output)
		fmt.Println("------------")
	}
}

func readPlainText() string {
	text := ""
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text = scanner.Text()

	return text
}

type pair [2]int

func encrypt(message string, pb polybiusSquare) string {
	xCoords := []int{}
	yCoords := []int{}

	for _, v := range message {
		x := findXInPolybius(string(v), pb)
		y := findYInPolybius(string(v), pb)
		xCoords = append(xCoords, x)
		yCoords = append(yCoords, y)
	}

	combined := append(xCoords, yCoords...)
	grouped := groupPairs(combined)

	var cipherBuilder strings.Builder
	for _, pair := range grouped {
		cipherBuilder.WriteString(findLetterByXY(pair, pb))
	}
	return cipherBuilder.String()
}

func decrypt(cipher string, pb polybiusSquare) string {
	combined := []int{}
	for _, v := range cipher {
		x := findXInPolybius(string(v), pb)
		y := findYInPolybius(string(v), pb)
		combined = append(combined, x, y)
	}

	halfLength := len(combined) / 2
	xCoords := combined[0:halfLength]
	yCoords := combined[halfLength:]

	var plaintextBuilder strings.Builder
	for i := 0; i < halfLength; i++ {
		x := xCoords[i]
		y := yCoords[i]
		plaintextBuilder.WriteString(pb[x][y])
	}
	return plaintextBuilder.String()
}

func groupPairs(codes []int) []pair {
	grouped := []pair{}

	for i := 0; i < len(codes); i += 2 {
		x := codes[i]
		y := codes[i+1]
		grouped = append(grouped, pair{x, y})
	}
	return grouped
}

func findLetterByXY(coords pair, pb polybiusSquare) string {
	x := pb[coords[0]]
	letter := x[coords[1]]
	return letter
}

func findXInPolybius(target string, pb polybiusSquare) int {
	for i, row := range pb {
		for _, v := range row {
			if string(v) == target {
				return i
			}
		}
	}
	return -1
}

func findYInPolybius(target string, pb polybiusSquare) int {
	for _, row := range pb {
		for j, v := range row {
			if string(v) == target {
				return j
			}
		}
	}
	return -1
}
