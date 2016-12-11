package mackerel

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelExpressionMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelExpressionMonitorCreate,
		Read:   resourceMackerelExpressionMonitorRead,
		Update: resourceMackerelExpressionMonitorUpdate,
		Delete: resourceMackerelExpressionMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"warning": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
			"critical": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
			"notification_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"is_mute": &schema.Schema{
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
		Warning:              d.Get("warning").(float64),
		Critical:             d.Get("critical").(float64),
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
			d.Set("id", mon.MonitorID())
			d.Set("name", mon.MonitorName())
			d.Set("expression", mon.Expression)
			d.Set("operator", mon.Operator)
			d.Set("warning", mon.Warning)
			d.Set("critical", mon.Critical)
			d.Set("notification_interval", mon.NotificationInterval)
			d.Set("is_mute", mon.IsMute)
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
		Warning:              d.Get("warning").(float64),
		Critical:             d.Get("critical").(float64),
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
