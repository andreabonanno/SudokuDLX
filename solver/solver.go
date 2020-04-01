package main

import (
	"../backtrack"
	"../dlx"
	"fmt"
	"time"
)

func main() {

	sudoku := [9][9]int{
		{0, 6, 1, 0, 0, 7, 0, 0, 3},
		{0, 9, 2, 0, 0, 3, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 5, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 4},
		{5, 0, 0, 0, 0, 8, 0, 0, 0},
		{0, 4, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 1, 6, 0, 8, 0, 0},
		{6, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	start := time.Now()
	m := encodeConstraints(sudoku)
	m.Solve(0)
	res := decodeExactCoverSolution(m.GetExactCover())
	elapsed := time.Since(start)
	printSudoku(res)
	fmt.Printf("Tempo esecuzione DLX: %s\n", elapsed)
	print("\n")

	start = time.Now()
	board := backtrack.Board{Cells: sudoku}
	board.Backtrack()
	elapsed = time.Since(start)
	print(board.String())
	fmt.Printf("Tempo esecuzione Backtrack: %s\n", elapsed)
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
func encodeConstraints(s [9][9]int) dlx.Matrix {
	m := dlx.NewMatrix(324)

	for row, position := 0, 0; row < 9; row++ {
		for column := 0; column < 9; column, position = column+1, position+1 {
			region := row/3*3 + column/3
			digit := s[row][column] - 1
			if digit > 0 && digit < 9 {
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
