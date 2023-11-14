package folders

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Improvements made to the code:
// 1. Removed unused variables and error check all return values from function calls.
// 2. There is a bug in the 2nd for loop in GetAllFolders. A pointer to v1 is being
//    appended to the slice, since v1's value changes, all the pointers in fp
//    effectively point to the last element of f. I fixed this by making a copy
//    of v1 and appending the address of the copy (copy := v1). copy is a variable
//    stored on the stack. A pointer to copy is valid because the golang compiler
//    performs escape analysis, which may promote copy to the heap making it a
//    valid reference.
// 3. The 1st for loop in GetAllFolders is redundant, copy can store the dereferenced value
//    of v instead (copy := *v) and we can append this to fp.
// 4. Used the make function for pre-allocating the fp slice. This saves the time spent on
//    re-allocations. We no longer need to append, and can simply assign the pointer instead.
// 5. Skipped over any nil pointer elements in FetchAllFoldersByOrgID.
// 6. Merged variable declaration and assignment of ffr. Once again ffr is a valid pointer
//    due to go compiler's escape analysis.
// 7. Changed variable names.

// The GetAllFolders function takes in a pointer to FetchFolderRequest as input
// and returns a pointer to FetchFolderResponse and error.
// - Calls the FetchAllFoldersByOrgID function with the OrgId from the request,
//   which returns a slice of pointers to folders.
// - Iterates over the result and copies the folders pointed to by the
//   elements of the result to slice f.
// - Iterates over the slice f and appends a pointer to each of its elements to
//   fp.
// - Constructs a FetchFolderResponse with fp and returns it along with
//   nil as the error.
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	if req == nil {
		return nil, errors.New("req cannot be nil")
	}

	orgFolders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, errors.New("FetchAllFoldersByOrgID failed")
	}

	foldersCopy := make([]*Folder, len(orgFolders))
	for idx, folder := range orgFolders {
		folderCopy := *folder
		foldersCopy[idx] = &folderCopy
	}

	fetchResponse := &FetchFolderResponse{Folders: foldersCopy}
	return fetchResponse, nil
}

// The FetchAllFoldersByOrgID takes in a orgID of type UUID as input
// and returns a slice of pointers to folders that belong to the orgID.
// - Gets a slice of pointers from calling GetSampleData.
// - Appends the folders that have the same orgID as the input to
//   resFolder which itself is a slice of pointers to folder.
// - Returns resFolder and nil as error.
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	folders := GetSampleData()

	resFolders := []*Folder{}
	for _, folder := range folders {
		if folder != nil && folder.OrgId == orgID {
			resFolders = append(resFolders, folder)
		}
	}

	return resFolders, nil
}
