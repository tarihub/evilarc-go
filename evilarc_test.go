package main

import (
	"fmt"
	"github.com/tarimoe/evilarc/pkg/helper"
	"log"
	"os"
	"strings"
	"testing"
)

func TestZip(t *testing.T) {
	of, err := os.Create("evil.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer of.Close()
	helper.CreateZip(of, []byte("test content"), "test_arc.txt", "test_passwd", "")
}

func TestZipWithPasswd(t *testing.T) {
	of, err := os.Create("evil.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer of.Close()
	helper.CreateZip(of, []byte("test content"), "test_arc.txt", "test_passwd", "standard")
}

func TestTarGZArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"travel", "--out", "xx.tar.gz"})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestZipArgsWithPasswd(t *testing.T) {
	rootCmd.SetArgs([]string{"travel", "--zip-enc", "STANDARD"})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestZipSymlinkWithNullArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"symlinks"})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestZipSymlinkWithPaaswd(t *testing.T) {
	rootCmd.SetArgs([]string{"symlinks", "--zip-enc", "STANDARD", "--zip-passwd", "1"})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestZipSymlink(t *testing.T) {
	rootCmd.SetArgs([]string{"symlinks", "--sym-file", "1.test", "--sym-target", "tmp", "--out", "sym-evil.zip"})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestTrimPrefix(t *testing.T) {
	fmt.Println(strings.TrimPrefix("etc", "/"))
	fmt.Println(strings.TrimPrefix("/etc", "/"))
	fmt.Println(strings.TrimPrefix("//etc", "/"))
}

func TestTarSymlink(t *testing.T) {
	rootCmd.SetArgs([]string{"symlinks", "--sym-file", "1.test", "--sym-target", "tmp", "--out", "sym-evil.tar"})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestTarGZSymlink(t *testing.T) {
	rootCmd.SetArgs([]string{"symlinks", "--sym-file", "1.test", "--sym-target", "tmp", "--out", "sym-evil.tar.gz"})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
