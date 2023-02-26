package main

import (
	"flag"
	"fmt"
	"github.com/obarbier/steganography/raster_image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	encodeCmd = flag.NewFlagSet("encode", flag.ExitOnError)
	decodeCmd = flag.NewFlagSet("decode", flag.ExitOnError)
	coverPath string
	secret    string
	k         int
)

func init() {
	// ENCODE
	encodeCmd.StringVar(&coverPath, "cover", "", "path to cover image")
	encodeCmd.StringVar(&secret, "secret", "", "secret message")
	encodeCmd.IntVar(&k, "k", 2, "LSB to be modified")

	// DECODE
	decodeCmd.StringVar(&coverPath, "cover", "", "path to cover image")

}

func main() {

	//steganography is the process of embedding secrets in plain text medium such as images
	//files etc.
	// LSB: change the least significant bit of an image with the secrets

	//create output by reading input byte and change the least significant bits 2 bits at a time
	//flag.Parse()
	switch os.Args[1] {
	case "encode":
		encodeCmd.Parse(os.Args[2:])
		encode()
	case "decode":
		decodeCmd.Parse(os.Args[2:])
		decode()
	default:
		log.Fatalf("[ERROR] unknown subcommand '%s', see help for more details.", os.Args[1])
	}
}

func decode() {
	log.Println("awesome steganography decoding")
	if len(coverPath) == 0 {
		log.Fatal("cover path cannot be empty")
	}

	log.Printf("coverPath=%s", coverPath)
	inputFile, errOpen := os.Open(coverPath)
	if errOpen != nil {
		log.Fatalf("error while opening the file: %s", errOpen)
		return
	}
	encoder, jpegInitErr := raster_image.New(inputFile)
	if jpegInitErr != nil {
		log.Fatalf("error while reading jpeg file: %s", jpegInitErr)
		return
	}

	secret, err := encoder.Decode()
	if err != nil {
		log.Fatalf("error while decoding jpeg file: %s", err)
		return
	}

	fmt.Printf("secret=%s", string(secret))
}

func encode() {
	log.Println("awesome steganography encoding")
	if len(secret) == 0 {
		log.Fatal("secret message cannot be empty")
	}
	log.Printf("secret=%s", secret)
	log.Printf("k=%d", k)
	if len(coverPath) == 0 {
		log.Fatal("cover path cannot be empty")
	}
	log.Printf("coverPath=%s", coverPath)
	inputFile, errOpen := os.Open(coverPath)
	if errOpen != nil {
		log.Fatalf("error while opening the file: %s", errOpen)
		return
	}
	encoder, jpegInitErr := raster_image.New(inputFile)
	if jpegInitErr != nil {
		log.Fatalf("error while reading jpeg file: %s", jpegInitErr)
		return
	}
	jpegEncodeErr := encoder.Encode(secret)
	if jpegEncodeErr != nil {
		log.Fatalf("error while embedding secrets into jpeg file: %s", jpegEncodeErr)
		return
	}

	baseNameWithExt := filepath.Base(coverPath)
	log.Printf("baseNameWithExt=%s", baseNameWithExt)
	baseName := strings.Split(baseNameWithExt, ".")
	outputFileName := strings.Join([]string{baseName[0], "png"}, "_embedded.")
	log.Printf("create output file: %s", outputFileName)
	embeddedFile, createFileErr := os.Create(outputFileName)
	if createFileErr != nil {
		log.Fatalf("error while creartingFile: %s", createFileErr)
		return
	}
	defer func(embeddedFile *os.File) {
		err := embeddedFile.Close()
		if err != nil {
			log.Printf("error occurred while closing file: %v", err)
			return
		}
	}(embeddedFile)
	writeBytesCounts, savingErr := encoder.Write(embeddedFile)
	if savingErr != nil {
		log.Printf("error occurred while closing file: %v", savingErr)
		return
	}
	log.Println(fmt.Sprintf("wrote %d bytes into %s", writeBytesCounts, outputFileName))
}
