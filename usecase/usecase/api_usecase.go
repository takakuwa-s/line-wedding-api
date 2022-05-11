package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
)

type ApiUsecase struct {
	ur igateway.IUserRepository
	lr igateway.ILineRepository
	fr igateway.IFileRepository
}

// Newコンストラクタ
func NewApiUsecase(ur igateway.IUserRepository, lr igateway.ILineRepository, fr igateway.IFileRepository) *ApiUsecase {
	return &ApiUsecase{ur: ur, lr: lr, fr: fr}
}

func (au *ApiUsecase) ValidateToken(token string) error {
	client := &http.Client{}
	resp, err := client.Get("https://api.line.me/oauth2/v2.1/verify?access_token=" + token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("token validation failed; %s", obj["error_description"])
	} else if obj["client_id"].(string) != os.Getenv("LIFF_CHANNEL_ID") {
		return fmt.Errorf("invalid access token audience")
	} else if obj["expires_in"].(float64) <= 0 {
		return fmt.Errorf("access token expired")
	}
	return nil
}

func (au *ApiUsecase) GetUser(id string) (*entity.User, error) {
	// Get user
	user, err := au.ur.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (au *ApiUsecase) UpdateUser(r *dto.UpdateUserRequest) (*entity.User, error) {
	// Get user
	user, err := au.ur.FindById(r.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = au.lr.GetUserProfileById(r.Id, dto.WeddingBotType)
		if err != nil {
			return nil, err
		}
	}
	user = r.ToUser(user)
	if err = au.ur.SaveUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (au *ApiUsecase) GetFileList(limit int, startAtId string) ([]entity.File, error) {
	if startAtId == "" {
		return au.fr.FindByLimit(limit)
	} else {
		return au.fr.FindByLimitAndStartAtId(limit, startAtId)
	}
}

func (au *ApiUsecase) DeleteFile(id string) error {
	return au.fr.DeleteFile(id, "api call")
}