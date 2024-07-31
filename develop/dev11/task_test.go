package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setup() {
	events = make(map[string]Event)
}

func TestCreateEvent(t *testing.T) {
	setup()

	form := "title=TestEvent&user_id=1&start_time=2024-07-25T15:00:00Z&end_time=2024-07-25T16:00:00Z"
	req, err := http.NewRequest("POST", "/create_event", strings.NewReader(form))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := `{"result":"event created"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUpdateEvent(t *testing.T) {
	setup()

	// First create an event to update
	event := Event{
		ID:        "12345",
		Title:     "TestEvent",
		UserID:    "1",
		StartTime: time.Date(2024, 7, 25, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 7, 25, 16, 0, 0, 0, time.UTC),
	}
	events[event.ID] = event

	form := "id=12345&title=UpdatedEvent&user_id=1&start_time=2024-07-25T15:00:00Z&end_time=2024-07-25T17:00:00Z"
	req, err := http.NewRequest("POST", "/update_event", strings.NewReader(form))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateEventHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"result":"event updated"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteEvent(t *testing.T) {
	setup()

	// First create an event to delete
	event := Event{
		ID:        "12345",
		Title:     "TestEvent",
		UserID:    "1",
		StartTime: time.Date(2024, 7, 25, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 7, 25, 16, 0, 0, 0, time.UTC),
	}
	events[event.ID] = event

	form := "id=12345"
	req, err := http.NewRequest("POST", "/delete_event", strings.NewReader(form))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteEventHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"result":"event deleted"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestEventsForDay(t *testing.T) {
	setup()

	// Create an event
	event := Event{
		ID:        "12345",
		Title:     "TestEvent",
		UserID:    "1",
		StartTime: time.Date(2024, 7, 25, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 7, 25, 16, 0, 0, 0, time.UTC),
	}
	events[event.ID] = event

	req, err := http.NewRequest("GET", "/events_for_day?user_id=1&date=2024-07-25", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(eventsForDayHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":"12345","title":"TestEvent","user_id":"1","start_time":"2024-07-25T15:00:00Z","end_time":"2024-07-25T16:00:00Z"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
