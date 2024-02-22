package clb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClbTargetGroupAttachmentsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupAttachments,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "load_balancer_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "associations.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbTargetGroupAttachmentsResource_target(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupAttachmentsTarget,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "target_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "associations.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbTargetGroupAttachmentsResource_sync(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupAttachmentsParallel,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "target_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "associations.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments1", "target_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments1", "associations.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbTargetGroupAttachmentsResource_tcp(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupAttachmentsTCP,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "target_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "associations.#"),
				),
			},
		},
	})
}

const testAccClbTargetGroupAttachments = `

resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach_example"
  vpc_id = "vpc-efc9vddt"
}

resource "tencentcloud_clb_listener" "public_listeners" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  #  protocol      = "HTTPS"
  #  port          = "443"
  protocol      = "HTTP"
  port          = "8090"
  listener_name = "iac-test-attach-2"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic2" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "baidu.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic3" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "tencent.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic4" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "aws.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  load_balancer_id = tencentcloud_clb_instance.clb_basic.id
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    target_group_id = "lbtg-f0od5ack"
    location_id = tencentcloud_clb_listener_rule.rule_basic.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    target_group_id = "lbtg-h7lzmdts"
    location_id = tencentcloud_clb_listener_rule.rule_basic2.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    target_group_id = "lbtg-577temh4"
    location_id = tencentcloud_clb_listener_rule.rule_basic3.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    target_group_id = "lbtg-119vqig8"
    location_id = tencentcloud_clb_listener_rule.rule_basic4.rule_id
  }
  depends_on = [tencentcloud_clb_listener.public_listeners,
    tencentcloud_clb_listener_rule.rule_basic4,
    tencentcloud_clb_listener_rule.rule_basic3,
    tencentcloud_clb_listener_rule.rule_basic2,
    tencentcloud_clb_listener_rule.rule_basic]
}

`
const testAccClbTargetGroupAttachmentsTarget = `

resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach_1"
  vpc_id = "vpc-efc9vddt"
}
resource "tencentcloud_clb_instance" "clb_basic2" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach_2"
  vpc_id = "vpc-efc9vddt"
}
resource "tencentcloud_clb_instance" "clb_basic3" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach_3"
  vpc_id = "vpc-efc9vddt"
}
resource "tencentcloud_clb_instance" "clb_basic4" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach_4"
  vpc_id = "vpc-efc9vddt"
}
resource "tencentcloud_clb_instance" "clb_basic5" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach_5"
  vpc_id = "vpc-efc9vddt"
}


resource "tencentcloud_clb_listener" "public_listeners" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  protocol      = "HTTP"
  port          = "8090"
  listener_name = "iac-test-attach-2"
}
resource "tencentcloud_clb_listener" "public_listeners2" {
  clb_id        = tencentcloud_clb_instance.clb_basic2.id
  protocol      = "HTTP"
  port          = "8090"
  listener_name = "iac-test-attach-3"
}
resource "tencentcloud_clb_listener" "public_listeners3" {
  clb_id        = tencentcloud_clb_instance.clb_basic3.id
  protocol      = "HTTP"
  port          = "8090"
  listener_name = "iac-test-attach-4"
}
resource "tencentcloud_clb_listener" "public_listeners4" {
  clb_id        = tencentcloud_clb_instance.clb_basic4.id
  protocol      = "HTTP"
  port          = "8090"
  listener_name = "iac-test-attach-5"
}
resource "tencentcloud_clb_listener" "public_listeners5" {
  clb_id        = tencentcloud_clb_instance.clb_basic5.id
  protocol      = "HTTP"
  port          = "8090"
  listener_name = "iac-test-attach-6"
}
resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic2" {
  clb_id              = tencentcloud_clb_instance.clb_basic2.id
  listener_id         = tencentcloud_clb_listener.public_listeners2.listener_id
  domain              = "baidu.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic3" {
  clb_id              = tencentcloud_clb_instance.clb_basic3.id
  listener_id         = tencentcloud_clb_listener.public_listeners3.listener_id
  domain              = "tencent.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic4" {
  clb_id              = tencentcloud_clb_instance.clb_basic4.id
  listener_id         = tencentcloud_clb_listener.public_listeners4.listener_id
  domain              = "aws.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic5" {
  clb_id              = tencentcloud_clb_instance.clb_basic5.id
  listener_id         = tencentcloud_clb_listener.public_listeners5.listener_id
  domain              = "aws2.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  target_group_id = "lbtg-nxd0dmcm"
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    load_balancer_id = tencentcloud_clb_instance.clb_basic.id
    location_id = tencentcloud_clb_listener_rule.rule_basic.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners2.listener_id
    load_balancer_id = tencentcloud_clb_instance.clb_basic2.id
    location_id = tencentcloud_clb_listener_rule.rule_basic2.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners3.listener_id
    load_balancer_id = tencentcloud_clb_instance.clb_basic3.id

    location_id = tencentcloud_clb_listener_rule.rule_basic3.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners4.listener_id
    load_balancer_id = tencentcloud_clb_instance.clb_basic4.id
    location_id = tencentcloud_clb_listener_rule.rule_basic4.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners5.listener_id
    load_balancer_id = tencentcloud_clb_instance.clb_basic5.id
    location_id = tencentcloud_clb_listener_rule.rule_basic5.rule_id
  }
  depends_on = [
    tencentcloud_clb_listener.public_listeners,
    tencentcloud_clb_listener.public_listeners2,
    tencentcloud_clb_listener.public_listeners3,
    tencentcloud_clb_listener.public_listeners4,
    tencentcloud_clb_listener.public_listeners5,
    tencentcloud_clb_listener_rule.rule_basic5,
    tencentcloud_clb_listener_rule.rule_basic4,
    tencentcloud_clb_listener_rule.rule_basic3,
    tencentcloud_clb_listener_rule.rule_basic2,
    tencentcloud_clb_listener_rule.rule_basic]
}
`
const testAccClbTargetGroupAttachmentsParallel = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach_1"
  vpc_id = "vpc-efc9vddt"
}
resource "tencentcloud_clb_listener" "public_listeners" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  protocol      = "HTTP"
  port          = "8090"
  listener_name = "iac-test-attach"
}
resource "tencentcloud_clb_listener" "public_listeners2" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  protocol      = "HTTP"
  port          = "8099"
  listener_name = "iac-test-attach-2"
}
resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic2" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners2.listener_id
  domain              = "baidu.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  target_group_id = "lbtg-fa1l7oh2"
  associations  {
    load_balancer_id = tencentcloud_clb_instance.clb_basic.id
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    location_id=tencentcloud_clb_listener_rule.rule_basic.rule_id
  }
  depends_on = [tencentcloud_clb_listener_rule.rule_basic,
    tencentcloud_clb_listener_rule.rule_basic,
    tencentcloud_clb_listener.public_listeners,
    tencentcloud_clb_listener.public_listeners2]
}
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments1" {
  target_group_id = "lbtg-aw5hicve"
  associations  {
    load_balancer_id = tencentcloud_clb_instance.clb_basic.id
    listener_id = tencentcloud_clb_listener.public_listeners2.listener_id
    location_id=tencentcloud_clb_listener_rule.rule_basic2.rule_id

  }
}
`
const testAccClbTargetGroupAttachmentsTCP = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach_1"
  vpc_id = "vpc-efc9vddt"
}
resource "tencentcloud_clb_listener" "public_listeners" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  protocol      = "TCP"
  port          = "8090"
  listener_name = "iac-test-attach"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  target_group_id = "lbtg-nxd0dmcm"
  associations  {
    load_balancer_id = tencentcloud_clb_instance.clb_basic.id
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
  }
}
`
