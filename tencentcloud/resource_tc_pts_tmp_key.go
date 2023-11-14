/*
Provides a resource to create a pts tmp_key

Example Usage

```hcl
resource "tencentcloud_pts_tmp_key" "tmp_key" {
  project_id = "project-abc"
  scenario_id = "scenario-abc"
}
```

Import

pts tmp_key can be imported using the id, e.g.

```
terraform import tencentcloud_pts_tmp_key.tmp_key tmp_key_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudPtsTmpKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsTmpKeyCreate,
		Read:   resourceTencentCloudPtsTmpKeyRead,
		Delete: resourceTencentCloudPtsTmpKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"scenario_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Scenario ID.",
			},
		},
	}
}

func resourceTencentCloudPtsTmpKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_tmp_key.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = pts.NewGenerateTmpKeyRequest()
		response   = pts.NewGenerateTmpKeyResponse()
		projectId  string
		scenarioId string
	)
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scenario_id"); ok {
		scenarioId = v.(string)
		request.ScenarioId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().GenerateTmpKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate pts tmpKey failed, reason:%+v", logId, err)
		return err
	}

	projectId = *response.Response.ProjectId
	d.SetId(strings.Join([]string{projectId, scenarioId}, FILED_SP))

	return resourceTencentCloudPtsTmpKeyRead(d, meta)
}

func resourceTencentCloudPtsTmpKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_tmp_key.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPtsTmpKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_tmp_key.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
