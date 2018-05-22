package merchant

import (
	"github.com/vegh1010/test/pkg/resperror"
)

// Validate valiates merchant request Data.
func (req *Request) Validate() error {
	// First check if data is present.
	if req.Data == nil {
		return resperror.ValidationRequired("request data")
	}

	if req.Data.Name == "" {
		return resperror.ValidationRequired("name")
	}
	if req.Data.ShortName == "" {
		return resperror.ValidationRequired("short_name")
	}
	if req.Data.DBAName == "" {
		return resperror.ValidationRequired("dba_name")
	}
	if req.Data.Country == "" {
		return resperror.ValidationRequired("country")
	}
	if req.Data.Timezone == "" {
		return resperror.ValidationRequired("timezone")
	}

	return nil
}
