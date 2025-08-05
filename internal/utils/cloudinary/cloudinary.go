package cloudinary

import (
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadOnCloudinary(file *multipart.FileHeader) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	safeFileName := filepath.Base(file.Filename)
	localFilePath := filepath.Join(path, "public", safeFileName)

	if err := os.MkdirAll(filepath.Join(path, "public"), os.ModePerm); err != nil {
		return "", err
	}

	if err := saveMultipartFile(file, localFilePath); err != nil {
		return "", err
	}

	defer func() {
		if err := os.Remove(localFilePath); err != nil {
			log.Println("Failed to delete temp file:", err)
		}
	}()

	Curl := os.Getenv("CLOUDINARY_URL")
	if Curl == "" {
		return "", errors.New("CLOUDINARY_URL is not set")
	}

	cld, err := cloudinary.NewFromURL(Curl)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	res, err := cld.Upload.Upload(ctx, localFilePath, uploader.UploadParams{
		PublicID: "user_avatar_" + safeFileName,
		ResourceType: "image",
		Folder: "Mhawk_images",
	})
	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}

func saveMultipartFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
