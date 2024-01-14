import ReactSlider from "react-slider";
import "./PosSlider.css";
import { useDebounce } from "usehooks-ts";
import { useEffect, useState } from "react";
import { APICall } from "../../lib/API.ts";

function fmtTime(rawtime: number): string {
  const secs = Math.floor(rawtime % 60);
  const mins = Math.floor(((rawtime - secs) % 3600) / 60);
  const hrs = Math.floor((rawtime - 60 * mins - secs) / 3600);
  const secs_str = secs.toString().padStart(2, "0");
  const mins_str = mins.toString().padStart(2, "0");
  const hrs_str = hrs === 0 ? "" : (hrs.toString() + ":");
  return `${hrs_str}${mins_str}:${secs_str}`;
}

function PosSlider(props: {
  maxPos: number;
  truePos: number | null;
  setTruePos: React.Dispatch<React.SetStateAction<number | null>>;
}) {
  const [pos, setPos] = useState<number | null>(null);
  const debouncedPos = useDebounce(pos, 500);

  useEffect(() => {
    if (pos !== null) {
      APICall.setPos(pos, false);
    }
  }, [debouncedPos]);

  return (
    <div className="slider">
      <span>
        {props.truePos !== null ? fmtTime(props.truePos) : fmtTime(pos || 0)}
      </span>
      <ReactSlider
        className="horizontal-slider"
        thumbClassName="thumb"
        trackClassName="track"
        value={props.truePos !== null ? props.truePos : (pos || 0)}
        max={props.maxPos}
        onChange={(val: number) => {
          props.setTruePos(null);
          setPos(val);
        }}
      />
    </div>
  );
}

export { PosSlider };
