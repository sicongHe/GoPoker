package poker

import (
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	t.Run("写入一个更短的值", func(t *testing.T) {
		file,clean := createTempfile(t, "12345")
		defer clean()
		tape := &Tape{file}
		tape.Write([]byte("abc"))
		file.Seek(0, 0)
		newFileContents, _ := ioutil.ReadAll(file)
		got := string(newFileContents)
		want := "abc"
		AssertString(t,got,want)
	})
}
