package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomainFast(b *testing.B) {
	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()

	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
}

func BenchmarkGetDomainSlow(b *testing.B) {
	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()

	for i := 0; i < b.N; i++ {
		getDomainStat(data, "biz")
	}
}
