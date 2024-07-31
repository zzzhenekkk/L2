package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Event представляет событие в календаре
type Event struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UserID    string    `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// JSONResponse представляет стандартный ответ в формате JSON
type JSONResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

var events = make(map[string]Event) // In-memory store for events

// toJSON сериализует объект в JSON
func toJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// fromJSON десериализует объект из JSON
func fromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// parseAndValidateEvent парсит и валидирует входные данные для события
func parseAndValidateEvent(r *http.Request) (Event, error) {
	var event Event
	if err := r.ParseForm(); err != nil {
		return event, fmt.Errorf("invalid request body: %v", err)
	}

	event.Title = r.FormValue("title")
	event.UserID = r.FormValue("user_id")
	startTimeStr := r.FormValue("start_time")
	endTimeStr := r.FormValue("end_time")

	if event.Title == "" || event.UserID == "" || startTimeStr == "" || endTimeStr == "" {
		return event, fmt.Errorf("missing required fields")
	}

	var err error
	event.StartTime, err = time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return event, fmt.Errorf("invalid start_time: %v", err)
	}

	event.EndTime, err = time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return event, fmt.Errorf("invalid end_time: %v", err)
	}

	if event.EndTime.Before(event.StartTime) {
		return event, fmt.Errorf("end_time cannot be before start_time")
	}

	return event, nil
}

// parseAndValidateID парсит и валидирует ID события
func parseAndValidateID(r *http.Request) (string, error) {
	id := r.FormValue("id")
	if id == "" {
		return "", fmt.Errorf("missing id")
	}
	return id, nil
}

// parseAndValidateUserIDAndDate парсит и валидирует user_id и дату
func parseAndValidateUserIDAndDate(r *http.Request) (string, time.Time, error) {
	userID := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	if userID == "" || dateStr == "" {
		return "", time.Time{}, fmt.Errorf("missing user_id or date")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid date: %v", err)
	}

	return userID, date, nil
}

// createEventHandler обработчик для создания события
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	event, err := parseAndValidateEvent(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	event.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	events[event.ID] = event

	response, _ := toJSON(JSONResponse{Result: "event created"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// updateEventHandler обработчик для обновления события
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseAndValidateID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	event, err := parseAndValidateEvent(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	if _, exists := events[id]; !exists {
		http.Error(w, `{"error": "event not found"}`, http.StatusNotFound)
		return
	}

	event.ID = id
	events[id] = event

	response, _ := toJSON(JSONResponse{Result: "event updated"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// deleteEventHandler обработчик для удаления события
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseAndValidateID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	if _, exists := events[id]; !exists {
		http.Error(w, `{"error": "event not found"}`, http.StatusNotFound)
		return
	}

	delete(events, id)

	response, _ := toJSON(JSONResponse{Result: "event deleted"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// eventsForDayHandler обработчик для получения событий за день
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseAndValidateUserIDAndDate(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	startOfDay := date
	endOfDay := date.Add(24 * time.Hour)

	var results []Event
	for _, event := range events {
		if event.UserID == userID && event.StartTime.After(startOfDay) && event.StartTime.Before(endOfDay) {
			results = append(results, event)
		}
	}

	response, _ := toJSON(results)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// eventsForWeekHandler обработчик для получения событий за неделю
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseAndValidateUserIDAndDate(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	startOfWeek := date
	endOfWeek := date.AddDate(0, 0, 7)

	var results []Event
	for _, event := range events {
		if event.UserID == userID && event.StartTime.After(startOfWeek) && event.StartTime.Before(endOfWeek) {
			results = append(results, event)
		}
	}

	response, _ := toJSON(results)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// eventsForMonthHandler обработчик для получения событий за месяц
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseAndValidateUserIDAndDate(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	startOfMonth := date
	endOfMonth := date.AddDate(0, 1, 0)

	var results []Event
	for _, event := range events {
		if event.UserID == userID && event.StartTime.After(startOfMonth) && event.StartTime.Before(endOfMonth) {
			results = append(results, event)
		}
	}

	response, _ := toJSON(results)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// loggingMiddleware middleware для логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// main основная функция, запускающая сервер
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", createEventHandler)
	mux.HandleFunc("/update_event", updateEventHandler)
	mux.HandleFunc("/delete_event", deleteEventHandler)
	mux.HandleFunc("/events_for_day", eventsForDayHandler)
	mux.HandleFunc("/events_for_week", eventsForWeekHandler)
	mux.HandleFunc("/events_for_month", eventsForMonthHandler)

	loggedMux := loggingMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: loggedMux,
	}

	log.Println("Starting server on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
