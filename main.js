import "./style.css";
import javascriptLogo from "./javascript.svg";
import viteLogo from "/vite.svg";

import * as zip from "@zip.js/zip.js";
import { encode } from "@msgpack/msgpack";

const fileselector = document.getElementById("file");

fileselector.addEventListener("change", async function (event) {
  var file = event.target.files[0];
  if (!file) {
    return;
  }

  // remove the file selector
  fileselector.remove();

  var processingNow = "Opening zip file...";

  const reader = new zip.ZipReader(new zip.BlobReader(file));

  const entries = await reader.getEntries();
  if (!entries.length) {
    return;
  }

  window.files = {};

  for (const entry of entries) {
    processingNow = entry.filename;
    const text = await entry.getData(new zip.TextWriter());
    window.files[entry.filename] = text;
  }

  // run the WASM
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch("main.wasm"),
    go.importObject
  );

  let t1 = performance.now();
  window.files = encode(window.files);
  let t2 = performance.now();
  console.log("Call to stringify took " + (t2 - t1) + " milliseconds.");
  go.run(result.instance);
});
