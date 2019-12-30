# TRACKARR

[![made-with-go](https://img.shields.io/badge/Made%20with-Go-blue.svg?style=flat-square)](https://golang.org/)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%203-blue.svg?style=flat-square)](https://github.com/l3uddz/plex_autoscan/blob/master/LICENSE.md)
[![Discord](https://img.shields.io/discord/381077432285003776.svg?colorB=177DC1&label=Discord&style=flat-square)](https://discord.io/cloudbox)
[![Contributing](https://img.shields.io/badge/Contributing-gray.svg?style=flat-square)](CONTRIBUTING.md)
[![Credits](https://img.shields.io/badge/Credits-gray.svg?style=flat-square)](CREDITS.md)
[![Donate](https://img.shields.io/badge/Donate-gray.svg?style=flat-square)](#donate)


## Basics

Trackarr monitors tracker announce IRC channel, parses the announcements, and forwards those announcements to ARR PVRs (eg Sonarr).


## Why use this over RSS feeds?

ARR PVRs by default monitor RSS feeds from torrent trackers that they will check on a configured interval. When a release is found, the PVR will pick it up for download. However, in this method, you are always limited to the interval set.

With Trackarr, the release is sent to the PVR as soon as it is released (i.e. when the announcement is made on the IRC channel).

Users who 'race' will find this useful to get the latest releases downloaded quickly, while maintaining a respectable seeding ratio, due to it being in the initial swarm.

Trackarr also allows you to filter out unwanted releases using 'expressions'.

## Why use this over AutoDL plugin for ruTorrent?

- Works independent of any torrent client.

- Robust filtering.

- Tight integration with Sonarr and Radarr.


## Supported Trackers

Supports these trackers: https://github.com/autodl-community/autodl-trackers/tree/master/trackers

# Installation

See [Wiki](https://gitlab.com/cloudb0x/trackarr/-/wikis/home).

# Donate

If you find this project helpful, feel free to make a small donation to the developer:

  - [GitHub Sponsor](https://github.com/sponsors/l3uddz) - Credit Cards, PayPal, etc
  - [Monzo](https://monzo.me/today): Credit Cards, Apple Pay, Google Pay

  - [Paypal: l3uddz@gmail.com](https://www.paypal.me/l3uddz)

  - BTC: 3CiHME1HZQsNNcDL6BArG7PbZLa8zUUgjL
