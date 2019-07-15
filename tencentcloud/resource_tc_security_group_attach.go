package tencentcloud

/*import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func resourceTencentCloudSecurityGroupAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupAttachCreate,
		Read:   resourceTencentCloudSecurityGroupAttachRead,
		Update: resourceTencentCloudSecurityGroupAttachUpdate,
		Delete: resourceTencentCloudSecurityGroupAttachDelete,

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "security group id",
			},

			"security_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group name",
			},

			"cvm_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "cvm ids",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validIdPrefix("ins-"),
				},
			},

			"seni_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "secondary eni ids",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validIdPrefix("eni-"),
				},
			},

			"mysql_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "mysql ids",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validIdPrefix("cdb-"),
				},
			},

			"clb_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "cloud load balance ids",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validIdPrefix("clb-"),
				},
			},
		},
	}
}

func resourceTencentCloudSecurityGroupAttachCreate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group_attach.create")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	sgId := d.Get("security_group_id").(string)

	client := m.(*TencentCloudClient).apiV3Conn

	var hasAttachId bool

	if cvmSet, ok := d.GetOk("cvm_ids"); ok {
		hasAttachId = true

		cvmIds := expandStringList(cvmSet.(*schema.Set).List())
		for _, cvmId := range cvmIds {
			if err := attachSecurityGroupToCvm(ctx, cvmId, sgId, client); err != nil {
				log.Printf("[CRITAL]%s attach cvm %s to security group %s error, reason: %v",
					logId, cvmId, sgId, err)
				return err
			}
		}
	}

	if seniSet, ok := d.GetOk("seni_ids"); ok {
		hasAttachId = true

		seniIds := expandStringList(seniSet.(*schema.Set).List())
		for _, seniId := range seniIds {
			if err := attachSecurityGroupToSeni(ctx, seniId, sgId, client); err != nil {
				log.Printf("[CRITAL]%s attach sencondary eni %s to security group %s error, reason: %v",
					logId, seniId, sgId, err)
				return err
			}
		}
	}

	if mysqlSet, ok := d.GetOk("mysql_ids"); ok {
		hasAttachId = true

		mysqlIds := expandStringList(mysqlSet.(*schema.Set).List())
		for _, mysqlId := range mysqlIds {
			if err := attachSecurityGroupToMysql(ctx, mysqlId, sgId, client); err != nil {
				log.Printf("[CRITAL]%s attach mysql %s to security group %s error, reason: %v",
					logId, mysqlId, sgId, err)
				return err
			}
		}
	}

	if clbSet, ok := d.GetOk("clb_ids"); ok {
		hasAttachId = true

		clbIds := expandStringList(clbSet.(*schema.Set).List())
		for _, clbId := range clbIds {
			if err := attachSecurityGroupToClb(ctx, clbId, sgId, client); err != nil {
				log.Printf("[CRITAL]%s attach cloud load balance %s to security group %s error, reason: %v",
					logId, clbId, sgId, err)
				return err
			}
		}
	}

	if !hasAttachId {
		return errors.New("should has cvm_ids, seni_ids, mysql_ids or clb_ids")
	}

	d.SetId(sgId)

	return resourceTencentCloudSecurityGroupAttachRead(d, m)
}

func resourceTencentCloudSecurityGroupAttachRead(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group_attach.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	sgId := d.Id()

	client := m.(*TencentCloudClient).apiV3Conn

	vpcService := VpcService{client: client}

	sgResp, has, err := vpcService.DescribeSecurityGroup(ctx, sgId)
	if err != nil {
		log.Printf("[DEBUG]%s get security group by id %s error, reason: %v",
			logId, sgId, err)
		return err
	}

	if has == 0 {
		d.SetId("")
		return nil
	}

	_ = d.Set("security_group_name", *sgResp.SecurityGroupName)

	// find out all cvm which attach this security group
	cvmService := CvmService{client: client}

	cvmInstances, err := cvmService.DescribeBySecurityGroups(ctx, sgId)
	if err != nil {
		return err
	}

	cvmIds := make([]string, 0, len(cvmInstances))
	for _, cvmInstance := range cvmInstances {
		cvmIds = append(cvmIds, *cvmInstance.InstanceId)
	}
	_ = d.Set("cvm_ids", cvmIds)

	// find out all secondary eni which attach this security group
	seniInstances, err := vpcService.DescribeNetworkInterfaces(ctx, nil, &sgId)
	if err != nil {
		return err
	}

	seniIds := make([]string, 0, len(seniInstances))
	for _, seniInstance := range seniInstances {
		seniIds = append(seniIds, *seniInstance.NetworkInterfaceId)
	}

	// find out all mysql which attach this security group
	mysqlService := MysqlService{client: client}

	mysqlInstances, err := mysqlService.DescribeDBInstancesBySecurityGroup(ctx, sgId)
	if err != nil {
		return err
	}

	mysqlIds := make([]string, 0, len(mysqlInstances))
	for _, mysqlInstance := range mysqlInstances {
		mysqlIds = append(mysqlIds, *mysqlInstance.InstanceId)
	}

	// find out all clb which attach this security group
	clbService := ClbService{client: client}

	clbInstances, err := clbService.DescribeLoadBalances(ctx, nil, &sgId)
	if err != nil {
		return err
	}

	clbIds := make([]string, 0, len(clbInstances))
	for _, clbInstance := range clbInstances {
		clbIds = append(clbIds, *clbInstance.LoadBalancerId)
	}

	if len(cvmIds) > 0 {
		_ = d.Set("cvm_ids", cvmIds)
	}

	if len(seniIds) > 0 {
		_ = d.Set("seni_ids", seniIds)
	}

	if len(mysqlIds) > 0 {
		_ = d.Set("mysql_ids", mysqlIds)
	}

	if len(clbIds) > 0 {
		_ = d.Set("clb_ids", clbIds)
	}

	return nil
}

func resourceTencentCloudSecurityGroupAttachUpdate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group_attach.update")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	sgId := d.Id()

	client := m.(*TencentCloudClient).apiV3Conn

	d.Partial(true)

	if d.HasChange("cvm_ids") {
		o, n := d.GetChange("cvm_ids")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		unattachCvmIds := expandStringList(oldSet.Difference(newSet).List())
		for _, cvmId := range unattachCvmIds {
			if err := unattachCvmFromSecurityGroup(ctx, cvmId, sgId, client); err != nil {
				return err
			}
		}

		attachCvmIds := expandStringList(newSet.Difference(oldSet).List())
		for _, cvmId := range attachCvmIds {
			if err := attachSecurityGroupToCvm(ctx, cvmId, sgId, client); err != nil {
				return err
			}
		}

		d.SetPartial("cvm_ids")
	}

	if d.HasChange("seni_ids") {
		o, n := d.GetChange("seni_ids")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		unattachSeniIds := expandStringList(oldSet.Difference(newSet).List())
		for _, seniId := range unattachSeniIds {
			if err := unattachSeniFromSecurityGroup(ctx, seniId, sgId, client); err != nil {
				return err
			}
		}

		attachSeniIds := expandStringList(newSet.Difference(oldSet).List())
		for _, seniId := range attachSeniIds {
			if err := attachSecurityGroupToSeni(ctx, seniId, sgId, client); err != nil {
				return err
			}
		}

		d.SetPartial("seni_ids")
	}

	if d.HasChange("mysql_ids") {
		o, n := d.GetChange("mysql_ids")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		unattachMysqlIds := expandStringList(oldSet.Difference(newSet).List())
		for _, mysqlId := range unattachMysqlIds {
			if err := unattachMysqlFromSecurityGroup(ctx, mysqlId, sgId, client); err != nil {
				return err
			}
		}

		attachMysqlIds := expandStringList(newSet.Difference(oldSet).List())
		for _, mysqlId := range attachMysqlIds {
			if err := attachSecurityGroupToMysql(ctx, mysqlId, sgId, client); err != nil {
				return err
			}
		}

		d.SetPartial("mysql_ids")
	}

	if d.HasChange("clb_ids") {
		o, n := d.GetChange("clb_ids")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		unattachClbIds := expandStringList(oldSet.Difference(newSet).List())
		for _, clbId := range unattachClbIds {
			if err := unattachClbFromSecurityGroup(ctx, clbId, sgId, client); err != nil {
				return err
			}
		}

		attachClbIds := expandStringList(newSet.Difference(oldSet).List())
		for _, clbId := range attachClbIds {
			if err := attachSecurityGroupToClb(ctx, clbId, sgId, client); err != nil {
				return err
			}
		}

		d.SetPartial("clb_ids")
	}

	d.Partial(false)

	return resourceTencentCloudSecurityGroupAttachRead(d, m)
}

func resourceTencentCloudSecurityGroupAttachDelete(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group_attach.delete")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	sgId := d.Id()

	client := m.(*TencentCloudClient).apiV3Conn

	if cvmSet, ok := d.GetOk("cvm_ids"); ok {
		cvmIds := expandStringList(cvmSet.(*schema.Set).List())
		for _, cvmId := range cvmIds {
			if err := unattachCvmFromSecurityGroup(ctx, cvmId, sgId, client); err != nil {
				return err
			}
		}
	}

	if seniSet, ok := d.GetOk("seni_ids"); ok {
		seniIds := expandStringList(seniSet.(*schema.Set).List())
		for _, seniId := range seniIds {
			if err := unattachSeniFromSecurityGroup(ctx, seniId, sgId, client); err != nil {
				return err
			}
		}
	}

	if mysqlSet, ok := d.GetOk("mysql_ids"); ok {
		mysqlIds := expandStringList(mysqlSet.(*schema.Set).List())
		for _, mysqlId := range mysqlIds {
			if err := unattachMysqlFromSecurityGroup(ctx, mysqlId, sgId, client); err != nil {
				return err
			}
		}
	}

	if clbSet, ok := d.GetOk("clb_ids"); ok {
		clbIds := expandStringList(clbSet.(*schema.Set).List())
		for _, clbId := range clbIds {
			if err := unattachClbFromSecurityGroup(ctx, clbId, sgId, client); err != nil {
				return err
			}
		}
	}

	return nil
}

func attachSecurityGroupToCvm(ctx context.Context, cvmId, sgId string, client *connectivity.TencentCloudClient) error {
	service := CvmService{client: client}

	existSgIds, err := service.DescribeAssociateSecurityGroups(ctx, cvmId)
	if err != nil {
		return err
	}

	existSgIds = addSecurityGroupToExistSecurityGroups(existSgIds, sgId)

	return service.ModifySecurityGroups(ctx, cvmId, existSgIds)
}

func attachSecurityGroupToSeni(ctx context.Context, seniId, sgId string, client *connectivity.TencentCloudClient) error {
	vpcService := VpcService{client: client}

	senis, err := vpcService.DescribeNetworkInterfaces(ctx, nil, &sgId)
	if err != nil {
		return err
	}

	if len(senis) != 1 {
		err := fmt.Errorf("eni id %s does not exist", seniId)
		return err
	}

	seniInstance := senis[0]
	existSgIds := make([]string, 0, len(seniInstance.GroupSet))
	for _, sgId := range seniInstance.GroupSet {
		existSgIds = append(existSgIds, *sgId)
	}

	existSgIds = addSecurityGroupToExistSecurityGroups(existSgIds, sgId)

	return vpcService.AttachEniToSecurityGroup(ctx, seniId, existSgIds)
}

func attachSecurityGroupToMysql(ctx context.Context, mysqlId, sgId string, client *connectivity.TencentCloudClient) error {
	service := MysqlService{client: client}

	existSgIds, err := service.DescribeDBSecurityGroups(ctx, mysqlId)
	if err != nil {
		return err
	}

	existSgIds = addSecurityGroupToExistSecurityGroups(existSgIds, sgId)

	if len(existSgIds) < 1 {
		panic("exit sg id len < 1")
	}

	if err := service.ModifyDBInstanceSecurityGroups(ctx, mysqlId, existSgIds); err != nil {
		log.Printf("[CRITAL]%s attach security group %s to mysql %s error, reason: %v", GetLogId(ctx), sgId, mysqlId, err)
		return err
	}
	return nil
}

func attachSecurityGroupToClb(ctx context.Context, clbId, sgId string, client *connectivity.TencentCloudClient) error {
	clbService := ClbService{client: client}

	clb, err := clbService.DescribeLoadBalances(ctx, []string{clbId}, nil)
	if err != nil {
		return err
	}

	if len(clb) != 1 {
		err := fmt.Errorf("clb id %s does not exist", clbId)
		return err
	}

	clbInstance := clb[0]
	existSgIds := make([]string, 0, len(clbInstance.SecureGroups))
	for _, sgId := range clbInstance.SecureGroups {
		existSgIds = append(existSgIds, *sgId)
	}

	existSgIds = addSecurityGroupToExistSecurityGroups(existSgIds, sgId)

	return clbService.ModifyLoadBalanceSecurityGroups(ctx, clbId, existSgIds)
}

func unattachCvmFromSecurityGroup(ctx context.Context, cvmId, sgId string, client *connectivity.TencentCloudClient) error {
	service := CvmService{client: client}

	existSgIds, err := service.DescribeAssociateSecurityGroups(ctx, cvmId)
	if err != nil {
		return err
	}

	for i := range existSgIds {
		if existSgIds[i] == sgId {
			existSgIds = append(existSgIds[:i], existSgIds[i+1:]...)
			break
		}
	}

	return service.ModifySecurityGroups(ctx, cvmId, existSgIds)
}

func unattachSeniFromSecurityGroup(ctx context.Context, seniId, sgId string, client *connectivity.TencentCloudClient) error {
	service := VpcService{client: client}

	enis, err := service.DescribeNetworkInterfaces(ctx, []string{seniId}, nil)
	if err != nil {
		return err
	}

	if len(enis) != 1 {
		err := fmt.Errorf("secondary eni %s does not exist", seniId)
		return err
	}

	seniInstance := enis[0]
	existSgIds := make([]string, 0, len(seniInstance.GroupSet))
	for _, existSgId := range seniInstance.GroupSet {
		if *existSgId == sgId {
			continue
		}

		existSgIds = append(existSgIds, *existSgId)
	}

	return service.AttachEniToSecurityGroup(ctx, seniId, existSgIds)
}

func unattachMysqlFromSecurityGroup(ctx context.Context, mysqlId, sgId string, client *connectivity.TencentCloudClient) error {
	service := MysqlService{client: client}

	if err := service.DisassociateSecurityGroup(ctx, mysqlId, sgId); err != nil {
		log.Printf("[CRITAL]%s unattach security group %s from mysql %s error, reason: %v", GetLogId(ctx), sgId, mysqlId, err)
		return err
	}
	return nil
}

func unattachClbFromSecurityGroup(ctx context.Context, clbId, sgId string, client *connectivity.TencentCloudClient) error {
	service := ClbService{client: client}

	clbs, err := service.DescribeLoadBalances(ctx, []string{clbId}, nil)
	if err != nil {
		return err
	}

	if len(clbs) != 1 {
		err := fmt.Errorf("get clb by id %s but return %d clb instances", clbId, len(clbs))
		return err
	}

	clbInstance := clbs[0]
	existSgIds := make([]string, 0, len(clbInstance.SecureGroups))
	for _, existSgId := range clbInstance.SecureGroups {
		if *existSgId == sgId {
			continue
		}

		existSgIds = append(existSgIds, *existSgId)
	}

	return service.ModifyLoadBalanceSecurityGroups(ctx, clbId, existSgIds)
}

func validIdPrefix(prefix string) schema.SchemaValidateFunc {
	return func(v interface{}, s string) (strs []string, errs []error) {
		value := v.(string)
		if !strings.HasPrefix(value, prefix) {
			errs = append(errs, errors.New("id is invalid"))
		}
		return
	}
}

// addSecurityGroupToExistSecurityGroups if id exists, will do nothing
func addSecurityGroupToExistSecurityGroups(existSecurityGroups []string, securityGroupId string) []string {
	var alreadyHasSecurityGroup bool
	for _, sgId := range existSecurityGroups {
		if sgId == securityGroupId {
			alreadyHasSecurityGroup = true
			break
		}
	}

	if !alreadyHasSecurityGroup {
		existSecurityGroups = append(existSecurityGroups, securityGroupId)
	}

	return existSecurityGroups
}*/
