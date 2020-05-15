package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Person struct {
	Name string
	Age uint
}

func main() {
	var liqi Person
	liqi.Name="李琪"
	liqi.Age=22
	var bufer bytes.Buffer
	encoder := gob.NewEncoder(&bufer)
	err := encoder.Encode(&liqi)
	if err != nil {
		log.Panic("编码失败")
	}
	fmt.Printf("编码后的李琪：%v\n",bufer.Bytes())

	decoder := gob.NewDecoder(bytes.NewReader(bufer.Bytes()))
	xiaoli := Person{}
	err = decoder.Decode(&xiaoli)
	if err != nil {
		log.Panic("解码失败")
	}
	fmt.Printf("解码后的小李：%v\n",xiaoli)
}
