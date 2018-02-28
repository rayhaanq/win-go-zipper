# win-go-zipper
Creates a zip file which contains all the files in a folder which can be used for aws lambda from windows. This is a modified version of the zipper provided by aws which also sets the file permissions to 777.

# Usage:
`win-go-zipper.exe -o <output-zip-file-name>.zip <folder-containing-zips>`

# Example:
`win-go-zipper.exe -o myservice.zip bin`

To use with the serverless framework edit your serverless.yml and add the generated zip file to the packages artifact. Read more here https://serverless.com/framework/docs/providers/aws/guide/packaging/

You can also create a bat file which builds all your binaries and then executes the zipper
