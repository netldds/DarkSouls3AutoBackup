package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

//https://stackoverflow.com/questions/37869793/how-do-i-zip-a-directory-containing-sub-directories-or-files-in-golang
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func checkFolderExt(fo string) (os.FileInfo, bool) {
	if info, err := os.Stat(fo); os.IsNotExist(err) {
		return nil, false
	} else {
		return info, true
	}
}

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Get Current User faild.")
		panic("Get Current User faild")
	}
	userPath := usr.HomeDir
	// fmt.Println(filepath.FromSlash(userPath))
	// C:\Users\Administrator\AppData\Roaming
	gamepath := filepath.Join(userPath, "/AppData/Roaming/DarkSoulsIII")
	gamepath = filepath.FromSlash(gamepath)
	fmt.Printf("Result of path :%v\n", gamepath)

	info, ok := checkFolderExt(gamepath)
	if !ok {
		return
	}

	if !info.IsDir() {
		fmt.Println("is not dir")
		return
	}

	outfile, err := os.Create("./DARKSOULSIII.zip")
	checkErr(err)

	w := zip.NewWriter(outfile)

	addfiles(w, gamepath+"/", "")
	defer outfile.Close()

}

func addfiles(w *zip.Writer, gamepath, baseInZip string) {

	files, err := ioutil.ReadDir(gamepath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(gamepath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(gamepath + file.Name())
			checkErr(err)

			f, err := w.Create(baseInZip + file.Name())
			checkErr(err)

			_, err = f.Write(dat)
			checkErr(err)
		} else if file.IsDir() {
			newBase := filepath.Join(gamepath, file.Name(), "/")
			fmt.Println("Recursing and Adding Subdir:" + newBase)
			addfiles(w, newBase, file.Name()+"/")
		}
	}
}
