## Overview

These are a collection of packages we need, or packages where we
want to control the exact version(s) of.

The packages are imported through the `nix-shell` automatically.
These nix files originally from Mayastor and heavily reduced
for use in openebs-e2e.

To update the nix packages
nix-env -iA nixpkgs.niv
niv update nixpkgs

Then raise a PR
