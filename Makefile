dev:
	export FIRESTORE_EMULATOR_HOST="localhost:8082"
	go build
	./legalaidcamp-go dev
prod:
	export GOOGLE_APPLICATION_CREDENTIALS="/Users/rohanbojja/Creds/legalaidcamp.json"
	go build
	./legalaidcamp-go prod