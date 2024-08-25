package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/icrxz/crm-api-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type tenantRepository struct {
	client *mongo.Client
}

func NewTenantRepository(client *mongo.Client) domain.TenantRepository {
	return &tenantRepository{
		client: client,
	}
}

func (db *tenantRepository) tenantCollection(ctx context.Context) *mongo.Collection {
	userCollection := GetCollection(db.client, "user")
	return userCollection
}

func (db *tenantRepository) Create(ctx context.Context, tenant domain.Tenant) (string, error) {
	tenantDTO := mapTenantToTenantDTO(tenant)

	_, err := db.tenantCollection(ctx).InsertOne(ctx, tenantDTO)

	if err != nil {
		return "", err
	}

	return tenant.TenantID, nil
}

func (db *tenantRepository) Delete(ctx context.Context, tenantID string) error {
	if tenantID == "" {
		return domain.NewValidationError("tenantID is required", map[string]any{"tenant_id": tenantID})
	}

	_, err := db.tenantCollection(ctx).DeleteOne(ctx, bson.D{{"_id", tenantID}})
	if err != nil {
		return err
	}

	return nil
}

func (db *tenantRepository) GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	var tenantDTO TenantDTO
	err := db.tenantCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", tenantID}}).Decode(&tenantDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no tenant found with this id", map[string]any{"tenant_id": tenantID})
		}
		return nil, err
	}

	tenant := mapTenantDTOToTenant(tenantDTO)

	return &tenant, nil
}

func (db *tenantRepository) Search(ctx context.Context, filters domain.TenantFilters) (domain.PagingResult[domain.Tenant], error) {

	filter := bson.M{}
	if len(filters.TenantID) > 0 {
		filter["_id"] = filters.TenantID
	}
	if len(filters.CompanyName) > 0 {
		filter["company_name"] = filters.CompanyName
	}
	if len(filters.Document) > 0 {
		filter["document"] = filters.Document
	}

	var tenantsResult TenantDTO
	err := db.tenantCollection(ctx).FindOne(ctx, filter).Decode(&tenantsResult)
	if err != nil {
		return domain.PagingResult[domain.Tenant]{}, err
	}
	tenants := mapTenantDTOsToTenants([]TenantDTO{tenantsResult})

	result := domain.PagingResult[domain.Tenant]{
		Result: tenants,
		Paging: domain.Paging{
			Total:  len(tenants),
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (db *tenantRepository) Update(ctx context.Context, tenant domain.Tenant) error {
	tenantDTO := mapTenantToTenantDTO(tenant)

	filter := bson.M{
		"_id": tenant.TenantID,
	}

	result := db.tenantCollection(ctx).FindOneAndUpdate(ctx, filter, tenantDTO)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
