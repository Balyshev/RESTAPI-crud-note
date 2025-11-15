package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Zametka struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var zametki = make(map[int]*Zametka) //хранилище заметок
var NextId = 1

func getAllZametki(w http.ResponseWriter, r *http.Request) {
	result := []*Zametka{}
	for _, i := range zametki {
		result = append(result, i)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func createZametka(w http.ResponseWriter, r *http.Request) {
	var novayaZametka Zametka

	err := json.NewDecoder(r.Body).Decode(&novayaZametka)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	novayaZametka.ID = NextId
	NextId++

	zametki[novayaZametka.ID] = &novayaZametka

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novayaZametka)
}
func getZametkaByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id должен быть числом", http.StatusBadRequest)
		return
	}
	zametka, i := zametki[id]
	if !i {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zametka)

}

func updateZametka(w http.ResponseWriter, r *http.Request) {
	//получаем id из url
	path := r.URL.Path                           //возвращаем строку из URL HTTP запроса
	idStr := strings.TrimPrefix(path, "/notes/") //убираем все лишнее, оставив число
	id, err := strconv.Atoi(idStr)               //делаем из строки "5" число 5
	if err != nil {
		http.Error(w, "id должно быть числом", http.StatusBadRequest)
		return
	}
	//-Проверяем существует ли заметка
	z, exists := zametki[id]
	if !exists {
		http.Error(w, "Заметка не найдена", http.StatusNotFound)
		return
	}
	//читаем новый данные от клиента
	var updateData Zametka
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	z.Title = updateData.Title
	z.Content = updateData.Content

	w.Header().Set("Content-Type", "application/json") //объявляем заголовок
	json.NewEncoder(w).Encode(z)

}

func main() {
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getAllZametki(w, r)
		} else if r.Method == "POST" {
			createZametka(w, r)
		} else {
			http.Error(w, "Метод не зарегистрирован", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Если GET - получить заметку по ID
			getZametkaByID(w, r)
		} else if r.Method == "PUT" {
			updateZametka(w, r)
		} else {
			// Если другой метод - ошибка
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	fmt.Println("📱 Веб-интерфейс: http://localhost:8090")
	fmt.Println("Доступные endpoints:")
	fmt.Println("  GET  /notes  - получить все заметки")
	fmt.Println("  POST /notes  - создать заметку")
	fmt.Println("  PUT    /notes/{id} - обновить заметку")

	// Раздача статических файлов (index.html)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
		} else {
			http.NotFound(w, r)
		}
	})

	log.Fatal(http.ListenAndServe(":8090", nil))
}
