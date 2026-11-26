package ssl

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func cert(id, domain string, subjectNames ...string) *ssl.Certificates {
	c := &ssl.Certificates{
		CertificateId: helper.String(id),
		Domain:        helper.String(domain),
	}
	for _, name := range subjectNames {
		c.SubjectAltName = append(c.SubjectAltName, helper.String(name))
	}
	return c
}

func TestSslCertificatesDataSourceSchema(t *testing.T) {
	res := DataSourceTencentCloudSslCertificates()

	domain, ok := res.Schema["domain"]
	assert.True(t, ok, "domain should be a top-level argument")
	assert.Equal(t, schema.TypeString, domain.Type)
	assert.True(t, domain.Optional)

	sortBy, ok := res.Schema["sort_by"]
	assert.True(t, ok, "sort_by should be a top-level argument")
	assert.True(t, sortBy.Optional)
	assert.NotNil(t, sortBy.ValidateFunc)

	_, errs := sortBy.ValidateFunc("end_time", "sort_by")
	assert.Empty(t, errs)
	_, errs = sortBy.ValidateFunc("begin_time", "sort_by")
	assert.NotEmpty(t, errs, "only end_time is supported")

	_, ok = res.Schema["order_by"]
	assert.False(t, ok, "order_by is not part of this feature")
}

func TestCertificateMatchesDomain(t *testing.T) {
	c := cert("cert-1", "example.com", "www.example.com", "api.example.com")

	assert.True(t, certificateMatchesDomain(c, "example.com"), "primary domain")
	assert.True(t, certificateMatchesDomain(c, "www.example.com"), "subject alternative name")
	assert.True(t, certificateMatchesDomain(c, "API.EXAMPLE.COM"), "case insensitive")
	assert.False(t, certificateMatchesDomain(c, "example.org"))
	assert.False(t, certificateMatchesDomain(c, "sub.www.example.com"), "exact match only, no suffix match")
	assert.False(t, certificateMatchesDomain(nil, "example.com"))
}

func TestGetCommonCertificates(t *testing.T) {
	a := cert("a", "a.example.com")
	b := cert("b", "b.example.com")
	c := cert("c", "c.example.com")

	// No filter: single list is returned as-is, order preserved.
	assert.Equal(t, []*ssl.Certificates{a, b}, GetCommonCertificates([]*ssl.Certificates{a, b}))

	// Empty list short-circuits to no results.
	assert.Empty(t, GetCommonCertificates([]*ssl.Certificates{}, []*ssl.Certificates{a}))

	// Two lists: intersection by certificate ID, order of the first list.
	got := GetCommonCertificates([]*ssl.Certificates{a, b, c}, []*ssl.Certificates{c, a})
	assert.Equal(t, []*ssl.Certificates{a, c}, got)

	// Three lists (id + name + domain filters all set): keep what all agree on.
	got = GetCommonCertificates([]*ssl.Certificates{a, b, c}, []*ssl.Certificates{b, c}, []*ssl.Certificates{c, b})
	assert.Equal(t, []*ssl.Certificates{b, c}, got)

	got = GetCommonCertificates([]*ssl.Certificates{a, b, c}, []*ssl.Certificates{b, c}, []*ssl.Certificates{c})
	assert.Equal(t, []*ssl.Certificates{c}, got)

	// Duplicated entries in the first list collapse.
	got = GetCommonCertificates([]*ssl.Certificates{a, a}, []*ssl.Certificates{a})
	assert.Equal(t, []*ssl.Certificates{a}, got)
}

func TestBuildDescribeCertificatesRequest(t *testing.T) {
	// No filters: nothing but the default request is set.
	req := buildDescribeCertificatesRequest(nil, nil, nil)
	assert.Nil(t, req.CertificateType)
	assert.Nil(t, req.SearchKey)
	assert.Nil(t, req.ExpirationSort, "API default order is kept when sort_by is omitted")

	// sort_by = end_time maps to ExpirationSort=DESC.
	req = buildDescribeCertificatesRequest(helper.String("SVR"), helper.String("example.com"), helper.String(SSL_CERTIFICATE_SORT_BY_END_TIME))
	assert.Equal(t, "SVR", *req.CertificateType)
	assert.Equal(t, "example.com", *req.SearchKey)
	assert.Equal(t, "DESC", *req.ExpirationSort)
}
