package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudFCFunction_basic(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudfcfunction-%d", rand)
	var basicMap = map[string]string{
		"service":     CHECKSET,
		"name":        name,
		"runtime":     "python2.7",
		"description": "tf",
		"handler":     "hello.handler",
		"oss_bucket":  CHECKSET,
		"oss_key":     CHECKSET,
	}
	resourceId := "alicloud_fc_function.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFCFunctionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service":     "${alicloud_fc_service.default.name}",
					"name":        "${var.name}",
					"runtime":     "python2.7",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  "${alicloud_oss_bucket.default.id}",
					"oss_key":     "${alicloud_oss_bucket_object.default.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "filename", "oss_bucket", "oss_key"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf unit test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf unit test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": map[string]string{
						"test":   "terraform",
						"prefix": "tfAcc",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.test":   "terraform",
						"environment_variables.prefix": "tfAcc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "512",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": "python3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "python3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service":     "${alicloud_fc_service.default.name}",
					"name":        "${var.name}",
					"runtime":     "python2.7",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  "${alicloud_oss_bucket.default.id}",
					"oss_key":     "${alicloud_oss_bucket_object.default.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func TestAccAlicloudFCFunctionMulti(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudfcfunction-%d", rand)
	var basicMap = map[string]string{
		"service":     CHECKSET,
		"name":        name + "-9",
		"runtime":     "python2.7",
		"description": "tf",
		"handler":     "hello.handler",
		"oss_bucket":  CHECKSET,
		"oss_key":     CHECKSET,
	}
	resourceId := "alicloud_fc_function.default.9"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFCFunctionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "10",
					"service":     "${alicloud_fc_service.default.name}",
					"name":        "${var.name}-${count.index}",
					"runtime":     "python2.7",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  "${alicloud_oss_bucket.default.id}",
					"oss_key":     "${alicloud_oss_bucket_object.default.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceFCFunctionConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%v"
}
resource "alicloud_log_project" "default" {
  name = "${var.name}"
  description = "tf unit test"
}

resource "alicloud_log_store" "default" {
  project = "${alicloud_log_project.default.name}"
  name = "${var.name}"
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_fc_service" "default" {
    name = "${var.name}"
    description = "tf unit test"
    log_config {
	project = "${alicloud_log_project.default.name}"
	logstore = "${alicloud_log_store.default.name}"
    }
    role = "${alicloud_ram_role.default.arn}"
    depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = "${alicloud_oss_bucket.default.id}"
  key = "fc/hello.zip"
  content = <<EOF
  	# -*- coding: utf-8 -*-
	def handler(event, context):
	    print "hello world"
	    return 'hello world'
  EOF
}

resource "alicloud_ram_role" "default" {
  name = "${var.name}"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
`, name, testFCRoleTemplate)
}
