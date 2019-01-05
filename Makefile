all:
	go run main.go

clean-testlib:
	rm -rf testlib/*

restore-import:
	cp -rf import-original/* import/
