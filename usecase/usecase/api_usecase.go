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
	lg igateway.ILineGateway
	fr igateway.IFileRepository
	br igateway.IBinaryRepository
}

// Newコンストラクタ
func NewApiUsecase(ur igateway.IUserRepository, lg igateway.ILineGateway, fr igateway.IFileRepository, br igateway.IBinaryRepository) *ApiUsecase {
	return &ApiUsecase{ur: ur, lg: lg, fr: fr, br: br}
}

func (au *ApiUsecase) ValidateToken(token string) error {
	client := &http.Client{}
	url := os.Getenv("LIFF_API_BASE_URL")
	resp, err := client.Get(url + "?access_token=" + token)
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
		user, err = au.lg.GetUserProfileById(r.Id, dto.WeddingBotType)
		if err != nil {
			return nil, err
		}
	}
	user = r.ToUser(user)
	user.IsRegistered = true
	if err = au.ur.SaveUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (au *ApiUsecase) GetFileList(limit int, startId, userId, orderBy string) ([]entity.File, error) {
	files, err :=  au.fr.FindByLimitAndStartIdAndUserId(limit, startId, userId, orderBy)
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
	if err := au.br.DeleteBinary(file.Name); err != nil {
		return err
	}
	if err := au.fr.DeleteFile(id); err != nil {
		return err
	}
	return nil
}