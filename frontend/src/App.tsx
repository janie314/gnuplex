import { useEffect, useState } from "react";
import { APICall } from "./lib/APICall";
import "./App.css";
import { Medialist } from "./components/Medialist";
import { TimeVolInput } from "./components/TimeVolInput";
import { CRUDPopup } from "./components/CRUDPopup";
import { MediaControls } from "./components/MediaControls";

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
        className="flex flex-row text-base font-sans pb-2/100"
        style={{ opacity: mediadirInputPopup ? "50%" : "100%" }}
      >
        <div className="flex-none flex-col w-1/4 p-1/100">
          <div className="logo-panel">
            <span className="logo">GNUPlex</span>
            <span className="version">{version}</span>
          </div>
          <MediaControls
            mediadirInputPopup={mediadirInputPopup}
            setMediadirInputPopup={setMediadirInputPopup}
            vol={vol}
            setVol={setVol}
            pos={pos}
            setPos={setPos}
            volPosToggle={volPosToggle}
            setVolPosToggle={setVolPosToggle}
          />
        </div>

        <div className="flex-col flex-none w-3/4 p-1/100">
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
