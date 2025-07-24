package store

import (
	"context"
	"evalevm/internal/datatype"
)

type Repository interface {
	AnalyzerList(ctx context.Context) []datatype.Analyzer
}
