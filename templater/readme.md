# Templater
A middleware to prepare templates for future middleware

## Installation
go get github.com/ScruffyProdigy/Middleware/templater

##  Documentation
http://godoc.org/github.com/ScruffyProdigy/Middleware/templater

## Usage
* Store all of your templates in one folder, or in subfolders of that folder
* Call templater.GetTemplates() and add the result to your rack
	* It needs to be added before anything that takes advantage of the following functions
* To find out whether a template exists, call (templater.V)(vars).Exists()
	* Pass it the name (including the subfolder names) of the template you're looking for
* To render a template, call (templater.V)(vars).Render()
	* The first parameter is the name of the template (including subfolder names)
	* The second parameter is where you want to render it to (any io.Writer will do)