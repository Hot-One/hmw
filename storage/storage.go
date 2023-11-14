package storage

import "app/models"

type StorageI interface {
	Close()
	Category() CategoryRepoI
	Product() ProductRepoI
}

type CategoryRepoI interface {
	Create(*models.CreateCategory) (string, error)
	GetByID(*models.CategoryPrimaryKey) (*models.Category, error)
	GetList(*models.CategoryGetListRequest) (*models.CategoryGetListResponse, error)
	Update(req *models.UpdateCategory) (*models.CategoryPrimaryKey, error)
	Delete(req *models.CategoryPrimaryKey) error
}

type ProductRepoI interface {
	Create(req *models.ProductCreate) (string, error)
	GetByID(req *models.ProductPrimaryKey) (*models.Product, error)
	GetList(req *models.ProductGetListRequest) (*models.ProductGetListResponse, error)
	Update(req *models.ProductUpdate) (*models.ProductPrimaryKey, error)
	Delete(req *models.ProductPrimaryKey) error
}
