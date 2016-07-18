# go-libgen
eBook conversion utility service with webfronten and email function, tbd

depends on calibres ebook-convert bin
link: ln -s /path/to/existing/ebook-convert /usr/local/bin/
test: echo $PATH
test: which ebook-convert
test: ebook-convert /path/to/rand/pdf.pdf /tmp/out.mobi


the programm still contains a user controlled remote triggerable local file overwrite bug/vulnerabillety
imput passed to ebook-convert should be sanetized to not hold any specail characters (.. ../ / . ^ &    and so on)



