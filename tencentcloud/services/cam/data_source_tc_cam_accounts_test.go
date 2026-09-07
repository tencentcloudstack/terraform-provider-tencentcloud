package cam_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	camv20190116 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cam"
)

type mockMetaForCamAccountsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCamAccountsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCamAccountsDS{}

func newMockMetaForCamAccountsDS() *mockMetaForCamAccountsDS {
	return &mockMetaForCamAccountsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCamAccountsDS(s string) *string { return &s }
func ptrInt64CamAccountsDS(v int64) *int64 { return &v }

// go test ./tencentcloud/services/cam/ -run "TestCamAccountsDataSource" -v -count=1 -gcflags="all=-l"

// TestCamAccountsDataSource_ReadBasic tests that the read handler flattens a populated
// ListAccounts response into the users list.
func TestCamAccountsDataSource_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&cam.CamService{}, "DescribeCamAccountsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*camv20190116.ListAllUser, error) {
		return []*camv20190116.ListAllUser{
			{
				Uin:          ptrInt64CamAccountsDS(100037718139),
				Name:         ptrStrCamAccountsDS("alice"),
				Uid:          ptrInt64CamAccountsDS(1034),
				Remark:       ptrStrCamAccountsDS("operator"),
				ConsoleLogin: ptrInt64CamAccountsDS(1),
				PhoneNum:     ptrStrCamAccountsDS("13800000000"),
				CountryCode:  ptrStrCamAccountsDS("86"),
				Email:        ptrStrCamAccountsDS("alice@example.com"),
				CreateTime:   ptrStrCamAccountsDS("2024-01-01 10:00:00"),
				UserType:     ptrStrCamAccountsDS("SubUser"),
			},
			{
				Uin:          ptrInt64CamAccountsDS(100037718140),
				Name:         ptrStrCamAccountsDS("bob"),
				Uid:          ptrInt64CamAccountsDS(1035),
				Remark:       ptrStrCamAccountsDS("collaborator"),
				ConsoleLogin: ptrInt64CamAccountsDS(0),
				PhoneNum:     ptrStrCamAccountsDS("13900000000"),
				CountryCode:  ptrStrCamAccountsDS("86"),
				Email:        ptrStrCamAccountsDS("bob@example.com"),
				CreateTime:   ptrStrCamAccountsDS("2024-02-02 11:00:00"),
				UserType:     ptrStrCamAccountsDS("Collaborator"),
			},
		}, nil
	})

	meta := newMockMetaForCamAccountsDS()
	res := cam.DataSourceTencentCloudCamAccounts()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	users := d.Get("users").([]interface{})
	assert.Len(t, users, 2)

	user0 := users[0].(map[string]interface{})
	assert.Equal(t, 100037718139, user0["uin"].(int))
	assert.Equal(t, "alice", user0["name"].(string))
	assert.Equal(t, 1034, user0["uid"].(int))
	assert.Equal(t, "operator", user0["remark"].(string))
	assert.Equal(t, 1, user0["console_login"].(int))
	assert.Equal(t, "13800000000", user0["phone_num"].(string))
	assert.Equal(t, "86", user0["country_code"].(string))
	assert.Equal(t, "alice@example.com", user0["email"].(string))
	assert.Equal(t, "2024-01-01 10:00:00", user0["create_time"].(string))
	assert.Equal(t, "SubUser", user0["user_type"].(string))

	user1 := users[1].(map[string]interface{})
	assert.Equal(t, 100037718140, user1["uin"].(int))
	assert.Equal(t, "bob", user1["name"].(string))
	assert.Equal(t, "Collaborator", user1["user_type"].(string))
	assert.Equal(t, 0, user1["console_login"].(int))
}

// TestCamAccountsDataSource_ReadWithNilFields tests that nil fields in the API response
// are safely skipped without panicking.
func TestCamAccountsDataSource_ReadWithNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&cam.CamService{}, "DescribeCamAccountsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*camv20190116.ListAllUser, error) {
		return []*camv20190116.ListAllUser{
			{
				Uin:        ptrInt64CamAccountsDS(100037718141),
				Name:       ptrStrCamAccountsDS("charlie"),
				Uid:        nil,
				Remark:     nil,
				PhoneNum:   nil,
				UserType:   ptrStrCamAccountsDS("MessageReceiver"),
				CreateTime: ptrStrCamAccountsDS("2024-03-03 12:00:00"),
			},
		}, nil
	})

	meta := newMockMetaForCamAccountsDS()
	res := cam.DataSourceTencentCloudCamAccounts()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	users := d.Get("users").([]interface{})
	assert.Len(t, users, 1)

	user0 := users[0].(map[string]interface{})
	assert.Equal(t, 100037718141, user0["uin"].(int))
	assert.Equal(t, "charlie", user0["name"].(string))
	assert.Equal(t, "MessageReceiver", user0["user_type"].(string))
	assert.Equal(t, "2024-03-03 12:00:00", user0["create_time"].(string))
	// nil fields default to zero values
	assert.Equal(t, 0, user0["uid"].(int))
	assert.Equal(t, "", user0["remark"].(string))
	assert.Equal(t, "", user0["phone_num"].(string))
}

// TestCamAccountsDataSource_ReadEmptyResponse tests that an empty Users list does not
// produce an error and the read handler still completes.
func TestCamAccountsDataSource_ReadEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&cam.CamService{}, "DescribeCamAccountsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*camv20190116.ListAllUser, error) {
		return []*camv20190116.ListAllUser{}, nil
	})

	meta := newMockMetaForCamAccountsDS()
	res := cam.DataSourceTencentCloudCamAccounts()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Get("users").([]interface{}))
}

// TestCamAccountsDataSource_Schema tests the schema definition includes the expected fields.
func TestCamAccountsDataSource_Schema(t *testing.T) {
	res := cam.DataSourceTencentCloudCamAccounts()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "user_type")
	assert.Contains(t, res.Schema, "users")
	assert.Contains(t, res.Schema, "result_output_file")
	// pagination is handled automatically in the service layer and must not be exposed
	assert.NotContains(t, res.Schema, "max_items")
	assert.NotContains(t, res.Schema, "marker")
	assert.NotContains(t, res.Schema, "is_truncated")

	userTypeSchema := res.Schema["user_type"]
	assert.Equal(t, schema.TypeString, userTypeSchema.Type)
	assert.True(t, userTypeSchema.Optional)
	assert.Nil(t, userTypeSchema.ValidateFunc)

	usersSchema := res.Schema["users"]
	assert.Equal(t, schema.TypeList, usersSchema.Type)
	assert.True(t, usersSchema.Computed)

	elemRes := usersSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "uin")
	assert.Contains(t, elemRes.Schema, "name")
	assert.Contains(t, elemRes.Schema, "uid")
	assert.Contains(t, elemRes.Schema, "remark")
	assert.Contains(t, elemRes.Schema, "console_login")
	assert.Contains(t, elemRes.Schema, "phone_num")
	assert.Contains(t, elemRes.Schema, "country_code")
	assert.Contains(t, elemRes.Schema, "email")
	assert.Contains(t, elemRes.Schema, "create_time")
	assert.Contains(t, elemRes.Schema, "user_type")

	consoleLoginSchema := elemRes.Schema["console_login"]
	assert.Equal(t, schema.TypeInt, consoleLoginSchema.Type)

	uinSchema := elemRes.Schema["uin"]
	assert.Equal(t, schema.TypeInt, uinSchema.Type)
}
