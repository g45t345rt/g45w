# G45W

`IN DEVELOPMENT`

A Dero Universal Wallet with mobile UI.
Cross-platform: Linux, Windows, MacOS, Android, IOS...

## Why

Utimately, my goal is to bring more users into the Dero Ecosystem with a wallet that is intuitive, easy to use
and don't compromise privacy.

Although there is already a fantastic wallet created by the DERO Foundation members <https://github.com/DEROFDN/Engram>,
I want to provide an alternative with a different user interface and additional features, such as smart contract interaction.
One notable technological difference is that this wallet uses Gio UI instead of Fyne for rendering components.

Additionally, I intend to discontinue the `Dero RPC Bridge` browser extension in favor of this wallet.
This change should remove a painful/unsafe step of connecting your wallet directly to
the platform. In the future, platforms should provide QR code for smart contract calls, allowing users to send
transactions directly from this wallet.

## Features

- Integrated fast-sync Dero Node for full privacy.
- Connect to a remote node for quick utilization.
- Multiple wallet management.
- Manage Dero tokens and NFTs.
- Easily send payments and receive with QR code.
- Fast registration implemented (by pieswap).
- Address book to manage contacts.
- Support as much language as possible.

## Demo

![Wallet app demo](https://github.com/g45t345rt/g45w/blob/master/g45w_demo.gif)

## Releases

You can build your own version by following build steps bellow or use available prebuilds here <https://github.com/g45t345rt/g45w/releases>.

## How to build

If you have Git installed, clone the repository.

```bash
git clone https://github.com/g45t345rt/g45w.git
```

Download and install GO.  
<https://go.dev/doc/install>

Install latest version of GioUI build tool.

```bash
go install gioui.org/cmd/gogio@latest
```

Check Github workflows for more build information.  
<https://github.com/g45t345rt/g45w/tree/master/.github/workflows>

### Android

Install Android SDK with NDK bundle!

Build the app.

```bash
./build_gio.sh android arm64
```

If your phone is connected and developer monde is on, install the APK directly with this command.

```bash
./adb_install.sh
```

### IOS

Xcode is required!  

Build the app.

```bash
./build_gio.sh ios arm64
```

### Linux

Install dependencies

Fedora 35+

```bash
dnf install gcc pkg-config wayland-devel libX11-devel libxkbcommon-x11-devel mesa-libGLES-devel mesa-libEGL-devel libXcursor-devel vulkan-headers
```

Ubuntu 18.04+

```bash
apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev
```

Build the app.

```bash
./build_go.sh linux amd64
```

### Windows

Build the app.

```bash
./build_gio.sh windows amd64
```

### MacOS

Xcode is required!  
Build the app.

``` bash
./build_gio.sh macos amd64
```

### Outputs

`/build/g45w_windows_amd64.exe`  
`/build/g45w_linux_amd64`  
`/build/g45w_macos_amd64.app`  
`/build/g45w_ios_arm64.app`  
`/build/g45w_android_arm64.apk`  
`/build/g45w_android_arm.apk`  

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).

The MIT License (MIT)  
Copyright (c) 2023 g45t345rt
