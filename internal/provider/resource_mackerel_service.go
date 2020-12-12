package provider

import (
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelService() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelServiceCreate,
		Read:   resourceMackerelServiceRead,
		Delete: resourceMackerelServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9\-\_]+$`), "must contain only alphanumeric characters, dashes, and underscores"),
			},
			"memo": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
		},
	}
}

func resourceMackerelServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.CreateServiceParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	}

	svc, err := client.CreateService(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel service %q created.", svc.Name)
	d.SetId(svc.Name)

	return resourceMackerelServiceRead(d, meta)
}

func resourceMackerelServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel service: %q", d.Id())
	services, err := client.FindServices()
	if err != nil {
		return err
	}
	for _, svc := range services {
		if svc.Name == d.Id() {
			d.SetId(svc.Name)
			_ = d.Set("name", svc.Name)
			_ = d.Set("memo", svc.Memo)
		}
	}

	return nil
}

func resourceMackerelServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteService(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel service %q deleted.", d.Id())

	return nil
}
