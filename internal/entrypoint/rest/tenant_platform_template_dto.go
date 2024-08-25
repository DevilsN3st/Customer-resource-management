package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type TenantPlatformTemplateDTO struct {
	URL           string `json:"url"`
	LoginName     string `json:"username"`
	LoginPassword string `json:"password"`
}

func mapTenantPlatformTemplateDTOToTenantPlatformTemplate(tenantPlatformTemplateDTO TenantPlatformTemplateDTO) (domain.TenantPlatformTemplate, error) {
	return domain.NewTenantPlatformTemplate(
		tenantPlatformTemplateDTO.URL,
		tenantPlatformTemplateDTO.LoginName,
		tenantPlatformTemplateDTO.LoginPassword,
	)
}
