package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectDatabaseOpensearchUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectDatabaseOpensearchUserCreate,
		Read:   resourceCloudProjectDatabaseOpensearchUserRead,
		Delete: resourceCloudProjectDatabaseOpensearchUserDelete,
		Update: resourceCloudProjectDatabaseOpensearchUserUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseOpensearchUserImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"acls": {
				Type:        schema.TypeSet,
				Description: "Acls of the user",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pattern": {
							Type:        schema.TypeString,
							Description: "Pattern of the ACL",
							Required:    true,
						},
						"permission": {
							Type:        schema.TypeString,
							Description: "Permission of the ACL",
							Required:    true,
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the user",
				ForceNew:    true,
				Required:    true,
			},
			"password_reset": {
				Type:        schema.TypeString,
				Description: "Arbitrary string to change to trigger a password update",
				Optional:    true,
			},

			//Computed
			"created_at": {
				Type:        schema.TypeString,
				Description: "Date of the creation of the user",
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password of the user",
				Sensitive:   true,
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the user",
				Computed:    true,
			},
		},
	}
}

func resourceCloudProjectDatabaseOpensearchUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return importCloudProjectDatabaseUser(d, meta)
}

func resourceCloudProjectDatabaseOpensearchUserCreate(d *schema.ResourceData, meta interface{}) error {
	f := func() interface{} {
		return (&CloudProjectDatabaseOpensearchUserCreateOpts{}).FromResource(d)
	}
	return postCloudProjectDatabaseUser(d, meta, "opensearch", dataSourceCloudProjectDatabaseOpensearchUserRead, resourceCloudProjectDatabaseOpensearchUserRead, resourceCloudProjectDatabaseOpensearchUserUpdate, f)
}

func resourceCloudProjectDatabaseOpensearchUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseOpensearchUserResponse{}

	log.Printf("[DEBUG] Will read user %s from cluster %s from project %s", id, clusterId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read user %+v", res)
	return nil
}

func resourceCloudProjectDatabaseOpensearchUserUpdate(d *schema.ResourceData, meta interface{}) error {
	f := func() interface{} {
		return (&CloudProjectDatabaseOpensearchUserUpdateOpts{}).FromResource(d)
	}
	return updateCloudProjectDatabaseUser(d, meta, "opensearch", resourceCloudProjectDatabaseOpensearchUserRead, f)
}

func resourceCloudProjectDatabaseOpensearchUserDelete(d *schema.ResourceData, meta interface{}) error {
	return deleteCloudProjectDatabaseUser(d, meta, "opensearch")
}
