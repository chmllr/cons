all:
	go run main.go

clean-testlib:
	rm -rf testdata/lib/*

restore-import:
	cp -rf testdata/import-original/* testdata/import/

restore:
	make clean-testlib
	make restore-import

show:
	tree testdata/lib
