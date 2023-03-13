package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDbaasLogsCluster() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return dataSourceDbaasLogsClusterRead(d, meta)
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			// Computed
			"cluster_type": {
				Type:        schema.TypeString,
				Description: "Cluster type",
				Computed:    true,
			},
			"dedicated_input_pem": {
				Type:        schema.TypeString,
				Description: "PEM for dedicated inputs",
				Computed:    true,
				Sensitive:   true,
			},
			"archive_allowed_networks": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Allowed networks for ARCHIVE flow type",
				Computed:    true,
			},
			"direct_input_allowed_networks": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Allowed networks for DIRECT_INPUT flow type",
				Computed:    true,
			},
			"direct_input_pem": {
				Type:        schema.TypeString,
				Description: "PEM for direct inputs",
				Computed:    true,
				Sensitive:   true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "hostname",
				Computed:    true,
			},
			"is_default": {
				Type:        schema.TypeBool,
				Description: "All content generated by given service will be placed on this cluster",
				Computed:    true,
			},
			"is_unlocked": {
				Type:        schema.TypeBool,
				Description: "Allow given service to perform advanced operations on cluster",
				Computed:    true,
			},
			"query_allowed_networks": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Allowed networks for QUERY flow type",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Data center localization",
				Computed:    true,
			},
		},
	}
}

func dbaasGetClusterID(config *Config, serviceName string) (string, error) {
	res := []string{}

	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/cluster",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return "", fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	return res[0], nil
}

func dataSourceDbaasLogsClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read dbaas logs cluster %s", serviceName)

	cluster_id, err := dbaasGetClusterID(config, serviceName)

	if err != nil {
		return fmt.Errorf("Error fetching info for %s:\n\t %q", serviceName, err)
	}

	d.SetId(cluster_id)

	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/cluster/%s",
		url.PathEscape(serviceName),
		url.PathEscape(cluster_id),
	)

	res := map[string]interface{}{}
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.Set("archive_allowed_networks", res["archiveAllowedNetworks"])
	d.Set("cluster_type", res["clusterType"])
	d.Set("dedicated_input_pem", res["dedicatedInputPEM"])
	d.Set("direct_input_allowed_networks", res["directInputAllowedNetworks"])
	d.Set("direct_input_pem", res["directInputPEM"])
	d.Set("hostname", res["hostname"])
	d.Set("is_default", res["isDefault"])
	d.Set("is_unlocked", res["isUnlocked"])
	d.Set("query_allowed_networks", res["queryAllowedNetworks"])
	d.Set("region", res["region"])

	return nil
}
