import { APICall } from "../lib/API.ts";
import {
  Book,
  IconoirProvider,
  LongArrowLeftUp,
  LongArrowRightDown,
  PauseSolid,
  PlaySolid,
} from "iconoir-react";
import "../App.css";
import "./CRUDPopup.css";
import "./MediaControls.css";

function MediaControls(props: {
  paused: boolean;
  media: string;
  setPaused: React.Dispatch<React.SetStateAction<boolean>>;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  return (
    <div className="mediacontrols">
      <span className="mediacontrol" onClick={() => APICall.incPos(-30)}>
        <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
          <LongArrowLeftUp />
        </IconoirProvider>
      </span>
      <span
        className="mediacontrol"
        onClick={() =>
          APICall.toggle().then((paused: boolean | null) => {
            if (paused !== null) {
              props.setPaused(paused);
            }
          })}
      >
        {props.paused ? <PlaySolid /> : <PauseSolid />}
      </span>
      <span className="mediacontrol" onClick={() => APICall.incPos(30)}>
        <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
          <LongArrowRightDown />
        </IconoirProvider>
      </span>
      <span
        className="mediacontrol small"
        onClick={() => {
          const url = window.prompt("YouTube URL:", "") || "";
          APICall.setOriginMedia(url);
        }}
      >
        Cast
      </span>
      <span className="mediacontrol small">
        Now Playing: {props.media.split("/").slice(-1).join("")}
      </span>
      <span
        className="mediacontrol rightjustify"
        onClick={() => {
          props.setMediadirInputPopup(true);
        }}
      >
        <Book />
      </span>
    </div>
  );
}

export { MediaControls };
