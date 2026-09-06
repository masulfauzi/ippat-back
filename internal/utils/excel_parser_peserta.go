package utils

type ExcelPesertaRow struct {
	RowIndex int
	Nama     string
	Username string
	Password string
}

// ParseExcelPesertaRow mengextract data peserta dari satu baris excel.
// Kolom: A=Nama, B=Username, C=Password (kolom "No" di file template diabaikan).
func ParseExcelPesertaRow(values []string, rowIndex int) *ExcelPesertaRow {
	getValueAt := func(idx int) string {
		if idx < len(values) {
			return toStringFromValue(values[idx])
		}
		return ""
	}

	return &ExcelPesertaRow{
		RowIndex: rowIndex,
		Nama:     getValueAt(1),
		Username: getValueAt(2),
		Password: getValueAt(3),
	}
}
