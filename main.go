package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// Создаем HTTP-клиента с настройками прокси-сервера
	proxyURL, err := url.Parse("http://64.225.8.132:10000") // Замените на свой прокси-сервер
	if err != nil {
		fmt.Println("Ошибка при разборе URL прокси-сервера:", err)
		os.Exit(1)
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("GET", "https://finance.yahoo.com/news/boeing-shares-fall-737-delivery-095832970.html", nil) // Замените на свой URL-адрес
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		os.Exit(1)
	}

	// Выполняем запрос через прокси-сервер
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Читаем и выводим ответ
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	traverseDiv(doc)
}

func traverseDiv(n *html.Node) {
	if n.Type == html.ElementNode && (n.Data == "time" || n.Data == "div") {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "caas-body" {
				divText := getTextContent(n)
				fmt.Println(divText)

			}
			if attr.Key == "class" && attr.Val == "caas-attr-meta-time" {
				divText := getTextContent(n)
				fmt.Println(divText)

			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseDiv(c)
	}

}

func getTextContent(n *html.Node) string {
	var textContent string

	if n.Type == html.TextNode {
		textContent = n.Data
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			textContent += getTextContent(c)
		}
	}

	return strings.TrimSpace(textContent)
}
