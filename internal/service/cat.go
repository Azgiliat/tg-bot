package service

import (
	aws2 "awesomeProject/internal/aws"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/types"
	"bytes"
	"fmt"
	"net/url"
	"time"
)

type CatService struct {
	repo       *repository.CatRepository
	catsBucket *aws2.S3Bucket
}

func NewCatService(repo *repository.CatRepository) *CatService {
	return &CatService{repo, aws2.GetCatsBucket()}
}

func (service *CatService) GetCatPhotosByTag(tag string) *types.CatPhoto {
	return service.repo.GetCatByTag(tag)
}

func (service *CatService) GetAvailableCats() []string {
	return service.repo.GetAvailableCats()
}

func generateCatImageKey(catName string, catExtension string) string {
	return fmt.Sprintf("%s_%s.%s", url.PathEscape(catName), time.Now().Format("20060102150405"), catExtension)
}

func (service *CatService) StoreCat(cat string, extension string, image []byte) error {
	catImageKey := generateCatImageKey(cat, extension)
	err := service.catsBucket.UploadImage(catImageKey, bytes.NewReader(image))

	if err != nil {
		return err
	}

	err = service.repo.StoreCat(cat, service.catsBucket.GenerateImageURL(catImageKey))

	return err
}
