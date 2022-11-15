{ lib, buildGoModule }:
buildGoModule rec {
  pname = "mumble_exporter";
  version = "1.0.0";

  src = lib.cleanSource ./.;

  vendorSha256 = "sha256-pMVmGuDJT9A+EDJsoyyIcbnzOGD2Qj5nfr7e/dmhM9w=";
  proxyVendor = true;

  meta = with lib; {
    description = "a prometheus exporter for mumble / murmur";
    homepage = "https://github.com/mguentner/mumble_exporter";
    license = licenses.mit;
    maintainers = with maintainers; [ mguentner ];
  };
}
