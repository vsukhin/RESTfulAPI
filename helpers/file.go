package helpers

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

type ImageType byte

const (
	IMAGE_TYPE_UNKNOWN ImageType = iota
	IMAGE_TYPE_GIF
	IMAGE_TYPE_JPEG
	IMAGE_TYPE_PNG
)

const (
	PARAM_NAME_KEY           = "key"
	PARAM_NAME_TYPE          = "type"
	PARAM_NAME_MODEID        = "modeId"
	PARAM_NAME_SIZE          = "size"
	PARAM_MODEID_ATTACH      = "attach"
	PARAM_MODEID_DATAURL     = "dataurl"
	PARAM_MODEID_DATAURLPNG  = "dataurlpng"
	PARAM_MODEID_DATAURLJPEG = "dataurljpeg"
	PARAM_MODEID_DATAURLGIF  = "dataurlgif"
	PARAM_MODEID_OBJECT      = "object"
	PARAM_MODEID_HEAD        = "head"
	PARAM_SIZE_NUMBER        = 2
	PARAM_SIZE_WIDTH         = 0
	PARAM_SIZE_HEIGHT        = 1
	IMAGE_GIF_COLORS         = 256
	IMAGE_JPEG_QUALITY       = 100
	CONTENT_TYPE_DEFAULT     = "application/octet-stream"
	CONTENT_ENCODING_DEFAULT = "utf-8"
)

var MimeTypes = map[string]string{
	"application/javascript":    ".js",
	"application/pdf":           ".pdf",
	"application/json":          ".json",
	"application/zip":           ".zip",
	"image/png":                 ".png",
	"image/gif":                 ".gif",
	"image/jpeg":                ".jpg",
	"text/plain; charset=utf-8": ".txt",
	"text/csv; charset=utf-8":   ".csv",
	"text/css; charset=utf-8":   ".css",
	"text/xml; charset=utf-8":   ".xml",
	"text/html; charset=utf-8":  ".html",
}

func DetectImage(data []byte) (imagetype ImageType, img image.Image, imgconfig image.Config) {
	imagetypes := []ImageType{IMAGE_TYPE_GIF, IMAGE_TYPE_JPEG, IMAGE_TYPE_PNG}
	for _, imagetype := range imagetypes {
		buf := bytes.NewBuffer(data)
		var err error

		switch imagetype {
		case IMAGE_TYPE_GIF:
			imgconfig, err = gif.DecodeConfig(buf)
		case IMAGE_TYPE_JPEG:
			imgconfig, err = jpeg.DecodeConfig(buf)
		case IMAGE_TYPE_PNG:
			imgconfig, err = png.DecodeConfig(buf)
		}
		if err == nil {
			buf = bytes.NewBuffer(data)
			switch imagetype {
			case IMAGE_TYPE_GIF:
				img, err = gif.Decode(buf)
			case IMAGE_TYPE_JPEG:
				img, err = jpeg.Decode(buf)
			case IMAGE_TYPE_PNG:
				img, err = png.Decode(buf)
			}
			if err == nil {
				return imagetype, img, imgconfig
			}
		}
	}

	return IMAGE_TYPE_UNKNOWN, nil, image.Config{}
}

func ConvertImage(image image.Image, imagetype ImageType) (data []byte, err error) {
	buf := new(bytes.Buffer)
	switch imagetype {
	case IMAGE_TYPE_GIF:
		err = gif.Encode(buf, image, &gif.Options{NumColors: IMAGE_GIF_COLORS})
	case IMAGE_TYPE_JPEG:
		err = jpeg.Encode(buf, image, &jpeg.Options{Quality: IMAGE_JPEG_QUALITY})
	case IMAGE_TYPE_PNG:
		err = png.Encode(buf, image)
	default:
		log.Error("Unknown image type %v", imagetype)
		return buf.Bytes(), errors.New("Unknown image type")
	}

	return buf.Bytes(), err
}
