package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenCSV(t *testing.T) {
	data := "totalIncome,wht,donation\n50000,5000,1000"
	fileHeader, err := MockFile(data, "file")
	assert.NoError(t, err)
	assert.NotNil(t, fileHeader)

	reader, closer, err := OpenCSV(fileHeader)
	assert.NoError(t, err)
	defer closer.Close()

	_, err = reader.Read()
	assert.NoError(t, err)
}

func TestValidateHeaders(t *testing.T) {
	expectedHeaders := []string{"totalIncome", "wht", "donation"}

	tests := []struct {
		name     string
		headers  []string
		expected error
	}{
		{"Correct Headers", []string{"totalIncome", "wht", "donation"}, nil},
		{"Incorrect Headers", []string{"income", "tax", "donation"}, fmt.Errorf("expected header totalIncome, got income")},
		{"Incomplete Headers", []string{"totalIncome", "wht"}, io.ErrUnexpectedEOF},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateHeaders(tc.headers, expectedHeaders)
			if tc.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expected.Error())
			}
		})
	}
}

func TestProcessRecords(t *testing.T) {
	data := "totalIncome,wht,donation\n50000,5000,1000\n60000,6000,1200"
	reader := csv.NewReader(strings.NewReader(data))
	_, _ = reader.Read() // Read off the header

	processMock := func(record []string) (interface{}, error) {
		if len(record) != 3 {
			return nil, io.ErrUnexpectedEOF
		}
		return record, nil
	}

	records, err := ProcessRecords(reader, processMock)
	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, []interface{}{[]string{"50000", "5000", "1000"}, []string{"60000", "6000", "1200"}}, records)
}

func TestMockFile(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{"Valid file", "totalIncome,wht,donation\n50000,5000,1000", "file", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFile, err := MockFile(tt.data, tt.want)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotFile)
			}
		})
	}
}
