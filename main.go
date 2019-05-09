package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
	"github.com/kawakami-o3/geam/erl"
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

	return chunk
}

const (
	tagLiteral    byte = 0   //  0 0 0 0 0 | 0 0 0
	tagInt             = 1   //  0 0 0 0 0 | 0 0 1
	tagAtom            = 2   //  0 0 0 0 0 | 0 1 0
	tagXreg            = 3   //  0 0 0 0 0 | 0 1 1
	tagYreg            = 4   //  0 0 0 0 0 | 1 0 0
	tagLabel           = 5   //  0 0 0 0 0 | 1 0 1
	tagChar            = 6   //  0 0 0 0 0 | 1 1 0
	tagZ               = 7   //  0 0 0 0 0 | 1 1 1  7
	tagFlaot           = 23  //  0 0 0 1 0 | 1 1 1  16 + 7
	tagList            = 39  //  0 0 1 0 0 | 1 1 1  32 + 7
	tagFloatReg        = 55  //  0 0 1 1 0 | 1 1 1  48 + 7
	tagAllocList       = 71  //  0 1 0 0 0 | 1 1 1  64 + 7
	tagExtLiteral      = 87  //  0 1 0 1 0 | 1 1 1  80 + 7
	tagFull            = 255 //  1 1 1 1 1 | 1 1 1  255
	tagErr             = 255 //  1 1 1 1 1 | 1 1 1  255

)

type Term struct {
	TermType byte
	Content  []byte
	Value    int64
}

type Instruction struct {
	Opcode erl.Opcode
	Args   []Term
}

type CodeChunk struct {
	Chunk

	SubSize    uint32
	Version    uint32
	MaxOpcode  uint32
	LabelCount uint32
	FunCount   uint32

	Instructions []Instruction
}

func LoadCode(buffer *bytes.Buffer, id string) *CodeChunk {
	chunk := &CodeChunk{}
	chunk.load(buffer, id)

	return chunk
}

type StrTChunk struct {
	Chunk
}

func LoadStrT(buffer *bytes.Buffer, id string) *StrTChunk {
	chunk := &StrTChunk{}
	chunk.load(buffer, id)

	return chunk
}

type ImportEntry struct {
	ModuleIndex   int
	ModuleName    string
	FunctionIndex int
	FunctionName  string
	Arity         uint32
}

type ImpTChunk struct {
	Chunk
	ImportCount uint32
	ImportTable []ImportEntry
}

func LoadImpT(buffer *bytes.Buffer, id string) *ImpTChunk {
	chunk := &ImpTChunk{}
	chunk.load(buffer, id)
	return chunk
}

type ExportEntry struct {
	Index        int
	FunctionName string
	Arity        uint32
	Label        uint32 // TODO
}

type ExpTChunk struct {
	Chunk
	ExportCount uint32
	ExportTable []ExportEntry
}

func LoadExpT(buffer *bytes.Buffer, id string) *ExpTChunk {
	chunk := &ExpTChunk{}
	chunk.load(buffer, id)

	return chunk
}

type LitTChunk struct {
	Chunk
}

func LoadLitT(buffer *bytes.Buffer, id string) *LitTChunk {
	chunk := &LitTChunk{}
	chunk.load(buffer, id)

	return chunk
}

type LocTChunk struct {
	Chunk
}

func LoadLocT(buffer *bytes.Buffer, id string) *LocTChunk {
	chunk := &LocTChunk{}
	chunk.load(buffer, id)

	return chunk
}

type AttrChunk struct {
	Chunk
}

func LoadAttr(buffer *bytes.Buffer, id string) *AttrChunk {
	chunk := &AttrChunk{}
	chunk.load(buffer, id)

	return chunk
}

type CInfChunk struct {
	Chunk
}

func LoadCInf(buffer *bytes.Buffer, id string) *CInfChunk {
	chunk := &CInfChunk{}
	chunk.load(buffer, id)

	return chunk
}

type DbgiChunk struct {
	Chunk
}

func LoadDbgi(buffer *bytes.Buffer, id string) *DbgiChunk {
	chunk := &DbgiChunk{}
	chunk.load(buffer, id)

	return chunk
}

type DocsChunk struct {
	Chunk
}

func LoadDocs(buffer *bytes.Buffer, id string) *DocsChunk {
	chunk := &DocsChunk{}
	chunk.load(buffer, id)

	return chunk
}

type ExDpChunk struct {
	Chunk
}

func LoadExDp(buffer *bytes.Buffer, id string) *ExDpChunk {
	chunk := &ExDpChunk{}
	chunk.load(buffer, id)

	return chunk
}

type LineChunk struct {
	Chunk
}

func LoadLine(buffer *bytes.Buffer, id string) *LineChunk {
	chunk := &LineChunk{}
	chunk.load(buffer, id)

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

func (this *BeamData) parseAtoms() {
	chunk := this.AtomChunk
	data := bytes.NewBuffer(chunk.Data)
	chunk.AtomCount = binary.BigEndian.Uint32(data.Next(4))
	labels := []string{}
	for i := 0; i < int(chunk.AtomCount); i++ {
		s := int(data.Next(1)[0])
		labels = append(labels, string(data.Next(s)))
	}
	chunk.Labels = labels

}

func (this *BeamData) parseExports() {
	//if this.AtomChunk == nil { return }

	chunk := this.ExpTChunk
	data := bytes.NewBuffer(chunk.Data)
	chunk.ExportCount = binary.BigEndian.Uint32(data.Next(4))

	info := []ExportEntry{}
	for i := 0; i < int(chunk.ExportCount); i++ {
		e := ExportEntry{}
		e.Index = int(binary.BigEndian.Uint32(data.Next(4)))
		e.FunctionName = this.AtomChunk.Labels[e.Index-1]
		e.Arity = binary.BigEndian.Uint32(data.Next(4))
		e.Label = binary.BigEndian.Uint32(data.Next(4))
		info = append(info, e)
	}

	chunk.ExportTable = info
}

func (this *BeamData) parseImports() {
	//if this.AtomChunk == nil { return }

	chunk := this.ImpTChunk
	data := bytes.NewBuffer(chunk.Data)
	chunk.ImportCount = binary.BigEndian.Uint32(data.Next(4))

	info := []ImportEntry{}
	for i := 0; i < int(chunk.ImportCount); i++ {
		e := ImportEntry{}
		e.ModuleIndex = int(binary.BigEndian.Uint32(data.Next(4)))
		e.FunctionName = this.AtomChunk.Labels[e.ModuleIndex-1]
		e.FunctionIndex = int(binary.BigEndian.Uint32(data.Next(4)))
		e.FunctionName = this.AtomChunk.Labels[e.FunctionIndex-1]
		e.Arity = binary.BigEndian.Uint32(data.Next(4))
		info = append(info, e)
	}

	chunk.ImportTable = info
}

func decodeCompactTerm(data *bytes.Buffer) Term {
	bs := data.Next(1)
	b := bs[0]

	ret := Term{tagErr, bs, 0}
	switch b & tagZ {
	case tagLiteral, tagInt, tagAtom, tagXreg, tagYreg, tagLabel, tagChar:
		tag := b & tagZ
		if (b>>3)&1 == 0 {
			ret = Term{tag, bs, int64(b >> 4)}
		} else if (b>>3)&3 == 1 {
			bs = append(bs, data.Next(1)[0])
			value := int64(bs[1]<<3) + int64(bs[0]>>5)
			ret = Term{tag, bs, value}
		} else if (b>>3)&3 == 3 {
			if b>>3 != 31 { // 16 + 8 + 4 + 2 + 1
				size := int(b>>5) + 2
				vs := data.Next(size)
				bs = append(bs, vs...)
				value := int64(0)
				base := int64(1)
				for i, v := range vs {
					if i == 0 {
						value += int64(v)
					} else {
						value += int64(v) * base
					}
					base *= 256
				}
				ret = Term{tag, bs, value}
			} else {
				panic("not implemented")
			}
		}
	}
	return ret
}

func (this *BeamData) parseCode() {
	chunk := this.CodeChunk

	data := bytes.NewBuffer(chunk.Data)
	chunk.SubSize = binary.BigEndian.Uint32(data.Next(4))
	chunk.Version = binary.BigEndian.Uint32(data.Next(4))
	chunk.MaxOpcode = binary.BigEndian.Uint32(data.Next(4))
	chunk.LabelCount = binary.BigEndian.Uint32(data.Next(4))
	chunk.FunCount = binary.BigEndian.Uint32(data.Next(4))

	// TODO
	insts := []Instruction{}
	for i := uint32(0); i < chunk.FunCount; i++ {
		opId := int(data.Next(1)[0])
		opc := findOpcode(opId)
		args := []Term{}
		l := 0
		for l < opc.Arity {
			t := decodeCompactTerm(data)
			l += len(t.Content)
			args = append(args, t)
		}
		insts = append(insts, Instruction{
			Opcode: opc,
			Args:   args,
		})
	}
	//pp.Println(insts)
	chunk.Instructions = insts
}

func findOpcode(id int) erl.Opcode {
	for _, o := range erl.Opcodes {
		if o.Id == id {
			return o
		}
	}
	return erl.Opcodes[0]
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

	data.parseAtoms()
	data.parseExports()
	data.parseImports()
	data.parseCode()
	return data, nil
}

func main() {
	beamPath := os.Args[1]

	data, err := LoadBeamFile(beamPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	pp.Println(data.CodeChunk)
	//pp.Println("debug", len([]*BeamData{data}))
}
