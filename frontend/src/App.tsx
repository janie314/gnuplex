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
      <div className="panel leftpanel">
        <span className="logo">GNUPlex</span>
        <div className="controlgroup">
          <input type="button" value="Play" onClick={APICall.play} />
          <input type="button" value="Pause" onClick={APICall.pause} />
        </div>
        <div className="controlgroup">
          <span>Pos</span>
          <input
            className="timeinput"
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
        <div className="controlgroup">
          <span>Vol</span>
          <input
            className="volinput"
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

      <div className="panel rightpanel">
        <div className="moviegroup">
          <span className="subtitle">Now Playing</span>
          <a className="mediafile" href="#">{media}</a>
        </div>
        <div className="moviegroup">
          <span className="subtitle">Recent</span>
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
        </div>
        <div className="moviegroup">
          <span className="subtitle">Library</span>
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
      </div>
    </div>
  );
}

export default App;
