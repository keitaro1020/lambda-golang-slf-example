package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"

	"github.com/keitaro1020/lambda-golang-slf-example/service/application"
	"github.com/keitaro1020/lambda-golang-slf-example/service/application/mocks"
)

func init() {

}

func Test_handler_SQSWorker(t *testing.T) {
	type args struct {
		ctx      context.Context
		sqsEvent events.SQSEvent
	}
	tests := []struct {
		name    string
		args    args
		app     func(ctrl *gomock.Controller) application.App
		wantErr bool
	}{
		{
			name: "ok_single",
			args: args{
				ctx: context.Background(),
				sqsEvent: events.SQSEvent{
					Records: []events.SQSMessage{
						{Body: "message_ok"},
					},
				},
			},
			app: func(ctrl *gomock.Controller) application.App {
				app := mocks.NewMockApp(ctrl)
				app.EXPECT().SQSWorker(gomock.Any(), "message_ok").Return(nil)
				return app
			},
			wantErr: false,
		},
		{
			name: "ok_multi",
			args: args{
				ctx: context.Background(),
				sqsEvent: events.SQSEvent{
					Records: []events.SQSMessage{
						{Body: "message_ok_1"},
						{Body: "message_ok_2"},
					},
				},
			},
			app: func(ctrl *gomock.Controller) application.App {
				app := mocks.NewMockApp(ctrl)
				app.EXPECT().SQSWorker(gomock.Any(), "message_ok_1").Return(nil)
				app.EXPECT().SQSWorker(gomock.Any(), "message_ok_2").Return(nil)
				return app
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				sqsEvent: events.SQSEvent{
					Records: []events.SQSMessage{
						{Body: "message_error"},
					},
				},
			},
			app: func(ctrl *gomock.Controller) application.App {
				app := mocks.NewMockApp(ctrl)
				app.EXPECT().SQSWorker(gomock.Any(), "message_error").Return(errors.New("error"))
				return app
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := &handler{
				app: tt.app(ctrl),
			}
			if err := h.SQSWorker(tt.args.ctx, tt.args.sqsEvent); (err != nil) != tt.wantErr {
				t.Errorf("SQSWorker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handler_S3Worker(t *testing.T) {
	type s3EventParam struct {
		bucketName string
		objectKey  string
	}
	s3Event := func(params []s3EventParam) events.S3Event {
		records := make([]events.S3EventRecord, len(params))
		for i := range params {
			records[i] = events.S3EventRecord{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{Name: params[i].bucketName},
					Object: events.S3Object{Key: params[i].objectKey},
				},
			}
		}
		return events.S3Event{
			Records: records,
		}
	}
	type args struct {
		ctx     context.Context
		s3Event events.S3Event
	}
	tests := []struct {
		name    string
		args    args
		app     func(ctrl *gomock.Controller) application.App
		wantErr bool
	}{
		{
			name: "ok_single",
			args: args{
				ctx:     context.Background(),
				s3Event: s3Event([]s3EventParam{{bucketName: "bucketName", objectKey: "objectKey"}}),
			},
			app: func(ctrl *gomock.Controller) application.App {
				app := mocks.NewMockApp(ctrl)
				app.EXPECT().S3Worker(gomock.Any(), "bucketName", "objectKey").Return(nil)
				return app
			},
			wantErr: false,
		},
		{
			name: "ok_multi",
			args: args{
				ctx: context.Background(),
				s3Event: s3Event([]s3EventParam{
					{bucketName: "bucketName_1", objectKey: "objectKey_1"},
					{bucketName: "bucketName_2", objectKey: "objectKey_2"},
				}),
			},
			app: func(ctrl *gomock.Controller) application.App {
				app := mocks.NewMockApp(ctrl)
				app.EXPECT().S3Worker(gomock.Any(), "bucketName_1", "objectKey_1").Return(nil)
				app.EXPECT().S3Worker(gomock.Any(), "bucketName_2", "objectKey_2").Return(nil)
				return app
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx:     context.Background(),
				s3Event: s3Event([]s3EventParam{{bucketName: "errorBucketName", objectKey: "errorObjectKey"}}),
			},
			app: func(ctrl *gomock.Controller) application.App {
				app := mocks.NewMockApp(ctrl)
				app.EXPECT().S3Worker(gomock.Any(), "errorBucketName", "errorObjectKey").Return(errors.New("error"))
				return app
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := &handler{
				app: tt.app(ctrl),
			}
			if err := h.S3Worker(tt.args.ctx, tt.args.s3Event); (err != nil) != tt.wantErr {
				t.Errorf("S3Worker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handler_Ping(t *testing.T) {
	type fields struct {
		app application.App
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Response
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				ctx: context.Background(),
			},
			want: Response{
				StatusCode:      200,
				IsBase64Encoded: false,
				Body:            "{\"message\":\"Okay so your other function also executed successfully!\"}",
				Headers: map[string]string{
					"Content-Type":           "application/json",
					"X-MyCompany-Func-Reply": "ping-cmd",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{
				app: tt.fields.app,
			}
			got, err := h.Ping(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ping() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//
//func Test_handler_GetCat(t *testing.T) {
//	type args struct {
//		ctx context.Context
//		req Request
//	}
//	tests := []struct {
//		name    string
//		args    args
//		app     func(ctrl *gomock.Controller) application.App
//		want    Response
//		wantErr bool
//	}{
//		{
//			name:    "ok_cat",
//			args:    args{},
//			app: func(ctrl *gomock.Controller) application.App {
//
//			},
//			want:    Response{
//
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			h := &handler{
//				app: tt.app(ctrl),
//			}
//			got, err := h.GetCat(tt.args.ctx, tt.args.req)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetCat() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetCat() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
