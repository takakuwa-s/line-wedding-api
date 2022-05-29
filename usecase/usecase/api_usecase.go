package usecase

import (
	"fmt"

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

func (au *ApiUsecase) GetInitialData(id string) (*dto.InitApiResponse, error) {
	// Get user
	user, err := au.ur.FindById(id)
	if err != nil {
		return nil, err
	}
	// Get file list
	files, err := au.GetFileList(12, "", "", "")
	if err != nil {
		return nil, err
	}
	return dto.NewInitApiResponse(user, files), nil
}

func (au *ApiUsecase) UpdateUser(r *dto.UpdateUserRequest) (*entity.User, error) {
	// Get user
	user, err := au.ur.FindById(r.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = au.lg.GetUserProfileById(r.Id)
		if err != nil {
			return nil, err
		}
	}
	user = r.ToUser(user)
	user.Registered = true
	if err = au.ur.SaveUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (au *ApiUsecase) PatchUser(userId, field string, val bool) error {
	// Check if user exists
	if _, err := au.GetUser(userId); err != nil {
		return err
	}

	return au.ur.UpdateBoolFieldById(userId, field, val)
}

func (au *ApiUsecase) GetUser(id string) (*entity.User, error) {
	// Get user
	user, err := au.ur.FindById(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("not found the user with id = %s", id)
	}
	return user, nil
}

func (au *ApiUsecase) GetUsers(limit int, startId, flag string, val bool) ([]entity.User, error) {
	// Get users
	users, err := au.ur.FindByFlagOrderByName(limit, startId, flag, val)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (au *ApiUsecase) GetFileList(limit int, startId, userId, orderBy string) ([]entity.File, error) {
	if orderBy == "" {
		orderBy = "UpdatedAt"
	}
	files, err := au.fr.FindByLimitAndStartIdAndUserId(limit, startId, userId, orderBy)
	if err != nil {
		return nil, err
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
	if err := au.fr.DeleteFileById(id); err != nil {
		return err
	}
	if err := au.br.DeleteBinary(file.Name); err != nil {
		return err
	}
	return nil
}

func (au *ApiUsecase) DeleteFileList(ids []string) error {
	files, err := au.fr.FindByIds(ids)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("not found the files with ids = %s", ids)
	}
	if len(files) != len(ids) {
		return fmt.Errorf("some ids is invalid and not found; ids = %s", ids)
	}
	for i, f := range files {
		if err := au.fr.DeleteFileById(f.Id); err != nil {
			return fmt.Errorf("successfully delete %d files, but failed to delete the file metadata of id = %s", i, f.Id)
		}
		if err := au.br.DeleteBinary(f.Name); err != nil {
			return fmt.Errorf("successfully delete %d files, but failed to delete the file binary of id = %s", i, f.Id)
		}
	}
	return nil
}
