# Mumble Prometheus Exporter

exports mumble / murmur server stats by using the mumble ping protocol

## Prometheus Configuration

The exporter expects a single HTTP parameter named `host` that must include a port, i.e. `mumble.example.com:1337`.

You can test if your the exporter works correctly with
```
curl localhost:8778?host=mumble.example.com:63748
```

Prometheus job example
```
scrape_config
 - job_name: 'mumble'
    scrape_interval: 120s
    scheme: http
    params:
      host: ['mumble.example.com:64738']
    static_configs:
      - targets:
        - 'localhost:8778'
```

## Credits

based on https://github.com/mumble-voip/mumble-scripts/blob/master/Non-RPC/mumble-ping.py

## Copyright & License

2022 Maximilian Güntner <code@mguentner.de>

MIT
