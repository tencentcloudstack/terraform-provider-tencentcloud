/*
Provides a resource to create a sqlserver instance_t_d_e

Example Usage

```hcl
resource "tencentcloud_sqlserver_instance_t_d_e" "instance_t_d_e" {
  instance_id = "mssql-i1z41iwd"
  certificate_attribution = ""
  quote_uin = ""
}
```

Import

sqlserver instance_t_d_e can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_instance_t_d_e.instance_t_d_e instance_t_d_e_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"log"
)

func resourceTencentCloudSqlserverInstanceTDE() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverInstanceTDECreate,
		Read:   resourceTencentCloudSqlserverInstanceTDERead,
		Update: resourceTencentCloudSqlserverInstanceTDEUpdate,
		Delete: resourceTencentCloudSqlserverInstanceTDEDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"certificate_attribution": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate attribution. self- means to use the account&amp;amp;#39;s own certificate, others- means to refer to the certificate of other accounts, and the default is self.",
			},

			"quote_uin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Other referenced main account IDs, required when CertificateAttribute is others.",
			},
		},
	}
}

func resourceTencentCloudSqlserverInstanceTDECreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_t_d_e.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverInstanceTDEUpdate(d, meta)
}

func resourceTencentCloudSqlserverInstanceTDERead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_t_d_e.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceTDEId := d.Id()

	instanceTDE, err := service.DescribeSqlserverInstanceTDEById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceTDE == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverInstanceTDE` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceTDE.InstanceId != nil {
		_ = d.Set("instance_id", instanceTDE.InstanceId)
	}

	if instanceTDE.CertificateAttribution != nil {
		_ = d.Set("certificate_attribution", instanceTDE.CertificateAttribution)
	}

	if instanceTDE.QuoteUin != nil {
		_ = d.Set("quote_uin", instanceTDE.QuoteUin)
	}

	return nil
}

func resourceTencentCloudSqlserverInstanceTDEUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_t_d_e.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyInstanceEncryptAttributesRequest()

	instanceTDEId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "certificate_attribution", "quote_uin"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyInstanceEncryptAttributes(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver instanceTDE failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverInstanceTDERead(d, meta)
}

func resourceTencentCloudSqlserverInstanceTDEDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_t_d_e.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
