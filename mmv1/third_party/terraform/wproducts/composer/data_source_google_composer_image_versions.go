package composerservice

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta"
)

func DataSourceGoogleComposerImageVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComposerImageVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported_python_versions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComposerImageVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*google.Config)
	userAgent, err := google.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := google.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := google.GetRegion(d, config)
	if err != nil {
		return err
	}

	url, err := google.ReplaceVars(d, config, "{{ComposerBasePath}}projects/{{project}}/locations/{{region}}/imageVersions")
	if err != nil {
		return err
	}

	versions, err := google.PaginatedListRequest(project, url, userAgent, config, flattenGoogleComposerImageVersions)
	if err != nil {
		return fmt.Errorf("Error listing Composer image versions: %s", err)
	}

	log.Printf("[DEBUG] Received Composer Image Versions: %q", versions)

	if err := d.Set("image_versions", versions); err != nil {
		return fmt.Errorf("Error setting image_versions: %s", err)
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/regions/%s", project, region))

	return nil
}

func flattenGoogleComposerImageVersions(resp map[string]interface{}) []interface{} {
	verObjList := resp["imageVersions"].([]interface{})
	versions := make([]interface{}, len(verObjList))
	for i, v := range verObjList {
		verObj := v.(map[string]interface{})
		versions[i] = map[string]interface{}{
			"image_version_id":          verObj["imageVersionId"],
			"supported_python_versions": verObj["supportedPythonVersions"],
		}
	}
	return versions
}
