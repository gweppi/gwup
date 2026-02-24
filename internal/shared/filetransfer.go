package shared

import (
	"io"
	"mime/multipart"
	"os"
)

func MultipartFilePipe(fileName string, file *os.File) (io.Reader, string, error) {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	go func() {
		defer w.Close()
		part, err := w.CreateFormFile(FieldName, fileName)
		if err != nil {
			pw.CloseWithError(err)
		}
		io.Copy(part, file)
	}()
	
	return pr, w.FormDataContentType(), nil
}
