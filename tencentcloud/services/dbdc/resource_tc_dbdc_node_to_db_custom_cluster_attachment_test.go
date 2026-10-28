package dbdc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

func TestAccTencentCloudDbdcNodeToDbCustomClusterAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbdcNodeToDbCustomClusterAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example", "cluster_id", "dbcc-xxxxxxxx"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example", "node_id", "dbcn-xxxxxxxx"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example", "status"),
				),
			},
			{
				ResourceName:            "tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_id", "login_settings", "labels", "taints", "host_name", "host_name_type"},
			},
		},
	})
}

const testAccDbdcNodeToDbCustomClusterAttachment = `
resource "tencentcloud_dbdc_node_to_db_custom_cluster_attachment" "example" {
  cluster_id = "dbcc-xxxxxxxx"
  node_id    = "dbcn-xxxxxxxx"
  image_id   = "img-xxxxxxxx"

  login_settings {
    password = "Passw0rd@2024"
  }

  labels {
    key   = "env"
    value = "prod"
  }

  taints {
    key    = "dedicated"
    effect = "NoSchedule"
    value  = "true"
  }

  host_name       = "node-{R:1}"
  host_name_type  = 1
}
`

// go test ./tencentcloud/services/dbdc/ -run "TestUnitDbdcNodeToDbCustomClusterAttachment" -v -count=1 -gcflags="all=-l"

func TestUnitDbdcNodeToDbCustomClusterAttachmentCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDbdcClusterNodesDS()
	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(meta.client, "UseDbdcV20201029Client", dbdcClient)

	createCalled := false
	patches.ApplyMethodFunc(dbdcClient, "AddNodesToDBCustomClusterWithContext", func(ctx context.Context, request *dbdcv20201029.AddNodesToDBCustomClusterRequest) (*dbdcv20201029.AddNodesToDBCustomClusterResponse, error) {
		createCalled = true
		assert.Equal(t, "dbcc-cluster123", *request.ClusterId)
		assert.Equal(t, "dbcn-node456", *request.NodeIds[0])
		assert.Equal(t, "img-rm13akp3", *request.ImageId)

		// verify login_settings
		assert.NotNil(t, request.LoginSettings)
		assert.Equal(t, "Passw0rd@2026", *request.LoginSettings.Password)

		// verify labels
		assert.Len(t, request.Labels, 1)
		assert.Equal(t, "env", *request.Labels[0].Key)
		assert.Equal(t, "prod", *request.Labels[0].Value)

		// verify taints
		assert.Len(t, request.Taints, 1)
		assert.Equal(t, "dedicated", *request.Taints[0].Key)
		assert.Equal(t, "NoSchedule", *request.Taints[0].Effect)
		assert.Equal(t, "true", *request.Taints[0].Value)

		// verify host_name and host_name_type
		assert.NotNil(t, request.HostName)
		assert.Equal(t, "node-{R:1}", *request.HostName)
		assert.NotNil(t, request.HostNameType)
		assert.Equal(t, int64(1), *request.HostNameType)

		resp := dbdcv20201029.NewAddNodesToDBCustomClusterResponse()
		resp.Response = &dbdcv20201029.AddNodesToDBCustomClusterResponseParams{
			RequestId: ptrStrCN("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodesWithContext", func(ctx context.Context, request *dbdcv20201029.DescribeDBCustomClusterNodesRequest) (*dbdcv20201029.DescribeDBCustomClusterNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodesResponseParams{
			TotalCount: ptrInt64CN(1),
			NodeSet: []*dbdcv20201029.DBCustomClusterNode{
				{
					NodeId:      ptrStrCN("dbcn-node456"),
					NodeName:    ptrStrCN("test-node"),
					LanIP:       ptrStrCN("10.0.0.1"),
					SSHEndpoint: ptrStrCN("10.0.0.1:22"),
					Status:      ptrStrCN("Running"),
					Zone:        ptrStrCN("ap-shanghai-5"),
					NodeType:    ptrStrCN("DB.AT5.8XLARGE128"),
					NetworkMode: ptrStrCN("cross_tenant_eni"),
					EniIP:       ptrStrCN("10.0.0.99"),
				},
			},
			RequestId: ptrStrCN("fake-request-id"),
		}
		return resp, nil
	})

	res := dbdc.ResourceTencentCloudDbdcNodeToDbCustomClusterAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-cluster123",
		"node_id":    "dbcn-node456",
		"image_id":   "img-rm13akp3",
		"login_settings": []interface{}{
			map[string]interface{}{
				"password": "Passw0rd@2026",
			},
		},
		"labels": []interface{}{
			map[string]interface{}{
				"key":   "env",
				"value": "prod",
			},
		},
		"taints": []interface{}{
			map[string]interface{}{
				"key":    "dedicated",
				"effect": "NoSchedule",
				"value":  "true",
			},
		},
		"host_name":      "node-{R:1}",
		"host_name_type": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.True(t, createCalled)
	assert.Equal(t, "dbcc-cluster123#dbcn-node456", d.Id())

	// verify computed fields populated by Read after Create
	assert.Equal(t, "cross_tenant_eni", d.Get("network_mode"))
	assert.Equal(t, "10.0.0.99", d.Get("eni_ip"))
}

func TestUnitDbdcNodeToDbCustomClusterAttachmentRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDbdcClusterNodesDS()
	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(meta.client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodesWithContext", func(ctx context.Context, request *dbdcv20201029.DescribeDBCustomClusterNodesRequest) (*dbdcv20201029.DescribeDBCustomClusterNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodesResponseParams{
			TotalCount: ptrInt64CN(1),
			NodeSet: []*dbdcv20201029.DBCustomClusterNode{
				{
					NodeId:      ptrStrCN("dbcn-node456"),
					NodeName:    ptrStrCN("test-node"),
					LanIP:       ptrStrCN("10.0.0.1"),
					SSHEndpoint: ptrStrCN("10.0.0.1:22"),
					Status:      ptrStrCN("Running"),
					Zone:        ptrStrCN("ap-shanghai-5"),
					NodeType:    ptrStrCN("DB.AT5.8XLARGE128"),
					NetworkMode: ptrStrCN("cross_tenant_eni"),
					EniIP:       ptrStrCN("10.0.0.99"),
				},
			},
			RequestId: ptrStrCN("fake-request-id"),
		}
		return resp, nil
	})

	res := dbdc.ResourceTencentCloudDbdcNodeToDbCustomClusterAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("dbcc-cluster123#dbcn-node456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "dbcc-cluster123", d.Get("cluster_id"))
	assert.Equal(t, "dbcn-node456", d.Get("node_id"))
	assert.Equal(t, "cross_tenant_eni", d.Get("network_mode"))
	assert.Equal(t, "10.0.0.99", d.Get("eni_ip"))
}

func TestUnitDbdcNodeToDbCustomClusterAttachmentReadNilNetworkFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDbdcClusterNodesDS()
	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(meta.client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodesWithContext", func(ctx context.Context, request *dbdcv20201029.DescribeDBCustomClusterNodesRequest) (*dbdcv20201029.DescribeDBCustomClusterNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodesResponseParams{
			TotalCount: ptrInt64CN(1),
			NodeSet: []*dbdcv20201029.DBCustomClusterNode{
				{
					NodeId:   ptrStrCN("dbcn-node456"),
					NodeName: ptrStrCN("test-node"),
					Status:   ptrStrCN("Running"),
				},
			},
			RequestId: ptrStrCN("fake-request-id"),
		}
		return resp, nil
	})

	res := dbdc.ResourceTencentCloudDbdcNodeToDbCustomClusterAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("dbcc-cluster123#dbcn-node456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// network_mode and eni_ip should be unset (empty) when nil
	assert.Equal(t, "", d.Get("network_mode"))
	assert.Equal(t, "", d.Get("eni_ip"))
}

func TestUnitDbdcNodeToDbCustomClusterAttachmentDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDbdcClusterNodesDS()
	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(meta.client, "UseDbdcV20201029Client", dbdcClient)

	deleteCalled := false
	patches.ApplyMethodFunc(dbdcClient, "RemoveNodesFromDBCustomClusterWithContext", func(ctx context.Context, request *dbdcv20201029.RemoveNodesFromDBCustomClusterRequest) (*dbdcv20201029.RemoveNodesFromDBCustomClusterResponse, error) {
		deleteCalled = true
		assert.Equal(t, "dbcc-cluster123", *request.ClusterId)
		assert.Equal(t, "dbcn-node456", *request.NodeIds[0])

		// verify login_settings passed on delete
		assert.NotNil(t, request.LoginSettings)
		assert.Equal(t, "Passw0rd@2026", *request.LoginSettings.Password)

		resp := dbdcv20201029.NewRemoveNodesFromDBCustomClusterResponse()
		resp.Response = &dbdcv20201029.RemoveNodesFromDBCustomClusterResponseParams{
			RequestId: ptrStrCN("fake-request-id"),
		}
		return resp, nil
	})

	res := dbdc.ResourceTencentCloudDbdcNodeToDbCustomClusterAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-cluster123",
		"node_id":    "dbcn-node456",
		"login_settings": []interface{}{
			map[string]interface{}{
				"password": "Passw0rd@2026",
			},
		},
	})
	d.SetId("dbcc-cluster123#dbcn-node456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}
