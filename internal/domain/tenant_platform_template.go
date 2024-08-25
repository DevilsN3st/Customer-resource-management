package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TenantPlatformTemplate struct {
	TenantPlatformTemplateID string
	URL                      string
	LoginName                string
	LoginPassword            string
	Fields                   map[string]map[string]string
}

func NewTenantPlatformTemplate(url, loginName, loginPassword string) (TenantPlatformTemplate, error) {
	tenantPlatformTemplateID, err := uuid.NewRandom()
	if err != nil {
		return TenantPlatformTemplate{}, err
	}

	encryptedLogin, err := encryptData(loginName)
	if err != nil {
		return TenantPlatformTemplate{}, err
	}
	encryptedPassword, err := encryptData(loginPassword)
	if err != nil {
		return TenantPlatformTemplate{}, err
	}

	return TenantPlatformTemplate{
		TenantPlatformTemplateID: tenantPlatformTemplateID.String(),
		URL:                      url,
		LoginName:                encryptedLogin,
		LoginPassword:            encryptedPassword,
	}, nil
}

func encryptData(input string) (string, error) {
	encryptedInput, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedInput), nil
}
