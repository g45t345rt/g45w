For icons, it's common to use https://github.com/google/iconvg with Gio icon widget.

The old material design icons is available here
https://pkg.go.dev/golang.org/x/exp/shiny@v0.0.0-20230817173708-d852ddb80c63/materialdesign/icons

When using svg from google fonts, make sure `filled` is selected instead of `outlined`.
https://github.com/google/material-design-icons/

An updated iconvg list available here by gennesseaux
https://github.com/gio-eui/md3-icons

The best platform to get svg icons with MIT license
https://www.svgrepo.com/

Use `go run ./svg_to_iconvg -fileName qr_code_scanner -name QRCodeScanner` to convert svg icons to iconvg format and copy bytes in icons.go.

!!NOTE!!
The SVG must use fills instead of strokes or the iconvg will not be rendered correctly.  
Use `svg-fixer` pkg to fix custom icons.  
https://iconly.io/tools/svg-convert-stroke-to-fill
https://github.com/oslllo/svg-fixer