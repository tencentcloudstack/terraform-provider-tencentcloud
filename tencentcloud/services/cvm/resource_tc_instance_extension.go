package cvm

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"

	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func internetChargeType(k, oldValue, newValue string, d *schema.ResourceData) bool {
	stopMode := d.Get("stopped_mode").(string)
	if stopMode != CVM_STOP_MODE_STOP_CHARGING || !d.HasChange("running_flag") {
		return oldValue == newValue
	}
	return oldValue == "" || newValue == ""
}

func keepImageLogin(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if newValue == "false" && oldValue == "" || oldValue == "false" && newValue == "" {
		return true
	} else {
		return oldValue == newValue
	}
}

func resourceTencentCloudInstanceCreatePostFillRequest0(ctx context.Context, req *cvm.RunInstancesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if v, ok := d.GetOk("instance_charge_type"); ok {
		instanceChargeType := v.(string)
		req.InstanceChargeType = &instanceChargeType
		if instanceChargeType == CVM_CHARGE_TYPE_PREPAID || instanceChargeType == CVM_CHARGE_TYPE_UNDERWRITE {
			req.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{}
			if period, ok := d.GetOk("instance_charge_type_prepaid_period"); ok {
				periodInt64 := int64(period.(int))
				req.InstanceChargePrepaid.Period = &periodInt64
			}
			if renewFlag, ok := d.GetOk("instance_charge_type_prepaid_renew_flag"); ok {
				req.InstanceChargePrepaid.RenewFlag = helper.String(renewFlag.(string))
			}
		}
		if instanceChargeType == CVM_CHARGE_TYPE_SPOTPAID {
			spotInstanceType, sitOk := d.GetOk("spot_instance_type")
			spotMaxPrice, smpOk := d.GetOk("spot_max_price")
			if sitOk || smpOk {
				req.InstanceMarketOptions = &cvm.InstanceMarketOptionsRequest{}
				req.InstanceMarketOptions.MarketType = helper.String(CVM_MARKET_TYPE_SPOT)
				req.InstanceMarketOptions.SpotOptions = &cvm.SpotMarketOptions{}
			}
			if sitOk {
				req.InstanceMarketOptions.SpotOptions.SpotInstanceType = helper.String(strings.ToLower(spotInstanceType.(string)))
			}
			if smpOk {
				req.InstanceMarketOptions.SpotOptions.MaxPrice = helper.String(spotMaxPrice.(string))
			}
		}
		if instanceChargeType == CVM_CHARGE_TYPE_CDHPAID {
			if v, ok := d.GetOk("cdh_instance_type"); ok {
				req.InstanceType = helper.String(v.(string))
			} else {
				return fmt.Errorf("cdh_instance_type can not be empty when instance_charge_type is %s", instanceChargeType)
			}
			if v, ok := d.GetOk("cdh_host_id"); ok {
				req.Placement.HostIds = append(req.Placement.HostIds, helper.String(v.(string)))
			} else {
				return fmt.Errorf("cdh_host_id can not be empty when instance_charge_type is %s", instanceChargeType)
			}
		}
	}

	if v, ok := d.GetOk("placement_group_id"); ok {
		req.DisasterRecoverGroupIds = []*string{helper.String(v.(string))}
	}

	// vpc
	if v, ok := d.GetOk("vpc_id"); ok {
		if v, ok = d.GetOk("private_ip"); ok {
			req.VirtualPrivateCloud.PrivateIpAddresses = []*string{helper.String(v.(string))}
		}
	}

	// login
	keyIds := d.Get("key_ids").(*schema.Set).List()
	if len(keyIds) > 0 {
		req.LoginSettings.KeyIds = helper.InterfacesStringsPoint(keyIds)
	} else if v, ok := d.GetOk("key_name"); ok {
		req.LoginSettings.KeyIds = []*string{helper.String(v.(string))}
	}
	v := d.Get("keep_image_login").(bool)
	if v {
		req.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN)
	} else {
		req.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN_NOT)
	}

	if v, ok := d.GetOk("user_data_raw"); ok {
		userData := base64.StdEncoding.EncodeToString([]byte(v.(string)))
		req.UserData = &userData
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		tags := make([]*cvm.Tag, 0)
		for tagKey, tagValue := range v {
			tag := cvm.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			tags = append(tags, &tag)
		}
		tagSpecification := cvm.TagSpecification{
			ResourceType: helper.String("instance"),
			Tags:         tags,
		}
		req.TagSpecification = append(req.TagSpecification, &tagSpecification)
	}

	return nil
}

func resourceTencentCloudInstanceCreatePreRequest0(ctx context.Context, req *cvm.RunInstancesRequest) *resource.RetryError {
	ratelimit.Check("create")
	return nil
}

func resourceTencentCloudInstanceCreateRequestOnError0(ctx context.Context, req *cvm.RunInstancesRequest, e error) *resource.RetryError {
	logId := ctx.Value(tccommon.LogIdKey)
	log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
		logId, req.GetAction(), req.ToJsonString(), e.Error())

	err, ok := e.(*sdkErrors.TencentCloudSDKError)
	if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, err.Code) {
		return resource.RetryableError(fmt.Errorf("cvm create error: %s, retrying", err.Error()))
	}
	return resource.NonRetryableError(e)
}

func resourceTencentCloudInstanceCreatePostHandleResponse0(ctx context.Context, resp *cvm.RunInstancesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId := *resp.Response.InstanceIdSet[0]

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance != nil && *instance.InstanceState == CVM_STATUS_LAUNCH_FAILED {
			//LatestOperationCodeMode
			return resource.NonRetryableError(fmt.Errorf("cvm instance %s launch failed, this resource will not be stored to tfstate and will auto removed\n.", *instance.InstanceId))
		}
		if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})

	if err != nil {
		return err
	}

	// Wait for the tags attached to the vm since tags attachment it's async while vm creation.
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			actualTags, e := tagService.DescribeResourceTags(ctx, "cvm", "instance", tcClient.Region, instanceId)
			if e != nil {
				return resource.RetryableError(e)
			}
			for tagKey, tagValue := range tags {
				if v, ok := actualTags[tagKey]; !ok || v != tagValue {
					return resource.RetryableError(fmt.Errorf("tag(%s, %s) modification is not completed", tagKey, tagValue))
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if !(d.Get("running_flag").(bool)) {
		stoppedMode := d.Get("stopped_mode").(string)
		err = cvmService.StopInstance(ctx, instanceId, stoppedMode)
		if err != nil {
			return err
		}

		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			if instance != nil && *instance.InstanceState == CVM_STATUS_STOPPED {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudInstanceReadRequestOnSuccess0(ctx context.Context, resp *cvm.Instance) *resource.RetryError {
	if resp != nil && resp.LatestOperationState != nil && *resp.LatestOperationState == "OPERATING" {
		return resource.RetryableError(fmt.Errorf("waiting for instance %s operation", *resp.InstanceId))
	}
	return nil
}

func resourceTencentCloudInstanceReadPreHandleResponse0(ctx context.Context, resp *cvm.Instance) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	instanceId := d.Id()

	if *resp.InstanceState == CVM_STATUS_LAUNCH_FAILED {
		d.SetId("")
		log.Printf("[CRITAL]instance %s not exist or launch failed", instanceId)
		return nil
	}

	var errRet error
	var cvmImages []string
	var response *cvm.DescribeImagesResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := cvm.NewDescribeImagesRequest()
		response, errRet = client.UseCvmClient().DescribeImages(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if *response.Response.TotalCount > 0 {
			for i := range response.Response.ImageSet {
				image := response.Response.ImageSet[i]
				cvmImages = append(cvmImages, *image.ImageId)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if d.Get("image_id").(string) == "" || resp.ImageId == nil || !tccommon.IsContains(cvmImages, *resp.ImageId) {
		_ = d.Set("image_id", resp.ImageId)
	}

	return nil
}

func resourceTencentCloudInstanceReadPostHandleResponse0(ctx context.Context, resp *cvm.Instance) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cbsService := svccbs.NewCbsService(client)

	instance := resp

	if *instance.InstanceChargeType == CVM_CHARGE_TYPE_CDHPAID {
		_ = d.Set("cdh_instance_type", instance.InstanceType)
	}

	if _, ok := d.GetOkExists("allocate_public_ip"); !ok {
		_ = d.Set("allocate_public_ip", len(instance.PublicIpAddresses) > 0)
	}

	tagService := svctag.NewTagService(client)
	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "instance", client.Region, d.Id())
	if err != nil {
		return err
	}
	// as attachment add tencentcloud:autoscaling:auto-scaling-group-id tag automatically
	// we should remove this tag, otherwise it will cause terraform state change
	delete(tags, "tencentcloud:autoscaling:auto-scaling-group-id")
	_ = d.Set("tags", tags)

	//set data_disks
	var hasDataDisks, isCombineDataDisks bool
	dataDiskList := make([]map[string]interface{}, 0, len(instance.DataDisks))
	diskSizeMap := map[string]*uint64{}
	diskOrderMap := make(map[string]int)

	if _, ok := d.GetOk("data_disks"); ok {
		hasDataDisks = true
	}
	if len(instance.DataDisks) > 0 {
		var diskIds []*string
		for i := range instance.DataDisks {
			id := instance.DataDisks[i].DiskId
			size := instance.DataDisks[i].DiskSize
			if id == nil {
				continue
			}
			if strings.HasPrefix(*id, "disk-") {
				diskIds = append(diskIds, id)
			} else {
				diskSizeMap[*id] = helper.Int64Uint64(*size)
			}
		}
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			disks, err := cbsService.DescribeDiskList(ctx, diskIds)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			for i := range disks {
				disk := disks[i]
				if *disk.DiskState == "EXPANDING" {
					return resource.RetryableError(fmt.Errorf("data_disk[%d] is expending", i))
				}
				diskSizeMap[*disk.DiskId] = disk.DiskSize
				if hasDataDisks {
					items := strings.Split(*disk.DiskName, "_")
					diskOrder := items[len(items)-1]
					diskOrderInt, err := strconv.Atoi(diskOrder)
					if err != nil {
						isCombineDataDisks = true
						continue
					}
					diskOrderMap[*disk.DiskId] = diskOrderInt
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	for _, disk := range instance.DataDisks {
		dataDisk := make(map[string]interface{}, 5)
		dataDisk["data_disk_id"] = disk.DiskId
		if disk.DiskId == nil {
			dataDisk["data_disk_size"] = disk.DiskSize
		} else if size, ok := diskSizeMap[*disk.DiskId]; ok {
			dataDisk["data_disk_size"] = size
		}
		dataDisk["data_disk_type"] = disk.DiskType
		dataDisk["data_disk_snapshot_id"] = disk.SnapshotId
		dataDisk["delete_with_instance"] = disk.DeleteWithInstance
		dataDisk["encrypt"] = disk.Encrypt
		dataDisk["throughput_performance"] = disk.ThroughputPerformance
		dataDiskList = append(dataDiskList, dataDisk)
	}
	if hasDataDisks && !isCombineDataDisks {
		sort.SliceStable(dataDiskList, func(idx1, idx2 int) bool {
			dataDiskIdIdx1 := *dataDiskList[idx1]["data_disk_id"].(*string)
			dataDiskIdIdx2 := *dataDiskList[idx2]["data_disk_id"].(*string)
			return diskOrderMap[dataDiskIdIdx1] < diskOrderMap[dataDiskIdIdx2]
		})
	}
	_ = d.Set("data_disks", dataDiskList)

	if len(instance.PrivateIpAddresses) > 0 {
		_ = d.Set("private_ip", instance.PrivateIpAddresses[0])
	}
	if len(instance.PublicIpAddresses) > 0 {
		_ = d.Set("public_ip", instance.PublicIpAddresses[0])
	}
	if len(instance.LoginSettings.KeyIds) > 0 {
		_ = d.Set("key_name", instance.LoginSettings.KeyIds[0])
		_ = d.Set("key_ids", instance.LoginSettings.KeyIds)
	} else {
		_ = d.Set("key_name", "")
		_ = d.Set("key_ids", []*string{})
	}
	if instance.LoginSettings.KeepImageLogin != nil {
		_ = d.Set("keep_image_login", *instance.LoginSettings.KeepImageLogin == CVM_IMAGE_LOGIN)
	}

	if *instance.InstanceState == CVM_STATUS_STOPPED {
		_ = d.Set("running_flag", false)
	} else {
		_ = d.Set("running_flag", true)
	}

	forceDelete := false
	if v, ok := d.GetOkExists("force_delete"); ok {
		forceDelete = v.(bool)
		_ = d.Set("force_delete", forceDelete)
	}

	return nil
}

func resourceTencentCloudInstanceUpdateOnStart(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	d.Partial(true)

	// Get the latest instance info from actual resource.
	instanceInfo, err := cvmService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	var (
		periodSet         = false
		renewFlagSet      = false
		expectChargeType  = CVM_CHARGE_TYPE_POSTPAID
		currentChargeType = *instanceInfo.InstanceChargeType
	)

	chargeType, chargeOk := d.GetOk("instance_charge_type")
	if chargeOk {
		expectChargeType = chargeType.(string)
	}

	if d.HasChange("instance_charge_type") && expectChargeType != currentChargeType {
		var (
			period    = -1
			renewFlag string
		)

		if v, ok := d.GetOk("instance_charge_type_prepaid_period"); ok {
			period = v.(int)
		}
		if v, ok := d.GetOk("instance_charge_type_prepaid_renew_flag"); ok {
			renewFlag = v.(string)
		}
		// change charge type
		err := cvmService.ModifyInstanceChargeType(ctx, instanceId, expectChargeType, period, renewFlag)
		if err != nil {
			return err
		}
		// query cvm status
		err = waitForOperationFinished(d, meta, 5*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}
		periodSet = true
		renewFlagSet = true
	}

	// When instance is prepaid but period was empty and set to 1, skip this case.
	op, np := d.GetChange("instance_charge_type_prepaid_period")
	if _, ok := op.(int); !ok && np.(int) == 1 {
		periodSet = true
	}
	if d.HasChange("instance_charge_type_prepaid_period") && !periodSet {
		chargeType := d.Get("instance_charge_type").(string)
		period := d.Get("instance_charge_type_prepaid_period").(int)
		renewFlag := ""

		if v, ok := d.GetOk("instance_charge_type_prepaid_renew_flag"); ok {
			renewFlag = v.(string)
		}
		err := cvmService.ModifyInstanceChargeType(ctx, instanceId, chargeType, period, renewFlag)
		if err != nil {
			return err
		}
		// query cvm status
		err = waitForOperationFinished(d, meta, 5*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}
		renewFlagSet = true
	}

	if d.HasChange("instance_charge_type_prepaid_renew_flag") && !renewFlagSet {
		//renew api
		err := cvmService.ModifyRenewParam(ctx, instanceId, d.Get("instance_charge_type_prepaid_renew_flag").(string))
		if err != nil {
			return err
		}

		//check success
		err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}

		time.Sleep(tccommon.ReadRetryTimeout)
	}

	// Reset Instance
	// Keep Login Info
	if d.HasChange("image_id") ||
		d.HasChange("hostname") ||
		d.HasChange("disable_security_service") ||
		d.HasChange("disable_monitor_service") ||
		d.HasChange("keep_image_login") {

		request := cvm.NewResetInstanceRequest()
		request.InstanceId = helper.String(d.Id())

		if v, ok := d.GetOk("image_id"); ok {
			request.ImageId = helper.String(v.(string))
		}
		if v, ok := d.GetOk("hostname"); ok {
			request.HostName = helper.String(v.(string))
		}

		// enhanced service
		request.EnhancedService = &cvm.EnhancedService{}
		if d.HasChange("disable_security_service") {
			v := d.Get("disable_security_service")
			securityService := v.(bool)
			request.EnhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
				Enabled: &securityService,
			}
		}

		if d.HasChange("disable_monitor_service") {
			v := d.Get("disable_monitor_service")
			monitorService := !(v.(bool))
			request.EnhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
				Enabled: &monitorService,
			}
		}

		// Modify or keep login info when instance reset
		request.LoginSettings = &cvm.LoginSettings{}

		if v, ok := d.GetOk("password"); ok {
			request.LoginSettings.Password = helper.String(v.(string))
		}

		if v, ok := d.GetOk("key_ids"); ok {
			request.LoginSettings.KeyIds = helper.InterfacesStringsPoint(v.(*schema.Set).List())
		} else if v, ok := d.GetOk("key_name"); ok {
			request.LoginSettings.KeyIds = []*string{helper.String(v.(string))}
		}

		if v := d.Get("keep_image_login").(bool); v {
			request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN)
		} else {
			request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN_NOT)
		}

		if err := cvmService.ResetInstance(ctx, request); err != nil {
			return err
		}

		// Modify Login Info Directly
	} else {
		if d.HasChange("password") {
			err := cvmService.ModifyPassword(ctx, instanceId, d.Get("password").(string))
			if err != nil {
				return err
			}
			err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
			if err != nil {
				return err
			}
		}

		if d.HasChange("key_name") {
			o, n := d.GetChange("key_name")
			oldKeyId := o.(string)
			keyId := n.(string)

			if oldKeyId != "" {
				err := cvmService.UnbindKeyPair(ctx, []*string{&oldKeyId}, []*string{&instanceId})
				if err != nil {
					return err
				}
				err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
				if err != nil {
					return err
				}
			}

			if keyId != "" {
				err = cvmService.BindKeyPair(ctx, []*string{&keyId}, instanceId)
				if err != nil {
					return err
				}
				err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
				if err != nil {
					return err
				}
			}
		}

		// support remove old `key_name` to `key_ids`, so do not follow "else"
		if d.HasChange("key_ids") {
			o, n := d.GetChange("key_ids")
			ov := o.(*schema.Set)

			nv := n.(*schema.Set)

			adds := nv.Difference(ov)
			removes := ov.Difference(nv)
			adds.Remove("")
			removes.Remove("")

			if removes.Len() > 0 {
				err := cvmService.UnbindKeyPair(ctx, helper.InterfacesStringsPoint(removes.List()), []*string{&instanceId})
				if err != nil {
					return err
				}
				err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
				if err != nil {
					return err
				}
			}
			if adds.Len() > 0 {
				err = cvmService.BindKeyPair(ctx, helper.InterfacesStringsPoint(adds.List()), instanceId)
				if err != nil {
					return err
				}
				err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("data_disks") {
		o, n := d.GetChange("data_disks")
		ov := o.([]interface{})
		nv := n.([]interface{})

		if len(ov) != len(nv) {
			return fmt.Errorf("error: data disk count has changed (%d -> %d) but doesn't support add or remove for now", len(ov), len(nv))
		}

		cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		for i := range nv {
			sizeKey := fmt.Sprintf("data_disks.%d.data_disk_size", i)
			idKey := fmt.Sprintf("data_disks.%d.data_disk_id", i)
			if !d.HasChange(sizeKey) {
				continue
			}
			size := d.Get(sizeKey).(int)
			diskId := d.Get(idKey).(string)

			err := cbsService.ResizeDisk(ctx, diskId, size)

			if err != nil {
				return fmt.Errorf("an error occurred when modifying %s, reason: %s", sizeKey, err.Error())
			}

		}
	}

	var flag bool
	if d.HasChange("running_flag") {
		flag = d.Get("running_flag").(bool)
		if err := switchInstance(&cvmService, ctx, d, flag); err != nil {
			return err
		}

	}

	return nil
}

func waitForOperationFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration, state string, immediately bool) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cvmService := CvmService{client}
	instanceId := d.Id()
	// We cannot catch LatestOperationState change immediately after modification returns, we must wait for LatestOperationState update to expected.
	if !immediately {
		time.Sleep(time.Second * 10)
	}

	err := resource.Retry(timeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance == nil {
			return resource.NonRetryableError(fmt.Errorf("%s not exists", instanceId))
		}
		if instance.LatestOperationState == nil {
			return resource.RetryableError(fmt.Errorf("wait for operation update"))
		}
		if *instance.LatestOperationState == state {
			return resource.RetryableError(fmt.Errorf("waiting for instance %s operation", instanceId))
		}
		if *instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_FAILED {
			return resource.NonRetryableError(fmt.Errorf("failed operation"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func switchInstance(cvmService *CvmService, ctx context.Context, d *schema.ResourceData, flag bool) (err error) {
	instanceId := d.Id()
	if flag {
		err = cvmService.StartInstance(ctx, instanceId)
		if err != nil {
			return err
		}
		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
		})
		if err != nil {
			return err
		}
	} else {
		stoppedMode := d.Get("stopped_mode").(string)
		skipStopApi := false
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			// when retry polling instance status, stop instance should skipped
			if !skipStopApi {
				err := cvmService.StopInstance(ctx, instanceId, stoppedMode)
				if err != nil {
					return resource.NonRetryableError(err)
				}
			}
			instance, err := cvmService.DescribeInstanceById(ctx, instanceId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if instance == nil {
				return resource.NonRetryableError(fmt.Errorf("instance %s not found", instanceId))
			}

			if instance.LatestOperationState != nil {
				operationState := *instance.LatestOperationState
				if operationState == "OPERATING" {
					skipStopApi = true
					return resource.RetryableError(fmt.Errorf("instance %s stop operating, retrying", instanceId))
				}
				if operationState == "FAILED" {
					skipStopApi = false
					return resource.RetryableError(fmt.Errorf("instance %s stop failed, retrying", instanceId))
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			if instance != nil && *instance.InstanceState == CVM_STATUS_STOPPED {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTencentCloudInstanceUpdatePreHandleResponse5(ctx context.Context, resp *cvm.ResizeInstanceDisksResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()
	size := d.Get("system_disk_size").(int)
	diskType := d.Get("system_disk_type").(string)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, err := cvmService.DescribeInstanceById(ctx, instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if instance != nil && instance.LatestOperationState != nil {
			if *instance.InstanceState == "FAILED" {
				return resource.NonRetryableError(fmt.Errorf("instance operation failed"))
			}
			if *instance.InstanceState == "OPERATING" {
				return resource.RetryableError(fmt.Errorf("instance operating"))
			}
		}
		if instance != nil && instance.SystemDisk != nil {
			//wait until disk result as expected
			if *instance.SystemDisk.DiskType != diskType || int(*instance.SystemDisk.DiskSize) != size {
				return resource.RetryableError(fmt.Errorf("waiting for expanding success"))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudInstanceUpdatePreHandleResponse6(ctx context.Context, resp *cvm.ResetInstancesTypeResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	err := waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudInstanceUpdatePreHandleResponse7(ctx context.Context, resp *cvm.ResetInstancesTypeResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	err := waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudInstanceUpdatePostFillRequest8(ctx context.Context, req *cvm.ModifyInstancesVpcAttributeRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	privateIp := d.Get("private_ip").(string)
	if privateIp != "" {
		req.VirtualPrivateCloud.PrivateIpAddresses = []*string{&privateIp}
	}
	return nil
}

func resourceTencentCloudInstanceUpdateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("cvm", "instance", region, instanceId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		//except instance ,system disk and data disk will be tagged
		//keep logical consistence with the console
		//tag system disk
		if systemDiskId, ok := d.GetOk("system_disk_id"); ok {
			if systemDiskId.(string) != "" {
				resourceName = tccommon.BuildTagResourceName("cvm", "volume", region, systemDiskId.(string))
				if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
					return err
				}
			}
		}
		//tag disk ids
		if dataDisks, ok := d.GetOk("data_disks"); ok {
			dataDiskList := dataDisks.([]interface{})
			for _, dataDisk := range dataDiskList {
				disk := dataDisk.(map[string]interface{})
				dataDiskId := disk["data_disk_id"].(string)
				resourceName = tccommon.BuildTagResourceName("cvm", "volume", region, dataDiskId)
				if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
					return err
				}
			}
		}

	}

	if d.HasChange("internet_max_bandwidth_out") {
		chargeType := d.Get("internet_charge_type").(string)
		bandWidthOut := int64(d.Get("internet_max_bandwidth_out").(int))
		if chargeType != "TRAFFIC_POSTPAID_BY_HOUR" && chargeType != "BANDWIDTH_POSTPAID_BY_HOUR" && chargeType != "BANDWIDTH_PACKAGE" {
			return fmt.Errorf("charge type should be one of `TRAFFIC_POSTPAID_BY_HOUR BANDWIDTH_POSTPAID_BY_HOUR BANDWIDTH_PACKAGE` when adjusting internet_max_bandwidth_out")
		}

		err := cvmService.ModifyInternetMaxBandwidthOut(ctx, instanceId, chargeType, bandWidthOut)
		if err != nil {
			return err
		}

		err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}

	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudInstanceDeletePostFillRequest0(ctx context.Context, req *cvm.TerminateInstancesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	_, err := cvmService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudInstanceDeletePreHandleResponse0(ctx context.Context, resp *cvm.TerminateInstancesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	//check is force delete or not
	forceDelete := d.Get("force_delete").(bool)

	//check recycling
	notExist := false

	var instance *cvm.Instance
	//check exist
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance == nil {
			notExist = true
			return nil
		}
		if *instance.InstanceState == CVM_STATUS_SHUTDOWN && *instance.LatestOperationState != CVM_LATEST_OPERATION_STATE_OPERATING {
			//in recycling
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})
	if err != nil {
		return err
	}

	vpcService := vpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	if notExist {
		err := waitIpRelease(ctx, vpcService, instance)
		if err != nil {
			return err
		}
		return nil
	}

	if !forceDelete {
		return nil
	}

	return nil
}

func resourceTencentCloudInstanceDeleteRequestOnError1(ctx context.Context, e error) *resource.RetryError {
	//check InvalidInstanceState.Terminating
	ee, ok := e.(*sdkErrors.TencentCloudSDKError)
	if !ok {
		return tccommon.RetryError(e)
	}
	if ee.Code == "InvalidInstanceState.Terminating" {
		return nil
	}
	return tccommon.RetryError(e, "OperationDenied.InstanceOperationInProgress")
}

func resourceTencentCloudInstanceDeletePreHandleResponse1(ctx context.Context, resp *cvm.TerminateInstancesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := ctx.Value(tccommon.LogIdKey)

	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	vpcService := vpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceId := d.Id()

	var instance *cvm.Instance
	//describe and check not exist
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("data_disks"); ok {
		dataDisks := v.([]interface{})
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskId := value["data_disk_id"].(string)
			deleteWithInstance := value["delete_with_instance"].(bool)
			if deleteWithInstance {
				cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
				err := resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
					diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}
					if *diskInfo.DiskState != svccbs.CBS_STORAGE_STATUS_UNATTACHED {
						return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
					return err
				}
				err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					e := cbsService.DeleteDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
					return err
				}
				err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
					diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}
					if *diskInfo.DiskState == svccbs.CBS_STORAGE_STATUS_TORECYCLE {
						return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s read cbs status failed, reason:%s\n ", logId, err.Error())
					return err
				}
				err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					e := cbsService.DeleteDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
					return err
				}
				err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
					diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}
					if diskInfo != nil {
						return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s read cbs status failed, reason:%s\n ", logId, err.Error())
					return err
				}
			}
		}
	}

	err = waitIpRelease(ctx, vpcService, instance)
	if err != nil {
		return err
	}
	return nil
}

func waitIpRelease(ctx context.Context, vpcService vpc.VpcService, instance *cvm.Instance) error {
	// wait ip release
	if len(instance.PrivateIpAddresses) > 0 {
		params := make(map[string]interface{})
		params["VpcId"] = instance.VirtualPrivateCloud.VpcId
		params["SubnetId"] = instance.VirtualPrivateCloud.SubnetId
		params["IpAddresses"] = instance.PrivateIpAddresses
		err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			usedIpAddress, errRet := vpcService.DescribeVpcUsedIpAddressByFilter(ctx, params)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			if len(usedIpAddress) > 0 {
				return resource.RetryableError(fmt.Errorf("wait cvm private ip release..."))
			}

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
