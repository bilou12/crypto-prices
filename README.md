# crypto-prices

* Build the project
```bash
go build
```

* Define the service file
```bash
[Unit]
Description=binance rt 1h

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/mnt/data/crypto-prices/crypto-prices

[Install]
WantedBy=multi-user.target
```

* Interact with the service
```bash
vi /lib/systemd/system/binance_rt_1h.service
service binance_rt_1h status
service binance_rt_1h start
service binance_rt_1h stop
systemctl daemon-reload
journalctl -u binance_rt_1h.service
```