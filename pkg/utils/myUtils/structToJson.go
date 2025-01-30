package myUtils

import (
	"encoding/json"
	"errors"
	"fmt"
)

// StructToJSON converts a struct to a JSON string.
func StructToJSON(input interface{}) (string, error) {
	// Marshal the input struct to JSON bytes
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return "", errors.New("failed to marshal struct to JSON: " + err.Error())
	}

	// Convert JSON bytes to string
	jsonString := string(jsonBytes)
	return jsonString, nil
}

// PrettyStructToJSON converts a struct to a pretty-printed JSON string.
func PrettyStructToJSON(input interface{}) (string, error) {
	// Marshal the input struct to JSON bytes with indentation
	jsonBytes, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", errors.New("failed to marshal struct to pretty JSON: " + err.Error())
	}

	// Convert JSON bytes to string
	jsonString := string(jsonBytes)
	return jsonString, nil
}

// JSONToStruct converts a JSON string to a struct.
func JSONToStruct(jsonString string, output interface{}) error {
	// Unmarshal the JSON string into the output struct
	err := json.Unmarshal([]byte(jsonString), output)
	if err != nil {
		return errors.New("failed to unmarshal JSON to struct: " + err.Error())
	}
	return nil
}

// ExampleStruct is an example struct for demonstration purposes.
type ExampleStruct struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}

// ExampleUsage demonstrates how to use the functions in this package.
func ExampleUsage() {
	// Create an example struct
	example := ExampleStruct{
		Name:    "John Doe",
		Age:     30,
		Email:   "john.doe@example.com",
		IsAdmin: true,
	}

	// Convert struct to JSON
	jsonString, err := StructToJSON(example)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("JSON String:", jsonString)

	// Convert struct to pretty-printed JSON
	prettyJSONString, err := PrettyStructToJSON(example)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Pretty JSON String:", prettyJSONString)

	// Convert JSON string back to struct
	var result ExampleStruct
	err = JSONToStruct(jsonString, &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Converted Struct:", result)
}
