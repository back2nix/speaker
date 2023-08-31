{ pkgs ? (
    let
      sources = import ./nix/sources.nix;
    in
    import sources.nixpkgs {
      overlays = [
        (import "${sources.gomod2nix}/overlay.nix")
      ];
    }
  )
,
}:
pkgs.buildGoApplication {
  pname = "speaker";
  version = "1.0.0";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;

  buildInputs = with pkgs; [
    translate-shell
    python310Packages.gtts
    mpg123
    libxkbcommon
    xorg.libX11.dev
    xorg.libXtst
  ];
}
