package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllFolders(t *testing.T) {
	t.Run("test nil request", func(t *testing.T) {
		res, err := folders.GetAllFolders(nil)
		if assert.Errorf(t, err, "GetAllFolders must fail") {
			assert.Nilf(t, res, "res must be nil")
		}
	})

	t.Run("test nil orgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.Nil,
		}

		res, err := folders.GetAllFolders(req)
		if assert.NoErrorf(t, err, "GetAllFolders must not fail") {
			assert.Emptyf(t, res.Folders, "res must be empty")
		}
	})

	t.Run("test valid orgId", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.Must(uuid.FromString(folders.DefaultOrgID)),
		}

		res, err := folders.GetAllFolders(req)
		if assert.NoErrorf(t, err, "GetAllFolders must not fail") {
			assert.NotEmptyf(t, res.Folders, "res must not be empty")

			for _, folder := range res.Folders {
				assert.NotNilf(t, folder, "folder must not be nil")
				assert.Equalf(
					t, req.OrgID, folder.OrgId,
					"folder orgID must be equal to request orgID",
				)
			}
		}
	})
}
