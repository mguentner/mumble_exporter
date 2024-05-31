{ lib, buildGoModule }:
buildGoModule rec {
  pname = "mumble_exporter";
  version = "1.0.0";

  src = lib.cleanSource ./.;

  vendorHash = "sha256-VQ4XCsS69psXxMp1iabcqcl0oK6kUb079IvZXU8kcBQ";
  proxyVendor = true;

  meta = with lib; {
    description = "a prometheus exporter for mumble / murmur";
    homepage = "https://github.com/mguentner/mumble_exporter";
    license = licenses.mit;
    maintainers = with maintainers; [ mguentner ];
  };
}
