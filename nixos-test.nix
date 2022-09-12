{ nixosTest }:

nixosTest {
  name = "mumble-exporter";

  nodes.server = { ... }: {
    imports = [ ./nixos-module.nix ];
    services.prometheus-mumble-exporter.enable = true;
    services.murmur.enable = true;
  };

  testScript = ''
    start_all()
    server.wait_for_unit("prometheus-mumble-exporter.service")
    server.succeed("curl http://localhost:8778/health")
    server.wait_for_open_port(64738)
    server.succeed("curl http://localhost:8778/metrics?host=localhost:64738 | grep -E '^mumble_current_users[^ ]* 0$'")
  '';
}
