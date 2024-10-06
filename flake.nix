{
  description = "Cix: minimal nix ci";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
  
  outputs = { self, nixpkgs }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      our_revision =
        if builtins.hasAttr "rev" self then
          self.rev
        else
          "uncommitted";
      
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
              go build -ldflags "-X github.com/steeleduncan/cix/version.BuildRevision=${our_revision}" -o cix github.com/steeleduncan/cix
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

