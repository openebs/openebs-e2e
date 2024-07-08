{ }:
let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs {
    overlays = [
      (_: _: { inherit sources; })
    ];
  };
in
with pkgs;
mkShell {

  buildInputs = [
    docker-compose
    kubectl
    kind
    docker
    cowsay
    e2fsprogs
    envsubst # for e2e tests
    gdb
    go
    golangci-lint
    google-cloud-sdk
    git
    kubernetes-helm
    nodejs-slim
    numactl
    meson
    ninja
    openssl
    pkg-config
    pre-commit
    procps
    python3
    utillinux
    xfsprogs
    hcloud
    yq
    protoc-gen-go-grpc
    protobuf
    protoc-gen-go
  ]
  ;

  shellHook = ''
    if [ -z "$CI" ]; then
      pre-commit install
      pre-commit install --hook commit-msg
    fi
  '';
}
