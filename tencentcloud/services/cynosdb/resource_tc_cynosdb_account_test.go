package cynosdb_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccynosdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbAccountResource_basic -v
func TestAccTencentCloudCynosdbAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCynosdbAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbAccountExists("tencentcloud_cynosdb_account.account"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account.account", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "cluster_id", tcacctest.DefaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "account_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "description", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "host", "%"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "max_user_connections", "1"),
				),
			},
			{
				ResourceName:            "tencentcloud_cynosdb_account.account",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccCynosdbAccountUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbAccountExists("tencentcloud_cynosdb_account.account"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account.account", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "cluster_id", tcacctest.DefaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "account_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "description", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "host", "%"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "max_user_connections", "2"),
				),
			},
		},
	})
}

func testAccCheckCynosdbAccountDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_account" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		accountName := idSplit[1]
		host := idSplit[2]

		has, err := cynosdbService.DescribeCynosdbAccountById(ctx, clusterId, accountName, host)
		if err != nil {
			return err
		}
		if has == nil {
			return nil
		}
		return fmt.Errorf("cynosdb cluster account still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb cluster account id is not set")
		}
		cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		accountName := idSplit[1]
		host := idSplit[2]

		has, err := cynosdbService.DescribeCynosdbAccountById(ctx, clusterId, accountName, host)
		if err != nil {
			return err
		}
		if has == nil {
			return fmt.Errorf("cynosdb cluster account doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCynosdbAccount = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_account" "account" {
	cluster_id = var.cynosdb_cluster_id
	account_name = "terraform_test"
	account_password = "Password@1234"
	host = "%"
	description = "test"
	max_user_connections = 1
}

`

const testAccCynosdbAccountUp = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_account" "account" {
	cluster_id = var.cynosdb_cluster_id
	account_name = "terraform_test"
	account_password = "Password@1234"
	host = "%"
	description = "terraform test"
	max_user_connections = 2
}

`

// ---- unit tests (gomonkey mock) ----

type mockMetaCynosdbAccount struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCynosdbAccount) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCynosdbAccount{}

func newMockMetaCynosdbAccount() *mockMetaCynosdbAccount {
	return &mockMetaCynosdbAccount{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringAccount(s string) *string { return &s }
func ptrInt64Account(i int64) *int64    { return &i }

// go test ./tencentcloud/services/cynosdb/ -run "TestUnitCynosdbAccount" -v -count=1 -gcflags="all=-l"

// TestUnitCynosdbAccount_Create_WithTaskPolling covers the Create flow where
// CreateAccounts returns a valid TaskId and DescribeTasks reports success.
func TestUnitCynosdbAccount_Create_WithTaskPolling(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaCynosdbAccount()
	cynosdbClient := &cynosdb.Client{}
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	// mock CreateAccounts to return a valid TaskId
	patches.ApplyMethodFunc(cynosdbClient, "CreateAccounts", func(request *cynosdb.CreateAccountsRequest) (*cynosdb.CreateAccountsResponse, error) {
		assert.Equal(t, "cynosdb-cid-test", *request.ClusterId)
		resp := &cynosdb.CreateAccountsResponse{}
		resp.Response = &cynosdb.CreateAccountsResponseParams{
			TaskId:    ptrInt64Account(123456),
			RequestId: ptrStringAccount("fake-request-id"),
		}
		return resp, nil
	})

	// mock DescribeTasks to report the task as success
	patches.ApplyMethodFunc(cynosdbClient, "DescribeTasks", func(request *cynosdb.DescribeTasksRequest) (*cynosdb.DescribeTasksResponse, error) {
		resp := &cynosdb.DescribeTasksResponse{}
		resp.Response = &cynosdb.DescribeTasksResponseParams{
			TotalCount: ptrInt64Account(1),
			TaskList: []*cynosdb.BizTaskInfo{
				{
					Status: ptrStringAccount("success"),
				},
			},
			RequestId: ptrStringAccount("fake-request-id"),
		}
		return resp, nil
	})

	// mock DescribeAccounts for the Read after task polling
	patches.ApplyMethodFunc(cynosdbClient, "DescribeAccounts", func(request *cynosdb.DescribeAccountsRequest) (*cynosdb.DescribeAccountsResponse, error) {
		resp := &cynosdb.DescribeAccountsResponse{}
		resp.Response = &cynosdb.DescribeAccountsResponseParams{
			AccountSet: []*cynosdb.Account{
				{
					AccountName:        ptrStringAccount("tf_test"),
					Host:               ptrStringAccount("%"),
					Description:        ptrStringAccount("test"),
					MaxUserConnections: ptrInt64Account(1),
				},
			},
			TotalCount: ptrInt64Account(1),
			RequestId:  ptrStringAccount("fake-request-id"),
		}
		return resp, nil
	})

	res := svccynosdb.ResourceTencentCloudCynosdbAccount()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":           "cynosdb-cid-test",
		"account_name":         "tf_test",
		"account_password":     "Password@1234",
		"host":                 "%",
		"description":          "test",
		"max_user_connections": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cynosdb-cid-test"+tccommon.FILED_SP+"tf_test"+tccommon.FILED_SP+"%", d.Id())
	assert.Equal(t, "cynosdb-cid-test", d.Get("cluster_id"))
	assert.Equal(t, "tf_test", d.Get("account_name"))
	assert.Equal(t, "%", d.Get("host"))
	assert.Equal(t, "test", d.Get("description"))
	assert.Equal(t, 1, d.Get("max_user_connections"))
}

// TestUnitCynosdbAccount_Create_TaskIdZeroSkipPolling covers the Create flow
// where CreateAccounts returns a TaskId of 0, polling is skipped and Read runs directly.
func TestUnitCynosdbAccount_Create_TaskIdZeroSkipPolling(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaCynosdbAccount()
	cynosdbClient := &cynosdb.Client{}
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	// mock CreateAccounts to return TaskId = 0
	patches.ApplyMethodFunc(cynosdbClient, "CreateAccounts", func(request *cynosdb.CreateAccountsRequest) (*cynosdb.CreateAccountsResponse, error) {
		resp := &cynosdb.CreateAccountsResponse{}
		resp.Response = &cynosdb.CreateAccountsResponseParams{
			TaskId:    ptrInt64Account(0),
			RequestId: ptrStringAccount("fake-request-id"),
		}
		return resp, nil
	})

	// DescribeTasks should never be called when TaskId is 0; if it is, fail the test.
	describeTasksCalled := false
	patches.ApplyMethodFunc(cynosdbClient, "DescribeTasks", func(request *cynosdb.DescribeTasksRequest) (*cynosdb.DescribeTasksResponse, error) {
		describeTasksCalled = true
		return nil, nil
	})

	// mock DescribeAccounts for the Read
	patches.ApplyMethodFunc(cynosdbClient, "DescribeAccounts", func(request *cynosdb.DescribeAccountsRequest) (*cynosdb.DescribeAccountsResponse, error) {
		resp := &cynosdb.DescribeAccountsResponse{}
		resp.Response = &cynosdb.DescribeAccountsResponseParams{
			AccountSet: []*cynosdb.Account{
				{
					AccountName:        ptrStringAccount("tf_test"),
					Host:               ptrStringAccount("%"),
					Description:        ptrStringAccount("test"),
					MaxUserConnections: ptrInt64Account(1),
				},
			},
			TotalCount: ptrInt64Account(1),
			RequestId:  ptrStringAccount("fake-request-id"),
		}
		return resp, nil
	})

	res := svccynosdb.ResourceTencentCloudCynosdbAccount()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":           "cynosdb-cid-test",
		"account_name":         "tf_test",
		"account_password":     "Password@1234",
		"host":                 "%",
		"description":          "test",
		"max_user_connections": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.False(t, describeTasksCalled, "DescribeTasks should not be called when TaskId is 0")
	assert.Equal(t, "cynosdb-cid-test"+tccommon.FILED_SP+"tf_test"+tccommon.FILED_SP+"%", d.Id())
}

// TestUnitCynosdbAccount_Read_NotFoundRetryExhausted covers the Read flow where
// the account is consistently not found and the retry is exhausted, clearing the id.
func TestUnitCynosdbAccount_Read_NotFoundRetryExhausted(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaCynosdbAccount()
	cynosdbClient := &cynosdb.Client{}
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	// mock DescribeAccounts to always return an empty AccountSet (account not found)
	patches.ApplyMethodFunc(cynosdbClient, "DescribeAccounts", func(request *cynosdb.DescribeAccountsRequest) (*cynosdb.DescribeAccountsResponse, error) {
		resp := &cynosdb.DescribeAccountsResponse{}
		resp.Response = &cynosdb.DescribeAccountsResponseParams{
			AccountSet: []*cynosdb.Account{},
			TotalCount: ptrInt64Account(0),
			RequestId:  ptrStringAccount("fake-request-id"),
		}
		return resp, nil
	})

	res := svccynosdb.ResourceTencentCloudCynosdbAccount()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":           "cynosdb-cid-test",
		"account_name":         "tf_test",
		"account_password":     "Password@1234",
		"host":                 "%",
		"description":          "test",
		"max_user_connections": 1,
	})
	d.SetId("cynosdb-cid-test" + tccommon.FILED_SP + "tf_test" + tccommon.FILED_SP + "%")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestUnitCynosdbAccount_Read_RetryThenSuccess covers the Read flow where the
// account is not found at first but is found on a subsequent retry.
func TestUnitCynosdbAccount_Read_RetryThenSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaCynosdbAccount()
	cynosdbClient := &cynosdb.Client{}
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	// mock DescribeAccounts: first call returns empty, subsequent calls return the account
	callCount := 0
	patches.ApplyMethodFunc(cynosdbClient, "DescribeAccounts", func(request *cynosdb.DescribeAccountsRequest) (*cynosdb.DescribeAccountsResponse, error) {
		callCount++
		if callCount == 1 {
			resp := &cynosdb.DescribeAccountsResponse{}
			resp.Response = &cynosdb.DescribeAccountsResponseParams{
				AccountSet: []*cynosdb.Account{},
				TotalCount: ptrInt64Account(0),
				RequestId:  ptrStringAccount("fake-request-id"),
			}
			return resp, nil
		}
		resp := &cynosdb.DescribeAccountsResponse{}
		resp.Response = &cynosdb.DescribeAccountsResponseParams{
			AccountSet: []*cynosdb.Account{
				{
					AccountName:        ptrStringAccount("tf_test"),
					Host:               ptrStringAccount("%"),
					Description:        ptrStringAccount("test"),
					MaxUserConnections: ptrInt64Account(1),
				},
			},
			TotalCount: ptrInt64Account(1),
			RequestId:  ptrStringAccount("fake-request-id"),
		}
		return resp, nil
	})

	res := svccynosdb.ResourceTencentCloudCynosdbAccount()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":           "cynosdb-cid-test",
		"account_name":         "tf_test",
		"account_password":     "Password@1234",
		"host":                 "%",
		"description":          "test",
		"max_user_connections": 1,
	})
	d.SetId("cynosdb-cid-test" + tccommon.FILED_SP + "tf_test" + tccommon.FILED_SP + "%")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cynosdb-cid-test"+tccommon.FILED_SP+"tf_test"+tccommon.FILED_SP+"%", d.Id())
	assert.Equal(t, "cynosdb-cid-test", d.Get("cluster_id"))
	assert.Equal(t, "tf_test", d.Get("account_name"))
	assert.Equal(t, "%", d.Get("host"))
	assert.Equal(t, "test", d.Get("description"))
	assert.Equal(t, 1, d.Get("max_user_connections"))
}
