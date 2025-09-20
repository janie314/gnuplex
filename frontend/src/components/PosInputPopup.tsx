import { useState } from "react";
import "./Popup.css";
import { API } from "../lib/API";

function PosInputPopup(props: {
  visible: boolean;
  setPosInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  currentPos: number;
  maxPos: number;
  setPos: React.Dispatch<React.SetStateAction<number>>;
}) {
  const [pos, setPos] = useState(props.currentPos);

  if (props.visible) {
    return (
      <div className="popup bg-white dark:bg-stone-800 m-5">
        <div className="flex flex-row mb-2 items-center">
          <span className="header mr-1">Position</span>
          <input
            type="number"
            className="bg-cyan-800 text-slate-100 border border-black text-sm font-mono font-bold p-1"
            value={pos}
            min={0}
            max={props.maxPos}
            onChange={(e) => {
              setPos(e.target.valueAsNumber);
            }}
          />
        </div>
        <div className="flex flex-row">
          <input
            type="button"
            value="OK"
            className="btn-standard mr-1"
            onClick={() => {
              props.setPos(pos);
              API.setPos(pos);
              props.setPosInputPopup(false);
            }}
          />
          <input
            type="button"
            value="Cancel"
            className="btn-standard"
            onClick={() => {
              setPos(props.currentPos);
              props.setPosInputPopup(false);
            }}
          />
        </div>
      </div>
    );
  }
  return null;
}

export { PosInputPopup };
