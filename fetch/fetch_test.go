package fetch

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"

	"github.com/pkg/errors"
)

func TestFetchUrlObject(t *testing.T) {
	defer gock.Off()

	gock.New("http://image.om").
		Get("/bar.jpg").
		Reply(200).
		BodyString("foo foo")

	gock.InterceptClient(client)
	buf, err := FetchFile("http://image.om/bar.jpg")

	assert.Nil(t, err)
	assert.Equal(t, string(buf), "foo foo")
}

func TestFetchUrlObjectErr(t *testing.T) {
	defer gock.Off()

	gock.New("http://image.om").
		Get("/bar.jpg").
		ReplyError(errors.New("error"))

	gock.InterceptClient(client)
	_, err := FetchFile("http://image.om/bar.jpg")

	assert.NotNil(t, err)
}
