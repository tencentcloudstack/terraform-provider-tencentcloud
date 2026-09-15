package bdrc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	bdrc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"
)

type mockMetaBdrc struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBdrc) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBdrc{}

func newMockMetaBdrc() *mockMetaBdrc {
	return &mockMetaBdrc{client: &connectivity.TencentCloudClient{}}
}

func ptrStringBdrc(s string) *string {
	return &s
}

func ptrInt64Bdrc(i int64) *int64 {
	return &i
}

func ptrBoolBdrc(b bool) *bool {
	return &b
}

func buildCreateTargetInstanceParameters() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"source_instance_id":   "ins-xxxxxxxx",
			"instance_charge_type": "POSTPAID_BY_HOUR",
			"instance_type":        "S5.MEDIUM4",
			"image_id":             "img-xxxxxxxx",
			"instance_name":        "tf-example-instance",
			"placement": []interface{}{
				map[string]interface{}{
					"zone":       "ap-guangzhou-3",
					"project_id": 0,
				},
			},
			"system_disk": []interface{}{
				map[string]interface{}{
					"disk_type":            "CLOUD_BSSD",
					"disk_size":            50,
					"delete_with_instance": true,
				},
			},
			"virtual_private_cloud": []interface{}{
				map[string]interface{}{
					"vpc_id":         "vpc-xxxxxxxx",
					"subnet_id":      "subnet-xxxxxxxx",
					"as_vpc_gateway": false,
				},
			},
			"internet_accessible": []interface{}{
				map[string]interface{}{
					"internet_charge_type":       "TRAFFIC_POSTPAID_BY_HOUR",
					"internet_max_bandwidth_out": 10,
					"public_ip_assigned":         false,
				},
			},
			"login_settings": []interface{}{
				map[string]interface{}{
					"password": "TFexample123",
				},
			},
			"enhanced_service": []interface{}{
				map[string]interface{}{
					"security_service": []interface{}{
						map[string]interface{}{
							"enabled": true,
						},
					},
					"monitor_service": []interface{}{
						map[string]interface{}{
							"enabled": true,
						},
					},
				},
			},
		},
	}
}

func TestAccBdrcInstanceCopyPair_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrc().client, "UseBdrcV20260330Client", bdrcClient)

	copyPairId := "cvmcopypair-xxxxxxxx"

	patches.ApplyMethodFunc(bdrcClient, "CreateInstanceCopyPairWithContext", func(_ context.Context, request *bdrcv20260330.CreateInstanceCopyPairRequest) (*bdrcv20260330.CreateInstanceCopyPairResponse, error) {
		assert.NotNil(t, request.ProtectGroupId)
		assert.Equal(t, "ProtectGroup-xxxxx", *request.ProtectGroupId)
		assert.NotNil(t, request.CreateTargetInstanceParameters)
		assert.Equal(t, 1, len(request.CreateTargetInstanceParameters))

		resp := bdrcv20260330.NewCreateInstanceCopyPairResponse()
		resp.Response = &bdrcv20260330.CreateInstanceCopyPairResponseParams{
			CopyPairIds: []*string{ptrStringBdrc(copyPairId)},
			RequestId:   ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&bdrc.BdrcService{}, "DescribeBdrcInstanceCopyPairById", func(_ context.Context, id string) (*bdrcv20260330.CopyPair, error) {
		assert.Equal(t, copyPairId, id)
		return &bdrcv20260330.CopyPair{
			CopyPairId:       ptrStringBdrc(copyPairId),
			CopyPairName:     ptrStringBdrc("tf-example-copy-pair"),
			CopyPairState:    ptrStringBdrc("RUNNING"),
			CopyPairType:     ptrStringBdrc("INSTANCE"),
			ProtectGroupId:   ptrStringBdrc("ProtectGroup-xxxxx"),
			SourceRegion:     ptrStringBdrc("ap-guangzhou"),
			TargetRegion:     ptrStringBdrc("ap-shanghai"),
			Percent:          ptrInt64Bdrc(100),
			CreateTime:       ptrStringBdrc("2024-01-01 00:00:00"),
			DeferredCreate:   ptrBoolBdrc(false),
			TargetCvmCreated: ptrBoolBdrc(false),
		}, nil
	})

	meta := newMockMetaBdrc()
	res := bdrc.ResourceTencentCloudBdrcInstanceCopyPair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"protect_group_id":                  "ProtectGroup-xxxxx",
		"instance_copy_pair_name":           "tf-example-copy-pair",
		"recovery_point_objective":          15,
		"create_target_instance_parameters": buildCreateTargetInstanceParameters(),
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, copyPairId, d.Id())
	assert.Equal(t, "RUNNING", d.Get("copy_pair_state").(string))
	assert.Equal(t, "ProtectGroup-xxxxx", d.Get("protect_group_id").(string))
}

func TestAccBdrcInstanceCopyPair_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrc().client, "UseBdrcV20260330Client", bdrcClient)

	copyPairId := "cvmcopypair-xxxxxxxx"

	patches.ApplyMethodFunc(&bdrc.BdrcService{}, "DescribeBdrcInstanceCopyPairById", func(_ context.Context, id string) (*bdrcv20260330.CopyPair, error) {
		assert.Equal(t, copyPairId, id)
		return &bdrcv20260330.CopyPair{
			CopyPairId:           ptrStringBdrc(copyPairId),
			CopyPairName:         ptrStringBdrc("tf-example-copy-pair"),
			CopyPairState:        ptrStringBdrc("RUNNING"),
			CopyPairType:         ptrStringBdrc("INSTANCE"),
			SitePairId:           ptrStringBdrc("sitepair-xxxxx"),
			ProtectGroupId:       ptrStringBdrc("ProtectGroup-xxxxx"),
			SourceRegion:         ptrStringBdrc("ap-guangzhou"),
			SourceZone:           ptrStringBdrc("ap-guangzhou-3"),
			TargetRegion:         ptrStringBdrc("ap-shanghai"),
			TargetZone:           ptrStringBdrc("ap-shanghai-2"),
			SourceResourceId:     ptrStringBdrc("ins-xxxxxxxx"),
			TargetResourceId:     ptrStringBdrc("ins-yyyyyyyy"),
			Percent:              ptrInt64Bdrc(100),
			CreateTime:           ptrStringBdrc("2024-01-01 00:00:00"),
			DataDirection:        ptrStringBdrc("POSITIVE"),
			DisasterRecoveryType: ptrStringBdrc("CROSS_REGION"),
			DeferredCreate:       ptrBoolBdrc(false),
			TargetCvmCreated:     ptrBoolBdrc(false),
			DiskCopyPairSet: []*bdrcv20260330.DiskCopyPairForCvm{
				{
					CopyPairId:       ptrStringBdrc("copypair-disk-001"),
					CopyPairName:     ptrStringBdrc("disk-copy-pair-001"),
					SourceResourceId: ptrStringBdrc("disk-xxxxxxxx"),
					TargetResourceId: ptrStringBdrc("disk-yyyyyyyy"),
					CreateTime:       ptrStringBdrc("2024-01-01 00:00:00"),
				},
			},
		}, nil
	})

	meta := newMockMetaBdrc()
	res := bdrc.ResourceTencentCloudBdrcInstanceCopyPair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId(copyPairId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, copyPairId, d.Id())
	assert.Equal(t, "RUNNING", d.Get("copy_pair_state").(string))
	assert.Equal(t, "INSTANCE", d.Get("copy_pair_type").(string))
	assert.Equal(t, "sitepair-xxxxx", d.Get("site_pair_id").(string))
	assert.Equal(t, "ap-guangzhou", d.Get("source_region").(string))
	assert.Equal(t, "ap-shanghai", d.Get("target_region").(string))
	assert.Equal(t, "ins-xxxxxxxx", d.Get("source_resource_id").(string))
	assert.Equal(t, "ins-yyyyyyyy", d.Get("target_resource_id").(string))
	assert.Equal(t, 100, d.Get("percent").(int))
	assert.Equal(t, "POSITIVE", d.Get("data_direction").(string))

	diskSet := d.Get("disk_copy_pair_set").([]interface{})
	assert.Equal(t, 1, len(diskSet))
	diskMap := diskSet[0].(map[string]interface{})
	assert.Equal(t, "copypair-disk-001", diskMap["copy_pair_id"].(string))
}

func TestAccBdrcInstanceCopyPair_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrc().client, "UseBdrcV20260330Client", bdrcClient)

	copyPairId := "cvmcopypair-notfound"

	patches.ApplyMethodFunc(&bdrc.BdrcService{}, "DescribeBdrcInstanceCopyPairById", func(_ context.Context, id string) (*bdrcv20260330.CopyPair, error) {
		return nil, nil
	})

	meta := newMockMetaBdrc()
	res := bdrc.ResourceTencentCloudBdrcInstanceCopyPair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId(copyPairId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestAccBdrcInstanceCopyPair_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrc().client, "UseBdrcV20260330Client", bdrcClient)

	copyPairId := "cvmcopypair-xxxxxxxx"

	patches.ApplyMethodFunc(bdrcClient, "ModifyCopyPairAttributeWithContext", func(_ context.Context, request *bdrcv20260330.ModifyCopyPairAttributeRequest) (*bdrcv20260330.ModifyCopyPairAttributeResponse, error) {
		assert.NotNil(t, request.CopyPairId)
		assert.Equal(t, copyPairId, *request.CopyPairId)
		assert.NotNil(t, request.CopyPairName)
		assert.Equal(t, "updated-copy-pair-name", *request.CopyPairName)

		resp := bdrcv20260330.NewModifyCopyPairAttributeResponse()
		resp.Response = &bdrcv20260330.ModifyCopyPairAttributeResponseParams{
			RequestId: ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&bdrc.BdrcService{}, "DescribeBdrcInstanceCopyPairById", func(_ context.Context, id string) (*bdrcv20260330.CopyPair, error) {
		return &bdrcv20260330.CopyPair{
			CopyPairId:     ptrStringBdrc(copyPairId),
			CopyPairName:   ptrStringBdrc("updated-copy-pair-name"),
			CopyPairState:  ptrStringBdrc("RUNNING"),
			CopyPairType:   ptrStringBdrc("INSTANCE"),
			ProtectGroupId: ptrStringBdrc("ProtectGroup-xxxxx"),
		}, nil
	})

	meta := newMockMetaBdrc()
	res := bdrc.ResourceTencentCloudBdrcInstanceCopyPair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"protect_group_id":                  "ProtectGroup-xxxxx",
		"instance_copy_pair_name":           "updated-copy-pair-name",
		"create_target_instance_parameters": buildCreateTargetInstanceParameters(),
	})
	d.SetId(copyPairId)

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, copyPairId, d.Id())
	assert.Equal(t, "updated-copy-pair-name", d.Get("copy_pair_name").(string))
}

func TestAccBdrcInstanceCopyPair_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrc().client, "UseBdrcV20260330Client", bdrcClient)

	copyPairId := "cvmcopypair-xxxxxxxx"

	patches.ApplyMethodFunc(bdrcClient, "DeleteCopyPairsWithContext", func(_ context.Context, request *bdrcv20260330.DeleteCopyPairsRequest) (*bdrcv20260330.DeleteCopyPairsResponse, error) {
		assert.NotNil(t, request.CopyPairIds)
		assert.Equal(t, 1, len(request.CopyPairIds))
		assert.Equal(t, copyPairId, *request.CopyPairIds[0])
		assert.NotNil(t, request.CopyPairType)
		assert.Equal(t, "INSTANCE", *request.CopyPairType)

		resp := bdrcv20260330.NewDeleteCopyPairsResponse()
		resp.Response = &bdrcv20260330.DeleteCopyPairsResponseParams{
			RequestId: ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrc()
	res := bdrc.ResourceTencentCloudBdrcInstanceCopyPair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"delete_target_resource": true,
	})
	d.SetId(copyPairId)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
