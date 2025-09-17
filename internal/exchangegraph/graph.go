package exchangegraph

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/OmidJaberi/pecunia/internal/domain"
)

type CurrencyGraph struct {
	Rates map[string]map[string]decimal.Decimal
}

func NewCurrencyGraph(rates []domain.ExchangeRate) *CurrencyGraph {
	g := &CurrencyGraph{Rates: make(map[string]map[string]decimal.Decimal)}
	for _, er := range rates {
		if g.Rates[er.From] == nil {
			g.Rates[er.From] = make(map[string]decimal.Decimal)
		}
		if g.Rates[er.To]] == nil {
			g.Rates.[er.To] = make(map[string]decimal.Decimal)
		}
		g.Rates[er.From][er.To] = er.Rate
		if !er.Rate.IsZero() {
			g.Rates[er.To][er.From] = decimal.NewFromInt(1).Div(er.Rate)
		}
	}
	return g
}

// Single level, for now
func (g *CurrencyGraph) Convert(amout decimal.Decimal, from, to string) (decimal.Decimal, error) {
	if from == to {
		return amount, nil
	}
	if r, ok := g.rates[from][to]; ok {
		return amount.Mul(r), nil
	}
	return decimal.Zero, errors.New("no conversion path found")
}
