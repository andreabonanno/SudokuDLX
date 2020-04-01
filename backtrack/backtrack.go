package backtrack

import "strconv"

type Board struct {
	Cells [9][9]int
}

// Restituisce la griglia in formato testuale.
func (b *Board) String() string {
	var output string

	for row, y := range b.Cells {
		for col, x := range y {
			output += strconv.Itoa(x)
			output += " "
			if col == 2 || col == 5 {
				output += "| "
			}
		}
		output += "\n"
		if (row+1)%3 == 0 {
			output += "\n"
		}
	}

	return output
}

// Procedura ricorsiva di backtracking. Il risultato indica se la soluzione e' stata trovata o meno.
func (b *Board) Backtrack() bool {
	nextRow, nextCol, hasEmptyCell := b.findEmptyCell()
	if !hasEmptyCell {
		return true
	}

	for candidate := 1; candidate <= 9; candidate++ {
		if b.isDigitValid(nextRow, nextCol, candidate) {
			b.Cells[nextRow][nextCol] = candidate

			if b.Backtrack() {
				return true
			}
			// reset
			b.Cells[nextRow][nextCol] = 0
		}
	}

	return false
}

// findEmptyCell controlla la presenza di caselle vuota, restituendo riga e colonna di una casella vuota, 0 altrimenti
func (b *Board) findEmptyCell() (int, int, bool) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if b.Cells[row][col] == 0 {
				return row, col, true
			}
		}
	}

	return 0, 0, false
}

// isDigitValid controlla se la cifra corrispondente soddifa i vincoli di riga colonna e gruppo del Sudoku
func (b *Board) isDigitValid(row, col, digit int) bool {
	// note: integer division 'rounds down' by ignoring all decimal places
	startRow := row / 3 * 3
	startCol := col / 3 * 3

	for i := 0; i < 9; i++ {
		// Controlla la corrispondente riga e colonna
		if b.Cells[row][i] == digit ||
			b.Cells[i][col] == digit ||
			// Controlla il gruppo 3x3
			b.Cells[startRow+i/3][startCol+i%3] == digit {
			return false
		}
	}

	return true
}
