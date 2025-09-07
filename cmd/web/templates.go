package main

import "snippetbox.adameury.dev/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
