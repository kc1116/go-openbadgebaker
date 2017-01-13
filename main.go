package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	/*for _, a := range os.Args[1:] {
		if err := do(a); err != nil {
			log.Fatal(err)
		}
	}*/
	log.Println(getFile(os.Args[1]))
}

func do(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Printf("%s\n", filename)
	r := bufio.NewReader(f)

	var tmp [8]byte
	const pngHeader = "\x89PNG\r\n\x1a\n"
	if _, err := io.ReadFull(r, tmp[:len(pngHeader)]); err != nil {
		return err
	}
	if string(tmp[:len(pngHeader)]) != pngHeader {
		fmt.Println(string(tmp[:len(pngHeader)]))
		//return fmt.Errorf("not a PNG file")
	}

	// Process each chunk.
	for {
		_, err := io.ReadFull(r, tmp[:8])
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		length := binary.BigEndian.Uint32(tmp[:4])
		chunkType := string(tmp[4:8])
		//fmt.Printf("\tChunkType=%q, Length=%d\n", chunkType, length)

		chunkData, err := ioutil.ReadAll(&io.LimitedReader{
			R: r,
			N: int64(length),
		})
		if err != nil {
			return err
		}

		//fmt.Printf("\tChunkData=%q\n\n\n", chunkData)

		if chunkType == "iTXt" {
			// TODO: parse the chunkData as per the table at
			// https://www.w3.org/TR/PNG/#11iTXt
			//
			// Note that the text payload may be compressed as per
			// compress/zlib.
			fmt.Printf("\t\t%q\n\t\t% x\n", chunkData, chunkData)
		}

		// Read and ignore the checksum.
		if _, err := io.ReadFull(r, tmp[:4]); err != nil {
			return err
		}
	}
}

func getFile(filepath string) bool {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return false
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		fmt.Println(err)
		return false
	}

	filetype := http.DetectContentType(buff)
	out, err := os.Create("temp.jpg")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(filetype)

	switch filetype {
	case "image/jpeg", "image/jpg":
		return true

	case "image/png":
		return true

	default:
		return false
	}
}
