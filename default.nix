{ lib, buildGoModule }:
buildGoModule rec {
  pname = "mumble_exporter";
  version = "1.0.0";

  src = ./.;

  vendorSha256 = "sha256-BlVNkGrDCR/nMS9uCrVfjlhtK1yApfYvBF74q/NdtXo";
  proxyVendor = true;

  meta = with lib; {
    description = "a prometheus exporter for mumble / murmur";
    homepage = "https://github.com/mguentner/mumble_exporter";
    license = licenses.mit;
    maintainers = with maintainers; [ mguentner ];
  };
}
