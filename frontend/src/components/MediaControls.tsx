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
import { VolSlider } from "./MediaControls/VolSlider.tsx";
import { useDebounce } from "usehooks-ts";
import { useEffect, useState } from "react";

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
  const [vol, setVol] = useState<number | null>(null);
  const debouncedVol = useDebounce(vol, 500);

  useEffect(() => {
    APICall.vol().then((vol) => {
      if (vol !== null) {
        setVol(vol);
      }
    });
  }, []);

  useEffect(() => {
    if (vol !== null) {
      APICall.setVol(debouncedVol as number);
    }
  }, [debouncedVol]);

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
        <VolSlider vol={vol} debouncedVol={debouncedVol} setVol={setVol} />
      </div>
    </div>
  );
}

export { MediaControls };
