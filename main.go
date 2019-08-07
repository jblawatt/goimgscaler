package main

import (
	"crypto/sha1"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"

	"github.com/disintegration/imaging"
	// https://golang.org/pkg/net/http/
	"net/http"
)

// https://yourbasic.org/golang/iota/
// https://github.com/nfnt/resize

func resizeImage() {

}

// MustInt converts a string to int. On error the
// default value will be returned.
func MustInt(input string, default_ int) int {
	if value, err := strconv.Atoi(input); err != nil {
		return default_
	} else {
		return value
	}
}

var InterpolationList = map[int]imaging.ResampleFilter{
	0: imaging.Lanczos,
	1: imaging.BSpline,
}

type ImageMethod int

const (
	MethodResize ImageMethod = iota
	MethodFill
	MethodFit
	endMethod
)

func validateImageMethod(value int) error {
	if value < int(endMethod) {
		return nil
	}
	return NewBadRequest("Invalid Method")
}

type Options struct {
	CacheDir      string
	ImageDir      string
	DefaultMethod ImageMethod
}

func GetOptions() Options {
	return Options{
		CacheDir:      viper.GetString("cache_dir"),
		ImageDir:      viper.GetString("image_dir"),
		DefaultMethod: MethodFill,
	}
}

func hashIt(filename string, method, height, width, anchor, interpolation int) string {
	h := sha1.New()
	io.WriteString(h, strconv.FormatInt(int64(height), 10))
	io.WriteString(h, strconv.FormatInt(int64(width), 10))
	io.WriteString(h, filename)
	io.WriteString(h, strconv.FormatInt(int64(interpolation), 10))
	io.WriteString(h, strconv.FormatInt(int64(method), 10))
	io.WriteString(h, strconv.FormatInt(int64(anchor), 10))
	return fmt.Sprintf("%x", string(h.Sum(nil)))
}

func MustCacheDir(cacheDir string) {
	if err := os.MkdirAll(cacheDir, os.ModeDir); err != nil {
		if os.IsExist(err) {
			log.Println("Cache dir already exists")
		} else {
			log.Fatal(err)
		}
	}
}

type ResampleFilter int

const (
	NearestNeighbor ResampleFilter = iota
	Box
	Linear
	Hermite
	MitchellNetravali
	CatmullRom
	BSpline
	Gaussian
	Bartlett
	Lanczos
	Hann
	Hamming
	Blackman
	Welch
	Cosine
	endRF
)

func GetResampleFilter(input ResampleFilter) (imaging.ResampleFilter, error) {
	if int(endRF) < int(input) {
		return imaging.ResampleFilter{}, NewBadRequest(fmt.Sprintf("Invalid resample filter: %s", input))
	}
	switch input {
	case NearestNeighbor:
		return imaging.NearestNeighbor, nil
	case Box:
		return imaging.Box, nil
	case Linear:
		return imaging.Linear, nil
	case Hermite:
		return imaging.Hermite, nil
	case MitchellNetravali:
		return imaging.MitchellNetravali, nil
	case CatmullRom:
		return imaging.CatmullRom, nil
	case BSpline:
		return imaging.BSpline, nil
	case Gaussian:
		return imaging.Gaussian, nil
	case Bartlett:
		return imaging.Bartlett, nil
	case Lanczos:
		return imaging.Lanczos, nil
	case Hann:
		return imaging.Hann, nil
	case Hamming:
		return imaging.Hamming, nil
	case Blackman:
		return imaging.Blackman, nil
	case Welch:
		return imaging.Welch, nil
	case Cosine:
		return imaging.Cosine, nil
	default:
		return imaging.NearestNeighbor, nil
	}
}

func applyImage(
	filename string, method ImageMethod,
	height, width int, interp ResampleFilter, anchor imaging.Anchor,
	opts Options) (image.Image, error) {

	// TODO: create safe path
	imagePath := path.Join(opts.ImageDir, filename)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, FileNotFound{filename}
	}

	ext := filepath.Ext(filename)
	cacheHash := hashIt(filename, int(method), height, width, int(anchor), int(interp))
	cacheName := fmt.Sprintf("%s%s", cacheHash, ext)

	MustCacheDir(opts.CacheDir)

	cacheFullPath := path.Join(opts.CacheDir, cacheName)

	if _, err := os.Stat(cacheFullPath); os.IsNotExist(err) {

		file, _ := os.Open(imagePath)
		img, _ := jpeg.Decode(file)
		file.Close()
		var m image.Image

		filter, err := GetResampleFilter(interp)
		if err != nil {
			return nil, err
		}

		switch method {
		case MethodResize:
			m = imaging.Resize(img, width, height, filter)
		case MethodFit:
			m = imaging.Fit(img, width, height, filter)
		case MethodFill:
			m = imaging.Fill(img, width, height, imaging.Anchor(anchor), filter)
		default:
			m = imaging.Resize(img, width, height, filter)
		}

		if out, oerr := os.Create(cacheFullPath); oerr != nil {
			log.Fatal(oerr)
		} else {
			defer out.Close()
			jpeg.Encode(out, m, nil)
			return m, nil
		}

	} else {
		file, _ := os.Open(cacheFullPath)
		img, _ := jpeg.Decode(file)
		defer file.Close()
		return img, nil
	}

	return nil, nil
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("f")

	anchor := MustInt(r.URL.Query().Get("a"), 0)
	height := MustInt(r.URL.Query().Get("h"), 0)
	width := MustInt(r.URL.Query().Get("w"), 0)
	interpolation := MustInt(r.URL.Query().Get("i"), 0)
	method := MustInt(r.URL.Query().Get("m"), 0)

	if img, err := applyImage(filename, ImageMethod(method), height, width, ResampleFilter(interpolation), imaging.Anchor(anchor), GetOptions()); err != nil {
		if _, ok := err.(BadRequest); ok {
			w.WriteHeader(http.StatusBadRequest)
		}
		if _, ok := err.(FileNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
		}
		io.WriteString(w, err.Error())
	} else {
		jpeg.Encode(w, img, nil)
	}

}

func main() {
	http.HandleFunc("/", imageHandler)
	// TODO: /list

	viper.SetDefault("cache_dir", "_cache")
	viper.SetDefault("image_dir", "input")
	viper.SetDefault("bind", "127.0.0.1:8080")
	viper.SetDefault("default_method", MethodResize)
	viper.SetDefault("default_filter", NearestNeighbor)
	viper.SetDefault("default_anchor", imaging.Center)

	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	log.Fatal(http.ListenAndServe(viper.GetString("bind"), nil))
}
