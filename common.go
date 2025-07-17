package main

import (
	"bufio"
	com_ss_ugc_tiktok "code.byted.org/tiktok/pb_builder/proto_gen"
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"io"
	"os"
)

func ReadMessagesFromFile(filePath string) ([]*com_ss_ugc_tiktok.AwemeV1AwemePostResponse, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建带缓冲的读取器
	reader := bufio.NewReader(file)

	var responses []*com_ss_ugc_tiktok.AwemeV1AwemePostResponse

	for {
		// 读取长度前缀
		var length uint64
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			if err == io.EOF {
				// 文件结束，正常退出
				break
			}
			return nil, err
		}

		// 读取消息数据
		data := make([]byte, length)
		if _, err := io.ReadFull(reader, data); err != nil {
			return nil, err
		}
		response := &com_ss_ugc_tiktok.AwemeV1AwemePostResponse{}
		// 反序列化消息
		err := proto.Unmarshal(data, response)
		if err != nil {
			return nil, err
		}

		responses = append(responses, response)
	}

	return responses, nil
}
