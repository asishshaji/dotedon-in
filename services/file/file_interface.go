package file_service

import (
	"context"
	"mime/multipart"
)

type IFileService interface {
	UploadFile(ctx context.Context, folderName string, file multipart.File) (string, error)
}
