{
  description = "Tools to scrape and archive postgres's HTML documentation.";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, flake-utils, nixpkgs }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = (import nixpkgs) { inherit system; };
      in
      rec {
        devShell = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go # 1.22.x
          ];
          buildInputs = with pkgs; [
            nixpkgs-fmt
            nil
            gopls
            gotools
            golangci-lint
            goreleaser
            gh
            helix
            bashInteractive
          ];
        };
      }
    );
}
