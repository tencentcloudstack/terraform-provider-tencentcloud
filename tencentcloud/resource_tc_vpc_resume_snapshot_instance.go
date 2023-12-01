/*
Provides a resource to create a vpc resume_snapshot_instance

Example Usage

Basic example

```hcl
resource "tencentcloud_vpc_resume_snapshot_instance" "resume_snapshot_instance" {
  snapshot_policy_id = "sspolicy-1t6cobbv"
  snapshot_file_id   = "ssfile-emtabuwu2z"
  instance_id        = "ntrgm89v"
}
```

Complete example

```hcl
data "tencentcloud_vpc_snapshot_files" "example" {
  business_type = "securitygroup"
  instance_id   = "sg-902tl7t7"
  start_date    = "2022-10-10 00:00:00"
  end_date      = "2023-10-30 00:00:00"
}

resource "tencentcloud_cos_bucket" "example" {
  bucket = "tf-example-1308919341"
  acl    = "private"
}

resource "tencentcloud_vpc_snapshot_policy" "example" {
  snapshot_policy_name = "tf-example"
  backup_type          = "time"
  cos_bucket           = tencentcloud_cos_bucket.example.bucket
  cos_region           = "ap-guangzhou"
  create_new_cos       = false
  keep_time            = 2

  backup_policies {
    backup_day  = "monday"
    backup_time = "00:00:00"
  }
  backup_policies {
    backup_day  = "tuesday"
    backup_time = "01:00:00"
  }
  backup_policies {
    backup_day  = "wednesday"
    backup_time = "02:00:00"
  }
}

resource "tencentcloud_vpc_resume_snapshot_instance" "example" {
  snapshot_policy_id = tencentcloud_vpc_snapshot_policy.example.id
  snapshot_file_id   = data.tencentcloud_vpc_snapshot_files.example.snapshot_file_set.0.snapshot_file_id
  instance_id        = "policy-1t6cob"
}
```

*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcResumeSnapshotInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcResumeSnapshotInstanceCreate,
		Read:   resourceTencentCloudVpcResumeSnapshotInstanceRead,
		Delete: resourceTencentCloudVpcResumeSnapshotInstanceDelete,
		Schema: map[string]*schema.Schema{
			"snapshot_policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Snapshot policy Id.",
			},

			"snapshot_file_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Snapshot file Id.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},
		},
	}
}

func resourceTencentCloudVpcResumeSnapshotInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_resume_snapshot_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = vpc.NewResumeSnapshotInstanceRequest()
		snapshotPolicyId string
		snapshotFileId   string
		instanceId       string
	)
	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		snapshotPolicyId = v.(string)
		request.SnapshotPolicyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_file_id"); ok {
		snapshotFileId = v.(string)
		request.SnapshotFileId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ResumeSnapshotInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc resumeSnapshotInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(snapshotPolicyId + FILED_SP + snapshotFileId + FILED_SP + instanceId)

	return resourceTencentCloudVpcResumeSnapshotInstanceRead(d, meta)
}

func resourceTencentCloudVpcResumeSnapshotInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_resume_snapshot_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcResumeSnapshotInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_resume_snapshot_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
