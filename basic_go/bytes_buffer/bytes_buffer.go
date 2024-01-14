package bytesbuffer

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func TestByteBuffer() {
	buf1 := bytes.NewBufferString("hello") // create from string
	fmt.Println(buf1)                      // hello

	buf2 := bytes.NewBuffer([]byte("hello"))
	fmt.Println(buf2) // hello
}

func TestByteBufferWrite() {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("hello") // write string
	fmt.Println(buf.String())

	buf.WriteByte(' ') // write byte

	buf.WriteRune('世')
	buf.WriteRune('界') // write rune

	buf.Write([]byte(" world")) // write []byte
	fmt.Println(buf.String())
}

func TestWriteToFile() {
	file, err := os.OpenFile("./text.txt", os.O_CREATE|os.O_RDWR, fs.ModeAppend)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBufferString("hello world")
	n, err := buf.WriteTo(file) // write to file
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("write %d bytes data to file, after write buf is %s\n", n, buf.String()) // after write buf is empty
	}
}

func TestRead() {
	buf := bytes.NewBufferString("hello world")
	var sRead = make([]byte, 3)
	n, err := buf.Read(sRead) // read to sread 3 byte
	if err != nil && err != io.EOF {
		fmt.Println("read failed, err: ", err)
	} else {
		fmt.Printf("read %d byte from buf, sread is %s\n", n, sRead)
	}
}

func TestReadByte() {
	buf := bytes.NewBufferString("hello world")
	b, _ := buf.ReadByte() // 读取一个byte
	fmt.Printf("read data= %c, after read one byte buf = %s\n", b, buf.String())

	bytes, _ := buf.ReadBytes(' ') // readbytes按分隔符读取
	fmt.Printf("read data= %s, after read one byte buf = %s\n", bytes, buf.String())
}

func TestReadString() {
	buf := bytes.NewBufferString("hello world")
	str, _ := buf.ReadString(' ') // 按分隔符读取string, str = hello
	fmt.Printf("read data= %s, after read one byte buf = %s\n", str, buf.String())

	bytes := buf.Next(buf.Len()) // Next读取固定长度的内容 bytes = world, buf为空
	fmt.Printf("read data= %s, after read one byte buf = %s\n", bytes, buf.String())
}
