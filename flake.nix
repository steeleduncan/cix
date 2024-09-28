{
  description = "Cix: minimal nix ci";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
  
  outputs = { self, nixpkgs }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      
    in rec {
      packages = {
        x86_64-linux.default =
          pkgs.stdenv.mkDerivation {
            name = "cix";
            buildInputs = [
              pkgs.go
            ];
            src = ./.;
            buildPhase = ''
              mkdir -p junk
              export HOME=$(pwd)/junk
              go build -o cix github.com/steeleduncan/cix
            '';

            installPhase = ''
              mkdir -p $out/bin
              mv cix $out/bin/
            '';
          };
      };
      checks = packages;
    };
}

