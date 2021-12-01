package metadata

type Repository interface {
	CreateImage(image *ImageMetadata) error
	GetImage(id string) (*ImageMetadata, error)
	GetAllImages() ([]ImageMetadata, error)
	UpdateImage(image *ImageMetadata) error
}
