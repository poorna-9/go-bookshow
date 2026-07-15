package config

import (
	"github.com/razorpay/razorpay-go"
)

func NewRazorpayClient(cfg *Config) *razorpay.Client {
	return razorpay.NewClient(cfg.RazorpayKeyID, cfg.RazorpayKeySecret)
}
