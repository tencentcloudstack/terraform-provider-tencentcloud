package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudVpcSnapshotPolicyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcSnapshotPolicyConfigCreate,
		Read:   resourceTencentCloudVpcSnapshotPolicyConfigRead,
		Update: resourceTencentCloudVpcSnapshotPolicyConfigUpdate,
		Delete: resourceTencentCloudVpcSnapshotPolicyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"snapshot_policy_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot policy Id.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "If enable snapshot policy.",
			},
		},
	}
}

func resourceTencentCloudVpcSnapshotPolicyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy_config.create")()
	defer inconsistentCheck(d, meta)()

	snapshotPolicyId := d.Get("snapshot_policy_id").(string)

	d.SetId(snapshotPolicyId)

	return resourceTencentCloudVpcSnapshotPolicyConfigUpdate(d, meta)
}

func resourceTencentCloudVpcSnapshotPolicyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	snapshotPolicyId := d.Id()

	snapshotPolicies, err := service.DescribeVpcSnapshotPoliciesById(ctx, snapshotPolicyId)
	if err != nil {
		return err
	}

	if len(snapshotPolicies) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcSnapshotPolicyConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	snapshotPolicy := snapshotPolicies[0]

	_ = d.Set("snapshot_policy_id", snapshotPolicyId)

	if snapshotPolicy.Enable != nil {
		_ = d.Set("enable", snapshotPolicy.Enable)
	}

	return nil
}

func resourceTencentCloudVpcSnapshotPolicyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy_config.update")()
	defer inconsistentCheck(d, meta)()

	var (
		enable         bool
		enableRequest  = vpc.NewEnableSnapshotPoliciesRequest()
		disableRequest = vpc.NewDisableSnapshotPoliciesRequest()
	)

	logId := getLogId(contextNil)

	snapshotPolicyId := d.Id()

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	}

	if enable {
		enableRequest.SnapshotPolicyIds = []*string{&snapshotPolicyId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().EnableSnapshotPolicies(enableRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc snapshotPolicyConfig failed, reason:%+v", logId, err)
			return err
		}
	} else {
		disableRequest.SnapshotPolicyIds = []*string{&snapshotPolicyId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DisableSnapshotPolicies(disableRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc snapshotPolicyConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudVpcSnapshotPolicyConfigRead(d, meta)
}

func resourceTencentCloudVpcSnapshotPolicyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
