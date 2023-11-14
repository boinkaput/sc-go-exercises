package folders

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Advantages of this design:
// - User gets to choose size of each pagination chunk.
// - Constant amortized time for adding, accessing
//   and deleting token - paginationState mappings.
// - Using pointers to avoid copying.
// - Chances of getting the same token for 2 different
//   chunks is extremely small with the use of UUID.

// Note: PaginationRequest and PaginationResponse types can be found in types.go

// Example:
// req := &folders.PaginationRequest{
// 	OrgID: uuid.FromStringOrNil(folders.DefaultOrgID),
// 	MaxFolders: 5,
// }

// for resp, err := folders.FoldersPagination(req); true; resp, err = folders.FoldersPagination(req) {
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 		return
// 	}
// 	folders.PrettyPrint(resp)
// 	req.Token = resp.Token
// 	if (req.Token == "") {
// 		break
// 	}
// }

type PaginationState struct {
	Folders []*Folder
	nextIdx int
}

var paginationMap = make(map[string]*PaginationState)

func FoldersPagination(req *PaginationRequest) (*PaginationResponse, error) {
	var pageState *PaginationState

	// Token would be an empty string for the first request.
	if req.Token == "" {
		req := &FetchFolderRequest{
			OrgID: req.OrgID,
		}
	
		// Get all folders with matching orgID
		res, err := GetAllFolders(req)
		if err != nil {
			return nil, errors.New("failed to get folders for the given orgID")
		}

		pageState = &PaginationState{
			Folders: res.Folders,
			nextIdx: 0,
		}
	} else {
		// Check if a mapping for the token exists in paginationMap
		var exists bool
		pageState, exists = paginationMap[req.Token]
		if !exists {
			return nil, errors.New("invalid token")
		}

		// Delete this entry a new entry would later be created for the next chunk.
		delete(paginationMap, req.Token)
	}

	return GetNextChunk(pageState, req.MaxFolders)
}

func GetNextChunk(pageState *PaginationState, maxFolders int) (*PaginationResponse, error) {
	// Chunk size would be the minimum of the remaining folders
	// in the slice and maxFolders given by the user.
	chunkSize := min(len(pageState.Folders) - pageState.nextIdx, maxFolders)
	folders := make([]*Folder, chunkSize)

	// Copy the folders (number of elements copied would be chunkSize).
	copy(folders, pageState.Folders[pageState.nextIdx:])
	pageState.nextIdx += chunkSize

	// Add a new entry in the paginationMap if the slice has not been exhausted yet.
	token, err := SaveState(pageState)
	if err != nil {
		return nil, err
	}

	resp := &PaginationResponse{
		Folders: folders,
		NumFolders: chunkSize,
		Token: token,
	}

	return resp, nil
}

func SaveState(pageState *PaginationState) (string, error) {
	if pageState.nextIdx >= len(pageState.Folders) {
		return "", nil
	}

	// Generate a uuid for the token.
	tokenUid, err := uuid.NewV4()
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	// Convert the uuid to a string to be used as a key in the map.
	token := tokenUid.String()
	paginationMap[token] = pageState

	return token, nil
}
