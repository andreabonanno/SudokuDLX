package main

import (
	"SudokuDLX/dlx"
	"fmt"
	"time"
)

func main() {

	sudoku := "___26_7_168__7__9_19___45__82_1___4___46_29___5___3_28__93___74_4__5__367_3_18___"
	start := time.Now()
	m := encodeConstraints(sudoku)
	m.Solve(0)
	res := decodeExactCoverSolution(m.GetExactCover())
	elapsed := time.Since(start)
	fmt.Printf("Tempo esecuzione DLX: %s\n", elapsed)
	printSudoku(res)
}

func printSudoku(s string) {
	for r, i := 0, 0; r < 9; r, i = r+1, i+9 {
		fmt.Printf("%c %c %c | %c %c %c | %c %c %c\n", s[i], s[i+1], s[i+2],
			s[i+3], s[i+4], s[i+5], s[i+6], s[i+7], s[i+8])
		if r == 2 || r == 5 {
			fmt.Printf("\n")
		}
	}
}

// Inizializza il problema di copertura
// Prende come argomento un problema sudoku in formato stringa di 81 caratteri
func encodeConstraints(s string) dlx.Matrix {
	m := dlx.NewMatrix(324)

	for row, position := 0, 0; row < 9; row++ {
		for column := 0; column < 9; column, position = column+1, position+1 {
			region := row/3*3 + column/3
			digit := int(s[position] - '1') // zero based digit
			if digit >= 0 && digit < 9 {
				m.AddRow(position, 81+row*9+digit, 162+column*9+digit, 243+region*9+digit)
			} else {
				for digit = 0; digit < 9; digit++ {
					m.AddRow(position, 81+row*9+digit, 162+column*9+digit, 243+region*9+digit)
				}
			}
		}
	}

	return m
}

// Effettua il parsing da soluzione di problema di copertura a soluzione del sudoku
func decodeExactCoverSolution(cs [][]int) string {
	b := make([]byte, len(cs))
	for _, row := range cs {
		position := row[0]
		value := row[1] % 9
		b[position] = byte(value) + '1'
	}
	return string(b)
}
