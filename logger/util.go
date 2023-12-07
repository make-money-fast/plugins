package logger

type ByteJson []byte

func RawJSON(b []byte) ByteJson {
	return ByteJson(b)
}

func (b *ByteJson) Marshal() ([]byte, error) {
	return *b, nil
}

func (b *ByteJson) UnMarshal(s []byte) error {
	*b = s
	return nil
}
