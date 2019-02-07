package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	beamPath := os.Args[1]

	file, _ := os.Open(beamPath)
	defer file.Close()

	bs, _ := ioutil.ReadAll(file)

	buffer := bytes.NewBuffer(bs)
	magic := string(buffer.Next(4))
	sizeBody := binary.BigEndian.Uint32(buffer.Next(4))
	formType := string(buffer.Next(4))

	fmt.Println("len: ", len(bs))
	fmt.Println("magic:", magic)
	fmt.Println("body size:", sizeBody)
	fmt.Println("form type:", formType)

	chunkId := binary.BigEndian.Uint32(buffer.Next(4))
	fmt.Println("chunk id:", chunkId)
	//chunkId := buffer.Next(4)
	//fmt.Println("chunk id:", binary.BigEndian.Uint32(chunkId), chunkId)

	chunkSize := binary.BigEndian.Uint32(buffer.Next(4))
	fmt.Println("chunk size:", chunkSize)

	numAtoms := binary.BigEndian.Uint32(buffer.Next(4))
	fmt.Println("number of atoms", numAtoms)

	//fmt.Println(buffer.Next(int(chunkSize - uint32(12))))
	//fmt.Println(string(buffer.Next(int(chunkSize - uint32(12)))))

	atoms := []string{}
	for i := 0; i < int(numAtoms); i++ {
		s := int(buffer.Next(1)[0])
		atoms = append(atoms, string(buffer.Next(s)))
	}
	fmt.Println("atoms:", atoms)


	exptId := binary.BigEndian.Uint32(buffer.Next(4))


}
