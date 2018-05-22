package merchant

import (
	"testing"

	"example.com/test/pkg/resperror"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	name := "name"
	shortName := "shortname"
	dbaName := "dbaname"
	country := "AU"
	timezone := "tz"
	status := "active"

	testCases := []struct {
		req      *Request
		expected error
	}{
		{
			// Empty Data.
			&Request{},
			resperror.ValidationRequired("request data"),
		},
		{
			// No errors.
			&Request{&Data{
				Name:      name,
				ShortName: shortName,
				DBAName:   dbaName,
				Country:   country,
				Timezone:  timezone,
				Status:    status,
			}},
			nil,
		},

		{
			// No Name.
			&Request{&Data{
				ShortName: shortName,
				DBAName:   dbaName,
				Country:   country,
				Timezone:  timezone,
				Status:    status,
			}},
			resperror.ValidationRequired("name"),
		},
		{
			// No ShortName.
			&Request{&Data{
				Name:     name,
				DBAName:  dbaName,
				Country:  country,
				Timezone: timezone,
				Status:   status,
			}},
			resperror.ValidationRequired("short_name"),
		},
		{
			// No DBAName.
			&Request{&Data{
				Name:      name,
				ShortName: shortName,
				Country:   country,
				Timezone:  timezone,
				Status:    status,
			}},
			resperror.ValidationRequired("dba_name"),
		},
		{
			// No Country.
			&Request{&Data{
				Name:      name,
				ShortName: shortName,
				DBAName:   dbaName,
				Timezone:  timezone,
				Status:    status,
			}},
			resperror.ValidationRequired("country"),
		},
		{
			// No Timezone.
			&Request{&Data{
				Name:      name,
				ShortName: shortName,
				DBAName:   dbaName,
				Country:   country,
				Status:    status,
			}},
			resperror.ValidationRequired("timezone"),
		},
	}

	for _, tc := range testCases {
		// Validate the Request.
		err := tc.req.Validate()

		if err == nil {
			assert.Nil(tc.expected)
			continue
		}

		assert.Equal(tc.expected.Error(), err.Error())
	}
}
