package packer

import (
	"encoding/binary"
	"io/ioutil"
)

const indexValueBitSize = 4
const numberOfValuesPerFile = 2

type PackedFile struct {
	NumberOfFiles int
	FileBounds    []FileBound
	IndexContent  []byte
	FileContents  []byte
}

type FileBound struct {
	Index  uint32
	Length uint32
}

type File struct {
	Index   int
	Content []byte
}

func (pf *PackedFile) Files() <-chan File {
	out := make(chan File)

	go func() {
		defer close(out)

		for i := 0; i < pf.NumberOfFiles; i++ {
			var f File
			f.Index = i
			f.Content = pf.LoadFile(i)
			out <- f
		}
	}()

	return out
}

func (pf *PackedFile) LoadFile(i int) []byte {
	b := pf.FileBounds[i]
	idxSize := indexSize(pf.NumberOfFiles)
	from := b.Index - uint32(idxSize)
	to := b.Index - uint32(idxSize) + b.Length

	return pf.FileContents[from:to]
}

func New(filenames []string) (*PackedFile, error) {
	var (
		pf          = PackedFile{}
		indexLength = indexSize(len(filenames))
	)

	for _, filename := range filenames {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		if pf.NumberOfFiles == 0 {
			pf.FileBounds = append(pf.FileBounds, FileBound{
				Index:  uint32(indexLength),
				Length: uint32(len(content)),
			})
		} else {
			prev := pf.FileBounds[len(pf.FileBounds)-1]
			pf.FileBounds = append(pf.FileBounds, FileBound{
				Index:  prev.Index + prev.Length,
				Length: uint32(len(content)),
			})
		}

		pf.NumberOfFiles++
		pf.FileBounds = append(pf.FileBounds)
		pf.FileContents = append(pf.FileContents, content...)
	}

	return &pf, nil
}

func indexSize(numFiles int) int {
	return numFiles*(indexValueBitSize*numberOfValuesPerFile) + indexValueBitSize
}

func (pf *PackedFile) Write(s string) error {
	pf.genIndex()
	content := append(pf.IndexContent, pf.FileContents...)

	return ioutil.WriteFile(s, content, 0644)
}

func (pf *PackedFile) genIndex() {
	var index []byte

	var numFiles = make([]byte, 4)
	binary.LittleEndian.PutUint32(numFiles, uint32(pf.NumberOfFiles))
	index = append(index, numFiles...)

	var idx = make([]byte, 4)
	var length = make([]byte, 4)
	for i := range pf.FileBounds {
		bounds := pf.FileBounds[i]

		binary.LittleEndian.PutUint32(idx, bounds.Index)
		binary.LittleEndian.PutUint32(length, bounds.Length)

		index = append(index, idx...)
		index = append(index, length...)
	}

	pf.IndexContent = index
}

func Read(filename string) (*PackedFile, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var pf = PackedFile{}
	const l = indexValueBitSize

	pf.NumberOfFiles = int(binary.LittleEndian.Uint32(content[0:l]))
	var offset = 0
	for i := 0; i < pf.NumberOfFiles; i++ {
		offset += l
		index := content[offset : offset+l]
		offset += l
		length := content[offset : offset+l]

		pf.FileBounds = append(pf.FileBounds, FileBound{
			Index:  binary.LittleEndian.Uint32(index),
			Length: binary.LittleEndian.Uint32(length),
		})
	}

	pf.IndexContent = content[:offset]
	pf.FileContents = content[offset+l:]

	return &pf, nil
}
