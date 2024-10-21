import { useEffect, useState } from "react";
import { API, type MediaItem } from "./lib/API";
import "./App.css";
import { CRUDPopup } from "./components/CRUDPopup";
import { MediaControls } from "./components/MediaControls";
import { Medialist } from "./components/Medialist";

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
  const [startPos, setStartPos] = useState(0);
  const [timeRemaining, setTimeRemaining] = useState(0);
  const [vol, setVol] = useState(0);
  const m: MediaItem = { ID: -1, Path: "", LastPlayed: "" };
  const [media, setMedia] = useState<MediaItem>(m);
  const [mediaItems, setMediaItems] = useState<MediaItem[]>([]);
  const [last25, setLast25] = useState<MediaItem[]>([]);
  const [mediadirInputPopup, setMediadirInputPopup] = useState(false);

  useEffect(() => {
    API.getVersion().then((version: string) => setVersion(version));
  }, []);

  useEffect(() => {
    API.getMedia().then((res) => setMedia(res));
    API.getMediaItems().then((res) => setMediaItems(res));
    API.getLast25Played().then((res) => setLast25(res));
  }, [mediaToggle]);

  useEffect(() => {
    API.getPos().then((res: number) => {
      setPos(res);
      setStartPos(res);
    });
    API.getTimeRemaining().then((res: number) => setTimeRemaining(res));
    API.getVol().then((res: number) => setVol(res));
  }, [media, volPosToggle]);

  return (
    <>
      <div
        className="flex flex-row flex-wrap max-w-full text-base font-sans pb-2/100"
        style={{ opacity: mediadirInputPopup ? "50%" : "100%" }}
      >
        <div className="sm:basis-1 md:basis-1/4 sm:max-w-full lg:max-w-sm grow flex-col px-1/100 pb-2 mb-1">
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
            startPos={startPos}
            timeRemaining={timeRemaining}
            setPos={setPos}
            volPosToggle={volPosToggle}
            setVolPosToggle={setVolPosToggle}
          />
        </div>

        <div className="sm:basis-1 md:basis-3/4 min-w-sm max-w-2xl shrink flex-col p-1/100">
          <Medialist
            mediaItems={[media]}
            subtitle="Now Playing"
            setMedia={setMedia}
          />
          <Medialist
            mediaItems={last25}
            subtitle="Recent"
            setMedia={setMedia}
          />
          <Medialist
            mediaItems={mediaItems}
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
