package helper

import (
	"bytes"
	"github.com/yeka/zip"
	"io"
	"os"
	"strings"
)

type ZipWriter struct {
	*zip.Writer
}

func NewZipWriter(file io.Writer) *ZipWriter {
	return &ZipWriter{zip.NewWriter(file)}
}

// CreateZipHeader
// mode: support symlink and default now
func (w *ZipWriter) CreateZipHeader(name, password, enc, mode string) (io.Writer, error) {
	var fh *zip.FileHeader
	if strings.ToLower(mode) == "symlink" {
		fh = &zip.FileHeader{
			Name:   name,
			Method: zip.Deflate,
		}
		fh.SetMode(os.ModeSymlink)
	} else {
		fh = &zip.FileHeader{
			Name:   name,
			Method: zip.Deflate,
		}
	}

	var zem zip.EncryptionMethod
	switch strings.ToUpper(enc) {
	case "STANDARD":
		zem = zip.StandardEncryption
	case "AES128":
		zem = zip.AES128Encryption
	case "AES192":
		zem = zip.AES192Encryption
	case "AES256":
		zem = zip.AES256Encryption
	case "":
		zem = 0
	default:
		panic("Only support Standard and AES128/192/256 encryption.")
	}
	if zem > 0 {
		fh.SetPassword(password)
		fh.SetEncryptionMethod(zem)
	}

	return w.CreateHeader(fh)
}

func (w *ZipWriter) AddToZip(name, password, enc, content, mode string) (io.Writer, error) {
	iow, err := w.CreateZipHeader(name, password, enc, mode)
	if err != nil {
		return nil, err
	}

	err = BytesToZipW(iow, []byte(content), mode)
	if err != nil {
		return nil, err
	}
	return iow, nil
}

func BytesToZipW(w io.Writer, content []byte, mode string) error {
	var err error = nil
	// To improve performance
	if strings.ToLower(mode) == "symlink" {
		_, err = w.Write(content)
	} else {
		_, err = io.Copy(w, bytes.NewReader(content))
	}
	return err
}

func CreateZip(of *os.File, content []byte, filename, password, enc string) {
	zipW := NewZipWriter(of)
	defer zipW.Close()
	_, err := zipW.AddToZip(filename, password, enc, string(content), "")
	if err != nil {
		panic(err)
	}
	zipW.Flush()
}

func CreateSymZip(of *os.File, content []byte, filename, symName, symTarget, password, enc string) {
	zipW := NewZipWriter(of)
	defer zipW.Close()

	_, err := zipW.AddToZip(symName, password, enc, symTarget, "symlink")
	if err != nil {
		panic(err)
	}

	_, err = zipW.AddToZip(filename, password, enc, string(content), "")
	if err != nil {
		panic(err)
	}

	zipW.Flush()
}
