package ports

import (
	"context"

	"github.com/cheersmas/jou/domains"
)

type JournalRepository interface {
	Create(ctx context.Context, content domains.Journal) (int, error)
	Read(ctx context.Context, journalId int) (domains.Journal, error)
	Update(ctx context.Context, id int, content string) (int, error)
	Delete(ctx context.Context, id int) (int, error)
}
