import { SyntheticEvent, useEffect, useState } from "react";
import { APICall } from "../lib/API";
import "./TimeInput.css";

const enum TimeUnit {
  secs = 0,
  mins,
  hrs,
}

function TimeInput(props: {
  rawtime: number;
  setRawtime: React.Dispatch<React.SetStateAction<number>>;
}) {
  const secs = props.rawtime % 60;
  const mins = ((props.rawtime - secs) % 3600) / 60;
  const hrs = (props.rawtime - 60 * mins - secs) / 3600;

  function numFmt(num: number, min: number, max: number) {
    const n = Number(num);
    if ((n < min) || (n > max) || (Number.isNaN(n))) {
      return "00";
    }
    return n.toString().slice(-2).padStart(2, "0");
  }

  function setTime(val: string, unit: TimeUnit) {
    const n = Number(val);
    if (unit === TimeUnit.hrs) {
      if (n < 0) {
        props.setRawtime(60 * mins + secs);
      } else {
        props.setRawtime(3600 * n + 60 * mins + secs);
      }
    } else if (unit === TimeUnit.mins) {
      if ((n < 0) || (n > 59)) {
        props.setRawtime(3600 * hrs + secs);
      } else {
        props.setRawtime(3600 * hrs + 60 * n + secs);
      }
    } else {
      if ((n < 0) || (n > 59)) {
        props.setRawtime(3600 * hrs + 60 * mins);
      } else {
        props.setRawtime(3600 * hrs + 60 * mins + n);
      }
    }
  }

  return (
    <div className="time-input">
      <span className="time-input-label">Pos</span>
      <input
        className="time-input-timenum"
        type="text"
        value={Number(hrs).toString().padStart(2, "0")}
        onChange={(e) => setTime(e.target.value, TimeUnit.hrs)}
      />
      <span className="time-input-timesep">:</span>
      <input
        className="time-input-timenum"
        type="text"
        value={Number(mins).toString().padStart(2, "0")}
        onChange={(e) => setTime(e.target.value, TimeUnit.mins)}
      />
      <span className="time-input-timesep">:</span>
      <input
        className="time-input-timenum"
        type="text"
        value={Number(secs).toString().padStart(2, "0")}
        onChange={(e) => setTime(e.target.value, TimeUnit.secs)}
      />
      <input
        type="button"
        className="time-input-button"
        value="Set"
        min={0}
        max={250}
        onClick={(e) => {
          APICall.setOriginPos(props.rawtime);
        }}
      />
    </div>
  );
}

export { TimeInput };
