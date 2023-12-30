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
import "./CRUDPopup.css";
import "./MediaControls.css";
import { Slider } from "./Slider.tsx";

function MediaControls(props: {
  paused: boolean;
  media: string;
  setPaused: React.Dispatch<React.SetStateAction<boolean>>;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  return (
    <div className="mediacontrols">
      <div className="mediacontrol" onClick={() => APICall.incPos(-30)}>
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
      <div className="mediacontrol" onClick={() => APICall.incPos(30)}>
        <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
          <LongArrowRightDown />
        </IconoirProvider>
      </div>
      <div
        className="mediacontrol"
        onClick={() => {
          const url = window.prompt("URL (YouTube, etc.):", "") || "";
          APICall.setOriginMedia(url);
        }}
      >
        <Chromecast />
      </div>
      <Slider />
      <div className="mediacontrol small">
        Now Playing: {props.media.split("/").slice(-1).join("")}
      </div>
      <div
        className="mediacontrol rightjustify"
        onClick={() => {
          props.setMediadirInputPopup(true);
        }}
      >
        <Book />
      </div>
    </div>
  );
}

export { MediaControls };
