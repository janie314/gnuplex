import { useEffect, useState } from "react";
import "./App.css";
import { mediafiles } from "./mediafiles";
import serverconfig from "../serverconfig.json";

interface IMPVRes {
  data?: number;
  request_id: number;
  error: string;
}

function App() {
  function play() {
    return fetch(
      `/api/play`,
      { method: "POST" },
    );
  }

  function pause() {
    return fetch(
      `/api/pause`,
      { method: "POST" },
    );
  }

  function queue(mediafile: string) {
    return fetch(
      `/api/queue?mediafile=${encodeURIComponent(mediafile)}`,
      { method: "POST" },
    );
  }

  /*
   * pos state
   */
  const [pos, setPos] = useState(0);
  async function getOriginPos() {
    fetch(
      `/api/pos`,
    ).then((res) => res.json()).then((res: IMPVRes) => {
      if (res.data !== undefined) {
        setPos(res.data);
      }
    });
  }
  async function setOriginPos(pos: number) {
    return await fetch(
      `http://${serverconfig.hostname}/api/pos?pos=${pos}`,
      { method: "POST" },
    ).then((res) => res.json());
  }

  /*
   * vol state
   */
  const [vol, setVol] = useState(0);
  async function getOriginVol() {
    fetch(
      `/api/vol`,
    ).then((res) => res.json()).then((res: IMPVRes) => {
      if (res.data !== undefined) {
        setVol(res.data);
      }
    });
  }
  async function setOriginVol(vol: number) {
    return await fetch(
      `/api/vol?vol=${vol}`,
      { method: "POST" },
    ).then((res) => res.json());
  }

  return (
    <div className="App">
      <span className="title">GNUPlex</span>
      <div className="controls">
        <input type="button" value="Play" onClick={play} />
        <input type="button" value="Pause" onClick={pause} />
      </div>
      <div className="controls">
        <span>Pos</span>
        <input
          type="number"
          value={pos}
          onChange={(e) => {
            setPos(Number(e.target.value));
          }}
        />
        <input
          type="button"
          value="Set"
          onClick={(e) => {
            setOriginPos(pos);
          }}
        />
      </div>
      <div className="controls">
        <span>Vol</span>
        <input
          type="number"
          value={vol}
          onChange={(e) => {
            setVol(Number(e.target.value));
          }}
        />
        <input
          type="button"
          value="Set"
          onClick={(e) => {
            setOriginVol(vol);
          }}
        />
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
