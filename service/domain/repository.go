package domain

type AllRepository struct {
	CatClient     CatClient
	CatRepository CatRepository
	S3Client      S3Client
}
