build: compile
	wasm-tools component embed ./wit manager.module.wasm --output manager.embed.wasm
	wasm-tools component new manager.embed.wasm -o manager.wasm --adapt wit/adapters/tier2/wasi_snapshot_preview1.wasm

bindings:
	wit-bindgen tiny-go --out-dir manager_out ./wit

compile: bindings
	tinygo build -target=wasi -o manager.module.wasm main.go

clean:
	rm -rf manager_out
	rm *.wasm

del-worker:
	golem-cli worker delete -p golem-poll -t manager --worker-name mgr-test-2

upload:
	golem-cli template update -p golem-poll -t manager manager.wasm
	golem-cli worker add -p golem-poll -t manager --worker-name mgr-test-2