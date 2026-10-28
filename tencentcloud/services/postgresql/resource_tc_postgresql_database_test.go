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

type mockMetaPostgresqlDatabase struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaPostgresqlDatabase) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaPostgresqlDatabase{}

func newMockMetaPostgresqlDatabase() *mockMetaPostgresqlDatabase {
	return &mockMetaPostgresqlDatabase{client: &connectivity.TencentCloudClient{}}
}

func ptrStringPgDb(s string) *string {
	return &s
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresqlDatabase" -v -count=1 -gcflags="all=-l"

// TestPostgresqlDatabase_Create tests the Create function
func TestPostgresqlDatabase_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	service := &svcpostgresql.PostgresqlService{}
	patches.ApplyMethodFunc(service, "CreatePostgresqlDatabase", func(_ context.Context, dbInstanceId, databaseName, databaseOwner, encoding, collate, ctype string) error {
		assert.Equal(t, "postgres-6fego161", dbInstanceId)
		assert.Equal(t, "test_db", databaseName)
		assert.Equal(t, "tcuser", databaseOwner)
		assert.Equal(t, "UTF8", encoding)
		assert.Equal(t, "C", collate)
		assert.Equal(t, "C", ctype)
		return nil
	})

	patches.ApplyMethodFunc(service, "DescribePostgresqlDatabaseById", func(_ context.Context, dbInstanceId, databaseName string) (*postgresql.Database, error) {
		assert.Equal(t, "postgres-6fego161", dbInstanceId)
		assert.Equal(t, "test_db", databaseName)
		return &postgresql.Database{
			DatabaseName:  ptrStringPgDb("test_db"),
			DatabaseOwner: ptrStringPgDb("tcuser"),
			Encoding:      ptrStringPgDb("UTF8"),
			Collate:       ptrStringPgDb("C"),
			Ctype:         ptrStringPgDb("C"),
		}, nil
	})

	meta := newMockMetaPostgresqlDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresqlDatabase()
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

// TestPostgresqlDatabase_Read tests the Read function
func TestPostgresqlDatabase_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	service := &svcpostgresql.PostgresqlService{}
	patches.ApplyMethodFunc(service, "DescribePostgresqlDatabaseById", func(_ context.Context, dbInstanceId, databaseName string) (*postgresql.Database, error) {
		assert.Equal(t, "postgres-6fego161", dbInstanceId)
		assert.Equal(t, "test_db", databaseName)
		return &postgresql.Database{
			DatabaseName:  ptrStringPgDb("test_db"),
			DatabaseOwner: ptrStringPgDb("tcuser"),
			Encoding:      ptrStringPgDb("UTF8"),
			Collate:       ptrStringPgDb("C"),
			Ctype:         ptrStringPgDb("C"),
		}, nil
	})

	meta := newMockMetaPostgresqlDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresqlDatabase()
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

// TestPostgresqlDatabase_Read_NotFound tests the Read function when the database is not found
func TestPostgresqlDatabase_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	service := &svcpostgresql.PostgresqlService{}
	patches.ApplyMethodFunc(service, "DescribePostgresqlDatabaseById", func(_ context.Context, dbInstanceId, databaseName string) (*postgresql.Database, error) {
		return nil, nil
	})

	meta := newMockMetaPostgresqlDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresqlDatabase()
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

// TestPostgresqlDatabase_Update tests the Update function
func TestPostgresqlDatabase_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	service := &svcpostgresql.PostgresqlService{}
	patches.ApplyMethodFunc(service, "ModifyPostgresqlDatabaseOwner", func(_ context.Context, dbInstanceId, databaseName, databaseOwner string) error {
		assert.Equal(t, "postgres-6fego161", dbInstanceId)
		assert.Equal(t, "test_db", databaseName)
		assert.Equal(t, "newuser", databaseOwner)
		return nil
	})

	patches.ApplyMethodFunc(service, "DescribePostgresqlDatabaseById", func(_ context.Context, dbInstanceId, databaseName string) (*postgresql.Database, error) {
		return &postgresql.Database{
			DatabaseName:  ptrStringPgDb("test_db"),
			DatabaseOwner: ptrStringPgDb("newuser"),
			Encoding:      ptrStringPgDb("UTF8"),
			Collate:       ptrStringPgDb("C"),
			Ctype:         ptrStringPgDb("C"),
		}, nil
	})

	meta := newMockMetaPostgresqlDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresqlDatabase()
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

// TestPostgresqlDatabase_Delete tests the Delete function
func TestPostgresqlDatabase_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	service := &svcpostgresql.PostgresqlService{}
	patches.ApplyMethodFunc(service, "DeletePostgresqlDatabaseById", func(_ context.Context, dbInstanceId, databaseName string) error {
		assert.Equal(t, "postgres-6fego161", dbInstanceId)
		assert.Equal(t, "test_db", databaseName)
		return nil
	})

	meta := newMockMetaPostgresqlDatabase()
	res := svcpostgresql.ResourceTencentCloudPostgresqlDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-6fego161",
		"database_name":  "test_db",
		"database_owner": "tcuser",
	})
	d.SetId("postgres-6fego161#test_db")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
