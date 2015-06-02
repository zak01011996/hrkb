package utils

import (
    "github.com/nfnt/resize"
    "github.com/astaxie/beego"
    "image/jpeg"
    "image/png"
    "image/gif"
    "os"
    "strings"
    "image"
    "mime/multipart"
)

func Resize(filepath string, fheader *multipart.FileHeader,  width uint, height uint) (string, string) {

	var err error

	file, err := os.Open(filepath)
	if err != nil {
		beego.Error(err)
		return "", ""
	}
	defer os.Remove(filepath)

	ftype := fheader.Header["Content-Type"]

	var img image.Image

	switch ftype[0] {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(file)
		if err != nil {
			beego.Error(err)
			return "", ""
		}
		file.Close()
	case "image/png":
		img, err = png.Decode(file)
		if err != nil {
			beego.Error(err)
			return "", ""
		}
		file.Close()
	case "image/gif":
		img, err = gif.Decode(file)
		if err != nil {
			beego.Error(err)
			return "", ""
		}
		file.Close()
	default:
		beego.Error("Unsupported file type")
		return "", ""
	}

	// and preserve aspect ratio
	m := resize.Resize(width, height, img, resize.Lanczos3)

	ftyps := strings.Split(fheader.Filename, ".")
	sftype := ftyps[len(ftyps) - 1]

	sname := RandString(10) + "." + sftype

	savepath := beego.AppConfig.String("tmp_dir") + sname

	out, err := os.Create(savepath)

	if err != nil {
		beego.Error(err)
		return "", ""
	}
	defer out.Close()

	// write new image to file
	switch ftype[0] {
	case "image/jpeg", "image/jpg":
		if err := jpeg.Encode(out, m, nil); err != nil {
			beego.Error(err)
		}
	case "image/png":
		if err := png.Encode(out, m); err != nil {
			beego.Error(err)
		}
	case "image/gif":
		if err := gif.Encode(out, m, nil); err != nil {
			beego.Error(err)
		}
	default:
		beego.Error("Invalid file type")
		return "", ""
	}

	return savepath, sname
}

