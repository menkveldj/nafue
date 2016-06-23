package nafue

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"os"
	"crypto/sha256"
	"github.com/menkveldj/nafue-api/models/display"
	"os/user"
	"path/filepath"
	"encoding/base64"
	"crypto/rand"
)


var fileIdRegex = regexp.MustCompile(`^.*file/(.*)$`)

func GetFile(url string) (*display.FileHeaderDisplay, string, error) {

	// get api url from share link
	aUrl := appifyUrl(url)

	// create temp file
	secureFile, err := createTempFile()
	if err != nil {
		return nil, "", err
	}
	defer secureFile.Close()

	// download file header info
	fileHeader, err := getFileHeader(aUrl)
	if err != nil {
		os.Remove(secureFile.Name())
		return nil, "", err
	}

	// download file body
	err = getFileBody(secureFile, fileHeader)
	if err != nil {
		os.Remove(secureFile.Name())
		return nil, "", err
	}

	return fileHeader, secureFile.Name(), nil
}

func UnsealFile(secureFileUri string, pass string, fileHeader *display.FileHeaderDisplay, pathToSave string) (string, error) {

	// open file for writing
	secureFile, err := os.OpenFile(secureFileUri, os.O_RDWR, 0777)
	if err != nil {
		return "", err
	}
	defer secureFile.Close()

	// decrypt to file
	fileName, err := Decrypt(secureFile, pass, fileHeader)
	if err != nil {
		if err != C_DECRYPT_UNAUTHENTICATED {
			os.Remove(secureFile.Name())
		}
		return "", err
	}

	// move file into place
	fp := filepath.Join(pathToSave, fileName)
	fileUri := filepath.Clean(fp)
	err = os.Rename(secureFileUri, fileUri)
	if err != nil {
		os.Remove(secureFile.Name())
		return "", err
	}

	return fileUri, nil
}

func SealShareFile(fileUri string, pass string) (string, error) {

	file, err := os.Open(fileUri)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// check file is under 50mb
	fstat, err := file.Stat()
	if err != nil {
		return"",  err
	}
	if fstat.Size() > (C.FILE_SIZE_LIMIT * 1024 * 1024) {
		err := errors.New("File is larger than " + strconv.FormatInt(C.FILE_SIZE_LIMIT, 10) + "mb.")
		return "", err
	}

	// create temp secure file
	secureFile, err := createTempFile()
	if err != nil {
		return "", err
	}

	// encrypt to temp file
	fileHeader, err := Encrypt(file, secureFile, pass)
	if err != nil {
		return "", err
	}

	//create checksum
	checksum := sha256.New()
	_, err = io.Copy(checksum, secureFile)
	if err != nil {
		return "", nil
	}

	err = putFileHeader(C.API_FILE_URL, fileHeader)
	if err != nil {
		return "", errors.New("PutFileHeader: " + err.Error())
	}


	// post body data
	err = putFileBody(fileHeader, secureFile)
	if err != nil {
		return "", errors.New("PutFileBody: " + err.Error())
	}

	// delete temp file
	err = os.Remove(secureFile.Name())
	if err != nil {
		return "", err
	}

	// provide share link
	shareLink := C.SHARE_LINK + fileHeader.ShortUrl
	return shareLink, nil
}

func createTempFile() (*os.File, error){
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	tmpDir := filepath.Join(usr.HomeDir, ".nafue")
	err = os.MkdirAll(tmpDir, os.ModeDir)
	if err != nil {
		return nil, err
	}

	// random file
	ran, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}

	w, err := os.Create(filepath.Join(tmpDir, ran + ".enn"))
	if err != nil {
		return nil, err
	}

	return w, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	if err != nil {
		return "", err
	}
	code := base64.URLEncoding.EncodeToString(b)
	return code[0:s], nil
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
