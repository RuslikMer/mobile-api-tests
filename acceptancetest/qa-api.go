package acceptancetest

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func DeleteUserById(t *testing.T, userId string) {
	data := url.Values{
		"accessKey": {os.Getenv("CL_QA_API_ACCESS_KEY")},
		"user_id":   {userId},
		"apiEnv":    {"stage"},
	}

	resp, err := http.PostForm(os.Getenv("CL_QA_API")+"user/delete/", data)

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}
