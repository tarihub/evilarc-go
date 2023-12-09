package helper

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path"
	"time"
)

func CreateTar(filename string, of *os.File, content []byte) {
	tw := tar.NewWriter(of)
	defer tw.Close()

	hdr := &tar.Header{
		Name:    filename,
		Mode:    int64(os.ModePerm),
		Size:    int64(len(content)),
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		panic(err)
	}
	if _, err := tw.Write(content); err != nil {
		panic(err)
	}
	if err := tw.Close(); err != nil {
		panic(err)
	}
}

func GetGzipBytes(data []byte) []byte {
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err := gzipWriter.Write(data)
	if err != nil {
		panic(err)
	}
	if err = gzipWriter.Close(); err != nil {
		panic(err)
	}
	return compressedData.Bytes()
}

func CreateTarGZ(filename string, of *os.File, content []byte) {
	CreateTar(filename, of, GetGzipBytes(content))
}

func CreateSymTar(of *os.File, content []byte, filename, symName, symTarget string) {
	tw := tar.NewWriter(of)
	defer tw.Close()

	linkHeader := &tar.Header{
		Name:     symName,
		Linkname: symTarget,
		Typeflag: tar.TypeSymlink,
		Mode:     0777,
	}
	if err := tw.WriteHeader(linkHeader); err != nil {
		panic(err)
	}

	hdr := &tar.Header{
		Name:    path.Join(symName, filename),
		Mode:    int64(os.ModePerm),
		Size:    int64(len(content)),
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		panic(err)
	}
	if _, err := tw.Write(content); err != nil {
		panic(err)
	}
	if err := tw.Close(); err != nil {
		panic(err)
	}
}

func CreateSymTarGZ(of *os.File, content []byte, filename, symName, symTarget string) {
	CreateSymTar(of, GetGzipBytes(content), filename, symName, symTarget)
}
