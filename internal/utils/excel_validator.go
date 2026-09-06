package utils

import (
	"strings"
)

// ValidateSoalRow memvalidasi satu row dari excel
func ValidateSoalRow(row *ExcelSoalRow) []string {
	var errors []string

	kunci := strings.TrimSpace(row.Kunci)
	if kunci == "" {
		errors = append(errors, "kunci tidak boleh kosong")
	} else if !IsValidKunci(kunci) {
		errors = append(errors, "kunci harus berupa A, B, C, D, atau E")
	}

	return errors
}

// IsValidKunci checks if kunci is valid (A-E)
func IsValidKunci(k string) bool {
	k = strings.ToUpper(strings.TrimSpace(k))
	return k == "A" || k == "B" || k == "C" || k == "D" || k == "E"
}

// ValidatePesertaRow memvalidasi satu row dari excel import peserta
func ValidatePesertaRow(row *ExcelPesertaRow) []string {
	var errors []string

	if strings.TrimSpace(row.Nama) == "" {
		errors = append(errors, "nama tidak boleh kosong")
	}
	if strings.TrimSpace(row.Username) == "" {
		errors = append(errors, "username tidak boleh kosong")
	}
	if strings.TrimSpace(row.Password) == "" {
		errors = append(errors, "password tidak boleh kosong")
	} else if len(row.Password) < 6 {
		errors = append(errors, "password minimal 6 karakter")
	}

	return errors
}
