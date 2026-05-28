package initialize

import "testing"

func TestGenerateSwaggerDocs(t *testing.T) {
	if err := GenerateSwaggerDocs(); err != nil {
		t.Fatal(err)
	}
}
