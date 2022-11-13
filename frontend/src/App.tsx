import { useEffect, useState } from "react";
import "./App.css";
import { mediafiles } from "./mediafiles";

interface IMPVRes {
  data?: number | string;
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

  /*
   * pos state
   */
  const [pos, setPos] = useState(0);
  async function getOriginPos() {
    fetch(
      `/api/pos`,
    ).then((res) => res.json()).then((res: IMPVRes) => {
      if (res.data !== undefined) {
        // @ts-ignore
        setPos(Math.floor(res.data));
      }
    });
  }
  async function setOriginPos(pos: number) {
    return await fetch(
      `/api/pos?pos=${pos}`,
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
        // @ts-ignore
        setVol(Math.floor(res.data));
      }
    });
  }
  async function setOriginVol(vol: number) {
    return await fetch(
      `/api/vol?vol=${vol}`,
      { method: "POST" },
    ).then((res) => res.json());
  }

  /*
   * media state
   */
  const [media, setMedia] = useState('');
  async function getOriginMedia() {
    fetch(
      `/api/media`,
    ).then((res) => res.json()).then((res: IMPVRes) => {
      if (res.data !== undefined) {
        // @ts-ignore
        setMedia(res.data);
      }
    });
  }
  async function setOriginMedia(mediafile: string) {
    return await fetch(
      `/api/media?mediafile=${mediafile}`,
      { method: "POST" },
    ).then((res) => res.json());
  }
  useEffect(() => {
    getOriginPos();
    getOriginVol();
    getOriginMedia();
  }, []);

  return (
    <div className="App">
      <span className="title">GNUPlex</span>
      <div className="controls">
        <span>Now playing: {media}</span>
      </div>
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
          onClick={() => setOriginMedia(mediafile)}
        >
          {mediafile}
        </a>
      ))}
    </div>
  );
}

export default App;
