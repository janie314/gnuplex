import { useState } from "react";
import playpause from "../assets/playpause.svg";
import skip from "../assets/skip.svg";
import { API, type SubTrack } from "../lib/API";
import { secondsToTimeComponents } from "../lib/Helpers";
import { PosInputPopup } from "./PosInputPopup";
import { SubSelector } from "./SubSelector";
import { VolInputPopup } from "./VolInputPopup";

function MediaControls(props: {
  mediadirInputPopup: boolean;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  setCastPopup: React.Dispatch<React.SetStateAction<boolean>>;
  setSettingsPopup: React.Dispatch<React.SetStateAction<boolean>>;
  startPos: number;
  pos: number;
  setPos: React.Dispatch<React.SetStateAction<number>>;
  timeRemaining: number;
  vol: number;
  setVol: React.Dispatch<React.SetStateAction<number>>;
  subs: SubTrack[] | null;
  dummyAudio: React.RefObject<HTMLAudioElement | null>;
  skipHook: () => void;
}) {
  const [posInputPopup, setPosInputPopup] = useState(false);
  const [volInputPopup, setVolInputPopup] = useState(false);
  return (
    <div className="flex flex-row flex-wrap items-center justify-center content-baseline p-1">
      <div className="mr-1">
        <button
          type="button"
          className="p-2 w-8 btn-standard"
          onClick={() => {
            API.playpause();
            if (props.dummyAudio.current) {
              if (props.dummyAudio.current.paused) {
                props.dummyAudio.current.play();
              } else {
                props.dummyAudio.current.pause();
              }
            }
          }}
        >
          <img src={playpause} alt="Play/Pause icon" />
        </button>
      </div>
      <div className="mr-2">
        <button
          type="button"
          className="p-2 w-8 btn-standard"
          onClick={() => {
            API.skip();
            props.skipHook();
          }}
        >
          <img src={skip} alt="Skip icon" />
        </button>
      </div>
      <div className="flex flex-col max-w-sm grow p-1">
        <div className="flex flex-row items-center">
          <span className="mx-1 dark:text-white">Pos</span>
          <input
            type="range"
            min={0}
            max={props.startPos + props.timeRemaining}
            value={props.pos}
            className="range range-xs dark:[--range-shdw:#0e7490]"
            onChange={(e) => props.setPos(e.target.valueAsNumber)}
            onMouseUp={() => API.setPos(props.pos)}
            onTouchCancel={() => API.setPos(props.pos)}
          />
          <button
            type="button"
            className="btn-subtle"
            onClick={() => setPosInputPopup(true)}
          >
            {timeFormat(props.pos)}
          </button>
        </div>

        <div className="flex flex-row items-center mt-3">
          <span className="mx-1 dark:text-white">Vol</span>
          <input
            type="range"
            min={0}
            max={120}
            value={props.vol}
            className="range range-xs dark:[--range-shdw:#0e7490]"
            onChange={(e) => props.setVol(e.target.valueAsNumber)}
            onMouseUp={() => API.setVol(props.vol)}
            onTouchCancel={() => API.setVol(props.vol)}
          />
          <button
            type="button"
            className="btn-subtle"
            onClick={() => setVolInputPopup(true)}
          >
            {props.vol}
          </button>
        </div>
      </div>
      <PosInputPopup
        visible={posInputPopup}
        setPosInputPopup={setPosInputPopup}
        currentPos={props.pos}
        maxPos={props.startPos + props.timeRemaining}
        setPos={props.setPos}
      />
      <VolInputPopup
        visible={volInputPopup}
        setVolInputPopup={setVolInputPopup}
        currentVol={props.vol}
        setVol={props.setVol}
      />
      <div className="flex flex-row justify-center mt-3 p-1">
        <SubSelector subs={props.subs} />
        <input
          type="button"
          className="mr-1 btn-standard"
          value="Manage Library"
          onClick={() => {
            props.setMediadirInputPopup(true);
          }}
        />
        <input
          type="button"
          className="btn-standard"
          value="Cast URL"
          onClick={() => props.setCastPopup(true)}
        />
        <input
          type="button"
          className="btn-standard ml-1"
          value="⚙"
          onClick={() => props.setSettingsPopup(true)}
        />
      </div>
    </div>
  );
}

function timeFormat(n: number) {
  const { hours, minutes, seconds } = secondsToTimeComponents(n);
  if (hours === 0) {
    return `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
  }
  return `${hours.toString().padStart(2, "0")}:${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
}

export { MediaControls };
