/**
 * @Time: 2023/1/5 22:24
 * @Author: solacowa@gmail.com
 * @File: service_wifi
 * @Software: GoLand
 */

package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/mdlayher/wifi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	wifiCmd = &cobra.Command{
		Use:               "wifi command <args> [flags]",
		Short:             "WIFI操作命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `## WIFI操作命令
可用的配置类型：
[test]

pi-kit wifi -h
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return prepare(ctx)
		},
	}

	wifiTestCmd = &cobra.Command{
		Use:               "test <args> [flags]",
		Short:             "Test wifi",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `## Test wifi
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			wifiCli, err := wifi.New()
			if err != nil {
				_ = level.Error(logger).Log("wifi", "New", "err", err.Error())
				return errors.Wrap(err, "wifi.New")
			}
			interfaces, err := wifiCli.Interfaces()
			if err != nil {
				_ = level.Error(logger).Log("wifiCli", "Interfaces", "err", err.Error())
				return errors.Wrap(err, "wifiCli.Interfaces")
			}
			for _, v := range interfaces {
				fmt.Println(v.Name, v.PHY, v.Frequency, v.Device, v.Index, v.Type.String(), v.HardwareAddr.String())
				bss, err := wifiCli.BSS(v)
				if err != nil {
					_ = level.Warn(logger).Log("wifiCli", "BSS", "name", v.Name, "err", err.Error())
					continue
				}
				fmt.Println("bss:", bss.BSSID.String(), bss.Frequency, bss.Status, bss.SSID, bss.BeaconInterval, bss.LastSeen)
				stationInfo, err := wifiCli.StationInfo(v)
				if err != nil {
					_ = level.Warn(logger).Log("wifiCli", "StationInfo", "name", v.Name, "err", err.Error())
					continue
				}
				for _, info := range stationInfo {
					b, _ := json.Marshal(info)
					fmt.Println("stationInfo:", string(b))
				}
			}
			return nil
		},
	}
)
