###this assumes you have your GOPATH and PATH set correctly
## for example
#export GOPATH=~/tropical.src/
#export PATH=~/tropical.src/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin
###and that you have done "go get" to install gopherjs 
##such as go get github.com/gopherjs/gopherjs

#HTML_DIR is where you want the html/js output to be placed. this directory
#must be visible through a web browser.  the setup here assumes that 
#you are putting things into the distributation directory for use with
#github pages.
HTML_DIR:=../iansmith.github.io/tropical

build:
	gopherjs install github.com/iansmith/mydumbapp
	cp $(GOPATH)/bin/mydumbapp.js ../iansmith.github.io/tropical
	cp $(GOPATH)/bin/mydumbapp.js.map ../iansmith.github.io/tropical

