package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "win-go-zipper"
	app.Usage = "Create zips for aws lambda from windows using folders"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "output file path for the zip. Defaults to the input file name.",
		},
	}

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			return errors.New("No input provided")
		}

		dir := c.Args().First()
		outputZip := c.String("output")
		if outputZip == "" {
			outputZip = fmt.Sprintf("%s.zip", filepath.Base(dir))
		}

		paths, err := getPaths(dir)

		if err != nil {
			return fmt.Errorf("Failed to get files from %s", dir)
		}

		err = compressFiles(paths, dir, outputZip)
		if err != nil {
			return fmt.Errorf("Failed to compress files")
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func getPaths(d string) ([]string, error) {
	paths := []string{}
	err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return paths, nil
}

func compressFiles(fPaths []string, path string, outputZip string) error {
	newfile, err := os.Create(outputZip)
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	for _, fp := range fPaths {

		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(fp)
		if err != nil {
			return err
		}

		writer, err := zipWriter.CreateHeader(&zip.FileHeader{
			CreatorVersion: 3 << 8,     // indicates Unix
			ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
			Name:           fp,
			Method:         zip.Deflate,
		})

		_, err = writer.Write(data)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully zipped %s\n", fp)
	}

	fmt.Println("Finished zipping")
	return nil
}
