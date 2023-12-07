# `secret-wallet`

`experimental`

A Dero Wallet with mobile UI.
Cross-platform: Linux, Windows, MacOS, Android, IOS...

## Why

The goal is to bring more users into the Dero Ecosystem with a wallet that is intuitive, easy to use and don't compromise privacy.

One notable technological difference is that this wallet uses Gio UI instead of Fyne for rendering components.

## Features

- Connect to a remote node for quick utilization.
- Multiple wallet management.
- Manage Dero and DERO tokens
- Easily send payments and receive with QR code.
- Address book to manage contacts.
- Support as much language as possible (11 and counting...).
- App color schemes (dark, light, blue).
- Draggable list items (for ordering items).

## TODO

- App release (iOS app store, Google play store & F-Droid).
- Expose daemon/wallet API for integrated node.


## Releases

You can build your own version by following build steps bellow or use available prebuilds here.

## How to build

### Setup

If you have Git installed, clone the repository.

```bash
git clone https://github.com/secretnamebasis/secret-wallet.git
```

Download and install GO.
<https://go.dev/doc/install>

Install latest version of GioUI build tool.

```bash
go install gioui.org/cmd/gogio@latest
```


### Android

Install Android SDK with NDK bundle!

Build the app.

```bash
./build_gio.sh android arm64
```

If your phone is connected and developer mode is on, install the APK directly with this command.

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

`/build/secret_wallet_windows_amd64.exe`
`/build/secret_wallet_linux_amd64`
`/build/secret_wallet_macos_amd64.app`
`/build/secret_wallet_ios_arm64.app`
`/build/secret_wallet_android_arm64.apk`
`/build/secret_wallet_android_arm.apk`

## Contributors

List of contributors. Thank you all!
If your alias is not here pls let me know.

### Translation

Most translations were done with ChatGPT and are far from being accurate.
If you notice any errors, pls update the values in `/assets/lang` folder and create a pull request.



### Testing

`secretnamebasis`, `chakipu`, `TheObjectiveAlpha`, `Ulmo`, `derionner`

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).

The MIT License (MIT)
Copyright (c) 2023 g45t345rt
