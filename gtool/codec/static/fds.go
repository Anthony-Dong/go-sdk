package static

type FileDescriptor interface {
	DecodeMessage(message string, content []byte) (out []byte, err error)
}
