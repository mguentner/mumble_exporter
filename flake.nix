{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs, ... }: let
    systems = [ "x86_64-linux" "aarch64-linux" ];
    overlay = final: prev: {
      mumble_exporter = final.callPackage ./default.nix {};
    };
    forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);
    forAllPkgs = f:
      forAllSystems (system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ overlay ];
        };
      in
        f pkgs);
  in {
    inherit overlay;

    packages = forAllPkgs (pkgs: {
      inherit (pkgs) mumble_exporter;
    });

    defaultPackage = forAllSystems (system: self.packages.${system}.mumble_exporter);
    devShell = forAllPkgs (pkgs:
      pkgs.mkShell {
        nativeBuildInputs = with pkgs; [ bashInteractive go ];
      });
  };
}
