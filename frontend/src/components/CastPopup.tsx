import { useState } from "react";
import "./Popup.css";
import { API } from "../lib/API";

function CastPopup(props: {
  visible: boolean;
  setCastPopup: React.Dispatch<React.SetStateAction<boolean>>;
  closeHook: () => void;
}) {
  const [url, setUrl] = useState("");
  const [addToLib, setAddToLib] = useState(false);
  if (props.visible) {
    return (
      <div className="popup">
        <div className="flex flex-row mb-2 items-center">
          <span className="mr-1">URL</span>
          <input
            type="text"
            className="border border-black text-sm font-mono p-1"
            value={url}
            onChange={(e) => {
              API.cast(url, !addToLib);
              setUrl(e.target.value);
            }}
          />
        </div>
        <div className="flex flex-row items-center mb-2">
          <label>
            Add to Library
            <input
              type="checkbox"
              id="boxy"
              className="ml-1"
              checked={addToLib}
              onChange={(e) => setAddToLib(e.target.checked)}
            />
          </label>
        </div>
        <div className="flex flex-row">
          <input
            type="button"
            value="OK"
            className="mr-1 p-1 border border-solid border-black hover:bg-cyan-200"
            onClick={() => {
              API.cast(url, !addToLib);
              setUrl("");
              setAddToLib(false);
              props.setCastPopup(false);
            }}
          />
          <input
            type="button"
            value="Cancel"
            className="p-1 border border-solid border-black hover:bg-cyan-200"
            onClick={() => props.setCastPopup(false)}
          />
        </div>
      </div>
    );
  }
  return null;
}

export { CastPopup };
