#!/usr/bin/env bash
set -x

REPO=$HOME/.local/share/flatpak/dev-repo

flatpak-builder --repo=$REPO --force-clean flatpak-build io.github.mjdiliscia.lemmeread.yml
flatpak build-update-repo $REPO
