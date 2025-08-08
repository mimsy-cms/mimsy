package media

type MediaResponse struct {
	Id           int64  `json:"id"`
	Uuid         string `json:"uuid"`
	Name         string `json:"name"`
	ContentType  string `json:"content_type"`
	CreatedAt    string `json:"created_at"`
	Size         int64  `json:"size"`
	UploadedById int64  `json:"uploaded_by_id"`
	URL          string `json:"url,omitempty"` // Optional URL for the media file
}

func NewMediaResponse(media *Media) MediaResponse {
	return MediaResponse{
		Id:           media.Id,
		Uuid:         media.Uuid.String(),
		Name:         media.Name,
		ContentType:  media.ContentType,
		CreatedAt:    media.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Size:         media.Size,
		UploadedById: media.UploadedById,
	}
}
