package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type tenantService struct {
	tenantRepository domain.TenantRepository
}

type TenantService interface {
	Create(ctx context.Context, tenant domain.Tenant) (string, error)
	GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error)
	Update(ctx context.Context, tenantID string, tenant domain.UpdateTenant) error
	Delete(ctx context.Context, tenantID string) error
	Search(ctx context.Context, filters domain.TenantFilters) (domain.PagingResult[domain.Tenant], error)
}

func NewTenantService(tenantRepository domain.TenantRepository) TenantService {
	return &tenantService{
		tenantRepository: tenantRepository,
	}
}

func (s *tenantService) Create(ctx context.Context, tenant domain.Tenant) (string, error) {
	return s.tenantRepository.Create(ctx, tenant)
}

func (s *tenantService) Delete(ctx context.Context, tenantID string) error {
	if tenantID == "" {
		return domain.NewValidationError("tenantID cannot be empty", nil)
	}

	return s.tenantRepository.Delete(ctx, tenantID)
}

func (s *tenantService) GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	if tenantID == "" {
		return nil, domain.NewValidationError("tenantID cannot be empty", nil)
	}

	return s.tenantRepository.GetByID(ctx, tenantID)
}

func (s *tenantService) Search(ctx context.Context, filters domain.TenantFilters) (domain.PagingResult[domain.Tenant], error) {
	return s.tenantRepository.Search(ctx, filters)
}

func (s *tenantService) Update(ctx context.Context, tenantID string, updateTenant domain.UpdateTenant) error {
	if tenantID == "" {
		return domain.NewValidationError("tenantID cannot be empty", nil)
	}

	tenant, err := s.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}

	tenant.MergeUpdate(updateTenant)

	return s.tenantRepository.Update(ctx, *tenant)
}
