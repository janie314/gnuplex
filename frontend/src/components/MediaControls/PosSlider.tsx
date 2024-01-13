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

function PosSlider(props: {
  flush: boolean;
}) {
  const [flush, setFlush] = useState(false);
  const [truePos, setTruePos] = useState<number | null>(null);
  const [pos, setPos] = useState<number | null>(null);
  const [maxPos, setMaxPos] = useState(0);
  const debouncedPos = useDebounce(pos, 500);

  useEffect(() => {
    APICall.pos().then((res: PosResponse | null) => {
      if (res !== null) {
        setTruePos(res.pos);
        setMaxPos(res.max_pos);
      }
    });
  }, [props.flush, flush]);

  useEffect(() => {
    if (pos !== null) {
      APICall.setPos(pos, false);
      setFlush(!flush);
    }
  }, [debouncedPos]);

  return (
    <div className="slider">
      <span>{truePos !== null ? fmtTime(truePos) : fmtTime(pos || 0)}</span>
      <ReactSlider
        className="horizontal-slider"
        thumbClassName="thumb"
        trackClassName="track"
        value={truePos !== null ? truePos : pos || 0}
        max={maxPos}
        onChange={(val: number) => {
          setTruePos(null);
          setPos(val);
        }}
      />
    </div>
  );
}

export { PosSlider };
