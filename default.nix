{
  # Чистая функция. Все аргументы передаются из flake.nix.
  pkgs,
  pkgsUnstable,
  lib ? pkgs.lib,
}:
pkgsUnstable.buildGoApplication {
  pname = "speaker";
  version = "1.0.0";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;

  buildInputs =
    (with pkgs; [
      libxkbcommon
      xorg.libX11
      xorg.libXtst
    ])
    ++ (with pkgsUnstable; [
      go
      delve
    ]);

  nativeBuildInputs = with pkgs; [ makeWrapper ];

  postInstall = ''
    cp -r sound $out/;

    wrapProgram "$out/bin/speaker" \
    --prefix PATH : ${lib.makeBinPath (
      [ pkgsUnstable.mpg123 ] ++
      (with pkgs; [
        xorg.libXtst
        libxkbcommon
        xorg.libX11
        translate-shell
        python312Packages.gtts
      ])
    )}:$out
  '';
}
