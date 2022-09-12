{ config, lib, pkgs, ... }:

let
  cfg = config.services.prometheus-mumble-exporter;
in
{
  options.services.prometheus-mumble-exporter = {
    enable = lib.mkEnableOption "mumble exporter";
    listenAddress = lib.mkOption {
      type = lib.types.str;
      default = ":8778";
      description = "Address and port to expose metrics";
    };
    metricsPath = lib.mkOption {
      type = lib.types.str;
      default = "/metrics";
      description = "Path under which to expose metrics";
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services."prometheus-mumble-exporter" = {
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      serviceConfig = {
        Restart = "always";
        DynamicUser = true;
        LockPersonality = true;
        MemoryDenyWriteExecute = true;
        PrivateDevices = true;
        ProtectSystem = "strict";
        RemoveIPC = true;
        RestrictAddressFamilies = [ "AF_INET" "AF_INET6" ];
        RestrictNamespaces = true;
        RestrictSUIDSGID = true;
        SystemCallArchitectures = "native";
        ExecStart = ''
          ${pkgs.callPackage ./. {}}/bin/mumble_exporter \
            --listenAddress=${cfg.listenAddress} \
            --metricsPath=${cfg.metricsPath}
        '';
      };
    };
  };
}
