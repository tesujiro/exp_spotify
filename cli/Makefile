all:
	go vet .
	go build .

search:
	./cli search artist つのだ
	./cli search album Secret
	./cli search playlist harakami

profile:
	./cli get profile

playlist:
	./cli get playlists

test:
	./test.sh ; echo RESULT=$$?
