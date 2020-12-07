package domain

type AllRepository struct {
	CatClient CatClient
	S3Client  S3Client
}
