import { useEffect, useState } from "react";
import { APICall } from "../lib/APICall.ts";
import "../App.css";
import "./CRUDPopup.css";

function MediaControls() {
  return (
    <div className="mediacontrols">
      <input type="button" value="⥀" />
      <input type="button" value="⏵" />
      <input type="button" value="⥁" />
    </div>
  );
}

export { MediaControls };
