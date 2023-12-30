import { useEffect, useState } from "react";
import { APICall } from "./lib/API.ts";
import "./App.css";
import { Medialist } from "./components/Medialist.tsx";
import { MediaControls } from "./components/MediaControls.tsx";
import { TimeVolInput } from "./components/TimeVolInput.tsx";
import { CRUDPopup } from "./components/CRUDPopup.tsx";

interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

function App() {
  const [version, setVersion] = useState("");
  const [mediaToggle, setMediaToggle] = useState(false);
  const [paused, setPaused] = useState(true);
  const [media, setMedia] = useState("");
  const [mediafiles, setMediafiles] = useState<string[]>([]);
  const [last25, setLast25] = useState<string[]>([]);
  const [mediadirInputPopup, setMediadirInputPopup] = useState(false);
  const [pos, setPos] = useState(0);

  useEffect(() => {
    APICall.version().then((version: string | null) => {
      if (version !== null) {
        setVersion(version);
      }
    });
    APICall.paused().then((paused: boolean | null) => {
      if (paused !== null) {
        setPaused(paused);
      }
    });
    APICall.media().then((media: string | null) => {
      if (media !== null) {
        setMedia(media);
      }
    });
  }, []);

  useEffect(() => {
    APICall.media().then((media: string | null) => {
      if (media !== null) {
        setMedia(media);
      }
    });
    APICall.getOriginMediafiles().then((res: string[]) => setMediafiles(res));
    APICall.getOriginLast25().then((res: string[]) => setLast25(res));
  }, [mediaToggle]);

  useEffect(() => {
    APICall.pos().then((res: number | null) => {
      if (res !== null) {
        setPos(res);
      }
    });
  }, [media]);

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
            <TimeVolInput rawtime={pos} setRawtime={setPos} type="time" />
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
      <MediaControls
        paused={paused}
        setPaused={setPaused}
        media={media}
        setMediadirInputPopup={setMediadirInputPopup}
      />
    </>
  );
}

export { App };
