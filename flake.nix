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
    overlays.default = overlay;

    packages = forAllPkgs (pkgs: {
      default = pkgs.mumble_exporter;
    });

    devShells = forAllPkgs (pkgs: {
      default = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [ bashInteractive go ];
      };
    });

    nixosModules.default = import ./nixos-module.nix;

    checks = forAllPkgs (pkgs: {
      nixos-test = pkgs.callPackage ./nixos-test.nix {};
    });
  };
}
