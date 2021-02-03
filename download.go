package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
)

//DownloadFromURL ...
func DownloadFromURL(url string, destination string) (int64, error) {

	output, err := os.Create(destination)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	response, err := http.Get(url)

	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return 0, err
	}

	return n, nil
}

//DownloadImage downloads the url to the destination path. The destPath
//should include the filename without the extension, which is calculated
//here.
func DownloadImage(url string, destPath string) (savedPath string, err error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	fileType := validateImageType(d)

	savedPath = destPath + "." + fileType

	//check if file exists
	if _, err := os.Stat(savedPath); !os.IsNotExist(err) {
		os.Remove(savedPath)
	}
	writeFile(savedPath, d)

	return savedPath, nil
}

func validateImageType(b []byte) string {
	m := http.DetectContentType(b)
	switch m {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	}
	return ""
}

//ReceiveNamedFile ...
func ReceiveNamedFile(file multipart.File, destination string) error {

	//check if file exists
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		os.Remove(destination)
	}

	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, file)

	if err != nil {
		return err
	}

	return nil
}

//ReceiveFile ...
func ReceiveFile(file multipart.File, header *multipart.FileHeader, requirePDF bool, overwrite bool, destinationFolder string) (fileName string, err error) {

	fileName = header.Filename
	filePath := path.Join(destinationFolder, fileName)

	//check if file exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		if overwrite {
			os.Remove(filePath)
		} else {
			//Construct a new filename by adding _1, _2, etc.
			//until the filename is unique
			fileName = getUniqueFilename(fileName, destinationFolder)
			//Update filePath
			filePath = path.Join(destinationFolder, fileName)
		}
	}

	output, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer output.Close()

	_, err = io.Copy(output, file)

	if err != nil {
		return "", err
	}

	if requirePDF {
		isPDF, _ := IsPDF(filePath)

		if !isPDF {
			os.Remove(filePath)
			return "", fmt.Errorf("file was not PDF")
		}
	}

	return fileName, nil
}

func getUniqueFilename(fileIn string, destinationFolder string) (fileOut string) {

	filePath := path.Join(destinationFolder, fileIn)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fileOut = fileIn
		return fileOut
	}

	fileNameRaw := fileIn //File name as it arrives

	fileName := ""
	extension := ""

	//Get the file parts
	dotLocation := strings.LastIndex(fileNameRaw, ".")

	if dotLocation > -1 {
		fileName = fileNameRaw[0:dotLocation]
		extension = fileNameRaw[dotLocation+1:]
	}
	//Add a number to the filename and check whether it exists
	i := 1
	for {
		fileOut = fmt.Sprintf("%s_%v.%s", fileName, i, extension)
		filePath := path.Join(destinationFolder, fileOut)
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			return fileOut
		}
		i++
	}

}

//IsPDF opens a file and checks whether it has a PDF signature
func IsPDF(filepath string) (bool, error) {

	// Open File
	f, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Get the content
	contentType, err := getFileContentType(f)
	if err != nil {
		return false, err
	}

	if strings.ToLower(contentType) == "application/pdf" {
		return true, nil
	}

	return false, fmt.Errorf("file is not PDF")

}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
