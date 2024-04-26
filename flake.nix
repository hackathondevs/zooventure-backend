{
  description = "Avanzu Backend dev environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          air
          delve
          go-migrate
          (go_1_22.overrideAttrs (prev: rec {
            version = "1.22.0";
            src = fetchurl {
              url = "https://go.dev/dl/go${version}.src.tar.gz";
              hash = "sha256-TRlsPUGg1sHfxk0E48wfYIsMQ2vYe3Bgzj4jI04fTVw=";
            };
          }))
          golangci-lint
          gopls
        ];
      };
    });
}
