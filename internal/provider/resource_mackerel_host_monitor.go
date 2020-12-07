package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelHostMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelHostMonitorCreate,
		Read:   resourceMackerelHostMonitorRead,
		Update: resourceMackerelHostMonitorUpdate,
		Delete: resourceMackerelHostMonitorDelete,
		Exists: resourceMackerelHostMonitorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"warning": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"critical": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_check_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Default:      1,
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"exclude_scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_mute": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceMackerelHostMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorHostMetric{
		Type:                 "host",
		Name:                 d.Get("name").(string),
		Duration:             uint64(d.Get("duration").(int)),
		Metric:               d.Get("metric").(string),
		Operator:             d.Get("operator").(string),
		Warning:              pfloat64(d.Get("warning").(float64)),
		Critical:             pfloat64(d.Get("critical").(float64)),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		MaxCheckAttempts:     uint64(d.Get("max_check_attempts").(int)),
		IsMute:               d.Get("is_mute").(bool),
	}

	if v, ok := d.GetOk("scopes"); ok {
		input.Scopes = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("exclude_scopes"); ok {
		input.ExcludeScopes = expandStringList(v.([]interface{}))
	}

	monitor, err := client.CreateMonitor(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q created.", monitor.MonitorID())
	d.SetId(monitor.MonitorID())

	return resourceMackerelHostMonitorRead(d, meta)
}

func resourceMackerelHostMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel monitor: %q", d.Id())
	monitors, err := client.FindMonitors()
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		if monitor.MonitorType() == "host" && monitor.MonitorID() == d.Id() {
			mon := monitor.(*mackerel.MonitorHostMetric)
			_ = d.Set("id", mon.ID)
			_ = d.Set("name", mon.Name)
			_ = d.Set("duration", mon.Duration)
			_ = d.Set("metric", mon.Metric)
			_ = d.Set("operator", mon.Operator)
			_ = d.Set("warning", mon.Warning)
			_ = d.Set("critical", mon.Critical)
			_ = d.Set("notification_interval", mon.NotificationInterval)
			_ = d.Set("max_check_attempts", mon.MaxCheckAttempts)
			_ = d.Set("scopes", flattenStringList(mon.Scopes))
			_ = d.Set("exclude_scopes", flattenStringList(mon.ExcludeScopes))
			_ = d.Set("is_mute", mon.IsMute)
			break
		}
	}

	return nil
}

func resourceMackerelHostMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorHostMetric{
		Type:                 "host",
		Name:                 d.Get("name").(string),
		Duration:             uint64(d.Get("duration").(int)),
		Metric:               d.Get("metric").(string),
		Operator:             d.Get("operator").(string),
		Warning:              pfloat64(d.Get("warning").(float64)),
		Critical:             pfloat64(d.Get("critical").(float64)),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		MaxCheckAttempts:     uint64(d.Get("max_check_attempts").(int)),
		IsMute:               d.Get("is_mute").(bool),
	}

	if v, ok := d.GetOk("scopes"); ok {
		input.Scopes = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("exclude_scopes"); ok {
		input.ExcludeScopes = expandStringList(v.([]interface{}))
	}

	_, err := client.UpdateMonitor(d.Id(), input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q updated.", d.Id())
	return resourceMackerelHostMonitorRead(d, meta)
}

func resourceMackerelHostMonitorExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*mackerel.Client)
	monitors, err := client.FindMonitors()
	if err != nil {
		return false, err
	}

	for _, m := range monitors {
		if m.MonitorType() == "host" && m.MonitorID() == d.Id() {
			return true, nil
		}
	}

	return false, nil
}

func resourceMackerelHostMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteMonitor(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q deleted.", d.Id())
	d.SetId("")

	return nil
}
