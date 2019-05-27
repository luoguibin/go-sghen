package controllers

import (
	"go-sghen/helper"
	"go-sghen/models"

	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

// PeotryController operations for Peotry
type PeotryController struct {
	BaseController
}

func (c *PeotryController) QueryPeotry() {
	data := c.GetResponseData()
	params := &getQueryPeotryParams{}

	if c.CheckFormParams(data, params) {
		if params.ID > 0 {
			peotry, err := models.QueryPeotryByID(params.ID)
			if err == nil {
				data[models.STR_DATA] = peotry
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "未查询到对应id的诗歌"
			}
		} else {
			list, err, count, totalPage, curPage, pageIsEnd := models.QueryPeotry(params.SID, params.Page, params.Limit, params.Content)
			if err == nil {
				if params.NeedComment {
					for _, peotry := range list {
						comments, e := models.QueryCommentByTypeID(peotry.ID)
						if e == nil {
							peotry.Comments = comments
						}
					}
				}

				data[models.STR_DATA] = list
				data["totalCount"] = count
				data["totalPage"] = totalPage
				data["curPage"] = curPage
				data["pageIsEnd"] = pageIsEnd
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "未查询到对应的诗歌"
			}
		}
	}
	c.respToJSON(data)
}

func (c *PeotryController) CreatePeotry() {
	data := c.GetResponseData()
	params := &getCreatePeotryParams{}
	uLevel, _ := strconv.Atoi(c.Ctx.Input.Context.Request.Form.Get("level"))

	if c.CheckPostParams(data, params) {
		set, err := models.QueryPeotrySetByID(params.SID)
		if err == nil {
			if set.UID == params.UID {
				imgDatas := make([]string, 0)
				fileNames := make([]string, 0)
				errDatas := make([]string, 0)

				err := json.Unmarshal(c.Ctx.Input.RequestBody, &imgDatas)

				if err == nil {
					for index, imgData := range imgDatas {
						if index > 9 {
							data[models.STR_MSG] = "诗歌图片超过10张，只保存前10张"
							break
						}
						fileName, err := savePeotryimage(imgData)
						if err == nil {
							fileNames = append(fileNames, fileName)
						} else {
							msg := "第" + strconv.Itoa(index+1) + "张图片保存失败：" + err.Error()
							errDatas = append(errDatas, msg)
						}
					}
				} else {
					data[models.STR_MSG] = "请求成功，未添加图片"
				}

				fileNameByte, _ := json.Marshal(fileNames)

				timeStr := helper.GetNowDateTime()
				var flag int
				if uLevel < 5 {
					flag = 2
				} else {
					flag = 1
				}
				pId, err := models.CreatePeotry(params.UID, params.SID, params.Title, timeStr, params.Content, params.End, string(fileNameByte[:]), flag)

				if err == nil {
					if len(errDatas) == 0 {
						data[models.STR_DATA] = pId
					} else {
						data[models.STR_CODE] = models.CODE_ERR
						data[models.STR_MSG] = "保存诗歌图片失败"
						data[models.STR_DATA] = errDatas
					}
				} else {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "创建诗歌失败"
				}
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "禁止在他人选集中创建个人诗歌"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "非本人创建的选集下不能创建新的诗歌"
		}
	}

	c.respToJSON(data)
}

func (c *PeotryController) UpdatePeotry() {
	data := c.GetResponseData()
	params := &getUpdatePeotryParams{}

	if c.CheckPostParams(data, params) {
		qPeotry, err := models.QueryPeotryByID(params.PID)
		if err == nil {
			if qPeotry.UID == params.UID {
				// 判断选集是否有更新
				if qPeotry.SID != params.SID {
					set, err := models.QueryPeotrySetByID(params.SID)
					if err == nil {
						if set.UID == params.UID {
							qPeotry.SID = params.SID
						} else {
							data[models.STR_CODE] = models.CODE_ERR
							data[models.STR_MSG] = "禁止在他人选集中更新个人诗歌"
							c.respToJSON(data)
							return
						}
					} else {
						data[models.STR_CODE] = models.CODE_ERR
						data[models.STR_MSG] = "未获取到相应选集id"
						c.respToJSON(data)
						return
					}
				}
				qPeotry.PTitle = params.Title
				qPeotry.PContent = params.Content
				qPeotry.PEnd = params.End
				// 更新时需要将这些附带的结构体置空
				qPeotry.UUser = nil
				qPeotry.SSet = nil
				qPeotry.PImage = nil

				err := models.UpdatePeotry(qPeotry)
				if err == nil {
					data[models.STR_DATA] = qPeotry.ID
				} else {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "更新诗歌失败"
				}
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "禁止更新他人诗歌"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "未获取到相应诗歌id"
		}
	}

	c.respToJSON(data)
}

func (c *PeotryController) DeletePeotry() {
	data := c.GetResponseData()
	params := &getDeletePeotryParams{}

	if c.CheckFormParams(data, params) {
		peotry, err := models.QueryPeotryByID(params.PID)
		if err == nil {
			if peotry.UID == params.UID {
				// err := models.DeletePeotry(params.PID)
				// do not delete peotry in reality, just make it unselectable
				peotry.Flag = -1
				err := models.UpdatePeotry(peotry)
				if err == nil {
					data[models.STR_DATA] = true
				} else {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "删除诗歌失败"
				}
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "禁止删除他人诗歌"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "未获取到相应诗歌id"
		}
	}

	c.respToJSON(data)
}

func savePeotryimage(baseStr string) (string, error) {
	if len(baseStr) == 0 {
		return "", errors.New("空数据")
	}

	baseIndex := strings.Index(baseStr, "base64")
	if baseIndex < 15 {
		return "", errors.New("数据错误")
	}

	format := baseStr[11 : baseIndex-1]

	data, err := base64.StdEncoding.DecodeString(baseStr[baseIndex+7:])
	if err != nil {
		return "", err
	}

	path := models.MConfig.PathTypeMap["peotry"]
	isExist, err := helper.PathExists(path)
	if !isExist {
		isMade := helper.MkdirAll(path)
		if !isMade {
			return "", err
		}
	}

	h := md5.New()
	h.Write([]byte(baseStr))
	fileRename := hex.EncodeToString(h.Sum(nil))
	fileName := fileRename + "." + format
	err2 := ioutil.WriteFile(path+fileName, data, 0666)
	if err2 != nil {
		return "", err2
	}

	return fileName, nil
}
