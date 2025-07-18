package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/klauspost/compress/zstd"
	"log"
	"runtime/debug"
)

// compress 使用Zstd压缩数据
func compress(data []byte) ([]byte, error) {
	// 创建编码器
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}
	defer encoder.Close()

	// 压缩数据
	return encoder.EncodeAll(data, make([]byte, 0, len(data))), nil
}

func performExperiment() {
	defer func() {
		if r := recover(); r != nil {
			// 打印 panic 信息和堆栈跟踪
			fmt.Printf("捕获到 panic: %v\n", r)
			log.Printf("堆栈跟踪: %s", string(debug.Stack()))
		}
	}()
	responses, err := ReadMessagesFromFile(RESPONSE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	if responses == nil {
		log.Fatal("responses is nil")
	}
	var totalResponseCount = len(responses)
	log.Printf("get %d responses", len(responses))
	v0Size := float64(0)
	v1Size := float64(0)
	for i := 0; i < totalResponseCount; i++ {
		// 序列化
		bytes, _ := proto.Marshal(responses[i])
		v1Bytes, _ := json.Marshal(responses[i])
		// zstd压缩
		bytes, _ = compress(bytes)
		v1Bytes, _ = compress(v1Bytes)

		v0Size += float64(len(bytes)) / float64(totalResponseCount)
		v1Size += float64(len(v1Bytes)) / float64(totalResponseCount)
		fmt.Println("success")
	}

	fmt.Printf("idc: %s\n", IDC_NAME)
	fmt.Printf("v0 size: %d byte\n", int(v0Size))
	fmt.Printf("v1 size: %d byte\n", int(v1Size))
}

func main() {
	performExperiment()
}
