package main

import (
	"flag"
	"golangify.com/snippetbox/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func main()  {
	// получаем порт из командной строки
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	flag.Parse()

	// go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log - перенаправление потока в файлы
	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)

	app := &config.Application{
		ErrLog: errLog,
		InfoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", home(app))
	mux.HandleFunc("/snippet", showSnippet(app))
	mux.HandleFunc("/snippet/create", createSnippet(app))

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  mux,
	}
	infoLog.Printf("Запуск веб-сервера на %s", *addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
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