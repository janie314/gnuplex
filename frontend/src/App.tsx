import { useEffect, useState } from "react";
import { APICall } from "./lib/API";
import "./App.css";
import { Medialist } from "./components/Medialist";
import { TimeVolInput } from "./components/TimeVolInput";

interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

function App() {
  const [toggle, setToggle] = useState(false);
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
  }, [media, toggle]);

  return (
    <div className="App">
      <div className="panel leftpanel">
        <span className="logo">GNUPlex</span>
        <div className="controlgroup">
          <input
            className="play-button"
            type="button"
            value="⏵"
            onClick={() =>
              APICall.play().then(() => APICall.sleep(2000)).then(() =>
                setToggle(!toggle)
              )}
          />
          <input
            className="pause-button"
            type="button"
            value="⏸"
            onClick={() =>
              APICall.pause().then(() => APICall.sleep(2000)).then(() =>
                setToggle(!toggle)
              )}
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
          <TimeVolInput rawtime={pos} setRawtime={setPos} type="time" />
        </div>
        <div className="controlgroup">
          <TimeVolInput vol={vol} setVol={setVol} type="vol" />
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
