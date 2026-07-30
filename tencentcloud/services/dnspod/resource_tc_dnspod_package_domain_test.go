package dnspod_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	dnspodService "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dnspod"
)

type mockMetaPackageDomain struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaPackageDomain) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaPackageDomain{}

func newMockMetaPackageDomain() *mockMetaPackageDomain {
	return &mockMetaPackageDomain{client: &connectivity.TencentCloudClient{}}
}

func ptrStr(s string) *string {
	return &s
}

func ptrUint64(i uint64) *uint64 {
	return &i
}

func ptrBool(b bool) *bool {
	return &b
}

func TestDnspodPackageDomainResource_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dnspodClient := &dnspod.Client{}
	patches.ApplyMethodReturn(newMockMetaPackageDomain().client, "UseDnsPodClient", dnspodClient)

	patches.ApplyMethodFunc(dnspodClient, "ModifyPackageDomainWithContext",
		func(_ context.Context, request *dnspod.ModifyPackageDomainRequest) (*dnspod.ModifyPackageDomainResponse, error) {
			assert.NotNil(t, request.Operation)
			assert.Equal(t, "bind", *request.Operation)
			assert.NotNil(t, request.ResourceId)
			assert.Equal(t, "res-test001", *request.ResourceId)
			assert.NotNil(t, request.NewDomainId)
			assert.Equal(t, uint64(12345), *request.NewDomainId)

			resp := dnspod.NewModifyPackageDomainResponse()
			resp.Response = &dnspod.ModifyPackageDomainResponseParams{
				RequestId: ptrStr("fake-request-id-create"),
			}
			return resp, nil
		})

	patches.ApplyMethodFunc(dnspodClient, "DescribeDomainVipListWithContext",
		func(_ context.Context, request *dnspod.DescribeDomainVipListRequest) (*dnspod.DescribeDomainVipListResponse, error) {
			resp := dnspod.NewDescribeDomainVipListResponse()
			resp.Response = &dnspod.DescribeDomainVipListResponseParams{
				TotalCount: ptrUint64(1),
				PackageList: []*dnspod.PackageListItem{
					{
						DomainId:      ptrUint64(12345),
						Domain:        ptrStr("example.com"),
						Grade:         ptrStr("DPG_PROFESSIONAL"),
						GradeTitle:    ptrStr("Professional"),
						VipStartAt:    ptrStr("2024-01-01 00:00:00"),
						VipEndAt:      ptrStr("2025-01-01 00:00:00"),
						VipAutoRenew:  ptrStr("YES"),
						RemainTimes:   ptrUint64(5),
						ResourceId:    ptrStr("res-test001"),
						GradeLevel:    ptrUint64(1),
						Status:        ptrStr("active"),
						IsGracePeriod: ptrStr("NO"),
						Downgrade:     ptrBool(false),
					},
				},
			}
			return resp, nil
		})

	meta := newMockMetaPackageDomain()
	res := dnspodService.ResourceTencentCloudDnspodPackageDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"resource_id": "res-test001",
		"domain_id":   12345,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "res-test001#12345", d.Id())
}

func TestDnspodPackageDomainResource_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dnspodClient := &dnspod.Client{}
	patches.ApplyMethodReturn(newMockMetaPackageDomain().client, "UseDnsPodClient", dnspodClient)

	patches.ApplyMethodFunc(dnspodClient, "DescribeDomainVipListWithContext",
		func(_ context.Context, request *dnspod.DescribeDomainVipListRequest) (*dnspod.DescribeDomainVipListResponse, error) {
			assert.NotNil(t, request.ResourceIdList)
			assert.Equal(t, 1, len(request.ResourceIdList))
			assert.Equal(t, "res-test001", *request.ResourceIdList[0])

			resp := dnspod.NewDescribeDomainVipListResponse()
			resp.Response = &dnspod.DescribeDomainVipListResponseParams{
				TotalCount: ptrUint64(1),
				PackageList: []*dnspod.PackageListItem{
					{
						DomainId:      ptrUint64(12345),
						Domain:        ptrStr("example.com"),
						Grade:         ptrStr("DPG_PROFESSIONAL"),
						GradeTitle:    ptrStr("Professional"),
						VipStartAt:    ptrStr("2024-01-01 00:00:00"),
						VipEndAt:      ptrStr("2025-01-01 00:00:00"),
						VipAutoRenew:  ptrStr("YES"),
						RemainTimes:   ptrUint64(5),
						ResourceId:    ptrStr("res-test001"),
						GradeLevel:    ptrUint64(1),
						Status:        ptrStr("active"),
						IsGracePeriod: ptrStr("NO"),
						Downgrade:     ptrBool(false),
					},
				},
			}
			return resp, nil
		})

	meta := newMockMetaPackageDomain()
	res := dnspodService.ResourceTencentCloudDnspodPackageDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"resource_id": "res-test001",
		"domain_id":   12345,
	})
	d.SetId("res-test001#12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "example.com", d.Get("domain"))
	assert.Equal(t, "DPG_PROFESSIONAL", d.Get("grade"))
	assert.Equal(t, "Professional", d.Get("grade_title"))
	assert.Equal(t, "2024-01-01 00:00:00", d.Get("vip_start_at"))
	assert.Equal(t, "2025-01-01 00:00:00", d.Get("vip_end_at"))
	assert.Equal(t, "YES", d.Get("vip_auto_renew"))
	assert.Equal(t, 5, d.Get("remain_times"))
	assert.Equal(t, 1, d.Get("grade_level"))
	assert.Equal(t, "active", d.Get("status"))
	assert.Equal(t, "NO", d.Get("is_grace_period"))
	assert.Equal(t, false, d.Get("downgrade"))
}

func TestDnspodPackageDomainResource_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dnspodClient := &dnspod.Client{}
	patches.ApplyMethodReturn(newMockMetaPackageDomain().client, "UseDnsPodClient", dnspodClient)

	patches.ApplyMethodFunc(dnspodClient, "DescribeDomainVipListWithContext",
		func(_ context.Context, request *dnspod.DescribeDomainVipListRequest) (*dnspod.DescribeDomainVipListResponse, error) {
			resp := dnspod.NewDescribeDomainVipListResponse()
			resp.Response = &dnspod.DescribeDomainVipListResponseParams{
				TotalCount:  ptrUint64(0),
				PackageList: []*dnspod.PackageListItem{},
			}
			return resp, nil
		})

	meta := newMockMetaPackageDomain()
	res := dnspodService.ResourceTencentCloudDnspodPackageDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"resource_id": "res-test001",
		"domain_id":   12345,
	})
	d.SetId("res-test001#12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestDnspodPackageDomainResource_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dnspodClient := &dnspod.Client{}
	patches.ApplyMethodReturn(newMockMetaPackageDomain().client, "UseDnsPodClient", dnspodClient)

	patches.ApplyMethodFunc(dnspodClient, "ModifyPackageDomainWithContext",
		func(_ context.Context, request *dnspod.ModifyPackageDomainRequest) (*dnspod.ModifyPackageDomainResponse, error) {
			assert.NotNil(t, request.Operation)
			assert.Equal(t, "change", *request.Operation)
			assert.NotNil(t, request.ResourceId)
			assert.Equal(t, "res-test001", *request.ResourceId)
			assert.NotNil(t, request.DomainId)
			assert.Equal(t, uint64(12345), *request.DomainId)
			assert.NotNil(t, request.NewDomainId)
			assert.Equal(t, uint64(67890), *request.NewDomainId)

			resp := dnspod.NewModifyPackageDomainResponse()
			resp.Response = &dnspod.ModifyPackageDomainResponseParams{
				RequestId: ptrStr("fake-request-id-update"),
			}
			return resp, nil
		})

	patches.ApplyMethodFunc(dnspodClient, "DescribeDomainVipListWithContext",
		func(_ context.Context, request *dnspod.DescribeDomainVipListRequest) (*dnspod.DescribeDomainVipListResponse, error) {
			resp := dnspod.NewDescribeDomainVipListResponse()
			resp.Response = &dnspod.DescribeDomainVipListResponseParams{
				TotalCount: ptrUint64(1),
				PackageList: []*dnspod.PackageListItem{
					{
						DomainId:      ptrUint64(67890),
						Domain:        ptrStr("newexample.com"),
						Grade:         ptrStr("DPG_PROFESSIONAL"),
						GradeTitle:    ptrStr("Professional"),
						VipStartAt:    ptrStr("2024-01-01 00:00:00"),
						VipEndAt:      ptrStr("2025-01-01 00:00:00"),
						VipAutoRenew:  ptrStr("YES"),
						RemainTimes:   ptrUint64(4),
						ResourceId:    ptrStr("res-test001"),
						GradeLevel:    ptrUint64(1),
						Status:        ptrStr("active"),
						IsGracePeriod: ptrStr("NO"),
						Downgrade:     ptrBool(false),
					},
				},
			}
			return resp, nil
		})

	meta := newMockMetaPackageDomain()
	res := dnspodService.ResourceTencentCloudDnspodPackageDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"resource_id": "res-test001",
		"domain_id":   67890,
	})
	d.SetId("res-test001#12345")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "res-test001#67890", d.Id())
	assert.Equal(t, "newexample.com", d.Get("domain"))
}

func TestDnspodPackageDomainResource_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dnspodClient := &dnspod.Client{}
	patches.ApplyMethodReturn(newMockMetaPackageDomain().client, "UseDnsPodClient", dnspodClient)

	patches.ApplyMethodFunc(dnspodClient, "ModifyPackageDomainWithContext",
		func(_ context.Context, request *dnspod.ModifyPackageDomainRequest) (*dnspod.ModifyPackageDomainResponse, error) {
			assert.NotNil(t, request.Operation)
			assert.Equal(t, "unbind", *request.Operation)
			assert.NotNil(t, request.ResourceId)
			assert.Equal(t, "res-test001", *request.ResourceId)
			assert.NotNil(t, request.DomainId)
			assert.Equal(t, uint64(12345), *request.DomainId)

			resp := dnspod.NewModifyPackageDomainResponse()
			resp.Response = &dnspod.ModifyPackageDomainResponseParams{
				RequestId: ptrStr("fake-request-id-delete"),
			}
			return resp, nil
		})

	meta := newMockMetaPackageDomain()
	res := dnspodService.ResourceTencentCloudDnspodPackageDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"resource_id": "res-test001",
		"domain_id":   12345,
	})
	d.SetId("res-test001#12345")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
