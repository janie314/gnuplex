import { useEffect, useState } from "react";
import { APICall } from "../lib/APICall.ts";
import "../App.css";
import "./CRUDPopup.css";
import "./MediaControls.css";

function MediaControls(props: {
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  return (
    <div className="mediacontrols">
      <span className="mediacontrol">⥀</span>
      <span
        className="mediacontrol"
        onClick={() => APICall.toggle()}
      >
        ⏵
      </span>
      <span className="mediacontrol">⥁</span>
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
