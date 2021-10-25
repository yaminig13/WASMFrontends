makeComponent = function(name, customEl) {
            const Element = class extends HTMLElement {
            constructor() {
                super();
                customEl(this);
            }
        };
            window.customElements.define(name, Element);
        
        };
async function init() {
    const go = new Go();
    let result = await WebAssembly.instantiateStreaming(fetch("controllers/app.wasm"), go.importObject)
    go.run(result.instance);
}
init();

