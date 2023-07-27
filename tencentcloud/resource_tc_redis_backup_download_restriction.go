/*
Provides a resource to create a redis backup_download_restriction

Example Usage

Modify the network information and address of the current region backup file download

```hcl
resource "tencentcloud_redis_backup_download_restriction" "foo" {
	limit_type = "Customize"
	vpc_comparison_symbol = "In"
	ip_comparison_symbol = "In"
	limit_vpc {
		  region = "ap-guangzhou"
		  vpc_list = [var.vpc_id]
	}
	limit_ip = ["10.1.1.12", "10.1.1.13"]
}
```

Import

redis backup_download_restriction can be imported using the region, e.g.

```
terraform import tencentcloud_redis_backup_download_restriction.foo ap-guangzhou
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisBackupDownloadRestriction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisBackupDownloadRestrictionCreate,
		Read:   resourceTencentCloudRedisBackupDownloadRestrictionRead,
		Update: resourceTencentCloudRedisBackupDownloadRestrictionUpdate,
		Delete: resourceTencentCloudRedisBackupDownloadRestrictionDelete,
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

func resourceTencentCloudRedisBackupDownloadRestrictionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_download_restriction.create")()
	defer inconsistentCheck(d, meta)()

	region := meta.(*TencentCloudClient).apiV3Conn.Region

	d.SetId(region)

	return resourceTencentCloudRedisBackupDownloadRestrictionUpdate(d, meta)
}

func resourceTencentCloudRedisBackupDownloadRestrictionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_download_restriction.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	backupDownloadRestriction, err := service.DescribeRedisBackupDownloadRestrictionById(ctx)
	if err != nil {
		return err
	}

	if backupDownloadRestriction == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisBackupDownloadRestriction` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupDownloadRestriction.LimitType != nil {
		_ = d.Set("limit_type", backupDownloadRestriction.LimitType)
	}

	if backupDownloadRestriction.VpcComparisonSymbol != nil {
		_ = d.Set("vpc_comparison_symbol", backupDownloadRestriction.VpcComparisonSymbol)
	}

	if backupDownloadRestriction.IpComparisonSymbol != nil {
		_ = d.Set("ip_comparison_symbol", backupDownloadRestriction.IpComparisonSymbol)
	}

	if backupDownloadRestriction.LimitVpc != nil {
		limitVpcList := []interface{}{}
		for _, limitVpc := range backupDownloadRestriction.LimitVpc {
			limitVpcMap := map[string]interface{}{}

			if limitVpc.Region != nil {
				limitVpcMap["region"] = limitVpc.Region
			}

			if limitVpc.VpcList != nil {
				limitVpcMap["vpc_list"] = limitVpc.VpcList
			}

			limitVpcList = append(limitVpcList, limitVpcMap)
		}

		_ = d.Set("limit_vpc", limitVpcList)

	}

	if backupDownloadRestriction.LimitIp != nil {
		_ = d.Set("limit_ip", backupDownloadRestriction.LimitIp)
	}

	return nil
}

func resourceTencentCloudRedisBackupDownloadRestrictionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_download_restriction.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyBackupDownloadRestrictionRequest()

	if v, ok := d.GetOk("limit_type"); ok {
		request.LimitType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_comparison_symbol"); ok {
		request.VpcComparisonSymbol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_comparison_symbol"); ok {
		request.IpComparisonSymbol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("limit_vpc"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			backupLimitVpcItem := redis.BackupLimitVpcItem{}
			if v, ok := dMap["region"]; ok {
				backupLimitVpcItem.Region = helper.String(v.(string))
			}
			if v, ok := dMap["vpc_list"]; ok {
				vpcListSet := v.(*schema.Set).List()
				for i := range vpcListSet {
					vpcList := vpcListSet[i].(string)
					backupLimitVpcItem.VpcList = append(backupLimitVpcItem.VpcList, &vpcList)
				}
			}
			request.LimitVpc = append(request.LimitVpc, &backupLimitVpcItem)
		}
	}

	if d.HasChange("limit_ip") {
		if v, ok := d.GetOk("limit_ip"); ok {
			limitIpSet := v.(*schema.Set).List()
			for i := range limitIpSet {
				limitIp := limitIpSet[i].(string)
				request.LimitIp = append(request.LimitIp, &limitIp)
			}
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
		log.Printf("[CRITAL]%s update redis backupDownloadRestriction failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisBackupDownloadRestrictionRead(d, meta)
}

func resourceTencentCloudRedisBackupDownloadRestrictionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_download_restriction.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
