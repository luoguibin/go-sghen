package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"go-sghen/helper"
	"go-sghen/models"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

// FileUploaderController ...
type FileUploaderController struct {
	BaseController
}

// UploadConfig ...
type UploadConfig struct {
	WaitArea  []*FileUploaderController
	IsWaiting bool
	WaitMax   int
	RateUnit  int64
	RateMax   int64
	RateTotal int64
}

func (c *UploadConfig) addWaitRequest(fc *FileUploaderController) bool {
	if c.RateTotal >= c.RateMax || len(c.WaitArea) > c.WaitMax {
		return false
	}

	// for _, v := range c.WaitArea {
	// 	if fc.Ctx == v.Ctx || fc.Ctx.Request == v.Ctx.Request {
	// 		fmt.Println("same Ctx or Request")
	// 		return false
	// 	}
	// }

	c.WaitArea = append(c.WaitArea, fc)
	return true
}
func (c *UploadConfig) removeWaitRequest(fc *FileUploaderController) bool {
	index := -1
	for i, v := range c.WaitArea {
		if v == fc {
			index = i
		}
	}
	if index == -1 {
		return false
	}

	c.WaitArea = append(c.WaitArea[:index], c.WaitArea[index+1:]...)
	c.RateTotal -= c.RateUnit
	c.IsWaiting = false
	return true
}
func (c *UploadConfig) callNextUpload() {
	if len(c.WaitArea) == 0 || c.RateTotal >= c.RateMax || c.IsWaiting {
		return
	}
	fc := c.WaitArea[0]
	c.RateTotal += c.RateUnit
	c.IsWaiting = true
	fc.FileUpload()
}

// MUploadConfig ...
var MUploadConfig = &UploadConfig{
	WaitArea:  make([]*FileUploaderController, 0),
	WaitMax:   10,
	RateUnit:  300,
	RateMax:   1024,
	RateTotal: 0,
}

func (c *FileUploaderController) uploadEnd(data ResponseData) {
	c.respToJSON(data)
	MUploadConfig.removeWaitRequest(c)

	MUploadConfig.callNextUpload()
}

// FileUploadQueue ...
func (c *FileUploaderController) FileUploadQueue() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	isOk = MUploadConfig.addWaitRequest(c)
	if !isOk {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "上传服务繁忙"
		c.respToJSON(data)
		return
	}
	MUploadConfig.callNextUpload()
}

// FileUpload 文件上传
func (c *FileUploaderController) FileUpload() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.uploadEnd(data)
		return
	}

	// 上传的文件存储在maxMemory大小的内存里面
	// 如果文件大小超过了maxMemory，那么剩下的部分将存储在系统的临时文件中
	c.Ctx.Request.ParseMultipartForm(32 << 20)
	// c.GetFile("file")	// 单文件
	// contentType := c.Ctx.Request.Header.Get("Content-Type")
	// "multipart/form-data"
	if c.Ctx.Request.MultipartForm == nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "获取上传信息失败"
		c.uploadEnd(data)
		return
	}
	fileHeaders, err := c.GetFiles("file") // 多文件

	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "获取上传文件列表失败"
		c.uploadEnd(data)
		return
	}

	// 检测上传目录是否存在
	pathType := c.GetString("pathType")
	if len(pathType) == 0 {
		pathType = "normal"
	}

	path, ok := models.MConfig.PathTypeMap[pathType]
	if !ok {
		path = models.MConfig.PathTypeMap["normal"]
	}

	isExist, _ := helper.PathExists(path)
	if !isExist {
		isMade := helper.MkdirAll(path)
		if !isMade {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "系统内部错误"
			c.uploadEnd(data)
			return
		}
	}

	if len(fileHeaders) > models.MConfig.MaxUploadCount {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "上传文件个数不能超过" + strconv.Itoa(models.MConfig.MaxUploadCount) + "个"
		c.uploadEnd(data)
		return
	}

	// 限制大文件上传
	for index, v := range fileHeaders {
		if len(v.Filename) <= 0 {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "上传文件不能为空"
			c.uploadEnd(data)
			return
		}
		sizeMB := int(v.Size / 1024 / 1024)
		if sizeMB > models.MConfig.MaxUploadSize {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "第" + strconv.Itoa(index+1) + "个文件:" + v.Filename + "文件大小超过" + strconv.Itoa(models.MConfig.MaxUploadSize) + "MB"
			c.uploadEnd(data)
			return
		}
	}

	// 遍历文件
	list := make([]string, 0)
	for index, v := range fileHeaders {
		fileName := v.Filename
		// fmt.Println(fileName)
		file, err := v.Open()
		defer file.Close()

		if err == nil {
			// 设置文件名字
			outputFilePath := path + fileName
			writer, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0666)

			if err == nil {
				// 保存文件内容
				io.Copy(writer, file)
				writer.Close()
				file.Seek(0, os.SEEK_SET)
				// 文件md5计算
				h := md5.New()
				io.Copy(h, file)
				fileRename := hex.EncodeToString(h.Sum(nil))
				thumbnailName := fileRename

				// 文件md5重命名
				dotIndex := strings.LastIndex(fileName, ".")
				if dotIndex != -1 && dotIndex != 0 {
					fileRename += fileName[dotIndex:]
					thumbnailName += "_100" + fileName[dotIndex:]
				}
				list = append(list, path+fileRename)

				err = os.Rename(outputFilePath, path+fileRename)
				if err != nil {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "第" + strconv.Itoa(index+1) + "文件重命名失败，请稍后再试"
					c.uploadEnd(data)
					return
				}
				if strings.HasSuffix(fileRename, "jpg") || strings.HasSuffix(fileRename, "jpeg") || strings.HasSuffix(fileRename, "png") {
					// decode jpeg into image.Image
					file.Seek(0, os.SEEK_SET)
					var img image.Image
					if strings.HasSuffix(fileRename, "png") {
						img, err = png.Decode(file)
					} else {
						img, err = jpeg.Decode(file)
					}

					if err != nil {
						data[models.STR_CODE] = models.CODE_ERR
						data[models.STR_MSG] = "文件上传失败"
						c.uploadEnd(data)
						return
					}

					// resize to width 100 using Lanczos resampling
					// and preserve aspect ratio
					m := resize.Resize(100, 0, img, resize.NearestNeighbor)

					out, err := os.Create(path + thumbnailName)
					if err != nil {
						log.Fatal(err)
					}
					defer out.Close()

					// write new image to file
					if strings.HasSuffix(fileRename, "png") {
						png.Encode(out, m)
					} else {
						jpeg.Encode(out, m, nil)
					}
				}
			} else {
				writer.Close()

				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "文件创建或打开失败"
				c.uploadEnd(data)
				return
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "上传文件打开失败"
			c.uploadEnd(data)
			return
		}
	}
	data[models.STR_DATA] = list

	c.uploadEnd(data)
}

// FileDownload ...
func (c *FileUploaderController) FileDownload() {
	filename := c.GetString("filename", "temp.txt")
	c.Ctx.Output.Download("./file/download/"+filename, filename)
}
