package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

type HeaderValidator func([]string) error
type RecordProcessor func([]string) (interface{}, error)

func OpenCSV(file *multipart.FileHeader) (*csv.Reader, io.Closer, error) {
	src, err := file.Open()
	if err != nil {
		return nil, nil, err
	}
	return csv.NewReader(src), src, nil
}

func ValidateHeaders(headers, expectedHeaders []string) error {
	if len(headers) != len(expectedHeaders) {
		return io.ErrUnexpectedEOF
	}

	for i, header := range headers {
		if header != expectedHeaders[i] {
			return fmt.Errorf("expected header %s, got %s", expectedHeaders[i], header)
		}
	}

	return nil
}

func ProcessRecords(csvReader *csv.Reader, process RecordProcessor) ([]interface{}, error) {
	var records []interface{}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		processedRecord, err := process(record)
		if err != nil {
			return nil, err
		}

		records = append(records, processedRecord)
	}
	return records, nil
}

func MockFile(data, fieldName string) (*multipart.FileHeader, error) {
	body := "--boundary\r\n"
	body += "Content-Disposition: form-data; name=\"" + fieldName + "\"; filename=\"testdata.csv\"\r\n"
	body += "Content-Type: text/csv\r\n\r\n"
	body += data + "\r\n--boundary--\r\n"

	r := multipart.NewReader(strings.NewReader(body), "boundary")
	form, err := r.ReadForm(1024)
	if err != nil {
		return nil, err
	}
	fileHeaders, ok := form.File[fieldName]
	if !ok || len(fileHeaders) == 0 {
		return nil, fmt.Errorf("no file part found with name %s", fieldName)
	}
	return fileHeaders[0], nil
}
