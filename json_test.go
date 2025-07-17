package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonSerialization(t *testing.T) {
	data := map[string]interface{}{
		"name":    "Bob",
		"age":     25,
		"hobbies": []string{"swimming", "hiking"},
		"active":  true,
	}
	// 序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
	// 输出：{"active":true,"age":25,"hobbies":["swimming","hiking"],"name":"Bob"}
}
