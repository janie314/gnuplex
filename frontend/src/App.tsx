import { useEffect, useState } from "react";
import { API, type MediaItem } from "./lib/API";
import "./App.css";
import { useDebounce } from "@uidotdev/usehooks";
import { CastPopup } from "./components/CastPopup";
import { MediaControls } from "./components/MediaControls";
import { MediadirsConfigPopup } from "./components/MediadirsConfigPopup";
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
  const [nowPlaying, setNowPlaying] = useState("");
  const [mediaItems, setMediaItems] = useState<MediaItem[]>([]);
  const [last25, setLast25] = useState<MediaItem[]>([]);
  const [mediadirInputPopup, setMediadirInputPopup] = useState(false);
  const [castPopup, setCastPopup] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const searchQueryDebounced = useDebounce(searchQuery, 1000);

  useEffect(() => {
    API.getVersion().then((version: string) => setVersion(version));
    const urlParams = new URLSearchParams(window.location.search);
    if ((urlParams.get("search") || "").length > 0) {
      setSearchQuery(urlParams.get("search") || "");
    }
  }, []);

  useEffect(() => {
    API.getNowPlaying().then((res) => setNowPlaying(res));
    API.getLast25Played().then((res) => setLast25(res));
  }, [mediaToggle]);

  useEffect(() => {
    API.getMediaItems(searchQueryDebounced).then((res) => setMediaItems(res));
  }, [searchQueryDebounced]);

  useEffect(() => {
    API.getPos().then((res: number) => {
      setPos(res);
      setStartPos(res);
    });
    API.getTimeRemaining().then((res: number) => setTimeRemaining(res));
    API.getVol().then((res: number) => setVol(res));
  }, [nowPlaying, volPosToggle]);

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    if (
      urlParams.get("search") !== searchQueryDebounced &&
      searchQueryDebounced.length !== 0
    ) {
      urlParams.set("search", searchQueryDebounced);
      window.location.search = urlParams.toString();
    }
  }, [searchQueryDebounced]);

  return (
    <>
      <div
        className="flex flex-row flex-wrap max-w-full text-base font-sans pb-2/100"
        style={{ opacity: mediadirInputPopup || castPopup ? "50%" : "100%" }}
      >
        <div className="sm:basis-1 md:basis-1/4 sm:max-w-full lg:max-w-sm grow flex-col px-1/100 pb-2 mb-1">
          <div className="logo-panel">
            <span className="logo">GNUPlex</span>
            <span className="version">{version}</span>
          </div>
          <MediaControls
            mediadirInputPopup={mediadirInputPopup}
            setMediadirInputPopup={setMediadirInputPopup}
            setCastPopup={setCastPopup}
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
        <div className="sm:basis-1 md:basis-3/4 min-w-sm max-w-2xl shrink flex-col p-1">
          <input
            type="text"
            className="grow mb-2 p-3 w-full border-2 border-gray-300 focus:bg-cyan-50"
            placeholder="Search"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          <Medialist mediaItems={[]} subtitle="Now Playing" />
          <Medialist mediaItems={last25} subtitle="Recent" />
          <Medialist mediaItems={mediaItems} subtitle="Library" />
        </div>
      </div>
      <MediadirsConfigPopup
        visible={mediadirInputPopup}
        setMediadirInputPopup={setMediadirInputPopup}
        closeHook={() => {
          setMediaToggle(!mediaToggle);
        }}
      />
      <CastPopup
        visible={castPopup}
        setCastPopup={setCastPopup}
        closeHook={() => {
          setMediaToggle(!mediaToggle);
        }}
      />
    </>
  );
}

export { App };
