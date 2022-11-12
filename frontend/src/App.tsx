import { useState } from "react";
import "./App.css";
import { mediafiles } from "./mediafiles";
import serverconfig from "../serverconfig.json";

function App() {
  function play() {
    return fetch(
      `http://${serverconfig.hostname}/api/play`,
      { method: "POST" },
    );
  }

  function pause() {
    return fetch(
      `http://${serverconfig.hostname}/api/pause`,
      { method: "POST" },
    );
  }

  function queue(mediafile: string) {
    return fetch(
      `http://${serverconfig.hostname}/api/queue?mediafile=${
        encodeURIComponent(mediafile)
      }`,
      { method: "POST" },
    );
  }
  
  async function getPos() {
    fetch(
      `http://${serverconfig.hostname}/api/getpos`
    ).then( res => res.json()).then( (res: Object) => {

    });
  }



  const [pos,setPos] = useState(0);

  fetch(``)

  return (
    <div className="App">
      <span className="title">GNUPlex</span>
      <div className="controls">
        <input type="button" value="Play" onClick={play} />
        <input type="button" value="Pause" onClick={pause} />
      </div>
      <div className="controls">
        <span>Pos</span>
        <input type="number" value="Play" onClick={play} />
      </div>
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
