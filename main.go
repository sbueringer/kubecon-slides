package main

import (
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// Setup via:
// dnf install ImageMagick-devel
func main() {

	folder := "/home/fedora/code/gopath/src/github.com/sbueringer/hugo-test/kubecon-slides/content/post/2018-kubecon-na"
	//pdfName := "/home/fedora/code/gopath/src/github.com/sbueringer/hugo-test/kubecon-slides/content/post/2018-kubecon-na/test.pdf"
	//imageName := "/home/fedora/code/gopath/src/github.com/sbueringer/hugo-test/kubecon-slides/content/post/2018-kubecon-na/test.jpg"

	remove := true

	fileInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	for _, f := range fileInfos {
		if strings.HasSuffix(f.Name(), ".pdf") {
			pdfName := path.Join(folder, f.Name())
			imageName := strings.TrimSuffix(pdfName, ".pdf") + ".jpg"

			if remove {
				os.Remove(imageName)
			}

			if _, err := os.Stat(imageName); os.IsNotExist(err) {
				fmt.Printf("Creating imge %s from pdf %s\n", imageName, pdfName)
				if err := ConvertPdfToJpg(pdfName, imageName); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

// ConvertPdfToJpg will take a filename of a pdf file and convert the file into an
// image which will be saved back to the same location. It will save the image as a
// high resolution jpg file with minimal compression.
func ConvertPdfToJpg(pdfName string, imageName string) error {

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Must be *before* ReadImageFile
	// Make sure our image is high quality
	if err := mw.SetResolution(300, 300); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(pdfName); err != nil {
		return err
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN); err != nil {
		return err
	}



	if err :=  mw.ResizeImage(mw.GetImageWidth()/5, mw.GetImageHeight()/5, imagick.FILTER_LANCZOS, 1); err != nil {
		return err
	}

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(95); err != nil {
		return err
	}

	// Select only first page of pdf
	mw.SetIteratorIndex(0)

	// Convert into JPG
	if err := mw.SetFormat("jpg"); err != nil {
		return err
	}

	// Save File
	return mw.WriteImage(imageName)
}
