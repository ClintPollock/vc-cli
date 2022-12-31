# Veracode Fix overview

This CLI allows demonstrating the Veracode Fix intelligent remediation feature. 

## Pre-requisites

1. Install Go on the machine that you'll use to demo the feature: 
	* [Installation instructions](https://go.dev/doc/install)
1. Make sure that you have valid Veracode API credentials in the US region configured in your `~/.veracode/credentials` file.

## Build and Run

At a terminal, cd to the `veracode-fix` directory created when you unzipped this package. Then execute the following from a Mac or Linux command line:

```
cd veracode
go build
cp ../demofiles/*.* .
export APIHOST=18.158.216.3:8080
./veracode-cli fix SecureLogger.java --apihost=$APIHOST
```

For a Windows machine you will need to substitute the appropriate shell commands above and use `veracode-cli.exe` instead.


## Notes

1. The first flaw that the CLI returns should work successfully. I've observed that the second one sometimes has an error. 
2. The CLI will automatically back up the source file (`SecureLogger.java`) when applying the patch. Don't forget to restore the source file from backup between demos.
3. I find it helpful to demo in the Visual Studio Code terminal so that I can show the patch file that the demo downloads.
4. You may not be able to do the demo the first time you try as the services may need to spin up. If this happens, wait about 15 minutes and try again.
