package petclinic_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseUrl = "http://petclinic.127.0.0.1.nip.io"
const registerOwnerRequest = `
{
    "id": "SR334451",
    "salutation": "Mr",
    "surname": "Smith",
    "name": "John",
    "email": "foobar@baz.com",
    "phone": "333-444-5555"
}`

func TestRegisterOwner(t *testing.T) {
	client := http.Client{}
	res, err := client.Post(baseUrl+"/v1.0/invoke/owner.petclinic/method/register", "application/json", strings.NewReader(registerOwnerRequest))
	require.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	fmt.Println(res.Header)
}
