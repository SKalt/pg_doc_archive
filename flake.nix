{
  description = "A very basic flake";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, flake-utils, nixpkgs }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = (import nixpkgs) {
          inherit system;
        };
      in
      rec {
        devShell = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go # 1.19.x
          ];
          buildInputs = with pkgs; [
            nixpkgs-fmt
            nil
            gopls
            gotools
            goreleaser
            helix
            bashInteractive
          ];
        };
      }
    );
}
