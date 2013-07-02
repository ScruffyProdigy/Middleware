# Staticer
A middleware that serves static files to the user

## Installation
`go get github.com/ScruffyProdigy/Middleware/staticer`

## Usage
* add the result of staticer.New() to your rack
	* first parameter should be the url route that the user must request from to get to the static files
	* the second parameter should be the location that this program can find the static files to serve back to the user