package history

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

func EncodeByGzip(src []byte) (string, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	// 将src数据写入gzip.Writer
	if _, err := gz.Write(src); err != nil {
		fmt.Printf("encode write failed, src: %s error: %s\n", src, err)
		return "", err
	}
	if err := gz.Flush(); err != nil {
		fmt.Printf("flush gzip writer failed, err: %s\n", err)
		return "", err
	}
	if err := gz.Close(); err != nil {
		fmt.Printf("close zip writer failed, err: %s\n", err)
		return "", err
	}
	// 返回buf.Bytes并编码成base64的字符串
	str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return str, nil
}

func DecodeByGzip(src string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(src) // 将src转成byte数组
	if err != nil {
		log.Printf("decode string failed, src: %s\b error: %s\n", src, err)
		return nil, err
	}

	rdata := bytes.NewReader(data)  // 创建bytes.Reader
	r, err := gzip.NewReader(rdata) // 创建gzip.Reader
	if err != nil {
		log.Printf("new gzip reader failed, error: %s\n", err)
		return nil, err
	}
	s, err := io.ReadAll(r) // 将压缩后的数据读到gzip.Reader 并将解压后的数据读出来
	if err != nil {
		log.Printf("read all failed, error: %s\n", err)
		return nil, err
	}
	return s, nil
}
