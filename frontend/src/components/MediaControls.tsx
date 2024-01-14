import { APICall } from "../lib/API.ts";
import {
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
import { useEffect, useState } from "react";
import { Button } from "@nextui-org/react";

function clipText(str: string, max: number) {
  if (str.length <= max) {
    return str;
  } else {
    return str.slice(0, 50) + "...";
  }
}

function MediaControls(props: {
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const [paused, setPaused] = useState(true);
  const [media, setMedia] = useState("");
  const [maxPos, setMaxPos] = useState<number>(0);
  // for these two, null represents "state being adjusted"
  // (true value not relevant for slider UI display)
  const [trueVol, setTrueVol] = useState<number | null>(0);
  const [truePos, setTruePos] = useState<number | null>(0);

  useEffect(() => {
    setInterval(() => {
      APICall.mediastate().then((res) => {
        if (res !== null) {
          setPaused(res.paused);
          setMedia(res.media);
          setTrueVol(res.vol);
          setTruePos(res.pos);
          setMaxPos(res.max_pos);
        }
      });
    }, 2000);
  }, []);

  return (
    <div className="mediacontrols">
      <div className="controlrow">
        <span className="nowplaying">
          {media.length === 0 ? "" : clipText(
            "Now Playing: " + media.split("/").slice(-1).join(""),
            50,
          )}
        </span>
      </div>
      <div className="controlrow">
        <PosSlider maxPos={maxPos} truePos={truePos} setTruePos={setTruePos} />
      </div>
      <div className="controlrow">
        <div
          className="mediacontrol"
          onClick={() => APICall.setPos(-30, true)}
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
                setPaused(paused);
              }
            })}
        >
          {paused ? <PlaySolid /> : <PauseSolid />}
        </div>
        <div
          className="mediacontrol"
          onClick={() => APICall.setPos(30, true)}
        >
          <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
            <LongArrowRightDown />
          </IconoirProvider>
        </div>
      </div>
      <div className="controlrow">
        <Button
          size="sm"
          color="primary"
          variant="flat"
          className="mediacontrol"
          onClick={() => {
            const url = window.prompt("URL (YouTube, etc.):", "") || "";
            APICall.setMedia(url);
          }}
        >
          Cast
        </Button>
        <Button
          size="sm"
          color="primary"
          variant="flat"
          className="mediacontrol"
          onClick={() => {
            props.setMediadirInputPopup(true);
          }}
        >
          Library
        </Button>
      </div>
      <div className="controlrow">
        <VolSlider trueVol={trueVol} setTrueVol={setTrueVol} />
      </div>
    </div>
  );
}

export { MediaControls };
