package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	NULL_CHUNKSIZE = 5
	PAYLOADCHUNK   = 10
	FILE_NAME1     = "ReducedEntropy.bin"
	FILE_NAME2     = "OriginalEntropy.bin"
)

func banner() {
	fmt.Println(`
	▓█████  ███▄    █ ▄▄▄█████▓ ██▀███   ▒█████   ██▓███ ▓██   ██▓  █████▒██▓▒██   ██▒
	▓█   ▀  ██ ▀█   █ ▓  ██▒ ▓▒▓██ ▒ ██▒▒██▒  ██▒▓██░  ██▒▒██  ██▒▓██   ▒▓██▒▒▒ █ █ ▒░
	▒███   ▓██  ▀█ ██▒▒ ▓██░ ▒░▓██ ░▄█ ▒▒██░  ██▒▓██░ ██▓▒ ▒██ ██░▒████ ░▒██▒░░  █   ░
	▒▓█  ▄ ▓██▒  ▐▌██▒░ ▓██▓ ░ ▒██▀▀█▄  ▒██   ██░▒██▄█▓▒ ▒ ░ ▐██▓░░▓█▒  ░░██░ ░ █ █ ▒ 
	░▒████▒▒██░   ▓██░  ▒██▒ ░ ░██▓ ▒██▒░ ████▓▒░▒██▒ ░  ░ ░ ██▒▓░░▒█░   ░██░▒██▒ ▒██▒
	░░ ▒░ ░░ ▒░   ▒ ▒   ▒ ░░   ░ ▒▓ ░▒▓░░ ▒░▒░▒░ ▒▓▒░ ░  ░  ██▒▒▒  ▒ ░   ░▓  ▒▒ ░ ░▓ ░
	░ ░  ░░ ░░   ░ ▒░    ░      ░▒ ░ ▒░  ░ ▒ ▒░ ░▒ ░     ▓██ ░▒░  ░      ▒ ░░░   ░▒ ░
	░      ░   ░ ░   ░        ░░   ░ ░ ░ ░ ▒  ░░       ▒ ▒ ░░   ░ ░    ▒ ░ ░    ░  
	░  ░         ░             ░         ░ ░           ░ ░             ░   ░    ░  
														░ ░                         
	@Auth: C1ph3rX13
	@Blog: https://c1ph3rx13.github.io
	@Note: EntropyFix - Go
	@Warn: 代码仅供学习使用，请勿用于其他用途`)
}

func ReduceEntropy(payload []byte) ([]byte, int) {
	loop := len(payload) / PAYLOADCHUNK
	remainder := len(payload) % NULL_CHUNKSIZE
	newPayloadSize := (3 * len(payload)) / 2

	newPayload := make([]byte, newPayloadSize+remainder)

	nPcntr, oPcntr := 0, 0

	for z := 0; z < loop; z++ {
		for i := 0; i < PAYLOADCHUNK; i++ {
			newPayload[nPcntr] = payload[oPcntr]
			nPcntr++
			oPcntr++
		}

		for j := 0; j < NULL_CHUNKSIZE; j++ {
			newPayload[nPcntr] = 0x00
			nPcntr++
		}
	}

	if remainder > 0 {
		for i := 0; i != remainder; i++ {
			newPayload[nPcntr] = payload[oPcntr]
			nPcntr++
			oPcntr++
		}
	}

	return newPayload, newPayloadSize
}

func ReverseEntropy(payload []byte) ([]byte, int) {
	payloadSize := len(payload) - 1
	remainder := payloadSize % NULL_CHUNKSIZE
	newPayloadSize := (payloadSize / 3) * 2
	loop := newPayloadSize / PAYLOADCHUNK

	newPayload := make([]byte, newPayloadSize+remainder)

	nPcntr, oPcntr := 0, 0

	for i := 0; i < loop; i++ {
		for j := 0; j < PAYLOADCHUNK; j++ {
			newPayload[nPcntr] = payload[oPcntr]
			nPcntr++
			oPcntr++
		}

		for z := 0; z < NULL_CHUNKSIZE; z++ {
			oPcntr++ // ignoring 5 bytes
		}
	}

	if remainder > 0 {
		for i := 0; i != remainder; i++ {
			newPayload[nPcntr] = payload[oPcntr]
			nPcntr++
			oPcntr++
		}
	}

	return newPayload, newPayloadSize
}

// ReadPayloadFile reads a payload from a file
func ReadPayloadFile(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WritePayloadFile writes the payload to a file
func WritePayloadFile(data []byte, fileName string) error {
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	banner()

	// msfvenom -p windows/x64/exec CMD=calc.exe -f raw -o calc.bin
	inputFile := flag.String("input", "", "input payload")
	option := flag.Int("option", 0, "Choose option 1(ReducedEntropy) or 2(ReverseEntropy)")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Please provide an input payload")
		return
	}

	payload, err := ReadPayloadFile(*inputFile)
	if err != nil {
		log.Panicf("Error reading input payload: %v", err)
	}

	var newPayload []byte
	var newFileName string

	switch *option {
	case 1:
		newPayload, _ = ReduceEntropy(payload)
		newFileName = FILE_NAME1
	case 2:
		newPayload, _ = ReverseEntropy(payload)
		newFileName = FILE_NAME2
	default:
		log.Panic("Please enter 1(ReducedEntropy) or 2(ReverseEntropy) only")
	}

	if err := WritePayloadFile(newPayload, newFileName); err != nil {
		log.Panicf("Error writing output file: %v", err)
	}

	outputFileName := map[bool]string{*option != 2: FILE_NAME1, *option == 2: FILE_NAME2}[true]
	fmt.Println("File Outputted Successfully To", outputFileName)
}
