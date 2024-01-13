import { APICall } from "../lib/API.ts";
import {
  Book,
  Chromecast,
  IconoirProvider,
  LongArrowLeftUp,
  LongArrowRightDown,
  PauseSolid,
  PlaySolid,
} from "iconoir-react";
import "../App.css";
import "./MediaControls.css";
import { VolSlider } from "./MediaControls/VolSlider.tsx";
import { PosSlider } from "./MediaControls/PosSlider.tsx";
import { useState } from "react";

function clipText(str: string, max: number) {
  if (str.length <= max) {
    return str;
  } else {
    return str.slice(0, 50) + "...";
  }
}

function MediaControls(props: {
  paused: boolean;
  setPaused: React.Dispatch<React.SetStateAction<boolean>>;
  media: string;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const [flushPos, setFlushPos] = useState(false);

  return (
    <div className="mediacontrols">
      <div className="controlrow">
        <span className="nowplaying">
          {clipText(
            "Now Playing: " + props.media.split("/").slice(-1).join(""),
            50,
          )}
        </span>
      </div>
      <div className="controlrow">
        <PosSlider flush={flushPos} />
      </div>
      <div className="controlrow">
        <div
          className="mediacontrol"
          onClick={() =>
            APICall.setPos(-30, true).then(() => setFlushPos(!flushPos))}
        >
          <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
            <LongArrowLeftUp />
          </IconoirProvider>
        </div>
        <div
          className="mediacontrol"
          onClick={() =>
            APICall.toggle().then((paused: boolean | null) => {
              if (paused !== null) {
                props.setPaused(paused);
              }
            })}
        >
          {props.paused ? <PlaySolid /> : <PauseSolid />}
        </div>
        <div
          className="mediacontrol"
          onClick={() =>
            APICall.setPos(30, true).then(() => setFlushPos(!flushPos))}
        >
          <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
            <LongArrowRightDown />
          </IconoirProvider>
        </div>
      </div>
      <div className="controlrow">
        <input
          type="button"
          className="mediacontrol"
          onClick={() => {
            const url = window.prompt("URL (YouTube, etc.):", "") || "";
            APICall.setMedia(url);
          }}
          value="Cast"
        />
        <input
          type="button"
          className="mediacontrol"
          onClick={() => {
            props.setMediadirInputPopup(true);
          }}
          value="Library"
        />
      </div>
      <div className="controlrow">
        <VolSlider />
      </div>
    </div>
  );
}

export { MediaControls };
