# G45W

`IN DEVELOPMENT`

A Dero Universal Wallet with mobile UI.
Cross-platform: Linux, Windows, MacOS, Android, IOS...

## Why

Utimately, my goal is to bring more users into the Dero Ecosystem with a wallet that is intuitive and easy to use.

The Dero Project already has a fantastic wallet made by the foundation members,
but I want to offer an alternative with a different UI and other features.
The most notable difference in technology is that this wallet uses Gio UI instead of Fyne for rendering components.

I also want to discontinue the `Dero RPC Bridge` browser extension in favor of this wallet.
This will remove a painful/unsafe step of connecting your wallet directly with
the platform. In the future platforms should provided QR code
of smart contract calls to directly send from wallet.

## Features

- Integrated fast-sync Dero Node for full privacy.
- Connect to a remote node for quick utilization.
- Multiple wallet management.
- Manage Dero tokens and NFTs.
- Easily send payments and receive with QR code.
- Fast registration implemented (by pieswap).
- Address book to manage contacts.
- Support as much language as possible.

## Screenshots

## Build

### Mobile

#### Android

Run `./build_android.sh` to compile go and create apk package in `/build/g45w_android.apk`

#### IOS

Xcode is required!

Run `./build_ios.sh` to compile go and create app in `/build/g45w_ios.app`

### Desktop

#### Linux

Install dependencies

Fedora 35+

```bash
dnf install gcc pkg-config wayland-devel libX11-devel libxkbcommon-x11-devel mesa-libGLES-devel mesa-libEGL-devel libXcursor-devel vulkan-headers
```

Ubuntu 18.04+

```bash
apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev
```

Run `./build.sh linux amd64` to compile and create file in `/build/g45w_linux_amd64`

#### Windows

Run `./build.sh windows amd64` to compile and create .exe in `/build/g45w_windows_amd64.exe`

#### MacOS

Xcode is required!

Run `./build.sh darwin amd64` to compile and create file in `/build/g45w_darwin_amd64`

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).

The MIT License (MIT)
Copyright (c) 2023 g45t345rt
