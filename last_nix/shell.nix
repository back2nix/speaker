# { pkgs ? import (<nixpkgs>) { } }:
{ sources ? import ./nix/sources.nix
, pkgs ? import sources.nixpkgs { }
}:

pkgs.mkShell {
  name = "go-shell";

  # nativeBuildInputs is usually what you want -- tools you need to run
  nativeBuildInputs = with pkgs; [
    nixpkgs-fmt
    rnix-lsp
    docker-client
    docker-compose
    gnumake

    # go development
    go
    go-outline
    gopls
    gopkgs
    go-tools
    delve
    translate-shell
    python310Packages.gtts
    mpg123

    libxkbcommon
    xorg.libX11.dev
    xorg.libXtst
  ];

  hardeningDisable = [ "all" ];
}
