package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/fs"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"gioui.org/gpu/headless"
	"gioui.org/layout"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
)

// From gio material theme
func HexColor(c uint32) color.NRGBA {
	return RGBA(0xff000000 | c)
}

func RGBA(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func CaptureFrame(gtx layout.Context, img *image.RGBA) error {
	rect := image.Rectangle{Max: gtx.Constraints.Max}
	w, err := headless.NewWindow(rect.Max.X, rect.Max.Y)
	if err != nil {
		return err
	}

	w.Frame(gtx.Ops)
	err = w.Screenshot(img)
	if err != nil {
		return err
	}

	return nil
}

func NewImageColor(size image.Point, color color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			img.SetRGBA(x, y, color)
		}
	}

	return img
}

func FormatHashRate(hash_rate uint64) string {
	hash_rate_string := ""

	switch {
	case hash_rate > 1000000000000:
		hash_rate_string = fmt.Sprintf("%.3f TH/s", float64(hash_rate)/1000000000000.0)
	case hash_rate > 1000000000:
		hash_rate_string = fmt.Sprintf("%.3f GH/s", float64(hash_rate)/1000000000.0)
	case hash_rate > 1000000:
		hash_rate_string = fmt.Sprintf("%.3f MH/s", float64(hash_rate)/1000000.0)
	case hash_rate > 1000:
		hash_rate_string = fmt.Sprintf("%.3f KH/s", float64(hash_rate)/1000.0)
	default:
		hash_rate_string = fmt.Sprintf("%d H/s", hash_rate)
	}
	return hash_rate_string
}

func FormatBytes(value int64) string {
	bytes_string := ""

	// math.Pow(1024, 2) for MB
	switch {
	case value > 1125899906842624:
		bytes_string = fmt.Sprintf("%.2f PB", float64(value)/1125899906842624)
	case value > 1099511627776:
		bytes_string = fmt.Sprintf("%.2f TB", float64(value)/1099511627776)
	case value > 1073741824:
		bytes_string = fmt.Sprintf("%.2f GB", float64(value)/1073741824)
	case value > 1048576:
		bytes_string = fmt.Sprintf("%.2f MB", float64(value)/1048576)
	case value > 1024:
		bytes_string = fmt.Sprintf("%.2f KB", float64(value)/1024)
	default:
		bytes_string = fmt.Sprintf("%d B", value)
	}

	return bytes_string
}

func GetFolderSize(folderPath string) (int64, error) {
	var size int64
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return nil
	})

	return size, err
}

func SplitString(s string, size int) []string {
	split := []string{}
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}

		split = append(split, s[i:end])
	}

	return split
}

func ReduceString(s string, maxLeft, maxRight int) string {
	if len(s) <= maxLeft+maxRight+3 {
		return s
	}
	start := s[:maxLeft]
	end := s[len(s)-maxRight:]
	return start + "..." + end
}

type ShiftNumber struct {
	Number   uint64
	Decimals int
}

func (s ShiftNumber) Value() float64 {
	return float64(s.Number) / math.Pow(10, float64(s.Decimals))
}

func (s ShiftNumber) Format() string {
	v := fmt.Sprintf("%%.%df", s.Decimals)
	return fmt.Sprintf(v, s.Value())
}

func ReduceAddr(addr string) string {
	return ReduceString(addr, 0, 10)
}

func ReduceTxId(txId string) string {
	return ReduceString(txId, 5, 5)
}

func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var destFile *os.File
	_, err = os.Stat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			destFile, err = os.Create(dest)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		destFile, err = os.Open(dest)
		if err != nil {
			return err
		}
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func DecodeString(value string) (string, error) {
	bytes, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func DecodeAddress(value string) (string, error) {
	p := new(crypto.Point)
	key, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}

	err = p.DecodeCompressed(key)
	if err != nil {
		return "", err
	}

	return rpc.NewAddressFromKeys(p).String(), nil
}

func HashSHA256(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))
	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum)
}

func FetchImage(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(res.Body)
	return img, err
}
