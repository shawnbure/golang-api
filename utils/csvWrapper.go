package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
)

type CsvWriterWrapper struct {
	writer *csv.Writer

	buffer *bytes.Buffer
}

func NewCsvWrapper() (*CsvWriterWrapper, error) {
	wrapper := &CsvWriterWrapper{}

	err := wrapper.create()
	if err != nil {
		return nil, err
	}

	return wrapper, nil
}

func (wrapper *CsvWriterWrapper) create() (err error) {
	wrapper.buffer = bytes.NewBufferString("")
	if wrapper.buffer == nil {
		return errors.New("Cannot create buffer")
	}

	wrapper.writer = csv.NewWriter(wrapper.buffer)
	//wrapper.reader = csv.NewReader(wrapper.buffer)
	return nil
}

func (wrapper *CsvWriterWrapper) WriteOneRecord(record []string) error {
	err := wrapper.writer.Write(record)
	if err != nil {
		return err
	}

	wrapper.writer.Flush()
	return nil
}

func (wrapper *CsvWriterWrapper) WriteBulkRecord(record [][]string) error {
	err := wrapper.writer.WriteAll(record)
	if err != nil {
		return err
	}

	wrapper.writer.Flush()
	return nil
}

func (wrapper *CsvWriterWrapper) Close() error {
	wrapper.writer.Flush()
	return nil
}

func (wrapper *CsvWriterWrapper) GetData() string {
	return wrapper.buffer.String()
}

func (wrapper *CsvWriterWrapper) GetBuffer() *bytes.Buffer {
	return wrapper.buffer
}
