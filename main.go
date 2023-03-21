package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

	router.POST("upload", func(c echo.Context) error {
		var response Response
		var fileType, fileName string
		var fileSize int64
		isSuccess := true

		file, err := c.FormFile("file")
		if err != nil {
			isSuccess = false
		} else {
			src, err := file.Open()
			if err != nil {
				isSuccess = false
			} else {
				fileBytes, _ := ioutil.ReadAll(src)
				fileType = http.DetectContentType(fileBytes)

				if fileType == "application/pdf" {
					fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".pdf"
				} else {
					fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
				}

				err = ioutil.WriteFile(fileName, fileBytes, 0777)
				if err != nil {
					isSuccess = false
				} else {
					fileSize = file.Size
				}
			}

			defer src.Close()
		}

		if isSuccess {
			response = Response{
				Message: "Sukses Mengupload file",
				Data: struct {
					FileName string
					Filetype string
					Filesize int64
				}{
					FileName: fileName,
					Filetype: fileType,
					Filesize: fileSize,
				},
			}
		} else {
			response = Response{
				Message: "Gagal Mengupload file",
			}
		}

		return c.JSON(http.StatusOK, response)
	})
	router.Start(":8080")
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
