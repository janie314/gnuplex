import { useEffect, useState } from "react";
import { play, pause, getOriginPos, setOriginPos } from "./lib/API"
import "./App.css";

interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

function App() {

  const [pos, setPos] = useState(0);


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
  const [media, setMedia] = useState("");
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

  const [mediafiles, setMediafiles] = useState([]);
  async function getOriginMediafiles() {
    fetch(`/api/medialist`).then((res) => res.json()).then((data) =>
      setMediafiles(data)
    );
  }

  const [last25, setLast25] = useState([]);
  async function getOriginLast25() {
    fetch(`/api/last25`).then((res) => res.json()).then((data) =>
      setLast25(data)
    );
  }

  useEffect(() => {
    getOriginPos().then((res: number) => setPos(res));
    getOriginVol();
    getOriginMedia();
    getOriginMediafiles();
    getOriginLast25();
  }, []);

  return (
    <div className="App">
      <span className="logo">GNUPlex</span>
      <div className="nowplaying">
        <span>Now playing: {media}</span>
      </div>
      <div className="controls">
      <div className="controls">
        <input type="button" value="Play" onClick={play} />
        <input type="button" value="Pause" onClick={pause} />
        </div>
      <div className="controls">
        <span>Pos</span>
        <input
          className="num-input"
          type="number"
          value={pos}
          min={0}
          onChange={(e) => {
            setPos(Number(e.target.value));
          }}
        />
        <input
          type="button"
          value="Set"
          min={0}
          max={250}
          onClick={(e) => {
            setOriginPos(pos);
          }}
        />
        </div>
      <div className="controls">
        <span>Vol</span>
        <input
          className="num-input"
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
        </div>
      <br />
      {last25.map((mediafile: string, i: number) => (
        <a
          className="mediafile"
          key={i}
          href="#"
          onClick={() => setOriginMedia(mediafile)}
        >
          {mediafile}
        </a>
      ))}
      <br />
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
