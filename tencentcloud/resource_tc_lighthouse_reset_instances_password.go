/*
Provides a resource to create a lighthouse reset_instances_password

Example Usage

```hcl
resource "tencentcloud_lighthouse_reset_instances_password" "reset_instances_password" {
  instance_ids =
  password = "xxxxx"
  user_name = "root"
}
```

Import

lighthouse reset_instances_password can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_reset_instances_password.reset_instances_password reset_instances_password_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudLighthouseResetInstancesPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseResetInstancesPasswordCreate,
		Read:   resourceTencentCloudLighthouseResetInstancesPasswordRead,
		Delete: resourceTencentCloudLighthouseResetInstancesPasswordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list.",
			},

			"password": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Login password of the instance(s). The password requirements vary among different operating systems:The password of a LINUX_UNIX instance must contain 8–30 characters (above 12 characters preferably) in at least three of the following types and cannot begin with /:Lowercase letters: [a–z]Uppercase letters: [A–Z]Digits: 0–9Special symbols: ()`!@#$%^&amp;amp;amp;*-+=_|{}[]:;&amp;amp;#39;&amp;amp;lt;&amp;amp;gt;,.?/The password of a WINDOWS instance must contain 12–30 characters in at least three of the following types and cannot begin with / or include the username:Lowercase letters: [a–z]Uppercase letters: [A–Z]Digits: 0–9Special symbols: ()`!@#$%^&amp;amp;amp;*-+=_|{}[]:;&amp;amp;#39; &amp;amp;lt;&amp;amp;gt;,.?/If both LINUX_UNIX and WINDOWS instances exist, the requirements for password complexity of WINDOWS instances shall prevail.",
			},

			"user_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "OS username of the instance for which you want to reset the password, which can contain up to 64 characters.",
			},
		},
	}
}

func resourceTencentCloudLighthouseResetInstancesPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_reset_instances_password.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewResetInstancesPasswordRequest()
		response   = lighthouse.NewResetInstancesPasswordResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ResetInstancesPassword(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse resetInstancesPassword failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseResetInstancesPasswordStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseResetInstancesPasswordRead(d, meta)
}

func resourceTencentCloudLighthouseResetInstancesPasswordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_reset_instances_password.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseResetInstancesPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_reset_instances_password.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
