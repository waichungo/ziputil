package main

// #include<stdlib.h>
//#include <string.h>
import "C"
import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"unsafe"
)

//export DownloadToBuffer
func DownloadToBuffer(linkStr *C.char, buffPtr unsafe.Pointer) C.int {
	link := C.GoString(linkStr)
	ret := 0
	if downloadToBuffer(link, buffPtr) {
		ret = 1
	}
	return C.int(ret)
}

//export DownloadToFile
func DownloadToFile(linkStr, fileStr *C.char) C.int {
	link := C.GoString(linkStr)
	file := C.GoString(fileStr)
	ret := 0
	if downloadToFile(link, file) {
		ret = 1
	}

	return C.int(ret)
}

//export ExtractArchive
func ExtractArchive(srcStr, destStr *C.char) C.int {
	src := C.GoString(srcStr)
	dest := C.GoString(destStr)
	ret := 0
	if extractArchive(src, dest) {
		ret = 1
	}

	return C.int(ret)
}
func main() {

}
func downloadToFile(link, file string) bool {
	resp, err := http.Get(link)
	if err == nil {
		defer resp.Body.Close()
		abs, _ := filepath.Abs(file)
		dir := filepath.Dir(abs)
		if !Exists(dir) {
			os.MkdirAll(dir, 0755)
		}
		fh, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err == nil {
			defer fh.Close()
			_, err = io.Copy(fh, resp.Body)
			return err == nil
		}
	}
	return false
}
func downloadToBuffer(link string, buffPtr unsafe.Pointer) bool {
	resp, err := http.Get(link)
	if err == nil {
		defer resp.Body.Close()
		data := make([]byte, 1024)
		buff := bytes.NewBuffer(data)
		if err == nil {
			_, err = io.Copy(buff, resp.Body)
			ptr := C.malloc(C.ulonglong(buff.Len()))
			C.memcpy(unsafe.Pointer(ptr), unsafe.Pointer(&buff.Bytes()[0]), C.ulonglong(buff.Len()))
			if err == nil {
				buffPtr = unsafe.Pointer(ptr)
			}
			return err == nil
		}
	}
	return false
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
func extractArchive(src, dest string) bool {

	zipReader, _ := zip.OpenReader(src)
	paths := []string{}
	success := true
	if !Exists(dest) {
		os.MkdirAll(dest, 0755)
	}
	for _, file := range zipReader.Reader.File {

		zippedFile, err := file.Open()
		if err != nil {
			success = false
			break
		}
		defer zippedFile.Close()

		extractedFilePath := filepath.Join(
			dest,
			file.Name,
		)
		paths = append(paths, extractedFilePath)
		if file.FileInfo().IsDir() {

			os.MkdirAll(extractedFilePath, file.Mode())
		} else {

			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				success = false
				break
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				success = false
				break
			}
		}
	}
	if !success {
		for _, path := range paths {
			st, err := os.Stat(path)
			if err == nil {

				if st.IsDir() {
					os.RemoveAll(path)
				} else {
					os.Remove(path)
				}
			}
		}
	}
	return success
}
