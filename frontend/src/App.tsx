import { useEffect, useRef, useState } from "react";
import { API, type MediaItem, type SubTrack } from "./lib/API";
import "./App.css";
import { CastPopup } from "./components/CastPopup";
import { MediaControls } from "./components/MediaControls";
import { MediadirsConfigPopup } from "./components/MediadirsConfigPopup";
import { Medialist } from "./components/Medialist";
import { QueuePopup } from "./components/QueuePopup";
import { useDebounce } from "./lib/useDebounce";
import { SettingsPopup } from "./components/SettingsPopup";

function App() {
  // Media player state info
  const [subs, setSubs] = useState<SubTrack[] | null>([]);
  const [pos, setPos] = useState(0);
  const [startPos, setStartPos] = useState(0);
  const [timeRemaining, setTimeRemaining] = useState(0);
  const [vol, setVol] = useState(0);
  const [nowPlaying, setNowPlaying] = useState<MediaItem[]>([]);
  const [mediaItems, setMediaItems] = useState<MediaItem[]>([]);
  const [mediaItemCount, setMediaItemCount] = useState(0);
  const [queueingTargetMediaItem, setQueueingTargetMediaItem] =
    useState<MediaItem | null>(null);
  const [paginationOffset, setPaginationOffset] = useState(
    Number(new URLSearchParams(window.location.search).get("offset") || 0) /
      1000,
  );
  const [last25, setLast25] = useState<MediaItem[]>([]);

  // UI popups' visibility
  const [mediaDirInputPopupVisible, setMediaDirInputPopupVisible] =
    useState(false);
  const [castPopupVisible, setCastPopupVisible] = useState(false);
  const [settingsPopupVisible, setSettingsPopupVisible] = useState(false);

  // URL params
  const [searchQuery, setSearchQuery] = useState(
    new URLSearchParams(window.location.search).get("search") || "",
  );
  const searchQueryDebounced = useDebounce(searchQuery, 1000);

  // Dummy audio component for media controls
  const dummyAudio = useRef<HTMLAudioElement>(null);

  // Whether or not we're on a mobile browser
  const mobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);

  useEffect(() => {
    // Escape key to quit out of windows
    window.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        setMediaDirInputPopupVisible(false);
        setCastPopupVisible(false);
        setQueueingTargetMediaItem(null);
        setSettingsPopupVisible(false);
      }
    });
    // Poll media player state from the backend
    window.setInterval(async () => {
      const res = await API.getNowPlaying();
      setNowPlaying(res || []);
      if (res === null) {
        setTimeRemaining(0);
        setVol(0);
        setSubs(null);
      } else {
        API.getTimeRemaining().then((res) => setTimeRemaining(res));
        API.getVol().then((res) => setVol(res));
        API.getPaused().then((res) => {
          if (!("mediaSession" in navigator) || !mobile) {
            return;
          }
          if (res) {
            dummyAudio.current?.pause();
            navigator.mediaSession.playbackState = "paused";
          } else {
            dummyAudio.current?.play();
            navigator.mediaSession.playbackState = "playing";
          }
        });
        API.getSubTracks().then((res) => setSubs(res));
        API.getPos().then((res) => {
          setPos(res);
          setStartPos(res);
        });
      }
    }, 2000);

    // Media Session API integration
    if (!("mediaSession" in navigator) || !mobile) {
      return;
    }
    navigator.mediaSession.playbackState = "playing";
    navigator.mediaSession.setActionHandler("play", () => {
      API.play();
      if (!/iPhone|iPad|iPod|Android/i.test(navigator.userAgent)) {
        return;
      }
      dummyAudio.current?.play();
      navigator.mediaSession.playbackState = "playing";
    });
    navigator.mediaSession.setActionHandler("pause", () => {
      API.pause();
      if (!/iPhone|iPad|iPod|Android/i.test(navigator.userAgent)) {
        return;
      }
      dummyAudio.current?.pause();
      navigator.mediaSession.playbackState = "paused";
    });
  }, []);

  useEffect(() => {
    API.getLast25Played().then((res) => setLast25(res));

    if (!("mediaSession" in navigator) || !mobile) {
      return;
    }

    if (nowPlaying?.[0]?.Path) {
      navigator.mediaSession.metadata = new MediaMetadata({
        title: nowPlaying[0].Path,
        artist: "GNUPlex",
      });
    } else {
      navigator.mediaSession.metadata = null;
      dummyAudio.current?.pause();
    }
  }, [nowPlaying]);

  // Refresh browser's search URL parameter when the search input changes
  function refreshMediaItems() {
    const urlParams = new URLSearchParams(window.location.search);
    let updateURL = false;
    if (urlParams.get("search") !== searchQueryDebounced) {
      urlParams.set("search", searchQueryDebounced);
      urlParams.set("offset", "0");
      updateURL = true;
    } else if (
      (Number(urlParams.get("offset")) || 0) / 1000 !==
      paginationOffset
    ) {
      urlParams.set("offset", paginationOffset.toString());
      updateURL = true;
    }
    if (updateURL) {
      const newURL = `${window.location.pathname}?${urlParams.toString()}`;
      window.history.pushState({}, "", newURL);
    }
    API.getMediaItems(searchQueryDebounced, paginationOffset * 1000).then(
      (res) => {
        setMediaItems(res.res);
        setMediaItemCount(res.count);
      },
    );
  }
  useEffect(() => {
    refreshMediaItems();
  }, [searchQueryDebounced, paginationOffset]);

  return (
    <>
      {mobile ? (
        <>
          {/** biome-ignore lint/a11y/useMediaCaption: just a dummy element to trigger mediacontrols */}
          <audio
            ref={dummyAudio}
            src="loop.ogg"
            autoPlay
            loop
            style={{ display: "none" }}
          />
        </>
      ) : null}
      <div
        className="flex flex-row flex-wrap max-w-full text-base font-sans pb-2/100 dark:bg-stone-950 text:white"
        style={{
          opacity:
            mediaDirInputPopupVisible || castPopupVisible ? "50%" : "100%",
        }}
      >
        <div className="sm:basis-1 md:basis-1/4 sm:max-w-full lg:max-w-sm grow flex-col px-1/100 pb-2 mb-1">
          <div className="logo-panel">
            <span className="logo">GNUPlex</span>
          </div>
          <MediaControls
            mediadirInputPopup={mediaDirInputPopupVisible}
            setMediadirInputPopup={setMediaDirInputPopupVisible}
            setCastPopup={setCastPopupVisible}
            setSettingsPopup={setSettingsPopupVisible}
            vol={vol}
            setVol={setVol}
            pos={pos}
            setPos={setPos}
            startPos={startPos}
            timeRemaining={timeRemaining}
            subs={subs}
            dummyAudio={dummyAudio}
          />
        </div>
        <div className="sm:basis-1 md:basis-3/4 min-w-sm max-w-2xl shrink flex-col p-1">
          <input
            type="text"
            className="mb-2 p-3 w-full border-2 border-stone-400 focus:bg-cyan-50 dark:bg-cyan-900 focus:dark:bg-cyan-700 dark:text-white"
            placeholder="Search"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          <Medialist
            mediaItems={nowPlaying}
            subtitle="Now Playing"
            mediaItemCount={null}
            paginationOffset={null}
            setPaginationOffset={null}
            setQueueingTargetMediaItem={setQueueingTargetMediaItem}
          />
          <Medialist
            mediaItems={last25}
            subtitle="Recent"
            mediaItemCount={null}
            paginationOffset={null}
            setPaginationOffset={null}
            setQueueingTargetMediaItem={setQueueingTargetMediaItem}
          />
          <Medialist
            mediaItems={mediaItems}
            subtitle="Library"
            mediaItemCount={mediaItemCount}
            paginationOffset={paginationOffset}
            setPaginationOffset={setPaginationOffset}
            setQueueingTargetMediaItem={setQueueingTargetMediaItem}
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
      <QueuePopup
        visible={queueingTargetMediaItem !== null}
        mediaItem={queueingTargetMediaItem}
        setQueueingTargetMediaItem={setQueueingTargetMediaItem}
        closeHook={() => {
          refreshMediaItems();
        }}
      />{" "}
      <SettingsPopup
        visible={settingsPopupVisible}
        closeHook={() => {
          setSettingsPopupVisible(false);
        }}
      />
    </>
  );
}

export { App };
