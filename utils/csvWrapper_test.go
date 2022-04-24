package utils

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCreateAndWriteCsvData(t *testing.T) {
	t.Run("Initialize Csv Wrapper and Write to it", func(t *testing.T) {
		wrapper, err := NewCsvWrapper()
		require.Nil(t, err)
		require.NotNil(t, wrapper)

		data := []string{"a", "b"}
		err = wrapper.WriteOneRecord(data)
		require.Nil(t, err)

		want := "a,b"
		got := wrapper.GetData()
		if strings.TrimSpace(got) != want {
			t.Fatalf("Want '%v', but got '%v'", want, got)
		}
	})

	t.Run("Saving multiple records", func(t *testing.T) {
		wrapper, err := NewCsvWrapper()
		require.Nil(t, err)
		require.NotNil(t, wrapper)

		data := [][]string{{"a", "b"}, {"c", "d"}}
		err = wrapper.WriteBulkRecord(data)
		require.Nil(t, err)

		want := `a,b
c,d`
		got := wrapper.GetData()
		if strings.TrimSpace(got) != want {
			t.Fatalf("Want '%v', but got '%v'", want, got)
		}
	})
}
