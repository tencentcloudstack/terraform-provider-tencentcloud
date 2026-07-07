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
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

type mockMetaPostgresqlInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaPostgresqlInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaPostgresqlInstance{}

func newMockMetaPostgresqlInstance() *mockMetaPostgresqlInstance {
	return &mockMetaPostgresqlInstance{client: &connectivity.TencentCloudClient{}}
}

func ptrStrInst(s string) *string {
	return &s
}

func ptrUint64Inst(i uint64) *uint64 {
	return &i
}

func ptrInt64Inst(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresqlInstanceBackupMethod_Read" -v -count=1 -gcflags="all=-l"
func TestPostgresqlInstanceBackupMethod_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgServiceType := svcpostgresql.PostgresqlService{}
	instanceId := "postgres-backup-method-test"

	// 1. DescribePostgresqlInstanceById -> running instance with 2 net infos
	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstanceById", func(ctx context.Context, id string) (*postgresql.DBInstance, bool, error) {
		assert.Equal(t, instanceId, id)
		running := "running"
		zone := "ap-guangzhou-3"
		netTypePrivate := "private"
		ip := "10.0.0.1"
		port := uint64(5432)
		netTypePublic := "public"
		statusClosed := "closed"
		pubHost := "1.2.3.4"
		pubPort := uint64(6432)
		vpcId := "vpc-test"
		subnetId := "subnet-test"
		instance := &postgresql.DBInstance{
			DBInstanceId:      &instanceId,
			DBInstanceName:    ptrStrInst("tf_postsql_instance"),
			DBInstanceStatus:  &running,
			Zone:              &zone,
			DBVersion:         ptrStrInst("13.3"),
			DBMajorVersion:    ptrStrInst("13"),
			DBKernelVersion:   ptrStrInst("v13.3_r1.1"),
			DBCharset:         ptrStrInst("UTF8"),
			ProjectId:         ptrUint64Inst(0),
			PayType:           ptrStrInst("postpaid"),
			AutoRenew:         ptrUint64Inst(0),
			DBInstanceMemory:  ptrUint64Inst(4 * 1024),
			DBInstanceStorage: ptrUint64Inst(100),
			DBInstanceCpu:     ptrUint64Inst(1),
			CreateTime:        ptrStrInst("2024-01-01 00:00:00"),
			VpcId:             &vpcId,
			SubnetId:          &subnetId,
			DBInstanceNetInfo: []*postgresql.DBInstanceNetInfo{
				{
					NetType:  &netTypePrivate,
					Ip:       &ip,
					Port:     &port,
					VpcId:    &vpcId,
					SubnetId: &subnetId,
				},
				{
					NetType: &netTypePublic,
					Status:  &statusClosed,
					Address: &pubHost,
					Port:    &pubPort,
				},
			},
		}
		return instance, true, nil
	})

	// 2. DescribeRootUser -> empty
	patches.ApplyMethodFunc(&pgServiceType, "DescribeRootUser", func(ctx context.Context, id string) ([]*postgresql.AccountInfo, error) {
		return []*postgresql.AccountInfo{}, nil
	})

	// 3. DescribeDBInstanceSecurityGroupsById -> empty
	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBInstanceSecurityGroupsById", func(ctx context.Context, id string) ([]string, error) {
		return []string{}, nil
	})

	// 4. DescribeDBInstanceAttribute -> ins with storage type and deletion protection
	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBInstanceAttribute", func(ctx context.Context, request *postgresql.DescribeDBInstanceAttributeRequest) (*postgresql.DBInstance, error) {
		storageType := "CLOUD_SSD"
		deletionProtection := false
		return &postgresql.DBInstance{
			DBInstanceStorageType: &storageType,
			DeletionProtection:    &deletionProtection,
			DBNodeSet:             []*postgresql.DBNode{},
		}, nil
	})

	// 5. DescribeDBEncryptionKeys -> has=false
	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBEncryptionKeys", func(ctx context.Context, request *postgresql.DescribeEncryptionKeysRequest) (bool, *postgresql.EncryptionKey, error) {
		return false, nil, nil
	})

	// 6. DescribePostgresqlInstances -> one instance with Uid
	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstances", func(ctx context.Context, filter []*postgresql.Filter) ([]*postgresql.DBInstance, error) {
		return []*postgresql.DBInstance{
			{
				DBInstanceId: &instanceId,
				Uid:          ptrUint64Inst(123456),
			},
		}, nil
	})

	// 7. DescribeBackupPlans -> week plan with BackupMethod
	patches.ApplyMethodFunc(&pgServiceType, "DescribeBackupPlans", func(ctx context.Context, request *postgresql.DescribeBackupPlansRequest) ([]*postgresql.BackupPlan, error) {
		periodTypeWeek := "week"
		minTime := "00:10:11"
		maxTime := "01:10:11"
		retention := uint64(7)
		backupPeriod := "[\"tuesday\",\"wednesday\"]"
		backupMethod := "logical"
		return []*postgresql.BackupPlan{
			{
				BackupPeriodType:          &periodTypeWeek,
				MinBackupStartTime:        &minTime,
				MaxBackupStartTime:        &maxTime,
				BaseBackupRetentionPeriod: &retention,
				BackupPeriod:              &backupPeriod,
				BackupMethod:              &backupMethod,
			},
		}, nil
	})

	// 8. DescribePgParams -> params map
	patches.ApplyMethodFunc(&pgServiceType, "DescribePgParams", func(ctx context.Context, id string) (map[string]string, error) {
		return map[string]string{
			"max_standby_streaming_delay": "10000",
			"max_standby_archive_delay":   "10000",
		}, nil
	})

	// 9. tag NewTagService + DescribeResourceTags
	tagServiceType := svctag.TagService{}
	patches.ApplyMethodFunc(&tagServiceType, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMetaPostgresqlInstance()
	res := svcpostgresql.ResourceTencentCloudPostgresqlInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":              "tf_postsql_instance",
		"vpc_id":            "vpc-test",
		"subnet_id":         "subnet-test",
		"availability_zone": "ap-guangzhou-3",
		"storage":           100,
		"memory":            4,
		"root_password":     "Root123$",
		"backup_plan": []interface{}{
			map[string]interface{}{
				"min_backup_start_time":        "00:10:11",
				"max_backup_start_time":        "01:10:11",
				"base_backup_retention_period": 7,
				"backup_period":                []interface{}{"tuesday", "wednesday"},
				"backup_method":                "logical",
			},
		},
	})
	d.SetId(instanceId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())

	// verify backup_method is read back from API
	backupPlan := d.Get("backup_plan").([]interface{})
	assert.Equal(t, 1, len(backupPlan))
	planMap := backupPlan[0].(map[string]interface{})
	assert.Equal(t, "logical", planMap["backup_method"])
	assert.Equal(t, "00:10:11", planMap["min_backup_start_time"])
	assert.Equal(t, "01:10:11", planMap["max_backup_start_time"])
	assert.Equal(t, 7, planMap["base_backup_retention_period"])
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresqlInstanceBackupMethod_ReadNil" -v -count=1 -gcflags="all=-l"
func TestPostgresqlInstanceBackupMethod_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgServiceType := svcpostgresql.PostgresqlService{}
	instanceId := "postgres-backup-method-nil-test"

	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstanceById", func(ctx context.Context, id string) (*postgresql.DBInstance, bool, error) {
		running := "running"
		zone := "ap-guangzhou-3"
		netTypePrivate := "private"
		ip := "10.0.0.1"
		port := uint64(5432)
		netTypePublic := "public"
		statusClosed := "closed"
		pubHost := "1.2.3.4"
		pubPort := uint64(6432)
		vpcId := "vpc-test"
		subnetId := "subnet-test"
		instance := &postgresql.DBInstance{
			DBInstanceId:     &instanceId,
			DBInstanceName:   ptrStrInst("tf_postsql_instance"),
			DBInstanceStatus: &running,
			Zone:             &zone,
			DBCharset:        ptrStrInst("UTF8"),
			ProjectId:        ptrUint64Inst(0),
			PayType:          ptrStrInst("postpaid"),
			VpcId:            &vpcId,
			SubnetId:         &subnetId,
			DBInstanceNetInfo: []*postgresql.DBInstanceNetInfo{
				{
					NetType:  &netTypePrivate,
					Ip:       &ip,
					Port:     &port,
					VpcId:    &vpcId,
					SubnetId: &subnetId,
				},
				{
					NetType: &netTypePublic,
					Status:  &statusClosed,
					Address: &pubHost,
					Port:    &pubPort,
				},
			},
		}
		return instance, true, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeRootUser", func(ctx context.Context, id string) ([]*postgresql.AccountInfo, error) {
		return []*postgresql.AccountInfo{}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBInstanceSecurityGroupsById", func(ctx context.Context, id string) ([]string, error) {
		return []string{}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBInstanceAttribute", func(ctx context.Context, request *postgresql.DescribeDBInstanceAttributeRequest) (*postgresql.DBInstance, error) {
		storageType := "CLOUD_SSD"
		deletionProtection := false
		return &postgresql.DBInstance{
			DBInstanceStorageType: &storageType,
			DeletionProtection:    &deletionProtection,
			DBNodeSet:             []*postgresql.DBNode{},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBEncryptionKeys", func(ctx context.Context, request *postgresql.DescribeEncryptionKeysRequest) (bool, *postgresql.EncryptionKey, error) {
		return false, nil, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstances", func(ctx context.Context, filter []*postgresql.Filter) ([]*postgresql.DBInstance, error) {
		return []*postgresql.DBInstance{
			{
				DBInstanceId: &instanceId,
				Uid:          ptrUint64Inst(123456),
			},
		}, nil
	})

	// BackupMethod is nil, should not be set
	patches.ApplyMethodFunc(&pgServiceType, "DescribeBackupPlans", func(ctx context.Context, request *postgresql.DescribeBackupPlansRequest) ([]*postgresql.BackupPlan, error) {
		periodTypeWeek := "week"
		minTime := "00:10:11"
		maxTime := "01:10:11"
		retention := uint64(7)
		backupPeriod := "[\"tuesday\",\"wednesday\"]"
		return []*postgresql.BackupPlan{
			{
				BackupPeriodType:          &periodTypeWeek,
				MinBackupStartTime:        &minTime,
				MaxBackupStartTime:        &maxTime,
				BaseBackupRetentionPeriod: &retention,
				BackupPeriod:              &backupPeriod,
			},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribePgParams", func(ctx context.Context, id string) (map[string]string, error) {
		return map[string]string{
			"max_standby_streaming_delay": "10000",
			"max_standby_archive_delay":   "10000",
		}, nil
	})

	tagServiceType := svctag.TagService{}
	patches.ApplyMethodFunc(&tagServiceType, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMetaPostgresqlInstance()
	res := svcpostgresql.ResourceTencentCloudPostgresqlInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":              "tf_postsql_instance",
		"vpc_id":            "vpc-test",
		"subnet_id":         "subnet-test",
		"availability_zone": "ap-guangzhou-3",
		"storage":           100,
		"memory":            4,
		"root_password":     "Root123$",
	})
	d.SetId(instanceId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())

	// backup_plan should be set, but backup_method is empty since API returned nil
	backupPlan := d.Get("backup_plan").([]interface{})
	assert.Equal(t, 1, len(backupPlan))
	planMap := backupPlan[0].(map[string]interface{})
	assert.Equal(t, "", planMap["backup_method"])
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresqlInstanceBackupMethod_Create" -v -count=1 -gcflags="all=-l"
func TestPostgresqlInstanceBackupMethod_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgServiceType := svcpostgresql.PostgresqlService{}
	instanceId := "postgres-backup-method-create-test"

	// capture ModifyBackupPlan request to verify BackupMethod is passed
	var capturedBackupMethod *string
	patches.ApplyMethodFunc(&pgServiceType, "ModifyBackupPlan", func(ctx context.Context, request *postgresql.ModifyBackupPlanRequest) error {
		capturedBackupMethod = request.BackupMethod
		return nil
	})

	// mock all other Create-flow dependencies
	patches.ApplyMethodFunc(&pgServiceType, "DescribeSpecinfos", func(ctx context.Context, zone string, storageType string) ([]*postgresql.SpecItemInfo, error) {
		majorVersion := "13"
		version := "13.3"
		memory := uint64(4 * 1024)
		cpu := uint64(1)
		specCode := "pgdb.test.1c4m"
		return []*postgresql.SpecItemInfo{
			{
				MajorVersion: &majorVersion,
				Version:      &version,
				Memory:       &memory,
				Cpu:          &cpu,
				SpecCode:     &specCode,
			},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "CreatePostgresqlInstance", func(ctx context.Context, args ...interface{}) (string, error) {
		return instanceId, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstanceById", func(ctx context.Context, id string) (*postgresql.DBInstance, bool, error) {
		running := "running"
		memory := uint64(4 * 1024)
		return &postgresql.DBInstance{
			DBInstanceId:     &instanceId,
			DBInstanceStatus: &running,
			DBInstanceMemory: &memory,
		}, true, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "CheckDBInstanceStatus", func(ctx context.Context, id string, retryMinutes ...int) error {
		return nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "ModifyPublicService", func(ctx context.Context, openInternet bool, id string) error {
		return nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "ModifyPostgresqlInstanceName", func(ctx context.Context, id, name string) error {
		return nil
	})

	// Read flow mocks (Create calls Read at the end)
	patches.ApplyMethodFunc(&pgServiceType, "DescribeRootUser", func(ctx context.Context, id string) ([]*postgresql.AccountInfo, error) {
		return []*postgresql.AccountInfo{}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBInstanceSecurityGroupsById", func(ctx context.Context, id string) ([]string, error) {
		return []string{}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBInstanceAttribute", func(ctx context.Context, request *postgresql.DescribeDBInstanceAttributeRequest) (*postgresql.DBInstance, error) {
		storageType := "CLOUD_SSD"
		deletionProtection := false
		return &postgresql.DBInstance{
			DBInstanceStorageType: &storageType,
			DeletionProtection:    &deletionProtection,
			DBNodeSet:             []*postgresql.DBNode{},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBEncryptionKeys", func(ctx context.Context, request *postgresql.DescribeEncryptionKeysRequest) (bool, *postgresql.EncryptionKey, error) {
		return false, nil, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstances", func(ctx context.Context, filter []*postgresql.Filter) ([]*postgresql.DBInstance, error) {
		return []*postgresql.DBInstance{
			{
				DBInstanceId: &instanceId,
				Uid:          ptrUint64Inst(123456),
			},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeBackupPlans", func(ctx context.Context, request *postgresql.DescribeBackupPlansRequest) ([]*postgresql.BackupPlan, error) {
		periodTypeWeek := "week"
		minTime := "00:10:11"
		maxTime := "01:10:11"
		retention := uint64(7)
		backupPeriod := "[\"tuesday\",\"wednesday\"]"
		backupMethod := "snapshot"
		return []*postgresql.BackupPlan{
			{
				BackupPeriodType:          &periodTypeWeek,
				MinBackupStartTime:        &minTime,
				MaxBackupStartTime:        &maxTime,
				BaseBackupRetentionPeriod: &retention,
				BackupPeriod:              &backupPeriod,
				BackupMethod:              &backupMethod,
			},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribePgParams", func(ctx context.Context, id string) (map[string]string, error) {
		return map[string]string{
			"max_standby_streaming_delay": "10000",
			"max_standby_archive_delay":   "10000",
		}, nil
	})

	tagServiceType := svctag.TagService{}
	patches.ApplyMethodFunc(&tagServiceType, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	// mock DescribePostgresqlInstanceById for Read (returns instance with net info)
	// Override the earlier patch to return full instance info for Read
	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstanceById", func(ctx context.Context, id string) (*postgresql.DBInstance, bool, error) {
		running := "running"
		zone := "ap-guangzhou-3"
		netTypePrivate := "private"
		ip := "10.0.0.1"
		port := uint64(5432)
		netTypePublic := "public"
		statusClosed := "closed"
		pubHost := "1.2.3.4"
		pubPort := uint64(6432)
		vpcId := "vpc-test"
		subnetId := "subnet-test"
		memory := uint64(4 * 1024)
		instance := &postgresql.DBInstance{
			DBInstanceId:     &instanceId,
			DBInstanceName:   ptrStrInst("tf_postsql_instance"),
			DBInstanceStatus: &running,
			Zone:             &zone,
			DBCharset:        ptrStrInst("UTF8"),
			ProjectId:        ptrUint64Inst(0),
			PayType:          ptrStrInst("postpaid"),
			VpcId:            &vpcId,
			SubnetId:         &subnetId,
			DBInstanceMemory: &memory,
			DBInstanceNetInfo: []*postgresql.DBInstanceNetInfo{
				{
					NetType:  &netTypePrivate,
					Ip:       &ip,
					Port:     &port,
					VpcId:    &vpcId,
					SubnetId: &subnetId,
				},
				{
					NetType: &netTypePublic,
					Status:  &statusClosed,
					Address: &pubHost,
					Port:    &pubPort,
				},
			},
		}
		return instance, true, nil
	})

	meta := newMockMetaPostgresqlInstance()
	res := svcpostgresql.ResourceTencentCloudPostgresqlInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":              "tf_postsql_instance",
		"vpc_id":            "vpc-test",
		"subnet_id":         "subnet-test",
		"availability_zone": "ap-guangzhou-3",
		"storage":           100,
		"memory":            4,
		"root_password":     "Root123$",
		"engine_version":    "13.3",
		"backup_plan": []interface{}{
			map[string]interface{}{
				"min_backup_start_time":        "00:10:11",
				"max_backup_start_time":        "01:10:11",
				"base_backup_retention_period": 7,
				"backup_period":                []interface{}{"tuesday", "wednesday"},
				"backup_method":                "snapshot",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())

	// verify BackupMethod was passed to ModifyBackupPlan during Create
	assert.NotNil(t, capturedBackupMethod, "BackupMethod should be passed to ModifyBackupPlan during Create")
	if capturedBackupMethod != nil {
		assert.Equal(t, "snapshot", *capturedBackupMethod)
	}
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresqlInstanceBackupMethod_CreateWithoutBackupMethod" -v -count=1 -gcflags="all=-l"
func TestPostgresqlInstanceBackupMethod_CreateWithoutBackupMethod(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgServiceType := svcpostgresql.PostgresqlService{}
	instanceId := "postgres-backup-method-create-empty-test"

	var capturedBackupMethod *string
	patches.ApplyMethodFunc(&pgServiceType, "ModifyBackupPlan", func(ctx context.Context, request *postgresql.ModifyBackupPlanRequest) error {
		capturedBackupMethod = request.BackupMethod
		return nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeSpecinfos", func(ctx context.Context, zone string, storageType string) ([]*postgresql.SpecItemInfo, error) {
		majorVersion := "13"
		version := "13.3"
		memory := uint64(4 * 1024)
		cpu := uint64(1)
		specCode := "pgdb.test.1c4m"
		return []*postgresql.SpecItemInfo{
			{
				MajorVersion: &majorVersion,
				Version:      &version,
				Memory:       &memory,
				Cpu:          &cpu,
				SpecCode:     &specCode,
			},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "CreatePostgresqlInstance", func(ctx context.Context, args ...interface{}) (string, error) {
		return instanceId, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstanceById", func(ctx context.Context, id string) (*postgresql.DBInstance, bool, error) {
		running := "running"
		zone := "ap-guangzhou-3"
		netTypePrivate := "private"
		ip := "10.0.0.1"
		port := uint64(5432)
		netTypePublic := "public"
		statusClosed := "closed"
		pubHost := "1.2.3.4"
		pubPort := uint64(6432)
		vpcId := "vpc-test"
		subnetId := "subnet-test"
		memory := uint64(4 * 1024)
		instance := &postgresql.DBInstance{
			DBInstanceId:     &instanceId,
			DBInstanceName:   ptrStrInst("tf_postsql_instance"),
			DBInstanceStatus: &running,
			Zone:             &zone,
			DBCharset:        ptrStrInst("UTF8"),
			ProjectId:        ptrUint64Inst(0),
			PayType:          ptrStrInst("postpaid"),
			VpcId:            &vpcId,
			SubnetId:         &subnetId,
			DBInstanceMemory: &memory,
			DBInstanceNetInfo: []*postgresql.DBInstanceNetInfo{
				{
					NetType:  &netTypePrivate,
					Ip:       &ip,
					Port:     &port,
					VpcId:    &vpcId,
					SubnetId: &subnetId,
				},
				{
					NetType: &netTypePublic,
					Status:  &statusClosed,
					Address: &pubHost,
					Port:    &pubPort,
				},
			},
		}
		return instance, true, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "CheckDBInstanceStatus", func(ctx context.Context, id string, retryMinutes ...int) error {
		return nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "ModifyPublicService", func(ctx context.Context, openInternet bool, id string) error {
		return nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "ModifyPostgresqlInstanceName", func(ctx context.Context, id, name string) error {
		return nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeRootUser", func(ctx context.Context, id string) ([]*postgresql.AccountInfo, error) {
		return []*postgresql.AccountInfo{}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBInstanceSecurityGroupsById", func(ctx context.Context, id string) ([]string, error) {
		return []string{}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBInstanceAttribute", func(ctx context.Context, request *postgresql.DescribeDBInstanceAttributeRequest) (*postgresql.DBInstance, error) {
		storageType := "CLOUD_SSD"
		deletionProtection := false
		return &postgresql.DBInstance{
			DBInstanceStorageType: &storageType,
			DeletionProtection:    &deletionProtection,
			DBNodeSet:             []*postgresql.DBNode{},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeDBEncryptionKeys", func(ctx context.Context, request *postgresql.DescribeEncryptionKeysRequest) (bool, *postgresql.EncryptionKey, error) {
		return false, nil, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribePostgresqlInstances", func(ctx context.Context, filter []*postgresql.Filter) ([]*postgresql.DBInstance, error) {
		return []*postgresql.DBInstance{
			{
				DBInstanceId: &instanceId,
				Uid:          ptrUint64Inst(123456),
			},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribeBackupPlans", func(ctx context.Context, request *postgresql.DescribeBackupPlansRequest) ([]*postgresql.BackupPlan, error) {
		periodTypeWeek := "week"
		minTime := "00:10:11"
		maxTime := "01:10:11"
		retention := uint64(7)
		backupPeriod := "[\"tuesday\",\"wednesday\"]"
		return []*postgresql.BackupPlan{
			{
				BackupPeriodType:          &periodTypeWeek,
				MinBackupStartTime:        &minTime,
				MaxBackupStartTime:        &maxTime,
				BaseBackupRetentionPeriod: &retention,
				BackupPeriod:              &backupPeriod,
			},
		}, nil
	})

	patches.ApplyMethodFunc(&pgServiceType, "DescribePgParams", func(ctx context.Context, id string) (map[string]string, error) {
		return map[string]string{
			"max_standby_streaming_delay": "10000",
			"max_standby_archive_delay":   "10000",
		}, nil
	})

	tagServiceType := svctag.TagService{}
	patches.ApplyMethodFunc(&tagServiceType, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMetaPostgresqlInstance()
	res := svcpostgresql.ResourceTencentCloudPostgresqlInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":              "tf_postsql_instance",
		"vpc_id":            "vpc-test",
		"subnet_id":         "subnet-test",
		"availability_zone": "ap-guangzhou-3",
		"storage":           100,
		"memory":            4,
		"root_password":     "Root123$",
		"engine_version":    "13.3",
		"backup_plan": []interface{}{
			map[string]interface{}{
				"min_backup_start_time":        "00:10:11",
				"max_backup_start_time":        "01:10:11",
				"base_backup_retention_period": 7,
				"backup_period":                []interface{}{"tuesday", "wednesday"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())

	// when backup_method is not configured, BackupMethod should NOT be set on the request
	assert.Nil(t, capturedBackupMethod, "BackupMethod should not be set when user does not configure backup_method")
}

// keep unused helpers to avoid lint issues; ptrInt64Inst reserved for future use
var _ = ptrInt64Inst
