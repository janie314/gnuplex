import { useEffect, useState } from "react";
import { APICall } from "./lib/API";
import "./App.css";
import { Medialist } from "./components/Medialist";
import { TimeInput } from "./components/TimeInput";

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
    APICall.getOriginMedia().then((res: string) => setMedia(res));
    APICall.getOriginMediafiles().then((res: string[]) => setMediafiles(res));
    APICall.getOriginLast25().then((res: string[]) => setLast25(res));
  }, []);

  useEffect(() => {
    APICall.getOriginPos().then((res: number) => setPos(res));
    APICall.getOriginVol().then((res: number) => setVol(res));
  }, [media]);

  return (
    <div className="App">
      <div className="panel leftpanel">
        <span className="logo">GNUPlex</span>
        <div className="controlgroup">
          <input
            className="play-button"
            type="button"
            value="⏵"
            onClick={APICall.play}
          />
          <input
            className="pause-button"
            type="button"
            value="⏸"
            onClick={APICall.pause}
          />
        </div>
        <div className="controlgroup">
          <input
            type="button"
            value="Refresh Library"
            onClick={APICall.refreshOriginMediafiles}
          />
        </div>
        <div className="controlgroup">
          <TimeInput rawtime={pos} setRawtime={setPos} />
        </div>
        <div className="controlgroup">
          <span id="vol-label">Vol</span>
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
            id="vol-set"
            value="Set"
            onClick={(e) => {
              APICall.setOriginVol(vol);
            }}
          />
        </div>
      </div>

      <div className="panel rightpanel">
        <Medialist
          medialist={[media]}
          subtitle="Now Playing"
          setMedia={setMedia}
        />
        <Medialist medialist={last25} subtitle="Recent" setMedia={setMedia} />
        <Medialist
          medialist={mediafiles}
          subtitle="Library"
          setMedia={setMedia}
        />
      </div>
    </div>
  );
}

export { App };
