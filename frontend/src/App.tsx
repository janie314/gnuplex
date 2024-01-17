import { useEffect, useState } from "react";
import { APICall } from "./lib/API.ts";
import "./App.css";
import "./index.css";
import { Medialist } from "./components/Medialist.tsx";
import { MediaControls } from "./components/MediaControls.tsx";
import { LibraryMgr } from "./components/LibraryMgr.tsx";

interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

function App() {
  const [version, setVersion] = useState("");
  const [mediafiles, setMediafiles] = useState<string[]>([]);
  const [last25, setLast25] = useState<string[]>([]);
  const [mediadirInputPopup, setMediadirInputPopup] = useState(false);

  useEffect(() => {
    APICall.version().then((version: string | null) => {
      if (version !== null) {
        setVersion(version);
      }
    });
    // TODO this should refresh when you refresh the library
    APICall.mediafiles().then((res: string[]) => setMediafiles(res));
    APICall.last25().then((res: string[]) => setLast25(res));
  }, []);

  return (
    <>
      <div
        className="flex flex-row"
        style={{ opacity: mediadirInputPopup ? "50%" : "100%" }}
      >
        <div className="flex-auto w-1/4 min-w-96">
          <div className="logo-panel">
            <span className="logo">GNUPlex</span>
            <span className="version">{version}</span>
          </div>
          <MediaControls
            setMediadirInputPopup={setMediadirInputPopup}
          />
        </div>
        <div className="flex-auto w-3/4">
          <Medialist medialist={last25} subtitle="Recent" />
          <Medialist
            medialist={mediafiles}
            subtitle="Library"
          />
        </div>
      </div>
      <LibraryMgr
        visible={mediadirInputPopup}
        setMediadirInputPopup={setMediadirInputPopup}
        closeHook={() => {
          setMediadirInputPopup(false);
        }}
      />
    </>
  );
}

export { App };
