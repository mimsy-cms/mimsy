package sync

type StatusResponse struct {
	Statuses   []SyncStatus `json:"statuses"`
	Repository string       `json:"repository"`
}

func NewStatusResponse(statuses []SyncStatus, repository string) *StatusResponse {
	return &StatusResponse{
		Statuses:   statuses,
		Repository: repository,
	}
}
