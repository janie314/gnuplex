import { useState } from "react";
import "./App.css";
import { mediafiles } from "./mediafiles";

function App() {
  
  function queue(mediafile: string) {
    return fetch(
      `http://localhost:40000/api/queue?mediafile=${encodeURIComponent(mediafile)}`,
      { method: "POST" },
    );
  }

  return (
    <div className="App">
      <span className="title">GNUPlex</span>
      {mediafiles.map((mediafile: string, i: number) => (
        <a
          className="mediafile"
          key={i}
          href="#"
          onClick={() => queue(mediafile)}
        >
          {mediafile}
        </a>
      ))}
    </div>
  );
}

export default App;
