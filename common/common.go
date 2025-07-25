package common

import (
	"bufio"
	com_ss_ugc_tiktok "code.byted.org/tiktok/pb_builder/proto_gen"
	"encoding/json"
	"fmt"
	"os"
)

func ReadMessagesFromFile(filePath string) ([]*com_ss_ugc_tiktok.AwemeV1AwemePostResponse, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var responses []*com_ss_ugc_tiktok.AwemeV1AwemePostResponse
	var lineBuilder []byte

	for {
		// 读取一行（可能是不完整的，需要拼接）
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}

		// 拼接行（处理超长行）
		lineBuilder = append(lineBuilder, line...)

		// 如果不是前缀，表示一行完整读取完毕
		if !isPrefix {
			// 尝试解析JSON
			response := &com_ss_ugc_tiktok.AwemeV1AwemePostResponse{}
			if err := json.Unmarshal(lineBuilder, response); err != nil {
				return nil, fmt.Errorf("解析JSON失败: %w, 行内容: %s", err, string(lineBuilder))
			}
			//EnsureNestedMessages(response)
			responses = append(responses, response)
			lineBuilder = nil // 重置缓冲区
		}
	}

	// 检查是否有未处理的最后一行
	if len(lineBuilder) > 0 {
		response := &com_ss_ugc_tiktok.AwemeV1AwemePostResponse{}
		if err := json.Unmarshal(lineBuilder, response); err != nil {
			return nil, fmt.Errorf("解析最后一行JSON失败: %w, 行内容: %s", err, string(lineBuilder))
		}
		//EnsureNestedMessages(response)
		responses = append(responses, response)
	}
	return responses, nil
}
