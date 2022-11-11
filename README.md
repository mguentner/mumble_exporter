# Mumble Prometheus Exporter

exports mumble / murmur server stats by using the mumble ping protocol

## Prometheus Configuration

The exporter expects a single HTTP parameter named `host` that must include a port, i.e. `mumble.example.com:1337`.

You can test if your the exporter works correctly with
```bash
curl "localhost:8778/metrics?host=mumble.example.com:63748"
```

Prometheus job example
```yml
scrape_config

 - job_name: 'mumble'
    scrape_interval: 120s
    scheme: http
    static_configs:
      - targets: # list of monitored mumble servers
        - 'mumble.example.com:64738' 
        - 'bumble.example.com:64738'
    relabel_configs:
    - source_labels: [__address__]
      target_label: __param_host
    - source_labels: [__param_host]
      target_label: instance
    - target_label: __address__
      replacement: "127.0.0.1:8778" # address of host running the exporter
```

## Credits

based on https://github.com/mumble-voip/mumble-scripts/blob/master/Non-RPC/mumble-ping.py

## Copyright & License

2022 Maximilian Güntner <code@mguentner.de>

MIT
