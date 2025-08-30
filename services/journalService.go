package services

import (
	"context"

	"github.com/cheersmas/jou/domains"
	"github.com/cheersmas/jou/ports"
)

type journalService struct {
	journalRepository ports.JournalRepository
}

func (js *journalService) Create(ctx context.Context, content domains.Journal) (int, error) {
	return js.journalRepository.Create(ctx, content)
}

func (js *journalService) Read(ctx context.Context, journalId int) (domains.Journal, error) {
	return js.journalRepository.Read(ctx, journalId)
}

func (js *journalService) Update(ctx context.Context, id int, content string) (int, error) {
	return js.journalRepository.Update(ctx, id, content)
}

func (js *journalService) Delete(ctx context.Context, id int) (int, error) {
	return js.journalRepository.Delete(ctx, id)
}

func (js *journalService) ListAll(ctx context.Context) ([]domains.Journal, error) {
	return js.journalRepository.ListAll(ctx)
}

func NewJournalService(js ports.JournalRepository) *journalService {
	return &journalService{
		journalRepository: js,
	}
}
