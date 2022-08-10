/*
Provides a datasource to copy a CBS snapshot.Snapshot replication across regions

Example Usage

```hcl
data "tencentcloud_cbs_copy_snapshot_cross_regions" "example" {
  destination_regions  = ["ap-beijing"]
  snapshot_id          = ""
  snapshot_name        = ""
}
```

*/
package tencentcloud

import (
	"crypto/md5"
	"fmt"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCbsCopySnapshotCrossRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCbsCopySnapshotCosssRegionsRead,

		Schema: map[string]*schema.Schema{
			"destination_regions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The target region to which the snapshot needs to be copied. The standard values of each region can be queried through the interface DescribeRegions, and can only be passed to regions that support snapshots.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"snapshot_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the snapshot.",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the the CBS which this snapshot created from.",
			},
			"snapshot_copy_result_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Snapshot results of cross-region replication.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The new snapshot ID copied to the target region.",
						},
						"destination_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target region for cross-replication.",
						},
						"code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Error code, the value is 'Success' when successful.",
						},
						"message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates specific error information, an empty string on success.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsCopySnapshotCosssRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cbs_copy_snapshot_cross_regions.read")()

	var (
		snapshotCopy []interface{}
	)

	logId := getLogId(contextNil)
	request := cbs.NewCopySnapshotCrossRegionsRequest()

	if v, ok := d.GetOk("destination_regions"); ok {
		regions := v.(*schema.Set).List()
		regionsArr := make([]*string, 0, len(regions))
		for _, receiver := range regions {
			regionsArr = append(regionsArr, helper.String(receiver.(string)))
		}
		request.DestinationRegions = regionsArr
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("snapshot_name"); ok {
		request.SnapshotName = helper.String(v.(string))
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().CopySnapshotCrossRegions(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err, InternalError)
		}
		snapshotCopyResultSet := response.Response.SnapshotCopyResultSet

		md := md5.New()
		id := fmt.Sprintf("%x", md.Sum(nil))
		d.SetId(id)

		for _, noticesItem := range snapshotCopyResultSet {
			resultItemMap := map[string]interface{}{
				"snapshot_id":        noticesItem.SnapshotId,
				"destination_region": noticesItem.DestinationRegion,
				"code":               noticesItem.Code,
				"message":            noticesItem.Message,
			}
			snapshotCopy = append(snapshotCopy, resultItemMap)
		}
		if err = d.Set("snapshot_copy_result_set", snapshotCopy); err != nil {
			return nil
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
