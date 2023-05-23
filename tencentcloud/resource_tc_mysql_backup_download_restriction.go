/*
Provides a resource to create a mysql backup_download_restriction

Example Usage

```hcl
resource "tencentcloud_mysql_backup_download_restriction" "backup_download_restriction" {
	limit_type = "Customize"
	vpc_comparison_symbol = "In"
	ip_comparison_symbol = "In"
	limit_vpc {
		  region = "ap-guangzhou"
		  vpc_list = ["vpc-4owdpnwr"]
	}
	limit_ip = ["127.0.0.1"]
}
```

Import

mysql backup_download_restriction can be imported using the "BackupDownloadRestriction", as follows.

```
terraform import tencentcloud_mysql_backup_download_restriction.backup_download_restriction BackupDownloadRestriction
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlBackupDownloadRestriction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlBackupDownloadRestrictionCreate,
		Read:   resourceTencentCloudMysqlBackupDownloadRestrictionRead,
		Update: resourceTencentCloudMysqlBackupDownloadRestrictionUpdate,
		Delete: resourceTencentCloudMysqlBackupDownloadRestrictionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"limit_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "NoLimit No limit, both internal and external networks can be downloaded; LimitOnlyIntranet Only intranet can be downloaded; Customize user-defined vpc:ip can be downloaded. LimitVpc and LimitIp can be set only when the value is Customize.",
			},

			"vpc_comparison_symbol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "This parameter only supports In, which means that the vpc specified by LimitVpc can be downloaded. The default is In.",
			},

			"ip_comparison_symbol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "In: The specified ip can be downloaded; NotIn: The specified ip cannot be downloaded. The default is In.",
			},

			"limit_vpc": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "vpc settings to limit downloads.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Restrict downloads from regions. Currently only the current region is supported.",
						},
						"vpc_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "List of vpcs to limit downloads.",
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
				Description: "ip settings to limit downloads.",
			},
		},
	}
}

func resourceTencentCloudMysqlBackupDownloadRestrictionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_download_restriction.create")()
	defer inconsistentCheck(d, meta)()

	d.SetId("BackupDownloadRestriction")

	return resourceTencentCloudMysqlBackupDownloadRestrictionUpdate(d, meta)
}

func resourceTencentCloudMysqlBackupDownloadRestrictionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_download_restriction.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupDownloadRestriction, err := service.DescribeMysqlBackupDownloadRestrictionById(ctx)
	if err != nil {
		return err
	}

	if backupDownloadRestriction == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlBackupDownloadRestriction` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudMysqlBackupDownloadRestrictionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_download_restriction.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mysql.NewModifyBackupDownloadRestrictionRequest()

	if v, ok := d.GetOk("limit_type"); ok {
		request.LimitType = helper.String(v.(string))
	}

	if d.HasChange("vpc_comparison_symbol") {
		if v, ok := d.GetOk("vpc_comparison_symbol"); ok {
			request.VpcComparisonSymbol = helper.String(v.(string))
		}
	}

	if d.HasChange("ip_comparison_symbol") {
		if v, ok := d.GetOk("ip_comparison_symbol"); ok {
			request.IpComparisonSymbol = helper.String(v.(string))
		}
	}

	if d.HasChange("limit_vpc") {
		if v, ok := d.GetOk("limit_vpc"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				limitVpcItem := mysql.BackupLimitVpcItem{}
				if v, ok := dMap["region"]; ok {
					limitVpcItem.Region = helper.String(v.(string))
				}
				if v, ok := dMap["vpc_list"]; ok {
					vpcListSet := v.(*schema.Set).List()
					limitVpcItem.VpcList = helper.InterfacesStringsPoint(vpcListSet)
				}
				request.LimitVpc = append(request.LimitVpc, &limitVpcItem)
			}
		}
	}

	if d.HasChange("limit_ip") {
		if v, ok := d.GetOk("limit_ip"); ok {
			limitIpSet := v.(*schema.Set).List()
			request.LimitIp = helper.InterfacesStringsPoint(limitIpSet)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyBackupDownloadRestriction(request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if strings.Contains(sdkerr.Code, "FailedOperation") {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql backupDownloadRestriction failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlBackupDownloadRestrictionRead(d, meta)
}

func resourceTencentCloudMysqlBackupDownloadRestrictionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_download_restriction.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
