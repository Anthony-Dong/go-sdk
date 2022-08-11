package commons

import "io"

func SkipReader(reader io.Reader, size int) error {
	_, err := ReadReader(reader, size)
	return err
}

func ReadReader(reader io.Reader, size int) ([]byte, error) {
	result := make([]byte, size)
	_, err := reader.Read(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
