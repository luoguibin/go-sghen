package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"go-sghen/helper"
	"go-sghen/models"
	"io"
	"os"
	"strconv"
	"strings"
)

type FileUploaderController struct {
	BaseController
}

func (c *FileUploaderController) FileUpload() {
	data := c.GetResponseData()

	// 上传的文件存储在maxMemory大小的内存里面，如果文件大小超过了maxMemory，那么剩下的部分将存储在系统的临时文件中
	c.Ctx.Request.ParseMultipartForm(32 << 20)
	// c.GetFile("file")	// 单文件
	fileHeaders, err := c.GetFiles("file") // 多文件

	if err != nil {
		// fmt.Println("file upload getFiles() err")
		// fmt.Println(err)
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "获取上传文件列表失败"
		c.respToJSON(data)
		return
	}

	// 检测上传目录是否存在
	pathType := c.GetString("pathType")
	if len(pathType) == 0 {
		pathType = "normal"
	}
	path, ok := models.MConfig.PathTypeMap[pathType]
	if !ok {
		path = models.MConfig.PathTypeMap["peotry"]
	}
	isExist, _ := helper.PathExists(path)
	if !isExist {
		isMade := helper.MkdirAll(path)
		if !isMade {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "创建文件夹失败"
			c.respToJSON(data)
			return
		}
	}

	if len(fileHeaders) > models.MConfig.MaxUploadCount {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "上传文件个数不能超过" + strconv.Itoa(models.MConfig.MaxUploadCount) + "个"
		c.respToJSON(data)
		return
	}

	// 限制大文件上传
	for index, v := range fileHeaders {
		if len(v.Filename) <= 0 {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "上传文件不能为空"
			c.respToJSON(data)
			return
		}
		sizeMB := int(v.Size / 1024 / 1024)
		if sizeMB > models.MConfig.MaxUploadSize {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "第" + strconv.Itoa(index+1) + "个文件:" + v.Filename + "文件大小超过" + strconv.Itoa(models.MConfig.MaxUploadSize) + "MB"
			c.respToJSON(data)
			return
		}
	}

	list := make([]string, 0)

	for index, v := range fileHeaders {
		fileName := v.Filename
		file, err := v.Open()
		defer file.Close()

		if err == nil {
			outputFilePath := path + fileName
			writer, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0666)

			if err == nil {
				io.Copy(writer, file)
				writer.Close()
				file.Seek(0, os.SEEK_SET)
				// 文件md5计算
				h := md5.New()
				io.Copy(h, file)
				fileRename := hex.EncodeToString(h.Sum(nil))

				// 文件md5重命名
				dotIndex := strings.LastIndex(fileName, ".")
				if dotIndex != -1 && dotIndex != 0 {
					fileRename += fileName[dotIndex:]
				}
				list = append(list, path+fileRename)

				err = os.Rename(outputFilePath, path+fileRename)
				if err != nil {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "第" + strconv.Itoa(index+1) + "文件重命名失败，请稍后再试"
					c.respToJSON(data)
					return
				}
			} else {
				writer.Close()

				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "文件创建或打开失败"
				c.respToJSON(data)
				return
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "上传文件打开失败"
			c.respToJSON(data)
			return
		}
	}
	if pathType == "peotry" {
		pId, err := c.GetInt64("pId")
		uId := c.Ctx.Input.GetData("uId").(int64)
		if err == nil && uId > 0 {
			err := checkSavePeotryImage(pId, uId, list)
			if err != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = err.Error()
				c.respToJSON(data)
				return
			}
		}
		data[models.STR_DATA] = list
	} else {
		data[models.STR_DATA] = true
	}

	c.respToJSON(data)
}

func checkSavePeotryImage(pId int64, uId int64, imgList []string) error {
	peotry, err := models.QueryPeotryByID(pId)
	if err == nil {
		if peotry.UID == uId {
			l := len(imgList)
			imgsByte, err := json.Marshal(imgList)
			if err != nil {
				return err
			}

			err = models.SavePeotryImage(pId, string(imgsByte[:]), l)
			if err != nil {
				return errors.New("更新诗歌的图片失败")
			} else {
				return err
			}
		} else {
			return errors.New("禁止更新非本人创建的诗歌的图片")
		}
	}
	return errors.New("未获取到对应id的诗歌")
}
