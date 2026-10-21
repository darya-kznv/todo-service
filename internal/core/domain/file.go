package domain

type File struct {
	buffer []byte
}

func NewFile(buffer []byte) File {
	return File{buffer: buffer}
}

func (f File) Buffer() []byte {
	return f.buffer
}
