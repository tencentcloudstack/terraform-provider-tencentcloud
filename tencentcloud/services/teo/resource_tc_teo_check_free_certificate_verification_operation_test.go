package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestAccTeoCheckFreeCertificateVerification" -v -count=1 -gcflags="all=-l"

// TestAccTeoCheckFreeCertificateVerification_Success tests successful certificate verification
func TestAccTeoCheckFreeCertificateVerification_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CheckFreeCertificateVerification", func(request *teov20220901.CheckFreeCertificateVerificationRequest) (*teov20220901.CheckFreeCertificateVerificationResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "example.com", *request.Domain)

		resp := teov20220901.NewCheckFreeCertificateVerificationResponse()
		resp.Response = &teov20220901.CheckFreeCertificateVerificationResponseParams{
			CommonName:         ptrString("example.com"),
			SignatureAlgorithm: ptrString("RSA 2048"),
			ExpireTime:         ptrString("2025-12-31T23:59:59Z"),
			RequestId:          ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckFreeCertificateVerificationOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"domain":  "example.com",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#example.com", d.Id())
	assert.Equal(t, "example.com", d.Get("common_name"))
	assert.Equal(t, "RSA 2048", d.Get("signature_algorithm"))
	assert.Equal(t, "2025-12-31T23:59:59Z", d.Get("expire_time"))
}

// TestAccTeoCheckFreeCertificateVerification_APIError tests handling API errors
func TestAccTeoCheckFreeCertificateVerification_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CheckFreeCertificateVerification", func(request *teov20220901.CheckFreeCertificateVerificationRequest) (*teov20220901.CheckFreeCertificateVerificationResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckFreeCertificateVerificationOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"domain":  "example.com",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestAccTeoCheckFreeCertificateVerification_NilResponseFields tests handling nil response fields
func TestAccTeoCheckFreeCertificateVerification_NilResponseFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CheckFreeCertificateVerification", func(request *teov20220901.CheckFreeCertificateVerificationRequest) (*teov20220901.CheckFreeCertificateVerificationResponse, error) {
		resp := teov20220901.NewCheckFreeCertificateVerificationResponse()
		resp.Response = &teov20220901.CheckFreeCertificateVerificationResponseParams{
			CommonName:         nil,
			SignatureAlgorithm: nil,
			ExpireTime:         nil,
			RequestId:          ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckFreeCertificateVerificationOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"domain":  "example.com",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#example.com", d.Id())
	assert.Equal(t, "", d.Get("common_name"))
	assert.Equal(t, "", d.Get("signature_algorithm"))
	assert.Equal(t, "", d.Get("expire_time"))
}

// TestAccTeoCheckFreeCertificateVerification_Read tests Read is no-op
func TestAccTeoCheckFreeCertificateVerification_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCheckFreeCertificateVerificationOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"domain":  "example.com",
	})
	d.SetId("zone-12345678#example.com")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestAccTeoCheckFreeCertificateVerification_Delete tests Delete is no-op
func TestAccTeoCheckFreeCertificateVerification_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCheckFreeCertificateVerificationOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"domain":  "example.com",
	})
	d.SetId("zone-12345678#example.com")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestAccTeoCheckFreeCertificateVerification_Schema validates schema definition
func TestAccTeoCheckFreeCertificateVerification_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCheckFreeCertificateVerificationOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "domain")
	assert.Contains(t, res.Schema, "common_name")
	assert.Contains(t, res.Schema, "signature_algorithm")
	assert.Contains(t, res.Schema, "expire_time")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	domain := res.Schema["domain"]
	assert.Equal(t, schema.TypeString, domain.Type)
	assert.True(t, domain.Required)
	assert.True(t, domain.ForceNew)

	commonName := res.Schema["common_name"]
	assert.Equal(t, schema.TypeString, commonName.Type)
	assert.True(t, commonName.Computed)

	signatureAlgorithm := res.Schema["signature_algorithm"]
	assert.Equal(t, schema.TypeString, signatureAlgorithm.Type)
	assert.True(t, signatureAlgorithm.Computed)

	expireTime := res.Schema["expire_time"]
	assert.Equal(t, schema.TypeString, expireTime.Type)
	assert.True(t, expireTime.Computed)
}
