package application

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/domain"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/domain/mocks"
)

func Test_app_SQSWorker(t *testing.T) {
	type args struct {
		ctx     context.Context
		message string
	}
	tests := []struct {
		name    string
		repos   func(ctrl *gomock.Controller) *domain.AllRepository
		config  *Config
		args    args
		wantErr bool
	}{
		{
			name: "ok_1件",
			repos: func(ctrl *gomock.Controller) *domain.AllRepository {
				catClient := mocks.NewMockCatClient(ctrl)
				catClient.EXPECT().Search(gomock.Any()).Return(domain.Cats{
					{ID: "123", URL: "http://hoge/huga", Width: 100, Height: 200},
				}, nil)

				file := bytes.NewBufferString(`{"id":"123","url":"http://hoge/huga","width":100,"height":200}`)
				s3Client := mocks.NewMockS3Client(ctrl)
				s3Client.EXPECT().Upload(
					gomock.Any(),
					"bucket_name",
					"message/123.txt",
					file,
				).Return(nil)

				return &domain.AllRepository{
					CatClient: catClient,
					S3Client:  s3Client,
				}
			},
			config: &Config{BucketName: "bucket_name"},
			args: args{
				ctx:     context.Background(),
				message: "message",
			},
			wantErr: false,
		},
		{
			name: "ok_2件",
			repos: func(ctrl *gomock.Controller) *domain.AllRepository {
				catClient := mocks.NewMockCatClient(ctrl)
				catClient.EXPECT().Search(gomock.Any()).Return(domain.Cats{
					{ID: "123", URL: "http://hoge/huga", Width: 100, Height: 200},
					{ID: "234", URL: "http://hoge/huga2", Width: 300, Height: 400},
				}, nil)

				file1 := bytes.NewBufferString(`{"id":"123","url":"http://hoge/huga","width":100,"height":200}`)
				file2 := bytes.NewBufferString(`{"id":"234","url":"http://hoge/huga2","width":300,"height":400}`)
				s3Client := mocks.NewMockS3Client(ctrl)
				s3Client.EXPECT().Upload(
					gomock.Any(), "bucket_name", "message/123.txt", file1,
				).Return(nil)
				s3Client.EXPECT().Upload(
					gomock.Any(), "bucket_name", "message/234.txt", file2,
				).Return(nil)

				return &domain.AllRepository{
					CatClient: catClient,
					S3Client:  s3Client,
				}
			},
			config: &Config{BucketName: "bucket_name"},
			args: args{
				ctx:     context.Background(),
				message: "message",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			app := &app{
				repos:  tt.repos(ctrl),
				config: tt.config,
			}
			if err := app.SQSWorker(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("SQSWorker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_app_S3Worker(t *testing.T) {
	type args struct {
		ctx      context.Context
		bucket   string
		filename string
	}
	tests := []struct {
		name    string
		repos   func(ctrl *gomock.Controller) *domain.AllRepository
		config  *Config
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			repos: func(ctrl *gomock.Controller) *domain.AllRepository {
				s3Client := mocks.NewMockS3Client(ctrl)
				s3Client.EXPECT().Download(gomock.Any(), "bucket_name", "file_name").Return([]byte(`{"id":"123","url":"http://hoge/huga","width":100,"height":200}`), nil)

				catRepository := mocks.NewMockCatRepository(ctrl)
				catRepository.EXPECT().CreateInTx(gomock.Any(), gomock.Any(), &domain.Cat{
					ID:     "123",
					URL:    "http://hoge/huga",
					Width:  100,
					Height: 200,
				}).Return(nil, nil)

				return &domain.AllRepository{
					S3Client:      s3Client,
					CatRepository: catRepository,
					Transaction: func(ctx context.Context, txFunc func(ctx context.Context, tx domain.Tx) error) (err error) {
						return txFunc(ctx, nil)
					},
				}
			},
			config: &Config{BucketName: "bucket_name"},
			args: args{
				ctx:      context.Background(),
				bucket:   "bucket_name",
				filename: "file_name",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			app := &app{
				repos:  tt.repos(ctrl),
				config: tt.config,
			}
			if err := app.S3Worker(tt.args.ctx, tt.args.bucket, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("S3Worker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_app_GetCat(t *testing.T) {
	type args struct {
		ctx context.Context
		id  domain.CatID
	}
	tests := []struct {
		name    string
		repos   func(ctrl *gomock.Controller) *domain.AllRepository
		config  *Config
		args    args
		want    *domain.Cat
		wantErr bool
	}{
		{
			name: "ok",
			repos: func(ctrl *gomock.Controller) *domain.AllRepository {
				catRepository := mocks.NewMockCatRepository(ctrl)
				catRepository.EXPECT().Get(gomock.Any(), domain.CatID("123")).Return(&domain.Cat{
					ID:     "123",
					URL:    "http://hoge/huga",
					Width:  100,
					Height: 200,
				}, nil)
				return &domain.AllRepository{CatRepository: catRepository}
			},
			config: &Config{BucketName: "bucket_name"},
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			want: &domain.Cat{
				ID:     "123",
				URL:    "http://hoge/huga",
				Width:  100,
				Height: 200,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			app := &app{
				repos:  tt.repos(ctrl),
				config: tt.config,
			}
			got, err := app.GetCat(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_app_GetCats(t *testing.T) {
	type args struct {
		ctx   context.Context
		first int64
	}
	tests := []struct {
		name    string
		repos   func(ctrl *gomock.Controller) *domain.AllRepository
		config  *Config
		args    args
		want    domain.Cats
		wantErr bool
	}{
		{
			name: "ok",
			repos: func(ctrl *gomock.Controller) *domain.AllRepository {
				catRepository := mocks.NewMockCatRepository(ctrl)
				catRepository.EXPECT().GetAll(gomock.Any(), int64(100)).Return(domain.Cats{
					{ID: "123", URL: "http://hoge/huga", Width: 100, Height: 200},
					{ID: "234", URL: "http://hoge/huga2", Width: 300, Height: 400},
				}, nil)
				return &domain.AllRepository{CatRepository: catRepository}

			},
			config: &Config{BucketName: "bucket_name"},
			args: args{
				ctx:   context.Background(),
				first: 100,
			},
			want: domain.Cats{
				{ID: "123", URL: "http://hoge/huga", Width: 100, Height: 200},
				{ID: "234", URL: "http://hoge/huga2", Width: 300, Height: 400},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			app := &app{
				repos:  tt.repos(ctrl),
				config: tt.config,
			}
			got, err := app.GetCats(tt.args.ctx, tt.args.first)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCats() got = %v, want %v", got, tt.want)
			}
		})
	}
}
