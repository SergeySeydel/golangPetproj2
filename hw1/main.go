package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

func main() {
	//NewServeMux возвращает указатель на "пустую структуру" ServeMux
	//newMux := http.NewServeMux()

	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/sub", subHandler)
	http.HandleFunc("/mul", mulHandler)
	http.HandleFunc("/div", divHandler)
	http.HandleFunc("/", rootHandler)

	fmt.Println("Server listening on 127.0.0.1:1234")
	http.ListenAndServe("127.0.0.1:1234", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	info := map[string]string{
		"message": "For more information type /info.",
	}
	jsonResponse(w, info)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	info := map[string]string{
		"info": "This is an API for basic arithmetic operations. Endpoints: /info, /first, /second, /add, /sub, /mul, /div.",
	}
	jsonResponse(w, info)
}

// Обработчик для эндпоинта /first, возвращает случайное число
func firstHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	number := rand.Intn(100)
	jsonResponse(w, map[string]int{"number": number})
}

// Обработчик для эндпоинта /second, возвращает случайное число
func secondHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	number := rand.Intn(100)
	jsonResponse(w, map[string]int{"number": number})
}

// Обработчик для эндпоинта /add, возвращает сумму двух случайных чисел
func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	number1 := rand.Intn(100)
	number2 := rand.Intn(100)
	result := number1 + number2
	jsonResponse(w, map[string]int{"first": number1, "second": number2, "result": result})
}

// Обработчик для эндпоинта /sub, возвращает разность двух случайных чисел
func subHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	number1 := rand.Intn(100)
	number2 := rand.Intn(100)
	result := number1 - number2
	jsonResponse(w, map[string]int{"first": number1, "second": number2, "result": result})
}

// Обработчик для эндпоинта /mul, возвращает произведение двух случайных чисел
func mulHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	number1 := rand.Intn(100)
	number2 := rand.Intn(100)
	result := number1 * number2
	jsonResponse(w, map[string]int{"first": number1, "second": number2, "result": result})
}

// Обработчик для эндпоинта /div, возвращает результат деления двух случайных чисел
func divHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	number1 := rand.Intn(100)
	number2 := rand.Intn(100)
	if number2 == 0 {
		number2 = 1
	}
	result := float64(number1) / float64(number2)
	jsonResponse(w, map[string]float64{"first": float64(number1), "second": float64(number1), "result": result})
}

// Функция для отправки JSON-ответа клиенту
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
