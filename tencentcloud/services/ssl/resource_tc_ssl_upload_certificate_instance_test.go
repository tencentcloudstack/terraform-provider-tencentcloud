package ssl_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcssl "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssl"
)

// go test ./tencentcloud/services/ssl/ -run "TestSslUploadCertificateInstanceResource" -v -count=1 -gcflags="all=-l"

type mockMetaUploadCertInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaUploadCertInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaUploadCertInstance{}

func newMockMetaUploadCertInstance() *mockMetaUploadCertInstance {
	return &mockMetaUploadCertInstance{client: &connectivity.TencentCloudClient{}}
}

func ptrStrUpload(s string) *string    { return &s }
func ptrUint64Upload(n uint64) *uint64 { return &n }
func ptrInt64Upload(n int64) *int64    { return &n }

// TestSslUploadCertificateInstanceResource_CreateSuccess tests create success with deployRecordId>0
func TestSslUploadCertificateInstanceResource_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sslClient := &ssl.Client{}
	patches.ApplyMethodReturn(newMockMetaUploadCertInstance().client, "UseSSLCertificateClient", sslClient)

	// mock UploadUpdateCertificateInstance returns DeployRecordId>0
	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstance", func(request *ssl.UploadUpdateCertificateInstanceRequest) (*ssl.UploadUpdateCertificateInstanceResponse, error) {
		resp := ssl.NewUploadUpdateCertificateInstanceResponse()
		resp.Response = &ssl.UploadUpdateCertificateInstanceResponseParams{
			DeployRecordId: ptrUint64Upload(12345),
			DeployStatus:   ptrInt64Upload(1),
			RequestId:      ptrStrUpload("fake-request-id"),
		}
		return resp, nil
	})

	// mock DescribeHostUploadUpdateRecord returns the matching record
	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecord", func(request *ssl.DescribeHostUploadUpdateRecordRequest) (*ssl.DescribeHostUploadUpdateRecordResponse, error) {
		resp := ssl.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &ssl.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64Upload(1),
			DeployRecordList: []*ssl.UploadUpdateRecordInfo{
				{
					Id:            ptrUint64Upload(12345),
					OldCertId:     ptrStrUpload("f8xnynbd"),
					ResourceTypes: []*string{ptrStrUpload("clb")},
					Status:        ptrUint64Upload(1),
					CreateTime:    ptrStrUpload("2024-01-01 00:00:00"),
					UpdateTime:    ptrStrUpload("2024-01-01 00:00:00"),
				},
			},
		}
		return resp, nil
	})

	// mock DescribeHostUploadUpdateRecordDetail returns detail
	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordDetail", func(request *ssl.DescribeHostUploadUpdateRecordDetailRequest) (*ssl.DescribeHostUploadUpdateRecordDetailResponse, error) {
		resp := ssl.NewDescribeHostUploadUpdateRecordDetailResponse()
		resp.Response = &ssl.DescribeHostUploadUpdateRecordDetailResponseParams{
			DeployRecordDetail: []*ssl.DeployRecord{
				{
					TotalCount:        ptrInt64Upload(1),
					SuccessTotalCount: ptrInt64Upload(1),
					Status:            ptrInt64Upload(1),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaUploadCertInstance()
	res := svcssl.ResourceTencentCloudSslUploadCertificateInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"old_certificate_id":      "f8xnynbd",
		"resource_types":          schema.NewSet(schema.HashString, []interface{}{"clb"}),
		"certificate_public_key":  "public-key-content",
		"certificate_private_key": "private-key-content",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "12345", d.Id())
	assert.Equal(t, 12345, d.Get("deploy_record_id"))
	assert.Equal(t, "f8xnynbd", d.Get("old_certificate_id"))
	assert.Equal(t, 1, d.Get("deploy_status"))
}

// TestSslUploadCertificateInstanceResource_CreateRetry tests deployRecordId=0 continues retry
func TestSslUploadCertificateInstanceResource_CreateRetry(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sslClient := &ssl.Client{}
	patches.ApplyMethodReturn(newMockMetaUploadCertInstance().client, "UseSSLCertificateClient", sslClient)

	callCount := 0
	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstance", func(request *ssl.UploadUpdateCertificateInstanceRequest) (*ssl.UploadUpdateCertificateInstanceResponse, error) {
		callCount++
		resp := ssl.NewUploadUpdateCertificateInstanceResponse()
		if callCount < 2 {
			// first call: still in progress
			resp.Response = &ssl.UploadUpdateCertificateInstanceResponseParams{
				DeployRecordId: ptrUint64Upload(0),
				DeployStatus:   ptrInt64Upload(0),
				RequestId:      ptrStrUpload("fake-request-id"),
			}
		} else {
			// second call: success
			resp.Response = &ssl.UploadUpdateCertificateInstanceResponseParams{
				DeployRecordId: ptrUint64Upload(12345),
				DeployStatus:   ptrInt64Upload(1),
				RequestId:      ptrStrUpload("fake-request-id"),
			}
		}
		return resp, nil
	})

	// mock read returns matching record
	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecord", func(request *ssl.DescribeHostUploadUpdateRecordRequest) (*ssl.DescribeHostUploadUpdateRecordResponse, error) {
		resp := ssl.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &ssl.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64Upload(1),
			DeployRecordList: []*ssl.UploadUpdateRecordInfo{
				{
					Id:        ptrUint64Upload(12345),
					OldCertId: ptrStrUpload("f8xnynbd"),
				},
			},
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordDetail", func(request *ssl.DescribeHostUploadUpdateRecordDetailRequest) (*ssl.DescribeHostUploadUpdateRecordDetailResponse, error) {
		resp := ssl.NewDescribeHostUploadUpdateRecordDetailResponse()
		resp.Response = &ssl.DescribeHostUploadUpdateRecordDetailResponseParams{
			DeployRecordDetail: []*ssl.DeployRecord{
				{
					Status: ptrInt64Upload(1),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaUploadCertInstance()
	res := svcssl.ResourceTencentCloudSslUploadCertificateInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"old_certificate_id":      "f8xnynbd",
		"resource_types":          schema.NewSet(schema.HashString, []interface{}{"clb"}),
		"certificate_public_key":  "public-key-content",
		"certificate_private_key": "private-key-content",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "12345", d.Id())
	assert.GreaterOrEqual(t, callCount, 1)
}

// TestSslUploadCertificateInstanceResource_ReadNotFound tests empty record list clears id
func TestSslUploadCertificateInstanceResource_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sslClient := &ssl.Client{}
	patches.ApplyMethodReturn(newMockMetaUploadCertInstance().client, "UseSSLCertificateClient", sslClient)

	// mock DescribeHostUploadUpdateRecord returns empty
	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecord", func(request *ssl.DescribeHostUploadUpdateRecordRequest) (*ssl.DescribeHostUploadUpdateRecordResponse, error) {
		resp := ssl.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &ssl.DescribeHostUploadUpdateRecordResponseParams{
			DeployRecordList: []*ssl.UploadUpdateRecordInfo{},
		}
		return resp, nil
	})

	meta := newMockMetaUploadCertInstance()
	res := svcssl.ResourceTencentCloudSslUploadCertificateInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"old_certificate_id": "f8xnynbd",
	})
	d.SetId("12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestSslUploadCertificateInstanceResource_ReadSuccess tests read populates fields
func TestSslUploadCertificateInstanceResource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sslClient := &ssl.Client{}
	patches.ApplyMethodReturn(newMockMetaUploadCertInstance().client, "UseSSLCertificateClient", sslClient)

	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecord", func(request *ssl.DescribeHostUploadUpdateRecordRequest) (*ssl.DescribeHostUploadUpdateRecordResponse, error) {
		resp := ssl.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &ssl.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64Upload(1),
			DeployRecordList: []*ssl.UploadUpdateRecordInfo{
				{
					Id:            ptrUint64Upload(12345),
					OldCertId:     ptrStrUpload("f8xnynbd"),
					ResourceTypes: []*string{ptrStrUpload("clb")},
					Status:        ptrUint64Upload(1),
					CreateTime:    ptrStrUpload("2024-01-01 00:00:00"),
					UpdateTime:    ptrStrUpload("2024-01-01 00:00:00"),
				},
			},
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordDetail", func(request *ssl.DescribeHostUploadUpdateRecordDetailRequest) (*ssl.DescribeHostUploadUpdateRecordDetailResponse, error) {
		resp := ssl.NewDescribeHostUploadUpdateRecordDetailResponse()
		resp.Response = &ssl.DescribeHostUploadUpdateRecordDetailResponseParams{
			DeployRecordDetail: []*ssl.DeployRecord{
				{
					Status: ptrInt64Upload(1),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaUploadCertInstance()
	res := svcssl.ResourceTencentCloudSslUploadCertificateInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"old_certificate_id": "f8xnynbd",
	})
	d.SetId("12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "12345", d.Id())
	assert.Equal(t, 12345, d.Get("deploy_record_id"))
	assert.Equal(t, 1, d.Get("deploy_status"))

	deployRecordList := d.Get("deploy_record_list").([]interface{})
	assert.Len(t, deployRecordList, 1)
	recordMap := deployRecordList[0].(map[string]interface{})
	assert.Equal(t, "f8xnynbd", recordMap["old_cert_id"])
	assert.Equal(t, 1, recordMap["status"])
}

// TestSslUploadCertificateInstanceResource_Delete tests delete calls rollback
func TestSslUploadCertificateInstanceResource_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sslClient := &ssl.Client{}
	patches.ApplyMethodReturn(newMockMetaUploadCertInstance().client, "UseSSLCertificateClient", sslClient)

	var capturedDeployRecordId *int64
	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateRecordRollback", func(request *ssl.UploadUpdateCertificateRecordRollbackRequest) (*ssl.UploadUpdateCertificateRecordRollbackResponse, error) {
		capturedDeployRecordId = request.DeployRecordId
		resp := ssl.NewUploadUpdateCertificateRecordRollbackResponse()
		resp.Response = &ssl.UploadUpdateCertificateRecordRollbackResponseParams{
			RequestId: ptrStrUpload("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaUploadCertInstance()
	res := svcssl.ResourceTencentCloudSslUploadCertificateInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("12345")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedDeployRecordId)
	assert.Equal(t, int64(12345), *capturedDeployRecordId)
}

// TestSslUploadCertificateInstanceResource_Schema validates schema definition
func TestSslUploadCertificateInstanceResource_Schema(t *testing.T) {
	res := svcssl.ResourceTencentCloudSslUploadCertificateInstance()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "old_certificate_id")
	assert.Contains(t, res.Schema, "resource_types")
	assert.Contains(t, res.Schema, "certificate_public_key")
	assert.Contains(t, res.Schema, "certificate_private_key")
	assert.Contains(t, res.Schema, "resource_types_regions")
	assert.Contains(t, res.Schema, "certificate_id")
	assert.Contains(t, res.Schema, "deploy_record_id")
	assert.Contains(t, res.Schema, "deploy_status")
	assert.Contains(t, res.Schema, "update_sync_progress")
	assert.Contains(t, res.Schema, "deploy_record_list")

	oldCertId := res.Schema["old_certificate_id"]
	assert.Equal(t, schema.TypeString, oldCertId.Type)
	assert.True(t, oldCertId.Required)

	resourceTypes := res.Schema["resource_types"]
	assert.Equal(t, schema.TypeSet, resourceTypes.Type)
	assert.True(t, resourceTypes.Required)

	publicKey := res.Schema["certificate_public_key"]
	assert.Equal(t, schema.TypeString, publicKey.Type)
	assert.True(t, publicKey.Required)
	assert.True(t, publicKey.Sensitive)

	deployRecordId := res.Schema["deploy_record_id"]
	assert.Equal(t, schema.TypeInt, deployRecordId.Type)
	assert.True(t, deployRecordId.Computed)
}
