package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// PolybiusSquare is an alias for a 5x5 matrix
type polybiusSquare [][]string

func main() {
	var configFilePath string = ""
	flag.StringVar(
		&configFilePath,
		"pb",
		"",
		"Specify the file of your 5x5 polybius square",
	)
	var mustGenRandomPolybius bool
	flag.BoolVar(
		&mustGenRandomPolybius,
		"gen",
		false,
		"Use this to random generate a new polybius square using atmospheric noise",
	)
	var apiKey string = ""
	flag.StringVar(
		&apiKey,
		"apiKey",
		"",
		"Your API key for the RANDOM.ORG service",
	)

	flag.Parse()

	if mustGenRandomPolybius {
		pb := generateRandomLatinPolybius(apiKey)
		filePath := writePolybiusToFile(pb)
		fmt.Println("------------")
		fmt.Println("🆗 Created new random polybius at path: " + filePath)
		fmt.Println("Re-run the program with the new polybius file")
		fmt.Println("------------")
		os.Exit(0)
	}

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
		input := readInput()
		operation := string(input[0])

		if operation == "+" {
			fmt.Println("🔐 Encrypting message...")
			output = encrypt(input[1:], polybius)
		} else if operation == "-" {
			fmt.Println("🔓 Decrypting message...")
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

func writePolybiusToFile(pb polybiusSquare) string {
	fileName := "polybius_" + strconv.Itoa(int(time.Now().Unix())) + ".json"
	filePath := "./generated_polybius_squares/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	textBytes, err := json.Marshal(pb)
	for _, b := range textBytes {
		writer.WriteByte(b)
	}
	writer.Flush()
	return filePath
}

func readInput() string {
	text := ""
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text = scanner.Text()

	return text
}

type intPair [2]int

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
	for _, intPair := range grouped {
		cipherBuilder.WriteString(findLetterByXY(intPair, pb))
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

func groupPairs(codes []int) []intPair {
	grouped := []intPair{}

	for i := 0; i < len(codes); i += 2 {
		x := codes[i]
		y := codes[i+1]
		grouped = append(grouped, intPair{x, y})
	}
	return grouped
}

func findLetterByXY(coords intPair, pb polybiusSquare) string {
	x := pb[coords[0]]
	letter := x[coords[1]]
	return letter
}

func findXInPolybius(target string, pb polybiusSquare) int {
	_target := target
	if target == "J" {
		_target = "I"
	}
	for i, row := range pb {
		for _, v := range row {
			if string(v) == _target {
				return i
			}
		}
	}
	return -1
}

func findYInPolybius(target string, pb polybiusSquare) int {
	_target := target
	if target == "J" {
		_target = "I"
	}
	for _, row := range pb {
		for j, v := range row {
			if string(v) == _target {
				return j
			}
		}
	}
	return -1
}

type params struct {
	APIKey      string `json:"apiKey"`
	N           int    `json:"n"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Replacement bool   `json:"replacement"`
}
type body struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  params `json:"params"`
	ID      int    `json:"id"`
}
type random struct {
	Data           []int  `json:"data"`
	CompletionTime string `json:"completionTime"`
}
type result struct {
	Random        random `json:"random"`
	BitsUsed      int    `json:"bitsUsed"`
	BitsLeft      int    `json:"bitsLeft"`
	RequestsLeft  int    `json:"requestsLeft"`
	AdvisoryDelay int    `json:"advisoryDelay"`
}
type response struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  result `json:"result"`
	ID      int    `json:"id"`
}

func generateRandomLatinPolybius(apiKey string) polybiusSquare {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	amount := len(alphabet)
	maxUpperBound := amount - 1
	randomNumbers := generateRandomNumbers(amount, 0, maxUpperBound, false, apiKey)
	randomNumbersWithJRemoved := []int{}
	// remove J's index in the alphabet because I and J refer to the same index in a polybius square
	for _, v := range randomNumbers {
		if v != 9 {
			randomNumbersWithJRemoved = append(randomNumbersWithJRemoved, v)
		}
	}

	pb := make([][]string, 5, 5)
	chunkSize := 5
	x := 0
	for i := 0; i < maxUpperBound; i += chunkSize {
		min := i + chunkSize
		if amount <= min {
			min = amount
		}
		randIndices := randomNumbersWithJRemoved[i:min]
		row := []string{}
		for _, v := range randIndices {
			row = append(row, string(alphabet[v]))
		}
		pb[x] = append(pb[x][:], row...)
		x++
	}
	return pb
}

func generateRandomNumbers(amount int, from int, to int, allowDuplicates bool, apiKey string) []int {
	body, err := json.Marshal(body{
		ID:      213,
		Jsonrpc: "2.0",
		Method:  "generateIntegers",
		Params: params{
			APIKey:      apiKey,
			N:           amount,
			Min:         from,
			Max:         to,
			Replacement: allowDuplicates,
		},
	})
	if err != nil {
		fmt.Print(err)
	}

	resp, err := http.Post("https://api.random.org/json-rpc/2/invoke", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("ERROR generating random sequence. %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR generating random sequence. %s", err.Error())
	}

	var data response
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		fmt.Printf("Cannot parse JSON response. %s", err.Error())
	}

	return data.Result.Random.Data
}
