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

	sourceRootFolder := "/home/sbuerin/code/src/github.com/sbueringer/kubecon-slides/slides"
	targetRootFolder := "/home/sbuerin/code/src/github.com/sbueringer/kubecon-slides/content"

	sourceFileInfos, err := ioutil.ReadDir(sourceRootFolder)
	if err != nil {
		panic(err)
	}

	var sourceFolders []string
	for _, f := range sourceFileInfos {
		if f.IsDir() && f.Name() ==  "2019-kubecon-na" {
			sourceFolders = append(sourceFolders, path.Join(sourceRootFolder, f.Name()))
		}
	}

	//remove := true
	remove := false

	for _, sourceFolder := range sourceFolders {
		targetFolder := path.Join(targetRootFolder, path.Base(sourceFolder))

		if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
			err := os.MkdirAll(targetFolder, 0755)
			if err != nil {
				panic(err)
			}
		}

		fileInfos, err := ioutil.ReadDir(sourceFolder)
		if err != nil {
			panic(err)
		}
		for _, f := range fileInfos {
			pdfName := path.Join(sourceFolder, f.Name())
			var imageName string
			if strings.HasSuffix(f.Name(), ".pdf") {
				imageName = path.Join(targetFolder, strings.TrimSuffix(f.Name(), ".pdf")+".png")
			} else {
				continue
			}
			//if strings.HasSuffix(f.Name(), ".pptx") {
			//	imageName = path.Join(targetFolder, strings.TrimSuffix(f.Name(), ".pptx")+".png")
			//}

			if remove {
				os.Remove(imageName)
			}

			if _, err := os.Stat(imageName); os.IsNotExist(err) {
				fmt.Printf("Creating image %s from pdf %s\n", imageName, pdfName)
				if err := ConvertPdfToJpg(pdfName, imageName); err != nil {
					log.Printf("error occured: %v", err)
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
	//if err := mw.SetResolution(300, 300); err != nil {
	//	return err
	//}

	// Load the image file into imagick
	if err := mw.ReadImage(pdfName); err != nil {
		return err
	}

	// Select only first page of pdf
	mw.SetIteratorIndex(0)

	mw.ThumbnailImage(500, 300)

	//mw.TrimImage()

	mw.SharpenImage(0, 1.0)

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	//if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN); err != nil {
	//	return err
	//}

	//if err := mw.ResizeImage(mw.GetImageWidth()/5, mw.GetImageHeight()/5, imagick.FILTER_LANCZOS, 1); err != nil {
	//	return err
	//}

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(95); err != nil {
		return err
	}

	// Convert into JPG
	if err := mw.SetFormat("png"); err != nil {
		return err
	}

	// Save File
	return mw.WriteImage(imageName)
}
