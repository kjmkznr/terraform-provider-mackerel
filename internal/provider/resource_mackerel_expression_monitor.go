package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelExpressionMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelExpressionMonitorCreate,
		Read:   resourceMackerelExpressionMonitorRead,
		Update: resourceMackerelExpressionMonitorUpdate,
		Delete: resourceMackerelExpressionMonitorDelete,
		Exists: resourceMackerelExpressionMonitorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": {
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
			"is_mute": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceMackerelExpressionMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorExpression{
		Type:                 "expression",
		Name:                 d.Get("name").(string),
		Expression:           d.Get("expression").(string),
		Operator:             d.Get("operator").(string),
		Warning:              pfloat64(d.Get("warning").(float64)),
		Critical:             pfloat64(d.Get("critical").(float64)),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		IsMute:               d.Get("is_mute").(bool),
	}

	monitor, err := client.CreateMonitor(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q created.", monitor.MonitorID())
	d.SetId(monitor.MonitorID())

	return resourceMackerelExpressionMonitorRead(d, meta)
}

func resourceMackerelExpressionMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel monitor: %q", d.Id())
	monitors, err := client.FindMonitors()
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		if monitor.MonitorType() == "expression" && monitor.MonitorID() == d.Id() {
			mon := monitor.(*mackerel.MonitorExpression)
			_ = d.Set("id", mon.MonitorID())
			_ = d.Set("name", mon.MonitorName())
			_ = d.Set("expression", mon.Expression)
			_ = d.Set("operator", mon.Operator)
			_ = d.Set("warning", mon.Warning)
			_ = d.Set("critical", mon.Critical)
			_ = d.Set("notification_interval", mon.NotificationInterval)
			_ = d.Set("is_mute", mon.IsMute)
			break
		}
	}

	return nil
}

func resourceMackerelExpressionMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorExpression{
		Type:                 "expression",
		Name:                 d.Get("name").(string),
		Expression:           d.Get("expression").(string),
		Operator:             d.Get("operator").(string),
		Warning:              pfloat64(d.Get("warning").(float64)),
		Critical:             pfloat64(d.Get("critical").(float64)),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		IsMute:               d.Get("is_mute").(bool),
	}

	_, err := client.UpdateMonitor(d.Id(), input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q updated.", d.Id())
	return resourceMackerelExpressionMonitorRead(d, meta)
}

func resourceMackerelExpressionMonitorExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*mackerel.Client)
	monitors, err := client.FindMonitors()
	if err != nil {
		return false, err
	}

	for _, m := range monitors {
		if m.MonitorType() == "expression" && m.MonitorID() == d.Id() {
			return true, nil
		}
	}

	return false, nil
}

func resourceMackerelExpressionMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteMonitor(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q deleted.", d.Id())
	d.SetId("")

	return nil
}
