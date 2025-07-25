package main

import (
	com_ss_ugc_tiktok "code.byted.org/tiktok/pb_builder/proto_gen"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/klauspost/compress/zstd"
	"log"
	"postlist-resp-field-statistic/common"
	"postlist-resp-field-statistic/config"
)

const KB = 1024

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

func performExperiment(responses []*com_ss_ugc_tiktok.AwemeV1AwemePostResponse) {
	responsesLen := len(responses)
	fmt.Printf("%d responses\n", responsesLen)
	v0Size := float64(0)
	v1Size := float64(0)
	for i := 0; i < responsesLen; i++ {
		// 序列化
		bytes, _ := json.Marshal(responses[i])
		v1Bytes, _ := proto.Marshal(responses[i])
		// zstd压缩
		bytes, _ = compress(bytes)
		v1Bytes, _ = compress(v1Bytes)

		v0Size += float64(len(bytes)) / float64(responsesLen)
		v1Size += float64(len(v1Bytes)) / float64(responsesLen)
	}

	fmt.Printf("idc: %s\n", config.IDC_NAME)
	fmt.Printf("v0 size: %d byte\n", int(v0Size))
	fmt.Printf("v1 size: %d byte\n", int(v1Size))
	fmt.Printf("平均减少%fKB，降低%f%%\n", (v0Size-v1Size)/float64(KB), (v0Size-v1Size)/v0Size*100)
}

func perform() {
	responses, err := common.ReadMessagesFromFile(config.RESPONSE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	if responses == nil {
		log.Fatal("responses is nil")
	}
	var totalResponseCount = len(responses)
	fmt.Printf("get %d responses int total\n", totalResponseCount)
	fmt.Printf("=================================\n")
	fmt.Printf("total response\n")
	performExperiment(responses)
	// 分类<20k、[20k, 40k)、[40k, 60k), >=60k
	var type1Responses []*com_ss_ugc_tiktok.AwemeV1AwemePostResponse
	var type2Responses []*com_ss_ugc_tiktok.AwemeV1AwemePostResponse
	var type3Responses []*com_ss_ugc_tiktok.AwemeV1AwemePostResponse
	var type4Responses []*com_ss_ugc_tiktok.AwemeV1AwemePostResponse
	for _, res := range responses {
		// proto序列化
		bytes, _ := proto.Marshal(res)
		// zstd压缩
		bytes, _ = compress(bytes)
		if len(bytes) < 20*KB {
			type1Responses = append(type1Responses, res)
		} else if len(bytes) >= 20*KB && len(bytes) < 40*KB {
			type2Responses = append(type2Responses, res)
		} else if len(bytes) >= 40*KB && len(bytes) < 60*KB {
			type3Responses = append(type3Responses, res)
		} else {
			type4Responses = append(type4Responses, res)
		}
	}
	fmt.Printf("==================================\n")
	fmt.Printf("<20k\n")
	performExperiment(type1Responses)
	fmt.Printf("==================================\n")
	fmt.Printf("[20k, 40k)\n")
	performExperiment(type2Responses)
	fmt.Printf("==================================\n")
	fmt.Printf("[40k, 60k)\n")
	performExperiment(type3Responses)
	fmt.Printf("==================================\n")
	fmt.Printf(">=60k\n")
	performExperiment(type4Responses)
}

func main() {
	perform()
}
