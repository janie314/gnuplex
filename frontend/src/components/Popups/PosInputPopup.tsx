import { useEffect, useState } from "react";
import "./Popup.css";
import { API } from "../../lib/API";
import {
  secondsToTimeComponents,
  timeComponentsToSeconds,
} from "../../lib/Helpers";

function PosInputPopup(props: {
  visible: boolean;
  setPosInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  currentPos: number;
  maxPos: number;
  setPos: React.Dispatch<React.SetStateAction<number>>;
}) {
  const [hours, setHours] = useState(0);
  const [minutes, setMinutes] = useState(0);
  const [seconds, setSeconds] = useState(0);

  useEffect(() => {
    const { hours, minutes, seconds } = secondsToTimeComponents(
      props.currentPos,
    );
    setHours(hours);
    setMinutes(minutes);
    setSeconds(seconds);
  }, [props.visible]);

  const getTotalSeconds = () => {
    return timeComponentsToSeconds({
      hours: hours,
      minutes: minutes,
      seconds: seconds,
    });
  };

  if (props.visible) {
    return (
      <div className="popup bg-white dark:bg-stone-800 m-5">
        <div className="flex flex-row mb-2 items-center">
          <span className="header mr-1">Position</span>
          <div className="flex gap-1">
            <input
              type="number"
              className="bg-cyan-800 text-slate-100 border border-black text-sm font-mono font-bold p-1 w-16"
              value={hours}
              min={0}
              onChange={(e) => {
                setHours(Math.max(0, e.target.valueAsNumber || 0));
              }}
              title="Hours"
            />
            <span className="header">:</span>
            <input
              type="number"
              className="bg-cyan-800 text-slate-100 border border-black text-sm font-mono font-bold p-1 w-16"
              value={minutes}
              min={0}
              max={59}
              onChange={(e) => {
                setMinutes(
                  Math.min(59, Math.max(0, e.target.valueAsNumber || 0)),
                );
              }}
              title="Minutes"
            />
            <span className="header">:</span>
            <input
              type="number"
              className="bg-cyan-800 text-slate-100 border border-black text-sm font-mono font-bold p-1 w-16"
              value={seconds}
              min={0}
              max={59}
              onChange={(e) => {
                setSeconds(
                  Math.min(59, Math.max(0, e.target.valueAsNumber || 0)),
                );
              }}
              title="Seconds"
            />
          </div>
        </div>
        <div className="flex flex-row">
          <input
            type="button"
            value="OK"
            className="btn-standard mr-1"
            onClick={() => {
              const totalSeconds = getTotalSeconds();
              if (totalSeconds <= props.maxPos) {
                props.setPos(totalSeconds);
                API.setPos(totalSeconds);
                props.setPosInputPopup(false);
              }
            }}
          />
          <input
            type="button"
            value="Cancel"
            className="btn-standard"
            onClick={() => {
              const totalSecs = props.currentPos;
              const { hours, minutes, seconds } =
                secondsToTimeComponents(totalSecs);
              setHours(hours);
              setMinutes(minutes);
              setSeconds(seconds);
              props.setPosInputPopup(false);
            }}
          />
        </div>
      </div>
    );
  }
  return null;
}

export { PosInputPopup };
