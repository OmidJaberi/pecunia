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
		if g.Rates[er.To] == nil {
			g.Rates[er.To] = make(map[string]decimal.Decimal)
		}
		g.Rates[er.From][er.To] = er.Rate
		if !er.Rate.IsZero() {
			g.Rates[er.To][er.From] = decimal.NewFromInt(1).Div(er.Rate)
		}
	}
	return g
}

func (g *CurrencyGraph) Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, error) {
	if from == to {
		return amount, nil
	}
	if r, ok := g.Rates[from][to]; ok {
		return amount.Mul(r), nil
	}

	type queue_element struct {
		currency string
		rate decimal.Decimal
	}

	var queue []queue_element
	queue = append(queue, queue_element{from, decimal.NewFromInt(1)})
	
	marked := make(map[string]bool)
	marked[from] = true
	
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		if v.currency == to {
			return amount.Mul(v.rate), nil
		}

		for child, rate := range g.Rates[v.currency] {
			if !marked[child] {
				marked[child] = true
				queue = append(queue, queue_element{child, v.rate.Mul(rate)})
			}
		} 
	}

	return decimal.Zero, errors.New("no conversion path found")
}
