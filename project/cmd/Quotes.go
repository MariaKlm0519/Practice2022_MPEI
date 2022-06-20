package main

import (
	"errors"
)

type Quote struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type QuotesStore struct {
	quotes map[int]Quote
	nextId int
}

var ErrNoRecord = errors.New("QuotesStore: подходящей записи не найдено")

func NewStore() *QuotesStore {
	ts := &QuotesStore{}
	ts.quotes = make(map[int]Quote)
	ts.nextId = 1
	return ts
}

// AddQuote создаёт новую запись в хранилище.
func (ts *QuotesStore) AddQuote(title string, text string) int {
	quote := Quote{
		ID:    ts.nextId,
		Title: title,
		Text:  text}

	ts.quotes[ts.nextId] = quote
	ts.nextId++
	return quote.ID
}

// GetQuote получает цитату из хранилища по ID. Если ID не существует - будет возвращена ошибка.
func (ts *QuotesStore) GetQuote(id int) (Quote, error) {
	t, ok := ts.quotes[id]
	if ok {
		return t, nil
	} else {
		return Quote{}, ErrNoRecord
	}
}

// DeleteQuote удаляет цитату с заданным ID. Если ID не существует - будет возвращена ошибка.
//func (ts *QuotesStore) DeleteQuote(id int) error {}

// DeleteAllQuotes удаляет из хранилища все цитаты.
//func (ts *QuotesStore) DeleteAllQuotes() error {}
