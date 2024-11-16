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
  // App info
  const [version, setVersion] = useState("");
  // Media player state info
  const [pos, setPos] = useState(0);
  const [startPos, setStartPos] = useState(0);
  const [timeRemaining, setTimeRemaining] = useState(0);
  const [vol, setVol] = useState(0);
  const [nowPlaying, setNowPlaying] = useState<MediaItem | null>(null);
  const [mediaItems, setMediaItems] = useState<MediaItem[]>([]);
  const [mediaItemCount, setMediaItemCount] = useState(0);
  const [paginationOffset, setPaginationOffset] = useState(
    Number(new URLSearchParams(window.location.search).get("offset") || 0) /
      1000,
  );
  const [last25, setLast25] = useState<MediaItem[]>([]);
  // UI popups' visibility
  const [mediaDirInputPopupVisible, setMediaDirInputPopupVisible] =
    useState(false);
  const [castPopupVisible, setCastPopupVisible] = useState(false);
  // URL params
  const [searchQuery, setSearchQuery] = useState(
    new URLSearchParams(window.location.search).get("search") || "",
  );
  const searchQueryDebounced = useDebounce(searchQuery, 1000);

  useEffect(() => {
    // Populate app version
    API.getVersion().then((version) => setVersion(version));
    // Poll media player state from the backend
    window.setInterval(() => {
      API.getPos().then((res) => {
        setPos(res);
        setStartPos(res);
      });
      API.getTimeRemaining().then((res) => setTimeRemaining(res));
      API.getVol().then((res) => setVol(res));
      API.getNowPlaying().then((res) => setNowPlaying(res));
    }, 2000);
  }, []);

  useEffect(() => {
    API.getLast25Played().then((res) => setLast25(res));
  }, [nowPlaying]);

  // Refresh browser's search URL parameter when the search input changes
  function refreshMediaItems() {
    const urlParams = new URLSearchParams(window.location.search);
    if (
      urlParams.get("search") !== searchQueryDebounced ||
      (Number(urlParams.get("offset")) || 0) / 1000 !== paginationOffset
    ) {
      urlParams.set("search", searchQueryDebounced);
      urlParams.set("offset", paginationOffset.toString());
      window.location.search = urlParams.toString();
    }
    API.getMediaItems(searchQueryDebounced, paginationOffset).then((res) => {
      setMediaItems(res.res);
      setMediaItemCount(res.count);
    });
  }

  useEffect(() => {
    refreshMediaItems();
  }, [searchQueryDebounced, paginationOffset]);

  return (
    <>
      <div
        className="flex flex-row flex-wrap max-w-full text-base font-sans pb-2/100"
        style={{
          opacity:
            mediaDirInputPopupVisible || castPopupVisible ? "50%" : "100%",
        }}
      >
        <div className="sm:basis-1 md:basis-1/4 sm:max-w-full lg:max-w-sm grow flex-col px-1/100 pb-2 mb-1">
          <div className="logo-panel">
            <span className="logo">GNUPlex</span>
            <span className="version">{version}</span>
          </div>
          <MediaControls
            mediadirInputPopup={mediaDirInputPopupVisible}
            setMediadirInputPopup={setMediaDirInputPopupVisible}
            setCastPopup={setCastPopupVisible}
            vol={vol}
            setVol={setVol}
            pos={pos}
            setPos={setPos}
            startPos={startPos}
            timeRemaining={timeRemaining}
          />
        </div>
        <div className="sm:basis-1 md:basis-3/4 min-w-sm max-w-2xl shrink flex-col p-1">
          <input
            type="text"
            className="mb-2 p-3 w-full border-2 border-gray-300 focus:bg-cyan-50"
            placeholder="Search"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          <Medialist
            mediaItems={[nowPlaying]}
            subtitle="Now Playing"
            mediaItemCount={null}
            paginationOffset={null}
            setPaginationOffset={null}
          />
          <Medialist
            mediaItems={last25}
            subtitle="Recent"
            mediaItemCount={null}
            paginationOffset={null}
            setPaginationOffset={null}
          />
          <Medialist
            mediaItems={mediaItems}
            subtitle="Library"
            mediaItemCount={mediaItemCount}
            paginationOffset={paginationOffset}
            setPaginationOffset={setPaginationOffset}
          />
        </div>
      </div>
      <MediadirsConfigPopup
        visible={mediaDirInputPopupVisible}
        setMediadirInputPopup={setMediaDirInputPopupVisible}
        closeHook={refreshMediaItems}
      />
      <CastPopup
        visible={castPopupVisible}
        setCastPopup={setCastPopupVisible}
        closeHook={refreshMediaItems}
      />
    </>
  );
}

export { App };
