package mackerel

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelServiceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelServiceMonitorCreate,
		Read:   resourceMackerelServiceMonitorRead,
		Update: resourceMackerelServiceMonitorUpdate,
		Delete: resourceMackerelServiceMonitorDelete,
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
			"service": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"duration": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric": &schema.Schema{
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

func resourceMackerelServiceMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorServiceMetric{
		Type:                 "service",
		Name:                 d.Get("name").(string),
		Service:              d.Get("service").(string),
		Duration:             uint64(d.Get("duration").(int)),
		Metric:               d.Get("metric").(string),
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

	return resourceMackerelServiceMonitorRead(d, meta)
}

func resourceMackerelServiceMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel monitor: %q", d.Id())
	monitors, err := client.FindMonitors()
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		if monitor.MonitorType() == "service" && monitor.MonitorID() == d.Id() {
			mon := monitor.(*mackerel.MonitorServiceMetric)
			d.Set("id", mon.ID)
			d.Set("name", mon.Name)
			d.Set("service", mon.Service)
			d.Set("duration", mon.Duration)
			d.Set("metric", mon.Metric)
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

func resourceMackerelServiceMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorServiceMetric{
		Type:                 "service",
		Name:                 d.Get("name").(string),
		Service:              d.Get("service").(string),
		Duration:             uint64(d.Get("duration").(int)),
		Metric:               d.Get("metric").(string),
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
	return resourceMackerelServiceMonitorRead(d, meta)
}

func resourceMackerelServiceMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteMonitor(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q deleted.", d.Id())
	d.SetId("")

	return nil
}
