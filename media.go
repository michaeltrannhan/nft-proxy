package nft_proxy

type Media struct {
	ID        uint   `json:"-" gorm:"primaryKey"`
	Mint      string `json:"mint" gorm:"uniqueIndex"`
	ImageUri  string `json:"imageUri"`
	ImageType string `json:"imageType"`
	MediaUri  string `json:"mediaUri,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
	LocalPath string `json:"-"`
}

type SolanaMedia struct {
	ID        uint   `json:"-" gorm:"primaryKey"`
	Mint      string `json:"mint" gorm:"uniqueIndex"`
	ImageUri  string `json:"imageUri"`
	ImageType string `json:"ImageType"`
	MediaUri  string `json:"mediaUri"`
	MediaType string `json:"mediaType"`
	LocalPath string `json:"-"`
}

func (m *SolanaMedia) Media() *Media {
	return &Media{
		ID:        m.ID,
		Mint:      m.Mint,
		ImageUri:  m.ImageUri,
		ImageType: m.ImageType,
		MediaUri:  m.MediaUri,
		MediaType: m.MediaType,
		LocalPath: m.LocalPath,
	}
}
