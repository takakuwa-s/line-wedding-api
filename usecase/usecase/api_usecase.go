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
	br igateway.IBinaryRepository
}

// Newコンストラクタ
func NewApiUsecase(ur igateway.IUserRepository, lr igateway.ILineRepository, fr igateway.IFileRepository, br igateway.IBinaryRepository) *ApiUsecase {
	return &ApiUsecase{ur: ur, lr: lr, fr: fr, br: br}
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

func (au *ApiUsecase) GetFileList(limit int, startId string) ([]entity.File, error) {
	var files []entity.File
	var err error
	if startId == "" {
		files, err = au.fr.FindByLimit(limit)
	} else {
		files, err =  au.fr.FindByLimitAndStartId(limit, startId)
	}
	if err != nil {
		return nil, err
	}
	if files == nil {
		return []entity.File{}, nil
	}
	return files, nil
}

func (au *ApiUsecase) DeleteFile(id string) error {
	file, err := au.fr.FindById(id)
	if err != nil {
		return err
	}
	if file == nil {
		return fmt.Errorf("not found the file with id = %s", id)
	}
	if err := au.br.DeleteBinary(file.FileId); err != nil {
		return err
	}
	if err := au.fr.DeleteFile(id); err != nil {
		return err
	}
	return nil
}