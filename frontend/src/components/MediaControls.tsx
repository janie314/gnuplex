import { APICall } from "../lib/APICall.ts";
import "../App.css";
import "./CRUDPopup.css";
import "./MediaControls.css";

function MediaControls(props: {
  paused: boolean;
  setPaused: React.Dispatch<React.SetStateAction<boolean>>;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  return (
    <div className="mediacontrols">
      <span className="mediacontrol" onClick={() => APICall.incPos(-30)}>
        ⥀
      </span>
      <span
        className="mediacontrol"
        onClick={() =>
          APICall.toggle().then((paused) => props.setPaused(paused))}
      >
        {props.paused ? "⏵" : "⏸"}
      </span>
      <span className="mediacontrol" onClick={() => APICall.incPos(30)}>
        ⥁
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
      <span>Now Playing: Night of the Living Dead</span>
      <span
        className="mediacontrol rightjustify"
        onClick={() => {
          props.setMediadirInputPopup(true);
        }}
      >
        🕮
      </span>
    </div>
  );
}

export { MediaControls };
