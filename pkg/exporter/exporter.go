package exporter

import (
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	namespace = "minecraft"
)

type Collector struct {
	describes []*prometheus.Desc
	sharedmem *sharedmem.SharedMem
	Logger    *logrus.Logger
}

func New(m *sharedmem.SharedMem, l *logrus.Logger) (Collector, error) {
	describes := []*prometheus.Desc{
		prometheus.NewDesc(
			"minecraft_user_info",
			"Minecraft Login Users",
			[]string{"username"},
			nil,
		),
		prometheus.NewDesc(
			"minecraft_health_gauge",
			"Minecraft User's Health",
			[]string{"username"},
			nil,
		),
	}

	return Collector{
		describes: describes,
		sharedmem: m,
		Logger:    l,
	}, nil
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	describeUserInfo := c.describes[0]
	describeHealthGauge := c.describes[1]

	data, err := c.sharedmem.ReadSharedMem()
	if err != nil {
		c.Logger.Warn(err)
		return
	}
	for _, user := range data.LogoutUsers {
		ch <- prometheus.MustNewConstMetric(
			describeUserInfo,
			prometheus.GaugeValue,
			0,
			user.Name,
		)
	}
	for _, user := range data.LoginUsers {
		ch <- prometheus.MustNewConstMetric(
			describeUserInfo,
			prometheus.GaugeValue,
			1,
			user.Name,
		)
	}

	for _, user := range data.AllUsers {
		ch <- prometheus.MustNewConstMetric(
			describeHealthGauge,
			prometheus.GaugeValue,
			float64(user.Health),
			user.Name,
		)
	}
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, describe := range c.describes {
		ch <- describe
	}
}
