import { APICall } from "../lib/APICall";
import "./TimeVolInput.css";

enum TimeUnit {
  secs = 0,
  mins = 1,
  hrs = 2,
}

function TimeVolInput(props: {
  rawtime?: number;
  setRawtime?: React.Dispatch<React.SetStateAction<number>>;
  vol?: number;
  setVol?: React.Dispatch<React.SetStateAction<number>>;
  type: string;
}) {
  if (
    props.type === "time" &&
    props.rawtime !== undefined &&
    props.setRawtime !== undefined
  ) {
    const secs = props.rawtime % 60;
    const mins = ((props.rawtime - secs) % 3600) / 60;
    const hrs = (props.rawtime - 60 * mins - secs) / 3600;

    function setTime(val: string, unit: TimeUnit) {
      const n = Number(val) || 0;
      if (props.rawtime === undefined || props.setRawtime === undefined) {
        return;
      }
      if (unit === TimeUnit.hrs) {
        if (n < 0) {
          props.setRawtime(60 * mins + secs);
        } else {
          props.setRawtime(3600 * n + 60 * mins + secs);
        }
      } else if (unit === TimeUnit.mins) {
        if (n < 0 || n > 59) {
          props.setRawtime(3600 * hrs + secs);
        } else {
          props.setRawtime(3600 * hrs + 60 * n + secs);
        }
      } else {
        if (n < 0 || n > 59) {
          props.setRawtime(3600 * hrs + 60 * mins);
        } else {
          props.setRawtime(3600 * hrs + 60 * mins + n);
        }
      }
    }

    return (
      <div className="timevol-input">
        <span className="timevol-input-label">Pos</span>
        <input
          className="timevol-input-timenum"
          type="text"
          value={Number(hrs).toString().padStart(2, "0")}
          onChange={(e) => setTime(e.target.value, TimeUnit.hrs)}
        />
        <span className="timevol-input-timesep">:</span>
        <input
          className="timevol-input-timenum"
          type="text"
          value={Number(mins).toString().padStart(2, "0")}
          onChange={(e) => setTime(e.target.value, TimeUnit.mins)}
        />
        <span className="timevol-input-timesep">:</span>
        <input
          className="timevol-input-timenum"
          type="text"
          value={Number(secs).toString().padStart(2, "0")}
          onChange={(e) => setTime(e.target.value, TimeUnit.secs)}
        />
        <input
          type="button"
          className="mediacontrol-button"
          value="Set"
          min={0}
          max={250}
          onClick={(e) => {
            APICall.setOriginPos(props.rawtime as number);
          }}
        />
      </div>
    );
  }
  if (
    props.type === "vol" &&
    props.vol !== undefined &&
    props.setVol !== undefined
  ) {
    return (
      <div className="timevol-input">
        <span className="timevol-input-label">Vol</span>
        <input
          className="w-[6ch]"
          type="number"
          value={props.vol}
          min={0}
          max={250}
          onChange={(e) =>
            // @ts-ignore
            props.setVol(Math.min(250, Math.max(e.target.valueAsNumber, 0)))
          }
        />
        <input
          type="button"
          className="mediacontrol-button"
          value="Set"
          onClick={(e) => {
            APICall.setOriginVol(props.vol as number);
          }}
        />
      </div>
    );
  }
  return <></>;
}

export { TimeVolInput };
