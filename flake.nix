{
  description = "flake speaker";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";
  inputs.nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixos-25.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs-unstable";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";

  outputs = {
    self,
    nixpkgs,
    nixpkgs-unstable,
    flake-utils,
    gomod2nix,
  }: (
    flake-utils.lib.eachDefaultSystem
    (system: let
      # Стабильный набор пакетов.
      pkgs = import nixpkgs { inherit system; };

      # Нестабильный набор пакетов, к которому применен оверлей gomod2nix.
      pkgsUnstable = import nixpkgs-unstable {
        inherit system;
        overlays = [ gomod2nix.overlays.default ];
      };

      callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;
    in {
      packages.default = callPackage ./default.nix {
        # Передаем оба набора.
        inherit pkgs pkgsUnstable;
      };
      devShells.default = callPackage ./shell.nix {
        # Передаем оба набора.
        inherit pkgs pkgsUnstable;
      };
    })
  );
}
