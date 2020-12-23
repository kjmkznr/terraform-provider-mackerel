package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelDashboardCreate,
		Read:   resourceMackerelDashboardRead,
		Update: resourceMackerelDashboardUpdate,
		Delete: resourceMackerelDashboardDelete,
		Exists: resourceMackerelDashboardExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"body_markdown": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url_path": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateUrlPathWordDiag(),
			},
		},
	}
}

func resourceMackerelDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.Dashboard{
		Title:        d.Get("title").(string),
		BodyMarkDown: d.Get("body_markdown").(string),
		URLPath:      d.Get("url_path").(string),
	}

	dashboard, err := client.CreateDashboard(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel dashboard %q created.", dashboard.ID)
	d.SetId(dashboard.ID)

	return resourceMackerelDashboardRead(d, meta)
}

func resourceMackerelDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel dashboard: %q", d.Id())
	dashboard, err := client.FindDashboard(d.Id())
	if err != nil {
		return err
	}

	d.SetId(dashboard.ID)
	_ = d.Set("title", dashboard.Title)
	_ = d.Set("body_markdown", dashboard.BodyMarkDown)
	_ = d.Set("url_path", dashboard.URLPath)

	return nil
}

func resourceMackerelDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.Dashboard{
		Title:        d.Get("title").(string),
		BodyMarkDown: d.Get("body_markdown").(string),
		URLPath:      d.Get("url_path").(string),
	}

	_, err := client.UpdateDashboard(d.Id(), input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel dashboard %q updated.", d.Id())
	return resourceMackerelDashboardRead(d, meta)
}

func resourceMackerelDashboardExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*mackerel.Client)
	dashboards, err := client.FindDashboards()
	if err != nil {
		return false, err
	}

	for _, dashboard := range dashboards {
		if dashboard.ID == d.Id() {
			return true, nil
		}
	}

	return false, nil
}

func resourceMackerelDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteDashboard(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel dashboard %q deleted.", d.Id())

	return nil
}
