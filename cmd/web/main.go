package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func main()  {
	// получаем порт из командной строки
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Запуск веб-сервера на %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// This returns an FileInfo type
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		index := filepath.Join(path, "index.html")
		//  Если файла нет, то err = os.ErrNorExist, которая преобразуется в 404 ответ на сервере
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}