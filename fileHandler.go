package nafue

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"fmt"
	"os"
)

var fileIdRegex = regexp.MustCompile(`^.*file/(.*)$`)


//func GetFile(url string) (*[]byte, *display.FileHeaderDisplay, error) {
//
//	// get api url from share link
//	aUrl := appifyUrl(url)
//
//	// dowload file header info
//	var fileHeader = display.FileHeaderDisplay{}
//	getFileHeader(aUrl, &fileHeader)
//
//	// dowload file body
//	fileBody, err := getFileBody(fileHeader.DownloadUrl)
//	if err != nil {
//		return nil, nil, err
//	}
//	return fileBody, &fileHeader, nil
//
//}
//
//func TryDecrypt(body *[]byte, header *display.FileHeaderDisplay, pass string) (io.Reader, string, error) {
//	fileBody, err := Decrypt(header, pass, body)
//	if err != nil {
//		return bytes.NewBufferString(""), "", err
//	}
//	return bytes.NewBuffer(fileBody.Content), fileBody.Name, nil
//}

func SealFile(reader io.ReaderAt, writer io.WriterAt, fileInfo os.FileInfo, name, pass string) error {

	// check file is under 50mb
	if fileInfo.Size() > (C.FILE_SIZE_LIMIT * 1024 * 1024) {
		err := errors.New("File is larger than " + strconv.FormatInt(C.FILE_SIZE_LIMIT, 10) + "mb.")
		return err
	}

	// encrypt to temp file
	fileHeader, err := Encrypt(reader, writer, fileInfo.Name(), pass)
	if err != nil {
		return err
	}

	fmt.Println("File Header: ", fileHeader)

	//err = putFileHeader(C.API_FILE_URL, &fileHeader)
	//if err != nil {
	//	return err
	//}

	//// create pbkdf2 key
	//key := getPbkdf2(pass, salt)
	//
	//// encrypt and write data
	//err = Encrypt(file, key, &fileHeader)
	//
	//// encrypt file with password
	////secureData, fileHeader, err := Encrypt(fileBody, pass)
	//if err != nil {
	//	return "", err
	//}
	//
	//
	//// post body data
	////putFileBody(fileHeader.UploadUrl, secureData)
	//
	//// provide share link
	//shareLink := C.SHARE_LINK + fileHeader.ShortUrl
	//return shareLink, nil
	return nil
}

//
//func getFileContentsFromReader(reader io.Reader, size int64, name string) (*models.FileBody, error) {
//
//
//	fileBytes, err := ioutil.ReadAll(reader)
//	if err != nil {
//		return nil , err
//	}
//
//	fbp := models.FileBody{
//		Name:    name,
//		Part:    0,
//		Content: fileBytes,
//	}
//
//	return &fbp, nil
//}
//
//func writeFileContentsToPath(fileBody *models.FileBody) error {
//
//	err := ioutil.WriteFile(fileBody.Name, fileBody.Content, 0644)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
