# go-libgen
eBook conversion utility service prototype with webfronten and email function

depends on calibres ebook-convert bin
``` bash
install dependancy: (pacman -S|apt install|brew install) calibre

# OSx only #
find bin: find / -name ebook-convert 2>/dev/null
link: ln -s /path/to/existing/ebook-convert /usr/local/bin/

# test/troubleshoot #
test: echo $PATH
test: which ebook-convert
test: ebook-convert /path/to/rand/pdf.pdf /tmp/out.mobi

``` 

Warning, don't use on public networks:
the programm still contains a user controlled remote triggerable local file overwrite bug/vulnerabillety.
Input passed to ebook-convert should be sanetized to not hold any specail characters (.. ../ / . ^ &    and so on)
Local code execution not unlikely. 


tbd:
- [ ] write tests
- [ ] add email function
- [ ] make sure all imput is sanetized
- [ ] escape dangerous characters in book file names
- [ ] find a way to exec ebook-convert with the least privileges.. Sandbox if possible.
- [ ] add convert_to_format choice
- [ ] much more




