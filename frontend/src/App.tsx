import { useEffect, useState } from "react";
import { APICall } from "./lib/APICall.ts";
import "./App.css";
import { Medialist } from "./components/Medialist.tsx";
import { TimeVolInput } from "./components/TimeVolInput.tsx";
import { CRUDPopup } from "./components/CRUDPopup.tsx";

interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

function App() {
  const [version, setVersion] = useState("");
  const [volPosToggle, setVolPosToggle] = useState(false);
  const [mediaToggle, setMediaToggle] = useState(false);
  const [pos, setPos] = useState(0);
  const [vol, setVol] = useState(0);
  const [media, setMedia] = useState("");
  const [mediafiles, setMediafiles] = useState<string[]>([]);
  const [last25, setLast25] = useState<string[]>([]);
  const [mediadirInputPopup, setMediadirInputPopup] = useState(false);

  useEffect(() => {
    APICall.getOriginVersion().then((version: string) => setVersion(version));
  }, []);

  useEffect(() => {
    APICall.getOriginMedia().then((res: string) => setMedia(res));
    APICall.getOriginMediafiles().then((res: string[]) => setMediafiles(res));
    APICall.getOriginLast25().then((res: string[]) => setLast25(res));
  }, [mediaToggle]);

  useEffect(() => {
    APICall.getOriginPos().then((res: number) => setPos(res));
    APICall.getOriginVol().then((res: number) => setVol(res));
  }, [media, volPosToggle]);

  return (
    <>
      <div
        className="App"
        style={{ opacity: mediadirInputPopup ? "50%" : "100%" }}
      >
        <div className="panel leftpanel">
          <div className="logo-panel">
            <span className="logo">GNUPlex</span>
            <span className="version">{version}</span>
          </div>
          <div className="controlgroup">
            <input
              className="play-button"
              type="button"
              value="⏵"
              onClick={() =>
                APICall.play().then(() => APICall.sleep(2000)).then(() =>
                  setVolPosToggle(!volPosToggle)
                )}
            />
            <input
              className="pause-button"
              type="button"
              value="⏸"
              onClick={() =>
                APICall.pause().then(() => APICall.sleep(2000)).then(() =>
                  setVolPosToggle(!volPosToggle)
                )}
            />
          </div>
          <div className="controlgroup">
            <input
              type="button"
              value="Manage Library"
              onClick={() => {
                setMediadirInputPopup(true);
              }}
            />
            <input
              type="button"
              value="Cast YouTube"
              onClick={() => {
                const url = window.prompt("YouTube URL:", "") || "";
                APICall.setOriginMedia(url);
              }}
            />
          </div>
          <div className="controlgroup">
            <TimeVolInput rawtime={pos} setRawtime={setPos} type="time" />
          </div>
          <div className="controlgroup">
            <TimeVolInput vol={vol} setVol={setVol} type="vol" />
          </div>
          <div className="controlgroup">
            <a
              href="https://gitlab.com/jane314/gnuplex/-/issues"
              target="_blank"
            >
              Bug?
            </a>
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
      <CRUDPopup
        visible={mediadirInputPopup}
        setMediadirInputPopup={setMediadirInputPopup}
        closeHook={() => {
          setMediaToggle(!mediaToggle);
        }}
      />
    </>
  );
}

export { App };
