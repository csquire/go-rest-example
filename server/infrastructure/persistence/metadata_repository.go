package persistence

import (
	"errors"

	"github.com/csquire/go-rest-example/domain/metadata"

	"github.com/jmoiron/sqlx"
)

func NewPostgresMetadataRepository(db *sqlx.DB) *PostgresMetadataRepository {
	return &PostgresMetadataRepository{
		db: db,
	}
}

type PostgresMetadataRepository struct {
	db *sqlx.DB
}

// CreateImage Create a new image metadata record
func (r *PostgresMetadataRepository) CreateImage(image *metadata.ImageMetadata) error {
	if image.Id == "" {
		return errors.New("id cannot be empty")
	}
	sqlStatement := `INSERT INTO metadata (id, name, base_image, approved) 
                     VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(sqlStatement, image.Id, image.Name, image.BaseImage, image.Approved)
	return err
}

// GetImage Gets one image metadata record by id
func (r *PostgresMetadataRepository) GetImage(id string) (*metadata.ImageMetadata, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	im := &metadata.ImageMetadata{}
	err := r.db.Get(im, "SELECT * FROM metadata WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return im, nil
}

// GetAllImages Gets all image metadata records in the database
// Caution: this will load all rows into memory, not suitable for a production situation
func (r *PostgresMetadataRepository) GetAllImages() ([]metadata.ImageMetadata, error) {
	ims := []metadata.ImageMetadata{}
	err := r.db.Select(&ims, "SELECT * FROM metadata")
	if err != nil {
		return nil, err
	}
	return ims, nil
}

// UpdateImage Updates an existing image metadata record
func (r *PostgresMetadataRepository) UpdateImage(image *metadata.ImageMetadata) error {
	sqlStatement := `UPDATE metadata
                     SET name = $2, base_image = $3, approved = $4
                     WHERE id = $1;`
	_, err := r.db.Exec(sqlStatement, image.Id, image.Name, image.BaseImage, image.Approved)
	return err
}
