package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
)

// parameter number detection function
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		// exit directly
		os.Exit(1)
	}
}

// Implement int64 to []byte
func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if nil != err {
		log.Panicf("int transact to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// Standard JSON format to slice
// need to add quotation marks under windows
// bc.exe send -from "[\"troytan\"]" -to "[\"Alice\"]" -amount "[\"5\"]"
// bc.exe send -from "[\"troytan\",\"Alice\"]" -to "[\"Alice\",\"troytan\"]" -amount "[\"5\", \"2\"]"
// troytan -> alice 5 --> alice 5, troytan 5
// alice -> troytan 2 --> alice 3, troytan 7
func JSONToSlice(jsonString string) []string {
	var strSlice []string
	// json
	if err := json.Unmarshal([]byte(jsonString), &strSlice); nil != err {
		log.Panicf("json to []string failed! %v\n", err)
	}
	return strSlice
}

// string to hash160
func StringToHash160(address string) []byte {
	pubKeyHash := Base58Decode([]byte(address))
	hash160 := pubKeyHash[:len(pubKeyHash)-addressCheckSumLen]
	return hash160
}

// get node ID
func GetEnvNodeId() string {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("NODE_ID is not set...")
		os.Exit(1)
	}
	return nodeID
}

// gob encoding
func gobEncode(data interface{}) []byte {
	var result bytes.Buffer
	enc := gob.NewEncoder(&result)
	err := enc.Encode(data)
	if nil != err {
		log.Panicf("encode the data failed! %v\n", err)
	}
	return result.Bytes()
}

// command converted to request ([]byte)
func commandToBytes(command string) []byte {
	var bytes [CMMAND_LENGTH]byte
	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

// Reverse parsing, parse the command in the request
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x00 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}
