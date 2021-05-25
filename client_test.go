package jsonplaceholder

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"gopkg.in/h2non/gock.v1"
)

func TestNew(t *testing.T) {
	httpClient := &http.Client{Timeout: time.Second * 30}

	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Default client",
			args: args{[]Option{}},
			want: func() *Client {
				c := &Client{
					httpClient: &http.Client{Timeout: time.Second * 10},
				}

				c.Post = newPostService(c)
				c.Comment = newCommentService(c)
				c.Photo = newPhotoService(c)
				c.Album = newAlbumService(c)
				c.Todo = newTodoService(c)
				c.User = newUserService(c)
				return c
			}(),
		},
		{
			name: "Custom client",
			args: args{opts: []Option{WithHTTPClient(httpClient)}},
			want: func() *Client {
				c := &Client{
					httpClient: httpClient,
				}

				c.Post = newPostService(c)
				c.Comment = newCommentService(c)
				c.Photo = newPhotoService(c)
				c.Album = newAlbumService(c)
				c.Todo = newTodoService(c)
				c.User = newUserService(c)
				return c
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeURL(t *testing.T) {
	type args struct {
		params []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty params",
			args: args{},
			want: apiURL,
		},
		{
			name: "One URI",
			args: args{[]interface{}{"foo"}},
			want: fmt.Sprintf("%s/foo", apiURL),
		},
		{
			name: "Two URI",
			args: args{[]interface{}{"foo", "bar"}},
			want: fmt.Sprintf("%s/foo/bar", apiURL),
		},
		{
			name: "URI with int",
			args: args{[]interface{}{"foo", "bar", 1}},
			want: fmt.Sprintf("%s/foo/bar/1", apiURL),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeURL(tt.args.params...); got != tt.want {
				t.Errorf("makeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_makeURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		want := apiURL + "/foo"
		if got := makeURL("foo"); got != want {
			b.Errorf("makeURL() = %v, want %v", got, want)
		}
	}
}

func TestClient_doFetch(t *testing.T) {
	defer gock.Off() // Disable HTTP interceptors

	// Invalid status code
	gock.New("https://foo.bar").
		Get("/status").
		Reply(500)

	// Invalid response body received
	gock.New("https://foo.bar").
		Get("/invalid").
		Persist().
		Reply(200).
		BodyString("foo")

	// Valid response body received
	gock.New("https://foo.bar").
		Get("/valid").
		Persist().
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	type fields struct {
		httpClient *http.Client
		Post       PostService
		Comment    CommentService
		Album      AlbumService
		Photo      PhotoService
		Todo       TodoService
		User       UserService
	}

	type args struct {
		ctx  context.Context
		url  string
		data interface{}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Nil context",
			fields: fields{
				httpClient: &http.Client{},
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "Invalid URL",
			fields: fields{
				httpClient: &http.Client{},
			},
			args: args{
				ctx:  context.Background(),
				url:  " https://foo.com",
				data: nil,
			},
			wantErr: true,
		},
		{
			name: "Missing URL",
			fields: fields{
				httpClient: &http.Client{},
			},
			args: args{
				ctx:  context.Background(),
				url:  "",
				data: nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid status code received",
			fields: fields{
				httpClient: &http.Client{},
			},
			args: args{
				ctx:  context.Background(),
				url:  "https://foo.bar/status",
				data: nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid response body",
			fields: fields{
				httpClient: &http.Client{},
			},
			args: args{
				ctx:  context.Background(),
				url:  "https://foo.bar/invalid",
				data: nil,
			},
			wantErr: true,
		},
		{
			name: "No pointer received",
			fields: fields{
				httpClient: &http.Client{},
			},
			args: args{
				ctx:  context.Background(),
				url:  "https://foo.bar/valid",
				data: *new(string),
			},
			wantErr: true,
		},
		{
			name: "Nil data provided",
			fields: fields{
				httpClient: &http.Client{},
			},
			args: args{
				ctx:  context.Background(),
				url:  "https://foo.bar/valid",
				data: nil,
			},
			wantErr: true,
		},
		{
			name: "Valid response body",
			fields: fields{
				httpClient: &http.Client{},
			},
			args: args{
				ctx:  context.Background(),
				url:  "https://foo.bar/valid",
				data: &map[string]string{"foo": "bar"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient: tt.fields.httpClient,
				Post:       tt.fields.Post,
				Comment:    tt.fields.Comment,
				Album:      tt.fields.Album,
				Photo:      tt.fields.Photo,
				Todo:       tt.fields.Todo,
				User:       tt.fields.User,
			}
			if err := c.doFetch(tt.args.ctx, tt.args.url, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("doFetch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
