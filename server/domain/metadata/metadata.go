package metadata

type ImageMetadata struct {
	Id        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	BaseImage string `json:"baseImage" db:"base_image"`
	Approved  bool   `json:"approved" db:"approved"`
}
