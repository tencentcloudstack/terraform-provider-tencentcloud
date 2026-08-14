package postgresql_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

type mockMetaPostgresDatabase struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaPostgresDatabase) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaPostgresDatabase{}

func newMockMetaPostgresDatabase() *mockMetaPostgresDatabase {
	return &mockMetaPostgresDatabase{client: &connectivity.TencentCloudClient{}}
}

func ptrStringPgDb(s string) *string {
	return &s
}

func ptrUint64PgDb(i uint64) *uint64 {
	return &i
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresDatabase" -v -count=1 -gcflags="all=-l"

// TestPostgresDatabase_Create tests the Create function
func TestPostgresDatabase_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaPostgresDatabase().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "CreateDatabaseWithContext", func(_ context.Context, request *postgresql.CreateDatabaseRequest) (*postgresql.CreateDatabaseResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-6fego161", *request.DBInstanceId)
		assert.NotNil(t, request.DatabaseName)
		assert.Equal(t, "test_db", *request.DatabaseName)
		assert.NotNil(t, request.DatabaseOwner)
		assert.Equal(t, "tcuser", *request.DatabaseOwner)
		assert.NotNil(t, request.Encoding)
		assert.Equal(t, "UTF8", *request.Encoding)
		assert.NotNil(t, request.Collate)
		assert.Equal(t, "C", *request.Collate)
		assert.NotNil(t, request.Ctype)
		assert.Equal(t, "C", *request.Ctype)

		resp := postgresql.NewCreateDatabaseResponse()
		resp.Response = &postgresql.CreateDatabaseResponseParams{
			RequestId: ptrStringPgDb("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeDatabasesWithContext", func(_ context.Context, request *postgresql.DescribeDatabasesRequest) (*postgresql.DescribeDatabasesResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-6fego161", *request.DBInstanceId)
		assert.NotNil(t, request.Filters)
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "database-name", *request.Filters[0].Name)

		resp := postgresql.NewDescribeDatabasesResponse()
		resp.Response = &postgresql.DescribeDatabasesResponseParams{
			TotalCount: ptrUint64PgDb(1),
			Databases: []*postgresql.Database{
				{
					DatabaseName:  ptrStringPgDb("test_db"),
					DatabaseOwner: ptrStringPgDb("tcuser"),
					Encoding:      ptrStringPgDb("UTF8"),
					Collate:       ptrStringPgDb("C"),
					Ctype:         ptrStringPgDb("C"),
				},
			},
			RequestId: ptrStringPgDb("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPostgresDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-6fego161",
		"database_name":  "test_db",
		"database_owner": "tcuser",
		"encoding":       "UTF8",
		"collate":        "C",
		"ctype":          "C",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-6fego161#test_db", d.Id())
	assert.Equal(t, "postgres-6fego161", d.Get("db_instance_id"))
	assert.Equal(t, "test_db", d.Get("database_name"))
	assert.Equal(t, "tcuser", d.Get("database_owner"))
	assert.Equal(t, "UTF8", d.Get("encoding"))
	assert.Equal(t, "C", d.Get("collate"))
	assert.Equal(t, "C", d.Get("ctype"))
}

// TestPostgresDatabase_Read tests the Read function
func TestPostgresDatabase_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaPostgresDatabase().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDatabasesWithContext", func(_ context.Context, request *postgresql.DescribeDatabasesRequest) (*postgresql.DescribeDatabasesResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-6fego161", *request.DBInstanceId)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, uint64(100), *request.Limit)
		assert.NotNil(t, request.Offset)
		assert.Equal(t, uint64(0), *request.Offset)

		resp := postgresql.NewDescribeDatabasesResponse()
		resp.Response = &postgresql.DescribeDatabasesResponseParams{
			TotalCount: ptrUint64PgDb(1),
			Databases: []*postgresql.Database{
				{
					DatabaseName:  ptrStringPgDb("test_db"),
					DatabaseOwner: ptrStringPgDb("tcuser"),
					Encoding:      ptrStringPgDb("UTF8"),
					Collate:       ptrStringPgDb("C"),
					Ctype:         ptrStringPgDb("C"),
				},
			},
			RequestId: ptrStringPgDb("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPostgresDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-6fego161",
		"database_name":  "test_db",
		"database_owner": "tcuser",
		"encoding":       "UTF8",
		"collate":        "C",
		"ctype":          "C",
	})
	d.SetId("postgres-6fego161#test_db")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-6fego161#test_db", d.Id())
	assert.Equal(t, "test_db", d.Get("database_name"))
	assert.Equal(t, "tcuser", d.Get("database_owner"))
	assert.Equal(t, "UTF8", d.Get("encoding"))
	assert.Equal(t, "C", d.Get("collate"))
	assert.Equal(t, "C", d.Get("ctype"))
}

// TestPostgresDatabase_Read_NotFound tests the Read function when the database is not found
func TestPostgresDatabase_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaPostgresDatabase().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDatabasesWithContext", func(_ context.Context, request *postgresql.DescribeDatabasesRequest) (*postgresql.DescribeDatabasesResponse, error) {
		resp := postgresql.NewDescribeDatabasesResponse()
		resp.Response = &postgresql.DescribeDatabasesResponseParams{
			TotalCount: ptrUint64PgDb(0),
			Databases:  []*postgresql.Database{},
			RequestId:  ptrStringPgDb("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPostgresDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-6fego161",
		"database_name":  "test_db",
		"database_owner": "tcuser",
	})
	d.SetId("postgres-6fego161#test_db")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestPostgresDatabase_Update tests the Update function
func TestPostgresDatabase_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaPostgresDatabase().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "ModifyDatabaseOwnerWithContext", func(_ context.Context, request *postgresql.ModifyDatabaseOwnerRequest) (*postgresql.ModifyDatabaseOwnerResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-6fego161", *request.DBInstanceId)
		assert.NotNil(t, request.DatabaseName)
		assert.Equal(t, "test_db", *request.DatabaseName)
		assert.NotNil(t, request.DatabaseOwner)
		assert.Equal(t, "newuser", *request.DatabaseOwner)

		resp := postgresql.NewModifyDatabaseOwnerResponse()
		resp.Response = &postgresql.ModifyDatabaseOwnerResponseParams{
			RequestId: ptrStringPgDb("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeDatabasesWithContext", func(_ context.Context, request *postgresql.DescribeDatabasesRequest) (*postgresql.DescribeDatabasesResponse, error) {
		resp := postgresql.NewDescribeDatabasesResponse()
		resp.Response = &postgresql.DescribeDatabasesResponseParams{
			TotalCount: ptrUint64PgDb(1),
			Databases: []*postgresql.Database{
				{
					DatabaseName:  ptrStringPgDb("test_db"),
					DatabaseOwner: ptrStringPgDb("newuser"),
					Encoding:      ptrStringPgDb("UTF8"),
					Collate:       ptrStringPgDb("C"),
					Ctype:         ptrStringPgDb("C"),
				},
			},
			RequestId: ptrStringPgDb("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPostgresDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-6fego161",
		"database_name":  "test_db",
		"database_owner": "newuser",
		"encoding":       "UTF8",
		"collate":        "C",
		"ctype":          "C",
	})
	d.SetId("postgres-6fego161#test_db")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-6fego161#test_db", d.Id())
	assert.Equal(t, "newuser", d.Get("database_owner"))
}

// TestPostgresDatabase_Delete tests the Delete function
func TestPostgresDatabase_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaPostgresDatabase().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DeleteDatabaseWithContext", func(_ context.Context, request *postgresql.DeleteDatabaseRequest) (*postgresql.DeleteDatabaseResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-6fego161", *request.DBInstanceId)
		assert.NotNil(t, request.DatabaseName)
		assert.Equal(t, "test_db", *request.DatabaseName)

		resp := postgresql.NewDeleteDatabaseResponse()
		resp.Response = &postgresql.DeleteDatabaseResponseParams{
			RequestId: ptrStringPgDb("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPostgresDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-6fego161",
		"database_name":  "test_db",
		"database_owner": "tcuser",
	})
	d.SetId("postgres-6fego161#test_db")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
