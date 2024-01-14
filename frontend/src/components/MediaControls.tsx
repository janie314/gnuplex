import { APICall } from "../lib/API.ts";
import {
  IconoirProvider,
  LongArrowLeftUp,
  LongArrowRightDown,
  PauseSolid,
  PlaySolid,
} from "iconoir-react";
import "../App.css";
import { VolSlider } from "./MediaControls/VolSlider.tsx";
import { PosSlider } from "./MediaControls/PosSlider.tsx";
import { useEffect, useState } from "react";
import { Button, Card, CardBody } from "@nextui-org/react";

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
    <Card isBlurred shadow="sm">
      <CardBody>
        <div className="flex flex-col justify-center">
          <div className="flex flex-row justify-center">
            <span className="nowplaying">
              {media.length === 0 ? "" : clipText(
                "Now Playing: " + media.split("/").slice(-1).join(""),
                50,
              )}
            </span>
          </div>
          <div className="flex flex-row justify-center">
            <PosSlider
              maxPos={maxPos}
              truePos={truePos}
              setTruePos={setTruePos}
            />
          </div>
          <div className="flex flex-row justify-center">
            <Button
              isIconOnly
              size="sm"
              color="primary"
              variant="faded"
              onClick={() => APICall.setPos(-30, true)}
            >
              <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
                <LongArrowLeftUp />
              </IconoirProvider>
            </Button>
            <Button
              isIconOnly
              size="sm"
              color="primary"
              variant="faded"
              onClick={() =>
                APICall.toggle().then((paused: boolean | null) => {
                  if (paused !== null) {
                    setPaused(paused);
                  }
                })}
            >
              {paused ? <PlaySolid /> : <PauseSolid />}
            </Button>
            <Button
              isIconOnly
              size="sm"
              color="primary"
              variant="faded"
              onClick={() => APICall.setPos(30, true)}
            >
              <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
                <LongArrowRightDown />
              </IconoirProvider>
            </Button>
          </div>
          <div className="flex flex-row justify-center">
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
          <div className="flex flex-row justify-center">
            <VolSlider trueVol={trueVol} setTrueVol={setTrueVol} />
          </div>
        </div>
      </CardBody>
    </Card>
  );
}

export { MediaControls };
