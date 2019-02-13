package main

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	idCInf = "CInf"
	idDbgi = "Dbgi"
	idDocs = "Docs"
	idExDp = "ExDp"
)

type Chunk struct {
	Id     string
	Length uint32
	Data   []byte
}

func (this *Chunk) load(buffer *bytes.Buffer, id string) {
	this.Id = id
	this.Length = binary.BigEndian.Uint32(buffer.Next(4))
	this.Data = buffer.Next(int(this.Length))

	buffer.Next(this.calRest())
}

func (this *Chunk) calRest() int {
	align := 4
	length := int(this.Length)
	return align*((length+align-1)/align) - length
}

type AtomChunk struct {
	Chunk

	AtomCount uint32
	Labels    []string
}

func LoadAtom(buffer *bytes.Buffer, id string) *AtomChunk {
	chunk := &AtomChunk{}
	chunk.load(buffer, id)

	data := bytes.NewBuffer(chunk.Data)
	chunk.AtomCount = binary.BigEndian.Uint32(data.Next(4))
	labels := []string{}
	for i := 0; i < int(chunk.AtomCount); i++ {
		s := int(data.Next(1)[0])
		labels = append(labels, string(data.Next(s)))
	}
	chunk.Labels = labels

	return chunk
}

type CodeChunk struct {
	Chunk

	Version    uint32
	MaxOpcode  uint32
	LabelCount uint32
	FunCount   uint32
}

func LoadCode(buffer *bytes.Buffer, id string) *CodeChunk {
	chunk := &CodeChunk{}
	chunk.load(buffer, id)

	data := bytes.NewBuffer(chunk.Data)
	chunk.Version = binary.BigEndian.Uint32(data.Next(4))
	chunk.MaxOpcode = binary.BigEndian.Uint32(data.Next(4))
	chunk.LabelCount = binary.BigEndian.Uint32(data.Next(4))
	chunk.FunCount = binary.BigEndian.Uint32(data.Next(4))

	// TODO
	//opcodes := []string{}
	/*
		for i := uint32(0); i < chunk.FunCount; i++ {
			fmt.Println(buffer.Next(1))
			return chunk
		}
	*/
	return chunk
}

type StrTChunk struct {
	Chunk
}

func LoadStrT(buffer *bytes.Buffer, id string) *StrTChunk {
	chunk := &StrTChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type ImpTChunk struct {
	Chunk
}

func LoadImpT(buffer *bytes.Buffer, id string) *ImpTChunk {
	chunk := &ImpTChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type ExpTChunk struct {
	Chunk
}

func LoadExpT(buffer *bytes.Buffer, id string) *ExpTChunk {
	chunk := &ExpTChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type LitTChunk struct {
	Chunk
}

func LoadLitT(buffer *bytes.Buffer, id string) *LitTChunk {
	chunk := &LitTChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type LocTChunk struct {
	Chunk
}

func LoadLocT(buffer *bytes.Buffer, id string) *LocTChunk {
	chunk := &LocTChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type AttrChunk struct {
	Chunk
}

func LoadAttr(buffer *bytes.Buffer, id string) *AttrChunk {
	chunk := &AttrChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type CInfChunk struct {
	Chunk
}

func LoadCInf(buffer *bytes.Buffer, id string) *CInfChunk {
	chunk := &CInfChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type DbgiChunk struct {
	Chunk
}

func LoadDbgi(buffer *bytes.Buffer, id string) *DbgiChunk {
	chunk := &DbgiChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type DocsChunk struct {
	Chunk
}

func LoadDocs(buffer *bytes.Buffer, id string) *DocsChunk {
	chunk := &DocsChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type ExDpChunk struct {
	Chunk
}

func LoadExDp(buffer *bytes.Buffer, id string) *ExDpChunk {
	chunk := &ExDpChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type LineChunk struct {
	Chunk
}

func LoadLine(buffer *bytes.Buffer, id string) *LineChunk {
	chunk := &LineChunk{}
	chunk.load(buffer, id)

	// TODO
	return chunk
}

type BeamData struct {
	Magic     string // 'FOR1'
	Length    uint32
	Type      string // 'BEAM'
	AtomChunk *AtomChunk
	CodeChunk *CodeChunk
	StrTChunk *StrTChunk
	ImpTChunk *ImpTChunk
	ExpTChunk *ExpTChunk
	LitTChunk *LitTChunk
	LocTChunk *LocTChunk
	AttrChunk *AttrChunk
	CInfChunk *CInfChunk
	DbgiChunk *DbgiChunk
	DocsChunk *DocsChunk
	ExDpChunk *ExDpChunk
	LineChunk *LineChunk
}

func LoadBeamFile(beamPath string) (*BeamData, error) {

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
		bs := buffer.Next(4)
		id := string(bs)
		fmt.Println(id)
		switch id {
		case idAtom, idAtU8:
			data.AtomChunk = LoadAtom(buffer, id)
		case idCode:
			data.CodeChunk = LoadCode(buffer, id)
		case idStrT:
			data.StrTChunk = LoadStrT(buffer, id)
		case idImpT:
			data.ImpTChunk = LoadImpT(buffer, id)
		case idExpT:
			data.ExpTChunk = LoadExpT(buffer, id)
		case idLitT:
			data.LitTChunk = LoadLitT(buffer, id)
		case idLocT:
			data.LocTChunk = LoadLocT(buffer, id)
		case idAttr:
			data.AttrChunk = LoadAttr(buffer, id)
		case idCInf:
			data.CInfChunk = LoadCInf(buffer, id)
		case idDbgi:
			data.DbgiChunk = LoadDbgi(buffer, id)
		case idDocs:
			data.DocsChunk = LoadDocs(buffer, id)
		case idExDp:
			data.ExDpChunk = LoadExDp(buffer, id)
		case idLine:
			data.LineChunk = LoadLine(buffer, id)
		default:
			return nil, errors.New(fmt.Sprintf("unknown id: %s %v", id, bs))
		}
	}

	return data, nil
}

func main() {
	beamPath := os.Args[1]

	data, err := LoadBeamFile(beamPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	pp.Println(data)
}
