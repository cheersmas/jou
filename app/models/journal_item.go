package models

import (
	"github.com/cheersmas/jou/domains"
)

type JournalItem struct {
	title string
	desc  string
}

func NewJournalItem(journal domains.Journal) JournalItem {
	return JournalItem{
		title: journal.CreatedAt.Format("2 Jan, 2006"),
		desc:  journal.Content,
	}
}

func (i JournalItem) Title() string       { return i.title }
func (i JournalItem) Description() string { return i.desc }
func (i JournalItem) FilterValue() string { return i.desc }
