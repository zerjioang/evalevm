package store

import (
	"context"
	"evalevm/internal/datatype"
)

type EmbedRepository struct {
	analyzerList []datatype.Analyzer
}

var _ Repository = (*EmbedRepository)(nil)

func NewEmbedRepository(analyzerList []datatype.Analyzer) *EmbedRepository {
	return &EmbedRepository{
		analyzerList: analyzerList,
	}
}

func (e EmbedRepository) AnalyzerList(ctx context.Context) []datatype.Analyzer {
	return e.analyzerList
}
