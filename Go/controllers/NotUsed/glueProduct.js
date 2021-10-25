async function init() {
    const go = new Go();
    let result = await WebAssembly.instantiateStreaming(fetch("controllers/product.wasm"), go.importObject)
    go.run(result.instance);
}
init();