package controller

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"go.uber.org/zap"
)

type UserApiController struct {
	au *usecase.ApiUsecase
}

// コンストラクタ
func NewUserApiController(au *usecase.ApiUsecase) *UserApiController {
	return &UserApiController{au: au}
}

func (uac *UserApiController) UpdateUser(c *gin.Context) {
	var request dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		conf.Log.Error("[UpdateUser] json convestion failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uac.au.UpdateUser(&request)
	if err != nil {
		conf.Log.Error("[UpdateUser] Updating user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uac *UserApiController) PatchUser(c *gin.Context) {
	var request dto.PatchUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		conf.Log.Error("[UpdateUser] json convestion failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	field, val, err := request.GetFieldAndVal()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uac.au.PatchUser(id, field, val); err != nil {
		conf.Log.Error("[PatchUser] Updating user field failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (uac *UserApiController) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	user, err := uac.au.GetUser(id)
	if err != nil {
		conf.Log.Error("[GetUser] Failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uac *UserApiController) GetUserList(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit is required and must be number"})
		return
	}
	startId := c.Query("startId")
	flag := c.Query("flag")
	var val bool
	if flag != "" {
		val, err = strconv.ParseBool(c.Query("val"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "val must be boolean"})
			return
		}
	}
	csvDownloadStr := c.Query("csvDownload")
	var csvDownload bool
	if csvDownloadStr != "" {
		csvDownload, err = strconv.ParseBool(csvDownloadStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "csvDownload must be boolean"})
			return
		}
	}
	users, err := uac.au.GetUsers(limit, startId, flag, val)
	if err != nil {
		conf.Log.Error("[GetUserList] Failed to get user list", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if csvDownload {
		data := uac.convertToCsv(users)
		c.Stream(func(w io.Writer) bool {
			c.Writer.Header().Set("Content-Disposition", "attachment; filename=user-list.csv")
			c.Writer.Header().Set("Content-Type", "attachment; text/csv")
			cw := csv.NewWriter(w)
			cw.WriteAll(data)
			cw.Flush()
			if err := cw.Error(); err != nil {
				conf.Log.Error("[GetUserList] Failed to download user list csv", zap.Error(err))
			}
			return false
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func (uac *UserApiController) convertToCsv(users []entity.User) [][]string {
	var data [][]string
	data = append(data, []string{
		"#",
		"LINE名",
		"名前",
		"かな",
		"管理者",
		"出席",
		"回答済",
		"フォロー",
		"ゲスト",
		"メールアドレス",
		"郵便番号",
		"住所",
		"アレルギー",
		"お連れ様1 名前",
		"お連れ様2 アレルギー",
		"お連れ様1 名前",
		"お連れ様2 アレルギー",
		"メッセージ",
		"管理メモ",
	})
	for i, u := range users {
		var guest string
		switch u.GuestType {
		case "GROOM":
			guest = "新郎側"
		case "BRIDE":
			guest = "新婦側"
		case "COMMON":
			guest = "共通"
		default:
			guest = ""
		}
		companionName1 := ""
		companionAllergy1 := ""
		companionName2 := ""
		companionAllergy2 := ""
		if len(u.Companions) > 0 {
			companionName1 = u.Companions[0].Name
			companionAllergy1 = u.Companions[0].Allergy
			companionName2 = u.Companions[1].Name
			companionAllergy2 = u.Companions[1].Allergy
		}

		data = append(data, []string{
			strconv.Itoa(i + 1),
			u.LineName,
			u.Name,
			u.NameKana,
			uac.convertBoolToStr(u.IsAdmin),
			uac.convertBoolToStr(u.Attendance),
			uac.convertBoolToStr(u.Registered),
			uac.convertBoolToStr(u.Follow),
			guest,
			u.Email,
			u.PostalCode,
			u.Address,
			u.Allergy,
			companionName1,
			companionAllergy1,
			companionName2,
			companionAllergy2,
			u.Message,
			u.Note,
		})
	}
	return data
}

func (uac *UserApiController) convertBoolToStr(b bool) string {
	if b {
		return "○"
	} else {
		return "×"
	}
}
