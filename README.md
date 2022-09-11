# Mumble Prometheus Exporter

exports mumble / murmur server stats by using the mumble ping protocol

## Prometheus Configuration

The exporter expects a single HTTP parameter named `host` that must include a port, i.e. `mumble.example.com:1337`.

## Credits

based on https://github.com/mumble-voip/mumble-scripts/blob/master/Non-RPC/mumble-ping.py

## Copyright & License

2022 Maximilian Güntner <code@mguentner.de>

MIT
