# Overview

This application manages patient data including creating, updating, and deleting patient data, diagnosis data, and patient document attatchments.

## Running

### With Docker

A dockerfile is included with this applicaiton to support easily running it on any platform.  To build, navigate to the root folder and run

`docker build -t mcg-app .`

once that is complete, run the application using

`docker run -p 8080:8080 --name mcg-app mcg-app`

This will run the application locally on port 8080

### Without Docker

Ensure that the latest version of golang is installed. Visit https://go.dev/doc/install for further instructions

#### On Windows
navigate to root and run 

`go build -o mcg-app.exe`

#### On MacOS/Linux

navigate to root and run 

`go build -o mcg-app`


This will create a runnable for your os.  Simply execute this runnable using the command line.  The application will run locally on port `8080`

## Authorization

This application requires a valid bearer token in order to access the patient management functionality.  To create one first POST to `/public/users` with a username and password.  Once completed, you can then get a bearer token by supplying the same username and password in a POST request to `/public/users/login`

## Documentation

The application hosts its own documentation at `/public/docs`.  The swagger ui is available, but will not work for authenticated routes.  Additionally, an openapi spec can be found at `/public/docs/openapi.json`