import { useEffect, useState } from "react";
import { APICall } from "./lib/API";
import "./App.css";

interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

function App() {
  const [pos, setPos] = useState(0);
  const [vol, setVol] = useState(0);
  const [media, setMedia] = useState("");
  const [mediafiles, setMediafiles] = useState<string[]>([]);
  const [last25, setLast25] = useState<string[]>([]);

  useEffect(() => {
    APICall.getOriginPos().then((res: number) => setPos(res));
    APICall.getOriginVol().then((res: number) => setVol(res));
    APICall.getOriginMedia().then((res: string) => setMedia(res));
    APICall.getOriginMediafiles().then((res: string[]) => setMediafiles(res));
    APICall.getOriginLast25().then((res: string[]) => setLast25(res));
  }, []);

  return (
    <div className="App">
      <span className="logo">GNUPlex</span>
      <div className="nowplaying">
        <span>Now playing: {media}</span>
      </div>
      <div className="controls">
        <div className="controls">
          <input type="button" value="Play" onClick={APICall.play} />
          <input type="button" value="Pause" onClick={APICall.pause} />
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
              APICall.setOriginPos(pos);
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
              APICall.setOriginVol(vol);
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
          onClick={() => APICall.setOriginMedia(mediafile)}
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
          onClick={() => APICall.setOriginMedia(mediafile)}
        >
          {mediafile}
        </a>
      ))}
    </div>
  );
}

export default App;
