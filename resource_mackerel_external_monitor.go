package mackerel

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelExternalMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelExternalMonitorCreate,
		Read:   resourceMackerelExternalMonitorRead,
		Update: resourceMackerelExternalMonitorUpdate,
		Delete: resourceMackerelExternalMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"response_time_warning": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"response_time_critical": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"response_time_duration": {
				Type:     schema.TypeFloat,
				Optional: true,
				//ValidateFunc: validateDurationTime 1 minute - 10 minute
			},
			"contains_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_check_attempts": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  1,
				//ValidateFunc: validateDurationTime 1 - 10
			},
			"certification_expiration_warning": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"certification_expiration_critical": {
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

func resourceMackerelExternalMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := getMackerelExternalMonitorInput(d)
	monitor, err := client.CreateMonitor(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q created.", monitor.MonitorID())
	d.SetId(monitor.MonitorID())

	return resourceMackerelExternalMonitorRead(d, meta)
}

func resourceMackerelExternalMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel monitor: %q", d.Id())
	monitors, err := client.FindMonitors()
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		if monitor.MonitorType() == "external" && monitor.MonitorID() == d.Id() {
			mon := monitor.(*mackerel.MonitorExternalHTTP)
			_ = d.Set("id", mon.MonitorID())
			_ = d.Set("name", mon.MonitorName())
			_ = d.Set("url", mon.URL)
			_ = d.Set("service", mon.Service)
			_ = d.Set("notification_interval", mon.NotificationInterval)
			_ = d.Set("response_time_duration", mon.ResponseTimeDuration)
			_ = d.Set("response_time_warning", mon.ResponseTimeWarning)
			_ = d.Set("response_time_critical", mon.ResponseTimeCritical)
			_ = d.Set("contains_string", mon.ContainsString)
			_ = d.Set("max_check_attempts", mon.MaxCheckAttempts)
			_ = d.Set("certification_expiration_warning", mon.CertificationExpirationWarning)
			_ = d.Set("certification_expiration_critical", mon.CertificationExpirationCritical)
			_ = d.Set("is_mute", mon.IsMute)
			break
		}
	}

	return nil
}

func resourceMackerelExternalMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := getMackerelExternalMonitorInput(d)
	_, err := client.UpdateMonitor(d.Id(), input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q updated.", d.Id())
	return resourceMackerelExternalMonitorRead(d, meta)
}

func resourceMackerelExternalMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteMonitor(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q deleted.", d.Id())
	d.SetId("")

	return nil
}

func getMackerelExternalMonitorInput(d *schema.ResourceData) *mackerel.MonitorExternalHTTP {
	input := &mackerel.MonitorExternalHTTP{
		Type: "external",
		Name: d.Get("name").(string),
		URL:  d.Get("url").(string),
	}

	if v, ok := d.GetOk("service"); ok {
		input.Service = v.(string)
	}
	if v, ok := d.GetOk("notification_interval"); ok {
		input.NotificationInterval = uint64(v.(int))
	}
	if v, ok := d.GetOk("response_time_duration"); ok {
		input.ResponseTimeDuration = puint64(uint64(v.(float64)))
	}
	if v, ok := d.GetOk("response_time_warning"); ok {
		input.ResponseTimeWarning = pfloat64(v.(float64))
	}
	if v, ok := d.GetOk("response_time_critical"); ok {
		input.ResponseTimeCritical = pfloat64(v.(float64))
	}
	if v, ok := d.GetOk("contains_string"); ok {
		input.ContainsString = v.(string)
	}
	if v, ok := d.GetOk("max_check_attempts"); ok {
		input.MaxCheckAttempts = uint64(v.(float64))
	}
	if v, ok := d.GetOk("certification_expiration_warning"); ok {
		input.CertificationExpirationWarning = puint64(uint64(v.(int)))
	}
	if v, ok := d.GetOk("certification_expiration_critical"); ok {
		input.CertificationExpirationCritical = puint64(uint64(v.(int)))
	}
	if v, ok := d.GetOk("is_mute"); ok {
		input.IsMute = v.(bool)
	}

	return input
}
