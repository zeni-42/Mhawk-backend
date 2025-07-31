package cloudinary

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadOnCloudinary(file *multipart.FileHeader)( string, error ){
	path, err := os.Getwd()
	if err != nil {
		return "", nil
	}

	safeFileName := filepath.Base(file.Filename)
	localFilePath := filepath.Join(path + "public" + safeFileName)

	defer func () {
		if err := os.Remove(path+"/public/"+file.Filename); err != nil {
			log.Println("Failed to delete temp file:", err)
		}
	} ()

	Curl := os.Getenv("CLOUDINARY_URL")
	if Curl == "" {
		return "", errors.New("CLOUDINARY_URL is not set")
	}


	cld, err := cloudinary.NewFromURL(Curl)
	if err != nil {
		return "", err
	}

	var ctx = context.Background()
	res, err := cld.Upload.Upload(ctx, localFilePath, uploader.UploadParams{
		PublicID: "user_avatar_" + safeFileName,
	} )
	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}