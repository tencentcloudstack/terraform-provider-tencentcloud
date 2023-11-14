/*
Provides a resource to create a redis backup

Example Usage

```hcl
resource "tencentcloud_redis_backup" "backup" {
  limit_type = "NoLimit"
  vpc_comparison_symbol = "In"
  ip_comparison_symbol = "In"
  limit_vpc {
		region = "ap-guangzhou"
		vpc_list =

  }
  limit_ip =
}
```

Import

redis backup can be imported using the id, e.g.

```
terraform import tencentcloud_redis_backup.backup backup_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"log"
)

func resourceTencentCloudRedisBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisBackupCreate,
		Read:   resourceTencentCloudRedisBackupRead,
		Update: resourceTencentCloudRedisBackupUpdate,
		Delete: resourceTencentCloudRedisBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"limit_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Types of network restrictions for downloading backup files:- NoLimit: There is no limit, and backup files can be downloaded from both Tencent Cloud and internal and external networks.- LimitOnlyIntranet: Only intranet addresses automatically assigned by Tencent Cloud can download backup files.- Customize: refers to a user-defined private network downloadable backup file.",
			},

			"vpc_comparison_symbol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "This parameter only supports entering In, which means that the custom LimitVpc can download the backup file.",
			},

			"ip_comparison_symbol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Identifies whether the customized LimitIP address can download the backup file.- In: Custom IP addresses are available for download.- NotIn: Custom IPs are not available for download.",
			},

			"limit_vpc": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "A custom VPC ID for a downloadable backup file.If the parameter LimitType is **Customize**, you need to configure this parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Customize the region of the VPC to which the backup file is downloaded.",
						},
						"vpc_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Customize the list of VPCs to download backup files.",
						},
					},
				},
			},

			"limit_ip": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A custom VPC IP address for downloadable backup files.If the parameter LimitType is **Customize**, you need to configure this parameter.",
			},
		},
	}
}

func resourceTencentCloudRedisBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup.create")()
	defer inconsistentCheck(d, meta)()

	var region string
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	d.SetId(region)

	return resourceTencentCloudRedisBackupUpdate(d, meta)
}

func resourceTencentCloudRedisBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupId := d.Id()

	backup, err := service.DescribeRedisBackupById(ctx, region)
	if err != nil {
		return err
	}

	if backup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisBackup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backup.LimitType != nil {
		_ = d.Set("limit_type", backup.LimitType)
	}

	if backup.VpcComparisonSymbol != nil {
		_ = d.Set("vpc_comparison_symbol", backup.VpcComparisonSymbol)
	}

	if backup.IpComparisonSymbol != nil {
		_ = d.Set("ip_comparison_symbol", backup.IpComparisonSymbol)
	}

	if backup.LimitVpc != nil {
		limitVpcList := []interface{}{}
		for _, limitVpc := range backup.LimitVpc {
			limitVpcMap := map[string]interface{}{}

			if backup.LimitVpc.Region != nil {
				limitVpcMap["region"] = backup.LimitVpc.Region
			}

			if backup.LimitVpc.VpcList != nil {
				limitVpcMap["vpc_list"] = backup.LimitVpc.VpcList
			}

			limitVpcList = append(limitVpcList, limitVpcMap)
		}

		_ = d.Set("limit_vpc", limitVpcList)

	}

	if backup.LimitIp != nil {
		_ = d.Set("limit_ip", backup.LimitIp)
	}

	return nil
}

func resourceTencentCloudRedisBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyBackupDownloadRestrictionRequest()

	backupId := d.Id()

	request.Region = &region

	immutableArgs := []string{"limit_type", "vpc_comparison_symbol", "ip_comparison_symbol", "limit_vpc", "limit_ip"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyBackupDownloadRestriction(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis backup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisBackupRead(d, meta)
}

func resourceTencentCloudRedisBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
