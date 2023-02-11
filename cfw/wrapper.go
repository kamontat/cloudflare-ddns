package cfw

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/logger"
)

type Wrapper struct {
	Config *models.Config
	ZoneId string

	context  context.Context
	cancelFn context.CancelFunc
	api      *cloudflare.API
	logger   *logger.Logger
}

// // t can be either A or AAAA
// func (w *Wrapper) ListRecords(t string) ([]cloudflare.DNSRecord, error) {
// 	records, _, err := w.api.ListDNSRecords(
// 		w.context,
// 		cloudflare.ZoneIdentifier(w.ZoneId),
// 		cloudflare.ListDNSRecordsParams{
// 			Type: t,
// 			ResultInfo: cloudflare.ResultInfo{
// 				Page:    1,
// 				PerPage: 100,
// 			},
// 		},
// 	)

// 	if err != nil {
// 		w.cancelFn()
// 		return records, err
// 	}

// 	return records, err
// }

func (w *Wrapper) UpdateRecords() {
	ipv4, err1 := GetPublicIPV4(w.Config.Settings.Ipv4)
	if err1 != nil {
		// Warning cannot get public ip
		// and skip it
		w.logger.Warn(err1)
	}

	ipv6, err1 := GetPublicIPV6(w.Config.Settings.Ipv6)
	if err1 != nil {
		// Warning cannot get public ip
		// and skip it
		w.logger.Warn(err1)
	}

	w.logger.Debug("ipv4: %s, ipv6: %s", ipv4, ipv6)
	err := w.updateRecordsA(ipv4)
	if err != nil {
		// Warning cannot get public ip
		// and skip it
		w.logger.Warn(err)
	}
	err = w.updateRecordsAAAA(ipv6)
	if err != nil {
		// Warning cannot get public ip
		// and skip it
		w.logger.Warn(err)
	}
}

func (w *Wrapper) updateRecordsA(ip string) (err error) {
	return
}

func (w *Wrapper) updateRecordsAAAA(ip string) (err error) {
	return
}
