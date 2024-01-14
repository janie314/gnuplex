import { APICall } from "../lib/API.ts";
import { Helpers } from "../lib/Helpers.ts";
import {
  IconoirProvider,
  LongArrowLeftUp,
  LongArrowRightDown,
  PauseSolid,
  PlaySolid,
  SoundHigh,
  SoundLow,
  SoundMin,
  SoundOff,
} from "iconoir-react";
import "../App.css";
import { useEffect, useState } from "react";
import { Button, Card, CardBody, Slider } from "@nextui-org/react";

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
  const [vol, setVol] = useState(0);
  const [pos, setPos] = useState(0);
  const [maxPos, setMaxPos] = useState(0);

  useEffect(() => {
    setInterval(() => {
      APICall.mediastate().then((res) => {
        if (res !== null) {
          setPaused(res.paused);
          setMedia(res.media);
          setVol(res.vol);
          setPos(res.pos);
          setMaxPos(res.max_pos);
        }
      });
    }, 2000);
  }, []);

  return (
    <Card isBlurred shadow="sm">
      <CardBody className="flex flex-col space-y-3 justify-center">
        <div className="flex flex-row justify-center space-x-1">
          <span className="nowplaying">
            {media.length === 0 ? "" : clipText(
              "Now Playing: " + media.split("/").slice(-1).join(""),
              50,
            )}
          </span>
        </div>
        <div className="flex flex-row justify-center space-x-1">
          <span>{Helpers.fmtTime(pos)}</span>
          <Slider
            radius="sm"
            maxValue={maxPos}
            value={pos}
            defaultValue={0}
            color="primary"
            onChange={(val: number | number[]) => {
              if (!Array.isArray(val)) {
                setPos(val);
              }
            }}
            onChangeEnd={(val: number | number[]) => {
              if (!Array.isArray(val)) {
                APICall.setPos(val, false);
              }
            }}
          />
        </div>
        <div className="flex flex-row justify-center space-x-1">
          <Button
            isIconOnly
            size="sm"
            color="primary"
            variant="flat"
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
            variant="flat"
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
            variant="flat"
            onClick={() => APICall.setPos(30, true)}
          >
            <IconoirProvider iconProps={{ transform: "rotate(-135)" }}>
              <LongArrowRightDown />
            </IconoirProvider>
          </Button>
        </div>
        <div className="flex flex-row justify-center space-x-1">
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
        <div className="flex flex-row justify-center space-x-1">
          {(vol === 0 || vol === null)
            ? <SoundOff />
            : (vol <= 33
              ? <SoundMin />
              : (vol <= 66 ? <SoundLow /> : <SoundHigh />))}
          <Slider
            radius="sm"
            maxValue={120}
            defaultValue={0}
            value={vol}
            onChange={(val: number | number[]) => {
              if (!Array.isArray(val)) {
                setVol(val);
              }
            }}
            onChangeEnd={(val: number | number[]) => {
              if (!Array.isArray(val)) {
                APICall.setVol(val);
              }
            }}
          />
        </div>
      </CardBody>
    </Card>
  );
}

export { MediaControls };
