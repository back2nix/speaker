{ 
pkgs ? (
    let
      sources = import ./nix/sources.nix;
    in
      import sources.nixpkgs {
        overlays = [
          (import "${sources.gomod2nix}/overlay.nix")
        ];
      }
  ),
}: 
let
  goEnv = pkgs.mkGoEnv {pwd = ./.;};
in
  pkgs.mkShell {
    name = "speaker-shell";
    packages = with pkgs; [
      # pkgs.gomod2nix
      goEnv
      go-tools
      # pkgs.niv
      translate-shell
      python310Packages.gtts
      mpg123
      libxkbcommon
      xorg.libX11.dev
      xorg.libXtst
    ];

    postShellHook = ''
      '';
  }
