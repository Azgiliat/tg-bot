package aws

import config "awesomeProject/internal/config"

var CatsClient *S3Bucket = nil

func GetCatsBucket() *S3Bucket {
	if CatsClient == nil {
		catsConfig := config.GetCatsBucketConfig()
		CatsClient = newClient(catsConfig.CatsBucketName, catsConfig.CatsBucketURL)
	}

	return CatsClient
}
