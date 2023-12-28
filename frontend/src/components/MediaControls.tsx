import { useEffect, useState } from "react";
import { APICall } from "../lib/APICall.ts";
import "../App.css";
import "./CRUDPopup.css";
import "./MediaControls.css";

function MediaControls() {
  return (
    <div className="mediacontrols">
      <input type="button" className="mediacontrol" value="⥀" />
      <input
        type="button"
        className="mediacontrol"
        value="⏵"
        onClick={() => APICall.toggle()}
      />
      <input type="button" className="mediacontrol" value="⥁" />
    </div>
  );
}

export { MediaControls };
