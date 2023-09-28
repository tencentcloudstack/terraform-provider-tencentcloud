/*
Provides a resource to create a pts tmp_key

Example Usage

```hcl
resource "tencentcloud_pts_tmp_key_generate" "tmp_key" {
  project_id = "project-1b0zqmhg"
  scenario_id = "scenario-abc"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPtsTmpKeyGenerate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsTmpKeyGenerateCreate,
		Read:   resourceTencentCloudPtsTmpKeyGenerateRead,
		Delete: resourceTencentCloudPtsTmpKeyGenerateDelete,

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

			"start_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The timestamp of the moment when the temporary access credential was obtained (in seconds).",
			},
			"expired_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Timestamp of temporary access credential timeout (in seconds).",
			},
			"credentials": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Temporary access credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tmp_secret_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Temporary secret ID.",
						},
						"tmp_secret_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Temporary secret key.",
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Temporary token.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPtsTmpKeyGenerateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_tmp_key_generate.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = pts.NewGenerateTmpKeyRequest()
		response  = pts.NewGenerateTmpKeyResponse()
		projectId string
	)
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scenario_id"); ok {
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

	d.SetId(projectId)

	if response != nil || response.Response != nil {
		credentials := response.Response.Credentials
		if credentials != nil {
			credentialsMap := map[string]interface{}{}
			if credentials.TmpSecretId != nil {
				credentialsMap["tmp_secret_id"] = credentials.TmpSecretId
			}

			if credentials.TmpSecretKey != nil {
				credentialsMap["tmp_secret_key"] = credentials.TmpSecretKey
			}

			if credentials.Token != nil {
				credentialsMap["token"] = credentials.Token
			}

			_ = d.Set("credentials", []interface{}{credentialsMap})
		}

		if response.Response.StartTime != nil {
			_ = d.Set("start_time", response.Response.StartTime)
		}

		if response.Response.ExpiredTime != nil {
			_ = d.Set("expired_time", response.Response.ExpiredTime)
		}
	}

	return resourceTencentCloudPtsTmpKeyGenerateRead(d, meta)
}

func resourceTencentCloudPtsTmpKeyGenerateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_tmp_key_generate.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPtsTmpKeyGenerateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_tmp_key_generate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
