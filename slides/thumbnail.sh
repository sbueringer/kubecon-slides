#!/bin/bash

set -e

sourceDir="/home/fedora/code/gopath/src/github.com/sbueringer/kubecon-slides/slides"
targetDir="/home/fedora/code/gopath/src/github.com/sbueringer/kubecon-slides/content/post"

for sourceFile in ${sourceDir}/2017-kubecon-eu/*.pdf; do
    echo $sourceFile
    targetFile=${sourceFile%.pdf}
    targetFile="${targetFile/$sourceDir/$targetDir}.png"
    echo "Convert from $sourceFile to $targetFile"

    echo "convert -density 150 -trim ${sourceFile}[0] -quality 100 -flatten -thumbnail 500x300 -sharpen 0x1.0 ${targetFile}"
    convert -density 150 -trim "${sourceFile}[0]" -quality 100 -flatten -thumbnail 500x300 -sharpen 0x1.0 "${targetFile}"
done