package media

import "io"
import "os"

// TmpReader is a file, and a reader for the file, where the file is deleted from disk
// when the reader is closed
func TempReader(f *os.File, r io.ReadCloser) *tempReader {
  return &tempReader{f, r}
}

type tempReader struct {
  *os.File
  io.ReadCloser
}

func (t *tempReader) Read(p []byte) (int, error) {
  return t.ReadCloser.Read(p)
}

func (t *tempReader) Close() error {
  t.ReadCloser.Close()
  t.File.Close()
  if err := os.Remove(t.File.Name()); err != nil {
    return err
  }
  return nil
}
