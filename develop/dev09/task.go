package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

// downloadFile скачивает файл по URL и сохраняет его в указанное место.
func downloadFile(filepath string, url string) error {
	// Создаем директорию для файла, если ее нет
	if err := os.MkdirAll(path.Dir(filepath), os.ModePerm); err != nil {
		return err
	}

	// Открываем файл для записи
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Выполняем HTTP-запрос
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Записываем тело ответа в файл
	_, err = io.Copy(out, resp.Body)
	return err
}

// convertLinks парсит HTML и заменяет ссылки на локальные пути.
func convertLinks(baseURL *url.URL, filepath string) error {
	// Открываем HTML файл для чтения
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return err
	}

	// Открываем HTML файл для записи
	outputFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var convert func(*html.Node)
	convert = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, attr := range n.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					link, err := baseURL.Parse(attr.Val)
					if err != nil {
						continue
					}
					if link.Host == baseURL.Host {
						relPath := link.Path
						n.Attr[i].Val = path.Join(baseURL.Path, relPath)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			convert(c)
		}
	}

	convert(doc)

	// Записываем измененный HTML в файл
	html.Render(outputFile, doc)

	return nil
}

// extractLinks парсит HTML и извлекает все ссылки.
func extractLinks(baseURL *url.URL, body io.Reader) ([]string, error) {
	var links []string
	tokenizer := html.NewTokenizer(body)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return links, nil
			}
			return nil, tokenizer.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "a", "link", "img", "script":
				for _, attr := range token.Attr {
					if attr.Key == "href" || attr.Key == "src" {
						link, err := baseURL.Parse(attr.Val)
						if err != nil {
							continue
						}
						links = append(links, link.String())
					}
				}
			}
		}
	}
}

// downloadPage скачивает страницу и все связанные с ней ресурсы.
func downloadPage(baseURL string, root string) error {
	// Выполняем HTTP-запрос для главной страницы
	resp, err := http.Get(baseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Сохраняем главную страницу
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	filepath := path.Join(root, u.Path)
	if strings.HasSuffix(filepath, "/") {
		filepath = path.Join(filepath, "index.html")
	}
	if err := downloadFile(filepath, baseURL); err != nil {
		return err
	}

	// Конвертируем ссылки в главной странице
	if err := convertLinks(u, filepath); err != nil {
		return err
	}

	// Извлекаем ссылки
	links, err := extractLinks(u, resp.Body)
	if err != nil {
		return err
	}

	// Скачиваем все найденные ссылки
	for _, link := range links {
		fmt.Println("Downloading", link)

		// Проверяем, что длина link больше длины baseURL
		if len(link) < len(baseURL) {
			fmt.Println("Skipping link (too short):", link)
			continue
		}

		relPath, err := url.PathUnescape(link[len(baseURL):])
		if err != nil {
			relPath = link[len(baseURL):]
		}
		if relPath == "" {
			relPath = "index.html"
		}
		if err := downloadFile(path.Join(root, relPath), link); err != nil {
			fmt.Println("Failed to download", link, ":", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wget URL [destination]")
		return
	}

	baseURL := os.Args[1]
	root := "site"
	if len(os.Args) > 2 {
		root = os.Args[2]
	}

	// Проверяем, существует ли root и является ли он директорией
	fileInfo, err := os.Stat(root)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(root, os.ModePerm); err != nil {
			fmt.Println("Error creating root directory:", err)
			return
		}
	} else if err != nil {
		fmt.Println("Error accessing root directory:", err)
		return
	} else if !fileInfo.IsDir() {
		fmt.Println("Error:", root, "is not a directory")
		return
	}

	if err := downloadPage(baseURL, root); err != nil {
		fmt.Println("Error:", err)
	}
}
