/*
Provides a resource to create a elasticsearch diagnose instance

Example Usage

```hcl
resource "tencentcloud_elasticsearch_diagnose_instance" "diagnose_instance" {
  instance_id = "es-xxxxxx"
  diagnose_jobs = ["cluster_health"]
  diagnose_indices = "*"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudElasticsearchDiagnoseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchDiagnoseInstanceCreate,
		Read:   resourceTencentCloudElasticsearchDiagnoseInstanceRead,
		Delete: resourceTencentCloudElasticsearchDiagnoseInstanceDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"diagnose_jobs": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Diagnostic items that need to be triggered.",
			},

			"diagnose_indices": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Indexes that need to be diagnosed. Wildcards are supported.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchDiagnoseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_diagnose_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = elasticsearch.NewDiagnoseInstanceRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("diagnose_jobs"); ok {
		diagnoseJobsSet := v.(*schema.Set).List()
		for i := range diagnoseJobsSet {
			diagnoseJobs := diagnoseJobsSet[i].(string)
			request.DiagnoseJobs = append(request.DiagnoseJobs, &diagnoseJobs)
		}
	}

	if v, ok := d.GetOk("diagnose_indices"); ok {
		request.DiagnoseIndices = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEsClient().DiagnoseInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate elasticsearch DiagnoseInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudElasticsearchDiagnoseInstanceRead(d, meta)
}

func resourceTencentCloudElasticsearchDiagnoseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_diagnose_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudElasticsearchDiagnoseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_diagnose_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
