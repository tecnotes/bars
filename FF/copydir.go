package main

import (
   "os"
   "flag"
   "fmt"
   "io"
   //"log"
   "crypto/md5"
   "encoding/hex"
)



func CopyFile(source string, dest string) (err error) {
    sourcefile, err := os.Open(source)
    if err != nil {
        return err
    }

    defer sourcefile.Close()
	 
    destfile, err := os.Create(dest)
    if err != nil {
		return err
    }

    defer destfile.Close()
	 
    _, err = io.Copy(destfile, sourcefile)
    if err == nil {
        sourceinfo, err := os.Stat(source)
        if err != nil {
            err = os.Chmod(dest, sourceinfo.Mode())
        }
    }
    return
}

func CopyDir(source string, dest string) (err error) {

    // get properties of source dir
    sourceinfo, err := os.Stat(source)
    if err != nil {
        return err
    }

    // create dest dir
    err = os.MkdirAll(dest, sourceinfo.Mode())
    if err != nil {
        return err
    }

    directory, _ := os.Open(source)

    objects, err := directory.Readdir(-1)

    for _, obj := range objects {
		obj := obj
      
		sourcefilepointer := source + "/" + obj.Name()
        destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
            // create sub-directories - recursively
            err = CopyDir(sourcefilepointer, destinationfilepointer)
			fmt.Println(sourcefilepointer)
            if err != nil {
                fmt.Println(err)
            }
        } else {
            // perform copy
            err = CopyFile(sourcefilepointer, destinationfilepointer)
			fmt.Println(sourcefilepointer)
			fmt.Println(err)
            if err != nil {
                fmt.Println(err)
            }
			//fmt.Println(destinationfilepointer)
			md1, err := md5sum(sourcefilepointer)
			fmt.Println(md1)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(destinationfilepointer)
			md2, err := md5sum(destinationfilepointer)
			fmt.Println(md2)
			if err != nil {
				fmt.Println(err)
			}
        }
    }
    return
}

func md5sum(filePath string) (result string, err error) {
    file, err := os.Open(filePath)
    if err != nil {
        return
    }
    defer file.Close()

    hash := md5.New()
    _, err = io.Copy(hash, file)
    if err != nil {
        return
    }

    result = hex.EncodeToString(hash.Sum(nil))
    return
}

func main() {
    flag.Parse() // get the source and destination directory

    source_dir := flag.Arg(0) // get the source directory from 1st argument

    dest_dir := flag.Arg(1) // get the destination directory from the 2nd argument

    fmt.Println("Source :" + source_dir)

    // check if the source dir exist
    src, err := os.Stat(source_dir)
    if err != nil {
      panic(err)
    }

    if !src.IsDir() {
     fmt.Println("Source is not a directory")
     os.Exit(1)
    }

    // create the destination directory
    fmt.Println("Destination :"+ dest_dir)

    _, err = os.Open(dest_dir)
    if !os.IsNotExist(err) {
        fmt.Println("Destination directory already exists. Abort!")
        os.Exit(1)
    }

    err = CopyDir(source_dir, dest_dir)
    if err != nil {
      fmt.Println(err)
    } else {
        fmt.Println("Directory copied")
    }
 }