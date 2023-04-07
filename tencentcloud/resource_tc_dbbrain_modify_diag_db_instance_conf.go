/*
Provides a resource to create a dbbrain modify_diag_db_instance_conf

Example Usage

```hcl
resource "tencentcloud_dbbrain_modify_diag_db_instance_conf" "modify_diag_db_instance_conf" {
  instance_confs {
		daily_inspection = ""
		overview_display = ""

  }
  regions = ""
  product = ""
  instance_ids =
}
```

Import

dbbrain modify_diag_db_instance_conf can be imported using the id, e.g.

```
terraform import tencentcloud_dbbrain_modify_diag_db_instance_conf.modify_diag_db_instance_conf modify_diag_db_instance_conf_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDbbrainModifyDiagDbInstanceConf() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbbrainModifyDiagDbInstanceConfCreate,
		Read:   resourceTencentCloudDbbrainModifyDiagDbInstanceConfRead,
		Delete: resourceTencentCloudDbbrainModifyDiagDbInstanceConfDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_confs": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Instance configuration, including inspection, overview switch, etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_inspection": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database inspection switch, Yes/No.",
						},
						"overview_display": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance overview switch, Yes/No.",
						},
					},
				},
			},

			"regions": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Effective instance region, the value is &amp;quot;All&amp;quot;, which means all regions.",
			},

			"product": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include： &amp;quot;mysql&amp;quot; - cloud database MySQL, &amp;quot;cynosdb&amp;quot; - cloud database CynosDB for MySQL.",
			},

			"instance_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Specifies the ID of the instance whose inspection status is changed.",
			},
		},
	}
}

func resourceTencentCloudDbbrainModifyDiagDbInstanceConfCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_modify_diag_db_instance_conf.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dbbrain.NewModifyDiagDBInstanceConfRequest()
		response   = dbbrain.NewModifyDiagDBInstanceConfResponse()
		instanceId uint64
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "instance_confs"); ok {
		instanceConfs := dbbrain.InstanceConfs{}
		if v, ok := dMap["daily_inspection"]; ok {
			instanceConfs.DailyInspection = helper.String(v.(string))
		}
		if v, ok := dMap["overview_display"]; ok {
			instanceConfs.OverviewDisplay = helper.String(v.(string))
		}
		request.InstanceConfs = &instanceConfs
	}

	if v, ok := d.GetOk("regions"); ok {
		request.Regions = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().ModifyDiagDBInstanceConf(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Println("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Println("[CRITAL]%s operate dbbrain modifyDiagDbInstanceConf failed, reason:%+v", logId, err)
		return nil
	}

	instanceId = *response.Response.InstanceId
	d.SetId(helper.UInt64ToStr(instanceId))

	return resourceTencentCloudDbbrainModifyDiagDbInstanceConfRead(d, meta)
}

func resourceTencentCloudDbbrainModifyDiagDbInstanceConfRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_modify_diag_db_instance_conf.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDbbrainModifyDiagDbInstanceConfDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_modify_diag_db_instance_conf.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
