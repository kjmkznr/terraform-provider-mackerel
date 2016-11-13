package mackerel

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelDashboardCreate,
		Read:   resourceMackerelDashboardRead,
		Update: resourceMackerelDashboardUpdate,
		Delete: resourceMackerelDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"body_markdown": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url_path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUrlPathWord,
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

	d.Set("id", dashboard.ID)
	d.Set("title", dashboard.Title)
	d.Set("body_markdown", dashboard.BodyMarkDown)
	d.Set("url_path", dashboard.URLPath)

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

func resourceMackerelDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteDashboard(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel dashboard %q deleted.", d.Id())
	d.SetId("")

	return nil
}
