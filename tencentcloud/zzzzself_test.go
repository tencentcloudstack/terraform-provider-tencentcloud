package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func TestDDDDDD(t *testing.T) {

	secretId := os.Getenv(PROVIDER_SECRET_ID)

	secretKey := os.Getenv(PROVIDER_SECRET_KEY)

	region := os.Getenv(PROVIDER_REGION)

	logId := GetLogId(nil)

	ctx := context.WithValue(context.TODO(), "logId", logId)

	apiV3Conn := connectivity.NewTencentCloudClient(secretId, secretKey, region)

	mysqlService := MysqlService{client: apiV3Conn}

	backup, err := mysqlService.CreateBackup(ctx, "cdb-n7j5ouir")

	fmt.Println(backup, err)

	ddd, err := mysqlService.DescribeBackupsByMysqlId(ctx, "cdb-n7j5ouir", 1000)
	_, _ = json.Marshal(ddd)
	for _, v := range ddd {
		fmt.Println(*v.BackupId, *v.Date, *v.FinishTime, *v.Status)
	}

}
