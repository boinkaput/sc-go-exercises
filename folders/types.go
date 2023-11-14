package folders

import "github.com/gofrs/uuid"

type FetchFolderRequest struct {
	OrgID uuid.UUID
}

type FetchFolderResponse struct {
	Folders []*Folder
}

type PaginationRequest struct {
	OrgID uuid.UUID
	MaxFolders int
	Token string
}

type PaginationResponse struct {
	Folders []*Folder
	NumFolders int
	Token string
}
