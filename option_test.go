package jsonplaceholder

import (
	"net/http"
	"reflect"
	"testing"
)

func TestWithHTTPClient(t *testing.T) {

	defaultClient := new(Client) // actually the defaultClient.httpClient is nil

	t.Run("HTTP client should be changed", func(t *testing.T) {
		want := &http.Client{}
		WithHTTPClient(want)(defaultClient)

		got := defaultClient.httpClient
		if !reflect.DeepEqual(got, want) {
			t.Errorf("WithHTTPClient() = %v, want %v", got, want)
		}
	})
}
