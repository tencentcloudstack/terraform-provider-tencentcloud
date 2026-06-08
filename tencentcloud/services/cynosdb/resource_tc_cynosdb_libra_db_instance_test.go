package cynosdb_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccynosdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"
)

type mockMetaLibraDB struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaLibraDB) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaLibraDB{}

func newMockMetaLibraDB() *mockMetaLibraDB {
	return &mockMetaLibraDB{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringLibraDB(s string) *string { return &s }
func ptrInt64LibraDB(i int64) *int64    { return &i }

// go test ./tencentcloud/services/cynosdb/ -run "TestUnitCynosdbLibraDbInstance" -v -count=1 -gcflags="all=-l"

func TestUnitCynosdbLibraDbInstance_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaLibraDB()
	cynosdbClient := &cynosdb.Client{}
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	patches.ApplyMethodFunc(cynosdbClient, "AddLibraDBInstances", func(request *cynosdb.AddLibraDBInstancesRequest) (*cynosdb.AddLibraDBInstancesResponse, error) {
		assert.Equal(t, "cynosdbmysql-12345678", *request.ClusterId)
		assert.Equal(t, "ap-guangzhou-3", *request.Zone)
		assert.Equal(t, int64(4), *request.Cpu)
		assert.Equal(t, int64(8), *request.Mem)
		assert.Equal(t, int64(100), *request.StorageSize)

		resp := &cynosdb.AddLibraDBInstancesResponse{}
		resp.Response = &cynosdb.AddLibraDBInstancesResponseParams{
			ResourceIds: []*string{ptrStringLibraDB("cynosdbmysql-ins-abcdefgh")},
			TranId:      ptrStringLibraDB("tran-123"),
			BigDealIds:  []*string{ptrStringLibraDB("deal-001")},
			DealNames:   []*string{ptrStringLibraDB("order-001")},
			RequestId:   ptrStringLibraDB("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cynosdbClient, "DescribeLibraDBInstanceDetail", func(request *cynosdb.DescribeLibraDBInstanceDetailRequest) (*cynosdb.DescribeLibraDBInstanceDetailResponse, error) {
		resp := &cynosdb.DescribeLibraDBInstanceDetailResponse{}
		resp.Response = &cynosdb.DescribeLibraDBInstanceDetailResponseParams{
			ClusterId:      ptrStringLibraDB("cynosdbmysql-12345678"),
			InstanceId:     ptrStringLibraDB("cynosdbmysql-ins-abcdefgh"),
			InstanceName:   ptrStringLibraDB("tf-test"),
			Zone:           ptrStringLibraDB("ap-guangzhou-3"),
			Cpu:            ptrInt64LibraDB(4),
			Memory:         ptrInt64LibraDB(8),
			Storage:        ptrInt64LibraDB(100),
			Status:         ptrStringLibraDB("running"),
			LibraDBVersion: ptrStringLibraDB("3.1.2"),
			VpcId:          ptrStringLibraDB("vpc-123"),
			SubnetId:       ptrStringLibraDB("subnet-123"),
			RequestId:      ptrStringLibraDB("fake-request-id"),
		}
		return resp, nil
	})

	res := svccynosdb.ResourceTencentCloudCynosdbLibraDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":   "cynosdbmysql-12345678",
		"zone":         "ap-guangzhou-3",
		"cpu":          4,
		"mem":          8,
		"storage_size": 100,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cynosdbmysql-12345678#cynosdbmysql-ins-abcdefgh", d.Id())
}

func TestUnitCynosdbLibraDbInstance_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaLibraDB()
	cynosdbClient := &cynosdb.Client{}
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	patches.ApplyMethodFunc(cynosdbClient, "DescribeLibraDBInstanceDetail", func(request *cynosdb.DescribeLibraDBInstanceDetailRequest) (*cynosdb.DescribeLibraDBInstanceDetailResponse, error) {
		assert.Equal(t, "cynosdbmysql-12345678", *request.ClusterId)
		assert.Equal(t, "cynosdbmysql-ins-abcdefgh", *request.InstanceId)

		resp := &cynosdb.DescribeLibraDBInstanceDetailResponse{}
		resp.Response = &cynosdb.DescribeLibraDBInstanceDetailResponseParams{
			ClusterId:      ptrStringLibraDB("cynosdbmysql-12345678"),
			InstanceId:     ptrStringLibraDB("cynosdbmysql-ins-abcdefgh"),
			InstanceName:   ptrStringLibraDB("tf-test"),
			Zone:           ptrStringLibraDB("ap-guangzhou-3"),
			Cpu:            ptrInt64LibraDB(4),
			Memory:         ptrInt64LibraDB(8),
			Storage:        ptrInt64LibraDB(100),
			Status:         ptrStringLibraDB("running"),
			PayMode:        ptrInt64LibraDB(0),
			InstanceType:   ptrStringLibraDB("Common"),
			StorageType:    ptrStringLibraDB("CLOUD_SSD"),
			LibraDBVersion: ptrStringLibraDB("3.1.2"),
			VpcId:          ptrStringLibraDB("vpc-123"),
			SubnetId:       ptrStringLibraDB("subnet-123"),
			RequestId:      ptrStringLibraDB("fake-request-id"),
		}
		return resp, nil
	})

	res := svccynosdb.ResourceTencentCloudCynosdbLibraDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":   "cynosdbmysql-12345678",
		"zone":         "ap-guangzhou-3",
		"cpu":          4,
		"mem":          8,
		"storage_size": 100,
	})
	d.SetId("cynosdbmysql-12345678#cynosdbmysql-ins-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ap-guangzhou-3", d.Get("zone"))
	assert.Equal(t, 4, d.Get("cpu"))
	assert.Equal(t, 8, d.Get("mem"))
	assert.Equal(t, 100, d.Get("storage_size"))
}

func TestUnitCynosdbLibraDbInstance_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaLibraDB()
	cynosdbClient := &cynosdb.Client{}
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	patches.ApplyMethodFunc(cynosdbClient, "IsolateLibraDBCluster", func(request *cynosdb.IsolateLibraDBClusterRequest) (*cynosdb.IsolateLibraDBClusterResponse, error) {
		assert.Equal(t, "cynosdbmysql-12345678", *request.ClusterId)

		resp := &cynosdb.IsolateLibraDBClusterResponse{}
		resp.Response = &cynosdb.IsolateLibraDBClusterResponseParams{
			FlowId:    ptrInt64LibraDB(12345),
			DealNames: []*string{ptrStringLibraDB("deal-001")},
			RequestId: ptrStringLibraDB("fake-request-id"),
		}
		return resp, nil
	})

	res := svccynosdb.ResourceTencentCloudCynosdbLibraDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":   "cynosdbmysql-12345678",
		"zone":         "ap-guangzhou-3",
		"cpu":          4,
		"mem":          8,
		"storage_size": 100,
	})
	d.SetId("cynosdbmysql-12345678#cynosdbmysql-ins-abcdefgh")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
