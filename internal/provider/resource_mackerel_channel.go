package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelChannelCreate,
		Read:   resourceMackerelChannelRead,
		Delete: resourceMackerelChannelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"email",
					"slack",
					"webhook",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"emails": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Field name may only contain lowercase alphanumeric characters & underscores.
			"user_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"events": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"emails", "user_ids"},
			},
			"mentions": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled_graph_image": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMackerelChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input, err := buildChannelParameter(d)
	if err != nil {
		return err
	}

	channel, err := client.CreateChannel(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel channel %q created.", channel.ID)
	d.SetId(channel.ID)

	return resourceMackerelChannelRead(d, meta)
}

func resourceMackerelChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel channel: %q", d.Id())
	channels, err := client.FindChannels()
	if err != nil {
		return err
	}

	for _, channel := range channels {
		if channel.ID == d.Id() {
			_ = d.Set("name", channel.Name)
			_ = d.Set("type", channel.Type)
			_ = d.Set("url", channel.URL)
			_ = d.Set("enabled_graph_image", channel.EnabledGraphImage)
			_ = d.Set("user_ids", channel.UserIDs)

			var mentions = map[string]string{}
			if ok := channel.Mentions.OK; ok != "" {
				mentions["ok"] = ok
			}
			if crit := channel.Mentions.Critical; crit != "" {
				mentions["critical"] = crit
			}
			if warn := channel.Mentions.Warning; warn != "" {
				mentions["warning"] = warn
			}
			_ = d.Set("mentions", mentions)
			_ = d.Set("events", channel.Events)
			break
		}
	}

	return nil
}

func resourceMackerelChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteChannel(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel channel %q deleted.", d.Id())

	return nil
}

func buildChannelParameter(d *schema.ResourceData) (*mackerel.Channel, error) {
	switch d.Get("type").(string) {
	case "email":
		return buildEmailParameter(d)
	case "slack":
		return buildSlackParameter(d)
	case "webhook":
		return buildWebhookParameter(d)
	default:
		return nil, fmt.Errorf("%v is not valid input for type", d.Get("type"))
	}
}

// build parameter for email
func buildEmailParameter(d *schema.ResourceData) (*mackerel.Channel, error) {
	input := &mackerel.Channel{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	if v, ok := d.GetOk("emails"); ok {
		tmp := expandStringList(v.([]interface{}))
		input.Emails = &tmp
	}

	if v, ok := d.GetOk("user_ids"); ok {
		tmp := expandStringList(v.([]interface{}))
		input.UserIDs = &tmp
	}

	if input.Emails == nil && input.UserIDs == nil {
		return nil, fmt.Errorf("emails or user_ids is required")
	}

	if v, ok := d.GetOk("events"); ok {
		tmp := expandStringList(v.([]interface{}))
		err := validateChannelEvent(tmp, []string{"alert", "alertGroup"})
		if err != nil {
			return nil, err
		}
		input.Events = &tmp
	}

	return input, nil
}

// build parameter for slack
func buildSlackParameter(d *schema.ResourceData) (*mackerel.Channel, error) {
	input := &mackerel.Channel{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
		URL:  d.Get("url").(string),
	}

	if v, ok := d.Get("enabled_graph_image").(bool); ok {
		input.EnabledGraphImage = &v
	}

	if v, ok := d.GetOk("mentions"); ok {
		// Convert from schema.TypeMap to mackerel.Mentions
		mentionJSON, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var mentions mackerel.Mentions
		err = json.Unmarshal(mentionJSON, &mentions)
		if err != nil {
			return nil, err
		}
		input.Mentions = mentions
	}

	if v, ok := d.GetOk("events"); ok {
		tmp := expandStringList(v.([]interface{}))
		err := validateChannelEvent(tmp, []string{"alert", "alertGroup", "hostStatus", "hostRegister", "hostRetire", "monitor"})
		if err != nil {
			return nil, err
		}
		input.Events = &tmp
	}

	return input, nil
}

// build parameter for webhook
func buildWebhookParameter(d *schema.ResourceData) (*mackerel.Channel, error) {
	input := &mackerel.Channel{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
		URL:  d.Get("url").(string),
	}

	if v, ok := d.Get("enabled_graph_image").(bool); ok {
		input.EnabledGraphImage = &v
	}

	if v, ok := d.GetOk("events"); ok {
		tmp := expandStringList(v.([]interface{}))
		err := validateChannelEvent(tmp, []string{"alert", "alertGroup", "hostStatus", "hostRegister", "hostRetire", "monitor"})
		if err != nil {
			return nil, err
		}
		input.Events = &tmp
	}

	return input, nil
}
