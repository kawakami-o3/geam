package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
)

const (
	idAtom = "Atom"
	idAtU8 = "AtU8"
	idCode = "Code"
	idAbst = "Abst"
	idCatT = "CatT"
	idFunT = "FunT"
	idExpT = "ExpT"
	idLitT = "LitT"
	idImpT = "ImpT"
	idLocT = "LocT"
	idLine = "Line"
	idStrT = "StrT"
	idAttr = "Attr"
)

func calcChunkRest(length uint32) int {
	align := 4
	return align*((int(length)+align-1)/align) - int(length)
}

type AtomChunk struct {
	Id        string
	Length    uint32
	AtomCount uint32
	Labels    []string
}

func ParseAtom(buffer *bytes.Buffer, id string) *AtomChunk {
	chunk := &AtomChunk{}
	//chunk.Id = string(buffer.Next(4))
	chunk.Id = id
	chunk.Length = binary.BigEndian.Uint32(buffer.Next(4))
	chunk.AtomCount = binary.BigEndian.Uint32(buffer.Next(4))

	length := 12

	labels := []string{}
	for i := 0; i < int(chunk.AtomCount); i++ {
		s := int(buffer.Next(1)[0])
		labels = append(labels, string(buffer.Next(s)))
		length += 1 + s
	}
	chunk.Labels = labels

	buffer.Next(calcChunkRest(chunk.Length))
	return chunk
}

type Term struct {
}

type CodeChunk struct {
	Id         string
	Length     uint32
	Version    uint32
	MaxOpcode  uint32
	LabelCount uint32
	FunCount   uint32
}

func ParseCode(buffer *bytes.Buffer, id string) *CodeChunk {
	chunk := &CodeChunk{}
	chunk.Id = id
	chunk.Length = binary.BigEndian.Uint32(buffer.Next(4))
	chunk.Version = binary.BigEndian.Uint32(buffer.Next(4))
	chunk.MaxOpcode = binary.BigEndian.Uint32(buffer.Next(4))
	chunk.LabelCount = binary.BigEndian.Uint32(buffer.Next(4))
	chunk.FunCount = binary.BigEndian.Uint32(buffer.Next(4))

	fmt.Println(buffer.Next(1))
	return chunk
}

type BeamData struct {
	Magic     string // 'FOR1'
	Length    uint32
	Type      string // 'BEAM'
	AtomChunk *AtomChunk
	CodeChunk *CodeChunk
}

func ParseBeamFile(beamPath string) (*BeamData, error) {

	file, err := os.Open(beamPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	data := &BeamData{}

	buffer := bytes.NewBuffer(bs)
	data.Magic = string(buffer.Next(4))
	data.Length = binary.BigEndian.Uint32(buffer.Next(4))
	data.Type = string(buffer.Next(4))

	for buffer.Len() > 0 {
		id := string(buffer.Next(4))
		switch id {
		case idAtom, idAtU8:
			data.AtomChunk = ParseAtom(buffer, id)
		case idCode:
			data.CodeChunk = ParseCode(buffer, id)
		default:
			//fmt.Println("id:", id)
			break
		}
	}

	return data, nil
}

func main() {
	beamPath := os.Args[1]

	data, _ := ParseBeamFile(beamPath)
	pp.Println(data)
}
