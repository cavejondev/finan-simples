// Package transaction define os tipos e regras de transações financeiras
package transaction

// Type representa o tipo da transação
type Type string

const (
	// Income representa entrada de dinheiro
	Income Type = "INCOME"

	// Expense representa saída de dinheiro
	Expense Type = "EXPENSE"

	// Transfer representa transferência entre contas
	Transfer Type = "TRANSFER"
)
