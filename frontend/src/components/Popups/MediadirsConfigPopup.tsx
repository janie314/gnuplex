import { useEffect, useState } from "react";
import { API } from "../../lib/API";
import "./Popup.css";
import { WorkingSpinnerTSX } from "../WorkingSpinner";

function MediadirsConfigPopup(props: {
  visible: boolean;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  closeHook: () => void;
}) {
  const [mediadirs, setMediadirs] = useState("");
  const [file_exts, setFileExts] = useState("");
  const [refreshLibraryWorking, setRefreshLibraryWorking] = useState(false);
  const [saveMediadirsWorking, setSaveMediadirsWorking] = useState(false);
  const [saveFileExtsWorking, setSaveFileExtsWorking] = useState(false);

  useEffect(() => {
    API.getMediadirs().then((res) => {
      setMediadirs(res.map((item) => item.Path).join("\n"));
    });
    API.getFileExts().then((res) => {
      setFileExts(res.map((item) => item.Extension).join("\n"));
    });
  }, [props.visible]);

  if (props.visible) {
    return (
      <div className="popup w-full sm:w-full md:w-full lg:w-3/5 bg-white dark:bg-stone-800 dark:text-white">
        <h1 className="header">Media Directories</h1>
        <textarea
          value={mediadirs}
          onChange={(e) => setMediadirs(e.target.value)}
          className="border border-solid border-black p-1 dark:bg-cyan-700 dark:text-white"
          rows={10}
          placeholder="/mnt/externalssd/tv/twilight_zone/eye_of_the_beholder.av1"
        />
        <h1 className="header">Excluded File Extensions</h1>
        <textarea
          value={file_exts}
          onChange={(e) => setFileExts(e.target.value)}
          className="border border-solid border-black p-1 dark:bg-cyan-700 dark:text-white"
          rows={10}
          placeholder=".pdf"
        />
        <div>
          <input
            type="button"
            className="btn-standard"
            value="Save Settings"
            onClick={() => {
              setSaveMediadirsWorking(true);
              setSaveFileExtsWorking(true);
              const arr1 = mediadirs
                .trim()
                .split("\n")
                .filter((line) => !/^\s*$/.test(line))
                .map((line) => line.trim());
              const arr2 = file_exts
                .trim()
                .split("\n")
                .filter((line) => !/^\s*$/.test(line))
                .map((line) => line.trim());
              API.setMediadirs(arr1).then(() => setSaveMediadirsWorking(false));
              API.setFileExts(arr2).then(() => setSaveFileExtsWorking(false));
            }}
          />
          <WorkingSpinnerTSX
            visible={saveFileExtsWorking || saveMediadirsWorking}
          />
        </div>
        <div>
          <input
            type="button"
            className="btn-standard"
            value="Refresh Library"
            onClick={() => {
              setRefreshLibraryWorking(true);
              API.scanLib().then(() => setRefreshLibraryWorking(false));
            }}
          />
          <WorkingSpinnerTSX visible={refreshLibraryWorking} />
        </div>
        <div className="okcancel">
          <input
            type="button"
            className="btn-standard"
            value="OK"
            onClick={() => {
              props.closeHook();
              props.setMediadirInputPopup(false);
            }}
          />
        </div>
      </div>
    );
  }
  return null;
}

export { MediadirsConfigPopup };
