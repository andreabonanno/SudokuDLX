package dlx

// Matrix e una matrice di valori binari che rappresenta il problema di copertura.
type Matrix interface {
	// AddRow aggiunge una riga di vincoli alla matrice.
	// AddRow prende come argomento gli indici colonna degli elementi non-zero della riga da aggiungere.
	AddRow(constraintsRow ...int)

	Solve(depth int) bool

	GetExactCover() [][]int
}

// Linked list toroidale per l'esecuzione dei dancing links
type matrixObject struct {
	columns         []columnObject
	head            *columnObject
	partialSolution []*dataObject
}

type dataObject struct {
	column                *columnObject
	up, down, left, right *dataObject
	rowStart              *dataObject
}

type columnObject struct {
	dataObject
	size  int
	index int
}

// NewMatrix crea una nuova Matrix avente nColumns colonne
func NewMatrix(nColumns int) Matrix {
	columns := make([]columnObject, nColumns+1)

	// inizializzo head puntatore all'intera struttura
	head := &columns[0]
	headDObj := &head.dataObject
	head.column = head
	head.left = &columns[nColumns].dataObject
	head.up = headDObj
	head.down = headDObj
	head.index = -1

	// l'ltima colonna si "riavvolge" su head
	columns[nColumns].right = headDObj
	prevColumn := head

	for i := range columns[1:] {
		column := &columns[i+1]
		columnDObj := &column.dataObject
		column.index = i
		column.column = column
		column.up = columnDObj
		column.down = columnDObj
		column.left = &prevColumn.dataObject
		prevColumn.right = columnDObj
		prevColumn = column
	}
	return &matrixObject{columns, head, nil}
}

// AddRow aggiunge una riga di vincoli al problema
func (m *matrixObject) AddRow(constraintsRow ...int) {
	constraintsCount := len(constraintsRow)
	if constraintsCount == 0 {
		return
	}
	rowDOs := make([]dataObject, constraintsCount)
	rowStart := &rowDOs[0]
	for i, c := range constraintsRow {
		column := &m.columns[c+1]
		column.size++

		// ottengo i puntatori dei nodi vicini
		prevUpDObj := column.up
		downDObj := &column.dataObject
		leftDObj := &rowDOs[(i+constraintsCount-1)%constraintsCount]
		rightDObj := &rowDOs[(i+1)%constraintsCount]

		// nuovo dataObject
		do := &rowDOs[i]
		do.column = column
		do.up = prevUpDObj
		do.down = downDObj
		do.left = leftDObj
		do.right = rightDObj
		do.rowStart = rowStart

		// aggiorno i puntatori dei nodi adiacenti
		prevUpDObj.down, downDObj.up, leftDObj.right, rightDObj.left = do, do, do, do
	}
}

func (m *matrixObject) Solve(depth int) bool {
	head := m.head
	headRight := head.right.column
	if headRight == head {
		return true
	}

	c := headRight
	minSize := headRight.size

	for jc := headRight.right.column; jc != head; jc = jc.right.column {
		jSize := jc.size
		if jSize >= minSize {
			continue
		}
		c, minSize = jc, jSize
	}

	coverColumn(c)

	stackSize := len(m.partialSolution)
	m.partialSolution = append(m.partialSolution, nil)

	for r := c.down; r != &c.dataObject; r = r.down {
		m.partialSolution[stackSize] = r

		for j := r.right; j != r; j = j.right {
			coverColumn(j.column)
		}

		if m.Solve(depth + 1) {
			return true
		}

		for j := r.left; j != r; j = j.left {
			unCoverColumn(j.column)
		}
	}
	m.partialSolution = m.partialSolution[:stackSize]
	unCoverColumn(c)

	return false
}

func (m *matrixObject) GetExactCover() [][]int {
	ec := make([][]int, len(m.partialSolution))

	for i, do := range m.partialSolution {
		row := ec[i]
		rowStart := do.rowStart
		row = append(row, rowStart.column.index)
		for j := rowStart.right; j != rowStart; j = j.right {
			row = append(row, j.column.index)
		}
		ec[i] = row
	}

	return ec
}

func coverColumn(c *columnObject) {
	cDObj := &c.dataObject
	c.right.left, c.left.right = c.left, c.right
	for i := c.down; i != cDObj; i = i.down {
		for j := i.right; j != i; j = j.right {
			j.down.up, j.up.down = j.up, j.down
			j.column.size--
		}
	}
}

func unCoverColumn(c *columnObject) {
	cDObj := &c.dataObject
	for i := c.up; i != cDObj; i = i.up {
		for j := i.left; j != i; j = j.left {
			j.column.size++
			j.down.up, j.up.down = j, j
		}
	}
	c.right.left, c.left.right = cDObj, cDObj
}
