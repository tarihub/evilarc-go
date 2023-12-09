package helper

import (
	"archive/tar"
	"github.com/dsnet/compress/bzip2"
	"os"
	"path"
	"time"
)

func CreateBZ2(filename string, of *os.File, content []byte) {
	bz2Writer, err := bzip2.NewWriter(of, &bzip2.WriterConfig{Level: 9})
	defer bz2Writer.Close()
	if err != nil {
		panic(err)
	}

	tarWriter := tar.NewWriter(bz2Writer)
	defer tarWriter.Close()

	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(content)),
	}
	if err = tarWriter.WriteHeader(hdr); err != nil {
		panic(err)
	}
	if _, err = tarWriter.Write(content); err != nil {
		panic(err)
	}
}

func CreateSymBZ2(of *os.File, content []byte, filename, symName, symTarget string) {
	bz2Writer, err := bzip2.NewWriter(of, &bzip2.WriterConfig{Level: 9})
	defer bz2Writer.Close()
	if err != nil {
		panic(err)
	}

	tarWriter := tar.NewWriter(bz2Writer)
	defer tarWriter.Close()

	linkHeader := &tar.Header{
		Name:     symName,
		Linkname: symTarget,
		Typeflag: tar.TypeSymlink,
		Mode:     0777,
	}
	if err = tarWriter.WriteHeader(linkHeader); err != nil {
		panic(err)
	}

	hdr := &tar.Header{
		Name:    path.Join(symName, filename),
		Mode:    int64(os.ModePerm),
		Size:    int64(len(content)),
		ModTime: time.Now(),
	}
	if err = tarWriter.WriteHeader(hdr); err != nil {
		panic(err)
	}
	if _, err = tarWriter.Write(content); err != nil {
		panic(err)
	}
}
