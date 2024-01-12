import ReactSlider from "react-slider";
import "./PosSlider.css";
import { useDebounce } from "usehooks-ts";
import { useEffect, useState } from "react";
import { APICall, PosResponse } from "../../lib/API.ts";

function fmtTime(rawtime: number): string {
  const secs = Math.floor(rawtime % 60);
  const mins = Math.floor(((rawtime - secs) % 3600) / 60);
  const hrs = Math.floor((rawtime - 60 * mins - secs) / 3600);
  const secs_str = secs.toString().padStart(2, "0");
  const mins_str = mins.toString().padStart(2, "0");
  const hrs_str = hrs === 0 ? "" : (hrs.toString() + ":");
  return `${hrs_str}${mins_str}:${secs_str}`;
}

function PosSlider() {
  const [pos, setPos] = useState(0);
  const [maxPos, setMaxPos] = useState(0);
  const debouncedPos = useDebounce(pos, 500);
  const [updateNum, setUpdateNum] = useState(0); // increment to flush update

  useEffect(() => {
    APICall.pos().then((pos: PosResponse | null) => {
      if (pos !== null) {
        setPos(pos.pos);
        setMaxPos(pos.max_pos);
      }
    });
  }, [updateNum]);

  useEffect(() => {
    if (updateNum !== 0) {
      APICall.setPos(pos, false).then(() => setUpdateNum(updateNum + 1));
    }
  }, [debouncedPos]);

  return (
    <div className="slider">
      <span>{fmtTime(pos)}</span>
      <ReactSlider
        className="horizontal-slider"
        thumbClassName="thumb"
        trackClassName="track"
        value={pos}
        max={maxPos}
        onChange={(val: number) => setPos(pos)}
      />
    </div>
  );
}

export { PosSlider };
