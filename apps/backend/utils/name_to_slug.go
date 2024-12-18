package utils

import (
	"regexp"
	"strings"
)

func Slugify(name string) string {
	// Ubah menjadi huruf kecil
	slug := strings.ToLower(name)

	// Hapus karakter non-alfanumerik dan ganti dengan tanda "-"
	re := regexp.MustCompile("[^a-z0-9]+")
	slug = re.ReplaceAllString(slug, "-")

	// Hapus tanda "-" berlebih di awal atau akhir
	slug = strings.Trim(slug, "-")

	return slug
}
