{
  description = "Knockdoor";

  nixConfig.bash-prompt = "\[develop\]$ ";

  inputs.nixpkgs.url = "nixpkgs/nixos-23.05";

  outputs = { self, nixpkgs }:
    let
      pname = "knockdoor";

      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {

      # Provide some binary packages for selected system types.
      packages = forAllSystems (system: rec {
        pkgs = nixpkgsFor.${system};
        knockdoor = pkgs.buildGoModule {
          inherit pname version;
          src = ./.;

          # remeber to bump this hash when your dependencies change.
          #vendorSha256 = pkgs.lib.fakeSha256;
          vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
        };
      });

      devShells = forAllSystems (system: rec {
        pkgs = nixpkgsFor.${system};
        default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gopls gotools go-tools ];

          shellHook = ''
            go mod tidy
            echo "next you can run 'go run main.go'"
          '';
        };
      });

      defaultPackage = forAllSystems (system: self.packages.${system}.${pname});
    };
}
