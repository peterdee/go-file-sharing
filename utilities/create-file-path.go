package utilities

import (
	"path/filepath"

	"file-sharing/constants"
)

func CreateFilePath(uid string) string {
	return filepath.Join(
		GetEnv(
			constants.ENV_NAMES.UplaodsDirectoryName,
			constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
		),
		uid,
	)
}
